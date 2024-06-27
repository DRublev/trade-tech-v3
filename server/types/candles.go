package types

import (
	"fmt"
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
	Close       Quant
	Time        time.Time
	LastTradeTS time.Time
	Volume      int64
}

func (o *OHLC) String() string {
	return fmt.Sprintf("OHLC: t %v; o %v; h %v; l %v; c %v; v %v", o.Time, o.Open, o.High, o.Low, o.Close, o.Volume)
}
