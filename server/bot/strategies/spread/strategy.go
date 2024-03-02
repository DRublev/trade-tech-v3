package spread

import (
	"fmt"
	"main/bot/orderbook"
	"main/bot/strategies"
	"main/types"
	"sync"
)

type Config struct {
	strategies.Config

	// Минимальная разница bid-ask, при которой выставлять ордер
	minSpread float32

	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
	nextOrderCooldownMs int32
}

type State struct {
	// Оставшееся количество денег
	balance float32

	// Количество лотов, купленных на данный момент
	holdingLots int32

	// Количество лотов, на которое выставлены ордера на покупку
	pendingBuyLots int32

	// Количество лотов, на которое выставлены ордера на продажу
	pendingSellLots int32
}

type SpreadStrategy struct {
	strategies.IStrategy
	strategies.Strategy
	config Config
	// Канал для стакана
	obCh  *chan *types.Orderbook
	state strategies.StrategyState[State]
}

func New() *SpreadStrategy {
	inst := &SpreadStrategy{}

	orderCh := make(chan types.PlaceOrder)
	inst.OrdersToPlaceCh = &orderCh

	return inst
}

func (s *SpreadStrategy) Start(config *strategies.Config) (bool, error) {
	// TODO: Нужен метод ConvertSerialsableToType[T](candidate) T, который конвертирует типы через json.Marshall
	s.config = ((any)(*config)).(Config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	obProvider := orderbook.NewOrederbookProvider()
	ch, err := obProvider.GetOrCreate(s.config.InstrumentId)
	if err != nil {
		return false, err
	}

	// Заполняем изначальное состояние
	s.state = strategies.StrategyState[State]{}
	err = s.state.Set(State{
		holdingLots:     0,
		pendingBuyLots:  0,
		pendingSellLots: 0,
		balance:         s.config.Balance,
	})
	if err != nil {
		return false, err
	}

	s.obCh = ch
	// Слушаем изменения в стакане
	go func(ch *chan *types.Orderbook) {
		for {
			select {
			case ob, ok := <-*ch:
				if !ok {
					fmt.Println("spread orderbook channel end")
				}

				go s.onOrderbook(ob)
			}
		}
	}(s.obCh)

	return true, nil
}

func (s *SpreadStrategy) Stop() (bool, error) {
	return false, nil
}

func (s *SpreadStrategy) onOrderbook(ob *types.Orderbook) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go s.buy(wg, ob)
	wg.Add(1)
	go s.sell(wg, ob)

	wg.Wait()
}

// TODO: Перенести в отдельный файл
func (s *SpreadStrategy) buy(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()

}

func (s *SpreadStrategy) sell(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()

}
