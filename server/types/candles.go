package types

import (
	"time"
)

type Interval byte

type Quant struct {
	// Целая часть цены
	Units int
	// Дробная часть цены
	Nano int
}

type OHLC struct {
	// Цена открытия
	Open Quant
	// Максимальная цена за интервал
	High Quant
	// Минимальная цена за интервал
	Low Quant
	// Цена закрытия
	Close Quant
	Time  time.Time
	Volume int64
}
