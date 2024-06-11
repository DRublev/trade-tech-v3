package rosshook

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"main/bot/candles"
	"main/bot/indicators"
	"main/bot/strategies"
	"main/types"
	"sync"
	"time"
)

type Config struct {
	strategies.Config
	// Доступный для торговли баланс
	Balance float32

	// Акция для торговли
	InstrumentID string

	MinProfit float32

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int32

	// Лотность инструмента
	LotSize int32

	// Если цена пошла ниже чем цена покупки - StopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	StopLossAfter float32
}

type State struct {
	// Оставшееся количество денег
	leftBalance float32

	// Сумма, которая должна списаться при выставлении ордера на покупку
	// Инкрементим когда хотим выставить бай ордер
	// Декрементим когда закрываем бай ордер
	notConfirmedBlockedMoney float32

	// Количество акций, купленных на данный момент
	holdingShares int32

	// Количество акций, на которое выставлены ордера на покупку
	pendingBuyShares int32

	// Количество акций, на которое выставлены ордера на продажу
	pendingSellShares int32

	lastBuyPrice float32

	placedOrders []types.OrderExecutionState
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type RossHookStrategy struct {
	strategies.IStrategy
	strategies.Strategy
	config Config
	// Канал для стакана
	obCh              *chan *types.Orderbook
	state             strategies.StrategyState[State]
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	macd indicators.MacdIndicator

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context
}

var cancelSwitch context.CancelFunc

func New() *RossHookStrategy {
	inst := &RossHookStrategy{}
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
		"strategy":   "ross_hook",
		"instrument": (*config)["InstrumentID"],
	})

	var res Config

	// TODO: Вынести в сущность конфига стратегии
	bts, err := json.Marshal(config)
	if err != nil {
		l.Error("Error parsing config %v", err)
		return false, err
	}

	err = json.Unmarshal(bts, &res)
	if err != nil {
		l.Error("Error parsing config %v", err)
		return false, err
	}
	s.config = res

	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	candlesProvider := candles.NewProvider()
	now := time.Now()

	ch, err := candlesProvider.GetOrCreate(s.config.InstrumentID, now, now)
	if err != nil {
		l.Errorf("Failed to get candles channel: %v", err)
		return false, err
	}

	go func(c *chan types.OHLC) {
		l.Info("Start listening latest candles")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case candle, ok := <-*c:
				l.Trace("New candle")
				if !ok {
					l.Trace("Candles channel closed")
					return
				}

				go s.onCandle(candle)
			}
		}
	}(ch)

	l.Trace("Setting state to empty")
	// Заполняем изначальное состояние
	s.state = strategies.StrategyState[State]{}
	err = s.state.Set(State{
		holdingShares:            0,
		pendingBuyShares:         0,
		pendingSellShares:        0,
		leftBalance:              s.config.Balance,
		notConfirmedBlockedMoney: 0,
		lastBuyPrice:             0,
	})
	if err != nil {
		l.Errorf("Failed to set strategy initial state: %v", err)
		return false, err
	}

	s.nextOrderCooldown = time.NewTimer(time.Duration(0) * time.Millisecond)

	return true, nil
}
func (s *RossHookStrategy) onCandle(c types.OHLC) {
	wg := &sync.WaitGroup{}

	close := c.Close.Float()
	s.macd.Update(close)
	allMacd, allSignals := s.macd.Get()

	if len(allMacd) < 2 {
		l.Infof("Not enough data for macd")
		return
	}

	state := s.state.Get()

	latestMacd := allMacd[len(allMacd)-1]

	latestSignal := allSignals[len(allSignals)-1]

	l.Tracef("Updating signal with new values: signal %v; macd: %v", latestSignal, latestMacd)

	s.state.Set(*state)

	wg.Add(1)
	go s.buy(wg, c)
	wg.Add(1)
	go s.sell(wg, c)

	wg.Wait()
}
func (s *RossHookStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *RossHookStrategy) sell(wg *sync.WaitGroup, c types.OHLC) {
	defer wg.Done()
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	state := *s.state.Get()

	// TODO: Выставить ордер на продажу

}

func (s *RossHookStrategy) buy(wg *sync.WaitGroup, c types.OHLC) {
	defer wg.Done()
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	state := s.state.Get()

	// TODO: Выставить ордер на покупку
}
