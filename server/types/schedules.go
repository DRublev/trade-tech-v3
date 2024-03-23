package types

import (
	"time"
)

type TradingDay struct {
	Date                           time.Time
	IsTradingDay                   bool
	StartTime                      time.Time
	EndTime                        time.Time
	OpeningAuctionStartTime        time.Time
	ClosingAuctionEndTime          time.Time
	EveningOpeningAuctionStartTime time.Time
	EveningStartTime               time.Time
	EveningEndTime                 time.Time
	ClearingStartTime              time.Time
	ClearingEndTime                time.Time
	PremarketStartTime             time.Time
	PremarketEndTime               time.Time
	ClosingAuctionStartTime        time.Time
	OpeningAuctionEndTime          time.Time
}

type TradingSchedule struct {
	Exchange string
	Days     []TradingDay
}
