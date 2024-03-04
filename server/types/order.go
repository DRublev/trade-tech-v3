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
