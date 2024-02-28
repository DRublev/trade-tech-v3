package orderbook

import "main/types"

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
	Id                 types.OrderId
	IdempodentId       types.IdempodentId
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
