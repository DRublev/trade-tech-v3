package spread

import (
	"context"
	"main/bot/orderbook"
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

	MinProfit float64

	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
	NextOrderCooldownMs int64

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int64

	// Лотность инструмента
	LotSize int64

	// Если цена пошла ниже чем цена покупки - StopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	StopLossAfter float64
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type SpreadStrategy struct {
	strategies.IStrategy
	strategies.Strategy[Config]

	provider       orderbook.BaseOrderbookProvider
	activityPubSub strategies.IStrategyActivityPubSub
	vault          strategies.Vault

	// Канал для стакана
	obCh              *chan *types.Orderbook
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context
}

var cancelSwitch context.CancelFunc

func New(provider orderbook.BaseOrderbookProvider, activityPubSub strategies.IStrategyActivityPubSub) *SpreadStrategy {
	inst := &SpreadStrategy{}
	inst.provider = provider
	inst.activityPubSub = activityPubSub
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())

	return inst
}

var l *log.Entry

func (s *SpreadStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "spread",
		"instrument": (*config)["InstrumentID"],
	})

	err := s.SetConfig(*config)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}
	l.Infof("Starting strategy with config: %v", s.Config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting orderbook channel")
	ch, err := s.provider.GetOrCreate(s.Config.InstrumentID)
	if err != nil {
		l.Errorf("Failed to get orderbook channel: %v", err)
		return false, err
	}

	s.vault = *strategies.NewVault(s.Config.LotSize, s.Config.Balance)

	s.nextOrderCooldown = time.NewTimer(time.Duration(0) * time.Millisecond)
	s.toPlaceOrders = *ordersToPlaceCh

	// стакан!
	s.obCh = ch

	go func() {
		l.Info("Start listening changes in orderbook")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case ob, ok := <-*ch:
				l.Trace("New orderbook change")
				if !ok {
					l.Trace("Orderbook channel closed")
					return
				}

				go s.onOrderbook(ob)
			}
		}
	}()

	go s.OnOrderSateChangeSubscribe(s.stopCtx, orderStateChangeCh, s.vault.OnOrderSateChange)

	return true, nil
}

func (s *SpreadStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *SpreadStrategy) onOrderbook(ob *types.Orderbook) {
	wg := &sync.WaitGroup{}

	if s.vault.PendingBuyShares == 0 {
		s.activityPubSub.Track("buyAt", "level", strategies.LevelActivityValue{
			DeleteFlag: true,
		})
	}
	if len(s.vault.PlacedBuyOrders) > 0 {
		lastBuyOrder := s.vault.PlacedBuyOrders[len(s.vault.PlacedBuyOrders)-1]

		s.activityPubSub.Track("buyAt", "level", strategies.LevelActivityValue{
			Level: s.vault.LastBuyPrice / float64(s.Config.LotSize),
			Text:  "buy at",
		})

		if lastBuyOrder.LotsExecuted == lastBuyOrder.LotsRequested {
			price := lastBuyOrder.ExecutedOrderPrice / float64(lastBuyOrder.LotsExecuted)

			s.activityPubSub.Track("stopLoss", "level", strategies.LevelActivityValue{
				DeleteFlag: true,
			})
			s.activityPubSub.Track("takeProfit", "level", strategies.LevelActivityValue{
				DeleteFlag: true,
			})
			s.activityPubSub.Track("stopLoss", "level", strategies.LevelActivityValue{
				Level: price - s.Config.StopLossAfter,
				Text:  "stop-loss",
			})
			s.activityPubSub.Track("takeProfit", "level", strategies.LevelActivityValue{
				Level: price + s.Config.MinProfit,
				Text:  "take-profit",
			})
		}
	}

	wg.Add(1)
	go s.checkForRottenBuys(wg, ob)
	wg.Add(1)
	go s.checkForRottenSells(wg, ob)

	wg.Add(1)
	go s.buy(wg, ob)
	wg.Add(1)
	go s.sell(wg, ob)

	wg.Wait()
}

// TODO: Перенести в отдельный файл
func (s *SpreadStrategy) buy(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	isHoldingMaxShares := s.vault.HoldingShares+s.vault.PendingBuyShares >= s.Config.MaxSharesToHold
	if isHoldingMaxShares {
		l.WithField("state", s.vault.String()).Tracef("Cannot buy, holding max shares")
		return
	}

	// Аукцион закрытия, только заявки на продажу
	if len(ob.Bids) == 0 {
		l.Trace("No bids")
		return
	}

	minBuyPrice := float64(ob.Bids[0].Price)
	l.Tracef("Min buy price: %v", minBuyPrice)
	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney
	if leftBalance < (minBuyPrice * float64(s.Config.LotSize)) {
		l.WithField("state", s.vault.String()).Tracef("Not enough money")
		return
	}

	canBuySharesAmount := int64(math.Abs(leftBalance / (minBuyPrice) * float64(s.Config.LotSize)))
	l.Tracef("First bid price: %v; Left money: %v; Can buy %v shares\n", minBuyPrice, leftBalance, canBuySharesAmount)
	if canBuySharesAmount <= 0 {
		l.WithField("state", s.vault.String()).Trace("Can buy 0 shares")
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
		Quantity:     canBuySharesAmount,
		Price:        types.Price(minBuyPrice),
		Direction:    types.Buy,
	}
	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	priceForAllShares := float64(s.Config.LotSize) * float64(order.Price)
	s.vault.PendingBuyShares += canBuySharesAmount
	s.vault.NotConfirmedBlockedMoney += float64(canBuySharesAmount) * priceForAllShares
	s.vault.LastBuyPrice = priceForAllShares
	l.WithField("state", s.vault.String()).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.activityPubSub.Track("buyAt", "level", strategies.LevelActivityValue{
			Level: float64(order.Price),
			Text:  "buy at",
		})

	s.toPlaceOrders <- order
}

func (s *SpreadStrategy) sell(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	if s.vault.HoldingShares-s.vault.PendingSellShares < 0 {
		l.WithField("state", s.vault.String()).Warn("Holding less than 0 shares")
		return
	}

	minAskPrice := float64(ob.Asks[0].Price)
	l.Tracef("Min ask price %v", minAskPrice)

	lastPrice := s.vault.LastBuyPrice / float64(s.Config.LotSize)
	isGoodPrice := minAskPrice-lastPrice >= s.Config.MinProfit
	hasStopLossBroken := s.vault.HoldingShares-s.vault.PendingSellShares > 0 &&
		s.Config.StopLossAfter != 0 &&
		float64(ob.Bids[0].Price) <= lastPrice-s.Config.StopLossAfter

	if s.vault.HoldingShares-s.vault.PendingSellShares == 0 && !hasStopLossBroken {
		l.WithField("state", s.vault.String()).Trace("Nothing to sell")
		return
	}

	shouldMakeSell := isGoodPrice || hasStopLossBroken
	if !shouldMakeSell {
		l.WithField("lastBuyPrice", lastPrice).Tracef("Not a good deal")
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

	price := minAskPrice
	l.Tracef("Selling price: %v", price)
	if hasStopLossBroken {
		price = float64(ob.Bids[0].Price)
		l.WithFields(log.Fields{
			"lastBuyPrice":  lastPrice,
			"stopLoss":      s.Config.StopLossAfter,
			"stopLossPrice": lastPrice - s.Config.StopLossAfter,
		}).Info("Stop loss broken")
	}

	order := &types.PlaceOrder{
		InstrumentID: s.Config.InstrumentID,
		Quantity:     int64(s.vault.HoldingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}

	if hasStopLossBroken {
		for _, o := range s.vault.PlacedSellOrders {
			if o.Status != types.New {
				order.CancelOrder = o.ID
				break
			}
		}
	}

	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	s.vault.PendingSellShares += s.vault.HoldingShares
	l.WithField("state", s.vault.String()).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	s.toPlaceOrders <- order

}

func (s *SpreadStrategy) checkForRottenBuys(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	// TODO: Чекать неаткуальные выставленные ордера и отменять их
	// TODO: Сбрасывать lastBuyPrice на предыдущий, если закрываем какой то бай ордер
}

func (s *SpreadStrategy) checkForRottenSells(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	// TODO: Чекать неаткуальные выставленные ордера и отменять их
}
