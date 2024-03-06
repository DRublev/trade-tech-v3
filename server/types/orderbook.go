package types

import "time"

type BidAsk struct {
	Quantity int64
	Price    float32
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
