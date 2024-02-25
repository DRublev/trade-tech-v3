package bot

import "main/types"

type IMIddleware[T any] interface {
	Do(arg T) error
}

type IPlaceOrderMiddleware interface {
	IMIddleware[types.Order]
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
	InitialOrderPrice  types.Money
	ExecutedOrderPrice types.Money
	InitialComission   types.Money
	ExecutedComission  types.Money
	Direction          types.OperationType
	InstrumentId       string
}

type IOrderWatcher interface {
	MarkAsSent(id IdempodentId, orderId OrderId) error
	Subscribe(orderId IdempodentId) error
	notify(orderId IdempodentId) error
}
