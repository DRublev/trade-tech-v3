package strategies

import (
	"main/identity"
	"main/metrics"
	"main/types"
	"sync"

	log "github.com/sirupsen/logrus"
)

type StatsCollector struct {
	source     *chan types.OrderExecutionState
	stopCh     chan struct{}
	strategy   string
	instrument string
	uid        string
	logger     *log.Entry
}

func NewStatsCollector(source *chan types.OrderExecutionState, strategy StrategyKey, instrument string) *StatsCollector {
	sc := &StatsCollector{
		source:     source,
		stopCh:     make(chan struct{}),
		logger:     log.WithField("metrics", strategy).WithField("instrument", instrument),
		strategy:   strategy.String(),
		instrument: instrument,
		uid:        identity.GetId(),
	}

	return sc
}

func (sc *StatsCollector) GetStats() {}

func (sc *StatsCollector) Start() {
	go func() {
		for {
			select {
			case <-sc.stopCh:
				return
			case state := <-*sc.source:
				sc.process(state)
			}
		}
	}()
}
func (sc *StatsCollector) Stop() {
	sc.stopCh <- struct{}{}
}

func (sc *StatsCollector) process(state types.OrderExecutionState) {
	// Профит
	// Оборот (уже есть дашборд по логам)
	// Кол-во совершенных сделок (с разрезом по их типу, купил или продал или стоп)
	// Кол-во отмененных/закрытых сделок
	// Кол-во прибыльных сделок
	// Кол-во не прибыльных сделок (общее число - прибыльные)
	sc.logger.Info("Got new state or metric")
	go sc.processProfitMetric(state)
	go sc.processTurnoverMetric(state)
	go sc.processDealsCountMetric(state)
	go sc.processDealsCancelledCountMetric(state)
	go sc.processDealsProfitableCountMetric(state)
	go sc.processDealsNonProfitableCountMetric(state)
}

var deals = make(map[string]types.OrderExecutionState)
var lastFilledBuys []types.OrderExecutionState = []types.OrderExecutionState{}
var dealsMX = sync.Mutex{}
var profit = 0
var profitMX = sync.Mutex{}

func (sc *StatsCollector) processProfitMetric(state types.OrderExecutionState) {
	if state.Status != types.Fill {
		return
	}

	if state.Direction == types.Buy {
		lastFilledBuys = append(lastFilledBuys, state)

		dealsMX.Lock()
		deals[string(state.ID)] = state
		dealsMX.Unlock()
		return
	}

	profit := 0.0

	// Продажа, когда ничего не куплено
	// TODO: Если будем делать торговлю в шорт, нужно и такой кейс учитывать
	if len(lastFilledBuys) == 0 {
		return
	}

	latestBuy := lastFilledBuys[len(lastFilledBuys)-1]

	if state.LotsExecuted == latestBuy.LotsExecuted {
		sellPrice := float64(state.LotsExecuted) * state.ExecutedOrderPrice
		buyPrice := latestBuy.ExecutedOrderPrice * float64(latestBuy.LotsExecuted)
		profit += sellPrice - buyPrice
		dealsMX.Lock()
		delete(deals, string(latestBuy.ID))
		dealsMX.Unlock()

		lastFilledBuys = lastFilledBuys[1:len(lastFilledBuys)-1]
	}

}

var turnover = 0.0
var turnoverMX = sync.Mutex{}

func (sc *StatsCollector) processTurnoverMetric(state types.OrderExecutionState) {
	if state.Status != types.Fill {
		return
	}

	turnoverMX.Lock()
	turnover += state.ExecutedOrderPrice * float64(state.LotsExecuted)
	turnoverMX.Unlock()
	sc.sendMetric("turnover", turnover)
}

var dealsCount = 0
var dealsCountMX = sync.Mutex{}

func (sc *StatsCollector) processDealsCountMetric(state types.OrderExecutionState) {
	if state.Status != types.Fill {
		return
	}

	dealsCountMX.Lock()
	dealsCount++
	dealsCountMX.Unlock()
	sc.sendMetric("deals_count", float64(dealsCount))
}

var dealsCancelledCount = 0
var dealsCancelledCountMX = sync.Mutex{}

func (sc *StatsCollector) processDealsCancelledCountMetric(state types.OrderExecutionState) {
	if state.Status != types.Cancelled {
		return
	}

	dealsCancelledCountMX.Lock()
	dealsCancelledCount++
	dealsCancelledCountMX.Unlock()
	sc.sendMetric("deals_cancelled_count", float64(dealsCancelledCount))
}
func (sc *StatsCollector) processDealsProfitableCountMetric(state types.OrderExecutionState) {

}
func (sc *StatsCollector) processDealsNonProfitableCountMetric(state types.OrderExecutionState) {

}

var metricsServer = metrics.NewMetricsService()

func (sc *StatsCollector) sendMetric(metric string, value float64) {
	sc.logger.WithField("metric-key", metric).WithField("value", value).Info("Updating metric")

	err := metricsServer.SendMetric(metric, value, map[string]string{"strategy": string(sc.strategy), "instrument_id": sc.instrument, "user_id": sc.uid})
	if err != nil {
		sc.logger.Warnf("Failed to send metric %v", err)
	}
}
