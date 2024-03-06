package types

type OperationType byte

const (
	Buy  OperationType = 1
	Sell OperationType = 2
)

type Order struct {
}

type Price float64

type PlaceOrder struct {
	IdempodentID IdempodentId

	// Количесто лотов в заявке
	Quantity int64

	// Цена за 1 лот
	Price Price

	Direction OperationType

	// ID  инструмента (акции, фонда и тп.)
	InstrumentID string
}

type IdempodentId string
type OrderId string


type ExecutionStatus byte

const (
	Unspecified   ExecutionStatus = 0
	Fill          ExecutionStatus = 1
	Rejected      ExecutionStatus = 2
	Cancelled     ExecutionStatus = 3
	New           ExecutionStatus = 4
	PartiallyFill ExecutionStatus = 5
)

type OrderExecutionState struct {
	Id                 OrderId
	IdempodentId       IdempodentId
	Status             ExecutionStatus
	LotsRequested      int
	LotsExecuted       int
	InitialOrderPrice  Money
	ExecutedOrderPrice Money
	InitialComission   Money
	ExecutedComission  Money
	Direction          OperationType
	InstrumentId       string
}
