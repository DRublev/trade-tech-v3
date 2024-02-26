package types

import "time"

type BidAsk struct {
	Quantiny int64
	Price    Quant
}

type Orderbook struct {
	InstrumentId string
	Depth        int32
	Time         time.Time
	LimitUp      Quant
	LimitDown    Quant
	Bids         []*BidAsk
	Asks         []*BidAsk
}
