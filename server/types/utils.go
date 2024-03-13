package types

// TODO: Выпилить, это деталь реализации Тинькофф
type Quant struct {
	// Целая часть цены
	Units int
	// Дробная часть цены
	Nano int
}

// TODO: Избавиться от Quant, завязаться на Price
type Money struct {
	Quant
	Currency string
}