package rosshook

import (
	"context"
	"encoding/json"
	"fmt"
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
	strategies.Strategy

	provider candles.BaseCandlesProvider
	config   Config
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
	targetGrow           *types.OHLC
	less                 *types.OHLC
	takeProfit           *types.OHLC
	lastBuyPendingCandle *types.OHLC
}

var cancelSwitch context.CancelFunc

func New(provider candles.BaseCandlesProvider) *RossHookStrategy {
	inst := &RossHookStrategy{}
	inst.provider = provider
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

	var res Config

	// TODO: Вынести в сущность конфига стратегии
	bts, err := json.Marshal(config)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}

	err = json.Unmarshal(bts, &res)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}
	s.config = res

	s.vault = *strategies.NewVault(s.config.LotSize, s.config.Balance)

	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.config.InstrumentID, now, now, false)
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

	go func() {
		l.Info("Start listening for orders")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case state, ok := <-*orderStateChangeCh:
				if !ok {
					l.Warn("Orders state channel closed")
					return
				}
				s.vault.OnOrderSateChange(state)
			}
		}
	}()

	s.nextOrderCooldown = time.NewTimer(time.Duration(0) * time.Millisecond)

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

func (s *RossHookStrategy) OnCandle(c types.OHLC) {
	if !s.isSelling.value {
		s.watchBuySignal(c)
	}
	if !s.isBuying.value {
		s.watchSellSignal(c)
	}
}

func isSameTF(candidate types.OHLC, toCompare types.OHLC) bool {
	cH, cM, _ := candidate.Time.Local().Clock()
	nH, nM, _ := toCompare.Time.Local().Clock()
	return cH == nH && cM == nM
}

func isCompletedTF(candle types.OHLC) bool {
	cH, cM, _ := candle.LastTradeTS.Clock()
	nH, nM, _ := time.Now().Clock()
	return nH >= cH && nM > cM
}
func isGreaterTF(candidate types.OHLC, toCompare types.OHLC) bool {
	cH, cM, _ := candidate.LastTradeTS.Local().Clock()
	nH, nM, _ := toCompare.LastTradeTS.Local().Clock()
	return cH >= nH && cM > nM
}

func (s *RossHookStrategy) watchBuySignal(c types.OHLC) {
	// Закрываем висящие на заявку покупки при поступлении новой свечи - мы проебали момент
	// Однако, если текущая цена равна цене, по которой выставляли заявку, есть шанс что еще исполнится
	if s.lastBuyPendingCandle != nil && c.High.Float() > s.lastBuyPendingCandle.High.Float() {
		s.closePendingBuys()
	}

	isLowDifTF := s.high != nil && isGreaterTF(c, *s.high)
	isTargetGrowDiffTF := isLowDifTF && s.low != nil && isGreaterTF(c, *s.low)
	isLessDiffTF := isTargetGrowDiffTF && s.targetGrow != nil && isGreaterTF(c, *s.targetGrow)

	if s.high == nil || s.high.High.Float() <= c.High.Float() {
		s.high = &c
		s.low = nil
		s.targetGrow = nil
		s.less = nil
		l.Infof("Set point 1. high: %v;", s.high.High.Float())
		return
	} else if (s.low == nil || s.low.Low.Float() >= c.Low.Float()) && isLowDifTF {
		s.low = &c
		s.targetGrow = nil
		s.less = nil
		l.Infof("Set point 2. high: %v; low: %v;", s.high.High.Float(), s.low.Low.Float())
		return
	} else if (s.targetGrow == nil || (s.targetGrow.High.Float() < c.High.Float() && s.less == nil)) && isTargetGrowDiffTF {
		s.targetGrow = &c
		s.less = nil
		s.takeProfit = &types.OHLC{
			Open:  s.high.Open,
			High:  s.high.High,
			Low:   s.high.Low,
			Close: s.high.Close,
			Time:  s.high.Time,
		}
		l.Infof("Set point 3. high: %v; low: %v; targetGrow: %v;", s.high.High.Float(), s.low.Low.Float(), s.targetGrow.High.Float())
		return
	} else if (s.less == nil || (s.less.Low.Float() > c.Low.Float() && s.less.Low.Float() < s.targetGrow.High.Float())) && isLessDiffTF {
		s.less = &c
		l.Infof("Set point 4. high: %v; low: %v; targetGrow: %v; less: %v;", s.high.High.Float(), s.low.Low.Float(), s.targetGrow.High.Float(), s.less.Low.Float())
		return
	} else {
		// l.Infof("None of price l %v; h: %v; t: %v; (high: %v; low: %v; targetGrow: %v; less: %v;)", c.Low.Float(), c.High.Float(), c.LastTradeTS.Local(), high, low, targetGrow, less)
	}

	if s.high != nil && s.low != nil && s.targetGrow != nil && s.less != nil {
		if s.targetGrow.High.Float()+0.0001 <= c.High.Float() {
			go s.buy(*&types.OHLC{
				Open:  s.targetGrow.Open,
				High:  s.targetGrow.High,
				Low:   s.targetGrow.Low,
				Close: s.targetGrow.High, // Точка входа
				Time:  s.targetGrow.Time,
			})
		}
	}
}

func (s *RossHookStrategy) watchSellSignal(c types.OHLC) {
	// Stop-loss
	if s.lastBuyPendingCandle != nil && s.low != nil &&
		isGreaterTF(c, *s.lastBuyPendingCandle) &&
		s.low.Low.Float()-s.config.StopLoss >= c.Close.Float() {
		l.Infof("Price reached stop-loss (low: %v; loss: %v; current: %v)", s.low.Low.Float(), s.config.StopLoss, c.Close.Float())
		go s.sell(c)
		return
	}

	if s.takeProfit == nil {
		return
	}
	if c.Close.Float() > s.takeProfit.Close.Float() {
		l.Infof("Updating take-profit high %v", c.Close.Float())
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
	if s.less != nil && s.lastBuyPendingCandle != nil && s.takeProfit.Close.Float()-float64(s.config.SaveProfit) >= c.Close.Float() && s.takeProfit.Close.Float() > s.lastBuyPendingCandle.High.Float() {
		l.Infof("Price reached take-profit (take: %v; save: %v; current: %v)", s.takeProfit.Close.Float(), s.config.SaveProfit, c.Close.Float())
		go s.sell(types.OHLC{
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.Close,
			Time:  c.Time,
		})
	} else if s.takeProfit.Close.Float()-float64(s.config.SaveProfit) <= c.Close.Float() {
		l.Infof("Price going up (take: %v; save: %v; current: %v)", s.takeProfit.Close.Float(), s.config.SaveProfit, c.Close.Float())
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
		InstrumentID: s.config.InstrumentID,
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
		// return
	}

	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney

	canBuySharesAmount := int64(math.Abs(leftBalance / (c.Close.Float() * float64(s.config.LotSize))))
	fmt.Printf("266 strategy lotSize %v; left balance %v; can buy %v \n", s.config.LotSize, leftBalance, canBuySharesAmount)
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
	if canBuySharesAmount > s.config.MaxSharesToHold {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = s.config.MaxSharesToHold
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentID,
		Quantity:     int64(canBuySharesAmount),
		Direction:    types.Buy,
		Price:        types.Price(c.Close.Float()),
	}

	l.Infof("Order to place: %v", order)

	s.vault.PendingBuyShares += int64(canBuySharesAmount)
	s.vault.NotConfirmedBlockedMoney += float64(canBuySharesAmount) * c.Close.Float()
	s.vault.LastBuyPrice = c.Close.Float()
	s.lastBuyPendingCandle = &c
	l.WithField("state", s.vault).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) closePendingBuys() {
	l.Infof("Pending buys: %v", len(s.vault.PlacedBuyOrders))
	for _, order := range s.vault.PlacedBuyOrders {
		o := &types.PlaceOrder{
			InstrumentID: s.config.InstrumentID,
			CancelOrder:  order.ID,
			Direction:    1,
			Quantity:     int64(order.LotsRequested),
		}
		s.toPlaceOrders <- o
	}
}
