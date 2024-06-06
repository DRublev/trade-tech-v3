package macd

import (
	"context"
	"encoding/json"
	"fmt"
	"main/bot/candles"
	"main/bot/indicators"
	"main/bot/strategies"
	"main/types"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	strategies.Config
	// Доступный для торговли баланс
	Balance float32

	// Акция для торговли
	InstrumentID string

	MinProfit float32

	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
	NextOrderCooldownMs int32

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

	// Храним пару последних значений индикаторов
	// Чтобы отслеживать их пересечения
	latestSignals []float64
	latestMacd    []float64

	placedOrders []types.OrderExecutionState
}

func (s *State) String() string {
	return fmt.Sprintf(
		"Holding %v\nLeft balance %v; Blocked money %v\nPending buy %v, sell %v\nLast buy price %v",
		s.holdingShares,
		s.leftBalance,
		s.notConfirmedBlockedMoney,
		s.pendingBuyShares,
		s.pendingSellShares,
		s.lastBuyPrice,
	)
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type MacdStrategy struct {
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

func New() *MacdStrategy {
	inst := &MacdStrategy{}
	inst.toPlaceOrders = make(chan *types.PlaceOrder)
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())
	inst.macd = *indicators.NewMacd(21, 16, 9)
	return inst
}

var l *log.Entry

func (s *MacdStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "macd",
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

	ch, err := candlesProvider.GetOrCreate(s.config.InstrumentID, now.Add(time.Duration(time.Minute) * 5 * 21))
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

func (s *MacdStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *MacdStrategy) onCandle(c types.OHLC) {
	wg := &sync.WaitGroup{}

	close := c.Close.Float()
	s.macd.Update([]float64{close})
	latestMacd, err := s.macd.Latest()
	if err != nil {
		l.Warnf("Cannot get latest MACD value: %v", err)
		return
	}

	state := s.state.Get()
	floatMacd, _ := latestMacd.Value.Float64()
	floatSignal, _ := latestMacd.Signal.Float64()
	l.Tracef("Updating signal with new values: signal %v; macd: %v", floatSignal, floatMacd)

	state.latestMacd = []float64{state.latestMacd[len(state.latestMacd)-1], floatMacd}
	state.latestSignals = []float64{state.latestSignals[len(state.latestSignals)-1], floatSignal}

	wg.Add(1)
	go s.buy(wg)
	wg.Add(1)
	go s.sell(wg)

	wg.Wait()
}

func (s *MacdStrategy) buy(wg *sync.WaitGroup) {
	defer wg.Done()
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	state := s.state.Get()
	isNowOver := state.latestMacd[1] > state.latestSignals[1]
	isPrevUnder := state.latestMacd[0] <= state.latestSignals[0]

	// Если дивергенция растет, то можно войти в позу
	shouldBuy := isNowOver && isPrevUnder

	if !shouldBuy {
		l.Infof("Not a good entry: macd %v, signal %v", state.latestMacd, state.latestSignals)
	}

	// TODO: Выставить ордер на покупку
}

func (s *MacdStrategy) sell(wg *sync.WaitGroup) {
	defer wg.Done()
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	state := *s.state.Get()

	isNowUnder := state.latestMacd[1] < state.latestSignals[1]
	isPrevOver := state.latestMacd[0] >= state.latestSignals[0]

	shouldSell := isNowUnder && isPrevOver

	if !shouldSell {
		l.Infof("Not a good exit: macd %v, signal %v", state.latestMacd, state.latestSignals)
	}

	// TODO: Выставить ордер на продажу

}

func (s *MacdStrategy) onOrderSateChange(state types.OrderExecutionState) {
	l.Infof("Order state changed %v", state)

	if state.Status == types.ErrorPlacing {
		l.Error("Order placing error. State restored")
	}

	newState := *s.state.Get()
	defer l.WithField("state", s.state.Get()).Info("State updated")

	s.state.Set(newState)
}
