package types

import "fmt"

// OperationType Тип операции (сделки, ордера) - покупка или продажа
type OperationType byte

func (o *OperationType) String() string {
	if *o == 1 {
		return "Buy"
	}
	if *o == 2 {
		return "Sell"
	}

	return "Unknown"
}

// Коллекция возможнных типов операций ордера
const (
	Buy  OperationType = 1
	Sell OperationType = 2
)

// Order Выставленная заявка
type Order struct {
}

// Price Абстракция для представления цены
type Price float64

// PlaceOrder Еще не выставленная заявка, содержит поля для выставления заявки
type PlaceOrder struct {
	IdempodentID IdempodentID

	// Количесто лотов в заявке
	Quantity int64

	// Цена за 1 лот
	Price Price

	Direction OperationType

	// ID  инструмента (акции, фонда и тп.)
	InstrumentID string
}

// IdempodentID Абстракция для ID идемподентности
// Он нужен для сопоставления генерируемых ботом PlaceOrder и выставленных брокером Order
type IdempodentID string

// OrderID Абстракция над id ордера со стороны брокера
type OrderID string

// ExecutionStatus Статус ордера
type ExecutionStatus byte

// Возможные состояния ордера
const (
	Unspecified ExecutionStatus = 0
	// Fill Исполнен
	Fill ExecutionStatus = 1
	// Rejected Отклонен брокером
	Rejected ExecutionStatus = 2
	// Cancelled Отменен пользователем
	Cancelled ExecutionStatus = 3
	New       ExecutionStatus = 4
	// PartiallyFill Исполнен не полностью, не все лоты проданы/куплены
	PartiallyFill ExecutionStatus = 5
)

func (s *ExecutionStatus) String() string {
	switch *s {
	case 1:
		return "Fill"
	case 2:
		return "Rejected"
	case 3:
		return "Cancelled"
	case 4:
		return "New"
	case 5:
		return "PartiallyFill"
	default:
		return "Unspecified"
	}
}

// OrderExecutionState Состояние исполнения ордера
type OrderExecutionState struct {
	ID                OrderID
	IdempodentID      IdempodentID
	Status            ExecutionStatus
	LotsRequested     int
	LotsExecuted      int
	InitialOrderPrice float64
	// Полная стоимость (цена за лот * лот * количество акций в сделке)
	ExecutedOrderPrice float64
	InitialComission   float64
	ExecutedComission  float64
	Direction          OperationType
	InstrumentID       string
}

func (s *OrderExecutionState) String() string {
	return fmt.Sprintf(
		"%v ID: %v; IdempodentID: %v; Status: %v;\nLots requested %v, executed %v\n;Price initial %v; executed %v;\nComission initial %v; executed %v",
		s.Direction,
		s.ID,
		s.IdempodentID,
		s.Status,
		s.LotsRequested,
		s.LotsExecuted,
		s.InitialOrderPrice,
		s.ExecutedOrderPrice,
		s.InitialComission,
		s.ExecutedComission,
	)
}
