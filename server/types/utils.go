package types

// TODO: Выпилить, это деталь реализации Тинькофф
type Quant struct {
	// Целая часть цены
	Units int
	// Дробная часть цены
	Nano int
}

type Money struct {
	Quant
	Currency string
}