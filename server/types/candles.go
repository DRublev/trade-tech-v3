package types

import (
	"time"
)

type Interval byte

type OHLC struct {
	// Цена открытия
	Open Quant
	// Максимальная цена за интервал
	High Quant
	// Минимальная цена за интервал
	Low Quant
	// Цена закрытия
	Close  Quant
	Time   time.Time
	Volume int64
}
