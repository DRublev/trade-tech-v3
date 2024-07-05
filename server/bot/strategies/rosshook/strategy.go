package rosshook

import (
	"context"
	"main/bot/candles"
	"main/bot/indicators"
	"main/bot/strategies"
	"main/types"
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	strategies.Config
	// Доступный для торговли баланс
	Balance float64

	// Акция для торговли
	InstrumentID string

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int64

	// Лотность инструмента
	LotSize int64

	// При падении ниже 2 точки минус этот парамер выставим продажу
	StopLoss float64

	// Нужен для Trailing take profit
	// При какой просадке от максимума выставить продажу
	SaveProfit float64
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type RossHookStrategy struct {
	strategies.IStrategy
	strategies.Strategy[Config]

	provider       candles.BaseCandlesProvider
	activityPubSub strategies.IStrategyActivityPubSub

	config Config
	// Канал для стакана
	obCh              *chan *types.Orderbook
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	macd indicators.MacdIndicator

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context

	vault strategies.Vault

	high                 *types.OHLC
	low                  *types.OHLC
	lowForStopLoss       *types.OHLC
	targetGrow           *types.OHLC
	less                 *types.OHLC
	takeProfit           *types.OHLC
	lastBuyPendingCandle *types.OHLC
}

var cancelSwitch context.CancelFunc

func New(provider candles.BaseCandlesProvider, activityPubSub strategies.IStrategyActivityPubSub) *RossHookStrategy {
	inst := &RossHookStrategy{}
	inst.provider = provider
	inst.activityPubSub = activityPubSub
	inst.toPlaceOrders = make(chan *types.PlaceOrder)
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())
	return inst
}

var l *log.Entry

func (s *RossHookStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "rosshook",
		"instrument": (*config)["InstrumentID"],
	})

	// Обнуляем, потому что при параллельном запуске тестов, значения запоминаются
	s.high = nil
	s.low = nil
	s.targetGrow = nil
	s.less = nil
	s.takeProfit = nil
	s.lastBuyPendingCandle = nil

	err := s.SetConfig(*config)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}

	s.vault = *strategies.NewVault(s.Config.LotSize, s.Config.Balance)

	l.Infof("Starting strategy with config: %v", s.Config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.Config.InstrumentID, now, now, false)
	if err != nil {
		l.Errorf("Failed to get candles channel: %v", err)
		return false, err
	}

	s.toPlaceOrders = *ordersToPlaceCh
	go func() {
		l.Info("Start listening latest candles")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case candle, ok := <-*ch:
				l.Trace("New candle")
				if !ok {
					l.Trace("Candles channel closed")
					return
				}

				go s.OnCandle(candle)
			}
		}
	}()

	go s.OnOrderSateChangeSubscribe(s.stopCtx, orderStateChangeCh, s.vault.OnOrderSateChange)

	return true, nil
}

func (s *RossHookStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *RossHookStrategy) mapAndSendState() {
	if s.high != nil {
		s.activityPubSub.Track("p1", "point", strategies.PointActivityValue[float64, time.Time]{
			X: s.high.High.Float(),
			Y: s.high.Time,
		})
	}
	if s.low != nil {
		s.activityPubSub.Track("p2", "point", strategies.PointActivityValue[float64, time.Time]{
			X: s.low.Low.Float(),
			Y: s.low.Time,
		})
	}
	if s.targetGrow != nil {
		s.activityPubSub.Track("p3", "point", strategies.PointActivityValue[float64, time.Time]{
			X: s.targetGrow.High.Float(),
			Y: s.targetGrow.Time,
		})
	}
	if s.less != nil {
		s.activityPubSub.Track("p4", "point", strategies.PointActivityValue[float64, time.Time]{
			X: s.less.Low.Float(),
			Y: s.less.Time,
		})
		s.activityPubSub.Track("buyAt", "level", s.targetGrow.High.Float())
	}
	if s.lowForStopLoss != nil {
		s.activityPubSub.Track("stopLoss", "level", s.lowForStopLoss.Low.Float()-s.Config.StopLoss)
	}
	if s.takeProfit != nil {
		s.activityPubSub.Track("takeProfit", "level", s.takeProfit.High.Float()-float64(s.Config.SaveProfit))
	}
}

func (s *RossHookStrategy) OnCandle(c types.OHLC) {
	defer s.mapAndSendState()

	if !s.isBuying.value {
		go s.watchBuySignal(c)
	}
	if !s.isSelling.value {
		go s.watchSellSignal(c)
	}
}

var utc, _ = time.LoadLocation("UTC")

func isSameTF(candidate types.OHLC, toCompare types.OHLC) bool {
	cH, cM, _ := candidate.Time.In(utc).Clock()
	nH, nM, _ := toCompare.Time.In(utc).Clock()
	return cH == nH && cM == nM
}

func isCompletedTF(candle types.OHLC) bool {
	cH, cM, _ := candle.LastTradeTS.In(utc).Clock()
	nH, nM, _ := time.Now().In(utc).Clock()
	return nH >= cH && nM > cM
}

// true if candidate > toCompare
func isGreaterTF(candidate types.OHLC, toCompare types.OHLC) bool {
	cH, cM, _ := candidate.LastTradeTS.In(utc).Clock()
	nH, nM, _ := toCompare.LastTradeTS.In(utc).Clock()
	return cH >= nH && cM > nM
}

func (s *RossHookStrategy) watchBuySignal(c types.OHLC) {
	// Закрываем висящие на заявку покупки при поступлении новой свечи - мы проебали момент
	// Однако, если текущая цена равна цене, по которой выставляли заявку, есть шанс что еще исполнится
	if s.lastBuyPendingCandle != nil && c.High.Float() > s.lastBuyPendingCandle.High.Float() && isGreaterTF(c, *s.lastBuyPendingCandle) {
		s.closePendingBuys()
	}

	isLowDiffTF := s.high != nil && isGreaterTF(c, *s.high)
	isTargetGrowDiffTF := s.low != nil && isGreaterTF(c, *s.low)
	isLessDiffTF := s.targetGrow != nil && isGreaterTF(c, *s.targetGrow)

	if s.high == nil || s.high.High.Float() <= c.High.Float() {
		s.high = &c
		s.low = nil
		s.targetGrow = nil
		s.less = nil
		s.takeProfit = nil
		l.Infof("Set point 1. high: %v;", s.high.High.Float())
		return
	} else if (s.low == nil || s.low.Low.Float() >= c.Low.Float()) && isLowDiffTF {
		s.low = &c
		s.targetGrow = nil
		s.less = nil
		l.Infof("Set point 2. high: %v; low: %v;", s.high.High.Float(), s.low.Low.Float())
		return
	} else if (s.targetGrow == nil || (s.targetGrow.High.Float() < c.High.Float() && s.less == nil)) && isTargetGrowDiffTF {
		s.targetGrow = &c
		s.less = nil
		l.Infof("Set point 3. high: %v; low: %v; targetGrow: %v;", s.high.High.Float(), s.low.Low.Float(), s.targetGrow.High.Float())
		return
	} else if s.targetGrow != nil && ((s.less == nil && c.Low.Float() < s.targetGrow.High.Float()) ||
		(s.less != nil && s.less.Low.Float() > c.Low.Float() && s.less.Low.Float() < s.targetGrow.High.Float())) &&
		isLessDiffTF {
		s.less = &c
		l.Infof("Set point 4. high: %v; low: %v; targetGrow: %v; less: %v;", s.high.High.Float(), s.low.Low.Float(), s.targetGrow.High.Float(), s.less.Low.Float())
		return
	} else {
		l.Tracef("None of price l %v; h: %v; t: %v; (high: %v; low: %v; targetGrow: %v; less: %v;)", c.Low.Float(), c.High.Float(), c.LastTradeTS.Local(), s.high, s.low, s.targetGrow, s.less)
	}

	if s.high != nil && s.low != nil && s.targetGrow != nil && s.less != nil {
		if s.targetGrow.High.Float() <= c.High.Float() {
			s.takeProfit = &types.OHLC{
				Open:        s.high.Open,
				High:        s.high.High,
				Low:         s.high.Low,
				Close:       s.high.Close,
				Time:        s.high.Time,
				LastTradeTS: s.high.LastTradeTS,
			}
			go s.buy(*&types.OHLC{
				Open:        s.targetGrow.Open,
				High:        s.targetGrow.High,
				Low:         s.targetGrow.Low,
				Close:       s.targetGrow.High, // Точка входа
				Time:        s.targetGrow.Time,
				LastTradeTS: s.targetGrow.LastTradeTS,
			})
		}
	}
}

func (s *RossHookStrategy) watchSellSignal(c types.OHLC) {
	// Stop-loss
	if s.lastBuyPendingCandle != nil && s.lowForStopLoss != nil &&
		isGreaterTF(c, *s.lastBuyPendingCandle) &&
		s.lowForStopLoss.Low.Float()-s.Config.StopLoss >= c.Close.Float() {
		l.Infof("Price reached stop-loss (low: %v; loss: %v; current: %v)", s.lowForStopLoss.Low.Float(), s.Config.StopLoss, c.Close.Float())
		go s.sell(c)
		return
	}

	if s.takeProfit == nil {
		return
	}
	if c.High.Float() > s.takeProfit.High.Float() {
		l.Infof("Updating take-profit high %v", c.High.Float())
		// Копируем свечу. Подозрение на баг, что свеча перезаписывается следующей, поэтомуне прокидываем просто &c
		s.takeProfit = &types.OHLC{
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.Close,
			Time:  c.Time,
		}
		return
	}
	if s.lastBuyPendingCandle == nil {
		return
	}
	if s.lastBuyPendingCandle != nil &&
		isCompletedTF(*s.takeProfit) &&
		// Не сразу после покупки
		isGreaterTF(c, *s.lastBuyPendingCandle) &&
		// Не ниже, чем цена покупки - за это отвечает стоп
		s.lastBuyPendingCandle.Close.Float() < c.High.Float() &&
		// Цена перестала расти и упала на Х - продаем
		s.takeProfit.High.Float()-float64(s.Config.SaveProfit) >= c.High.Float() {
		l.Infof("Price reached take-profit (take: %v; save: %v; current: %v)", s.takeProfit.Close.Float(), s.Config.SaveProfit, c.High.Float())
		go s.sell(types.OHLC{
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.High,
			Time:  c.Time,
		})
	}
}

func (s *RossHookStrategy) sell(c types.OHLC) {
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	if s.vault.HoldingShares-s.vault.PendingSellShares == 0 {
		l.WithField("state", s.vault.String()).Info("Nothing to sell")
		return
	}

	ok := s.isSelling.TryLock()
	if !ok {
		l.Warn("isSelling mutex cannot be locked")

		return
	}
	defer s.isSelling.Unlock()

	l.Trace("Set is selling")
	s.isSelling.value = true

	price := c.Close.Float()
	order := &types.PlaceOrder{
		InstrumentID: s.Config.InstrumentID,
		Quantity:     int64(s.vault.HoldingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}
	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	s.vault.PendingSellShares += s.vault.HoldingShares
	l.WithField("state", s.vault).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	s.high = nil
	s.low = nil
	s.lowForStopLoss = nil
	s.targetGrow = nil
	s.less = nil
	s.lastBuyPendingCandle = nil
	s.takeProfit = nil

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) buy(c types.OHLC) {
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev candle item for buy")
		return
	}

	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney

	canBuySharesAmount := int64(math.Abs(leftBalance / (c.Close.Float() * float64(s.Config.LotSize))))
	if canBuySharesAmount <= 0 {
		l.WithField("state", s.vault).Trace("Can buy 0 shares")
		return
	}

	ok := s.isBuying.TryLock()
	if !ok {
		l.Warn("IsBuiyng mutex cannot be locked")
		return
	}
	defer s.isBuying.Unlock()

	l.Trace("Set is buiyng")
	s.isBuying.value = true
	if canBuySharesAmount > s.Config.MaxSharesToHold {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = s.Config.MaxSharesToHold
	}

	order := &types.PlaceOrder{
		InstrumentID: s.Config.InstrumentID,
		Quantity:     int64(canBuySharesAmount),
		Direction:    types.Buy,
		Price:        types.Price(c.Close.Float()),
	}

	l.Infof("Order to place: %v", order)

	s.vault.PendingBuyShares += int64(canBuySharesAmount)
	s.vault.NotConfirmedBlockedMoney += float64(canBuySharesAmount) * c.Close.Float()
	s.vault.LastBuyPrice = c.Close.Float()
	s.lastBuyPendingCandle = &c
	l.Infof("Last buy pending candle TF: %v;", s.lastBuyPendingCandle.LastTradeTS)
	s.lowForStopLoss = s.low
	l.WithField("state", s.vault).Trace("State updated after place buy order")

	s.targetGrow = nil
	s.less = nil

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) closePendingBuys() {
	l.Infof("Pending buys: %v", len(s.vault.PlacedBuyOrders))
	for _, order := range s.vault.PlacedBuyOrders {
		o := &types.PlaceOrder{
			InstrumentID: s.Config.InstrumentID,
			CancelOrder:  order.ID,
			Direction:    1,
			Quantity:     int64(order.LotsRequested - order.LotsExecuted),
		}
		s.toPlaceOrders <- o
	}
}
