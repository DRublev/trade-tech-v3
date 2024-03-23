package types

import "fmt"

// TODO: Выпилить, это деталь реализации Тинькофф
type Quant struct {
	// Целая часть цены
	Units int
	// Дробная часть цены
	Nano int
}

const BILLION = 1_000_000_000

func quantToNumber(q Quant) float64 {
	return float64(q.Units) + (float64(q.Nano) / BILLION)
}
func (q *Quant) String() string {
	return fmt.Sprintf("%v", quantToNumber(*q))
}

// TODO: Избавиться от Quant, завязаться на Price
type Money struct {
	Quant
	Currency string
}
