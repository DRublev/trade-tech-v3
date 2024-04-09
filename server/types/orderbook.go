package types

import "time"

// BidAsk Предложение покупки или продажи из стакана
type BidAsk struct {
	Quantity int64
	Price    float32
}

// Orderbook Стакан
type Orderbook struct {
	InstrumentId string
	Depth        int32
	Time         time.Time
	LimitUp      Quant
	LimitDown    Quant
	Bids         []*BidAsk
	Asks         []*BidAsk
}
