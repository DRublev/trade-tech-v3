package types

type OperationType byte

const (
	Buy  OperationType = 1
	Sell OperationType = 2
)

type Order struct {
}

type PlaceOrder struct {
	IdempodentId IdempodentId
}

type IdempodentId string
type OrderId string
