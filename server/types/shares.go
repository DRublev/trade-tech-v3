package types

import (
	"time"
)

// InstrumentStatus Статус инструмента
type InstrumentStatus byte
// TradingStatus Статус торгуемости инструмента
type TradingStatus byte

// Share Акция или фонд
// TODO: Переименовать в Instrument
type Share struct {
	Figi                string
	Name                string
	Exchange            string
	Ticker              string
	Lot                 int32
	IpoDate             time.Time
	TradingStatus       TradingStatus
	MinPriceIncrement   Quant
	Uid                 string
	First1minCandleDate time.Time
	First1dayCandleDate time.Time
}
