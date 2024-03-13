package types

import (
	"time"
)

// Interval  Интервал на графике (минута, 5 минут, час и тп)
type Interval byte

// OHLC Педставление свечи
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
