package types

import (
	"time"
)

type InstrumentStatus byte
type TradingStatus byte

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
