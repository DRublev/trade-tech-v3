package strategies

import (
	"main/types"

	log "github.com/sirupsen/logrus"
)

var l = log.WithFields(log.Fields{
	"strategy": "???",
})

type Vault struct {
	// Оставшееся количество денег
	LeftBalance float64

	// Сумма, которая должна списаться при выставлении ордера на покупку
	// Инкрементим когда хотим выставить бай ордер
	// Декрементим когда закрываем бай ордер
	NotConfirmedBlockedMoney float64

	// Количество акций, купленных на данный момент
	HoldingShares int64

	// Количество акций, на которое выставлены ордера на покупку
	PendingBuyShares int64

	// Количество акций, на которое выставлены ордера на продажу
	PendingSellShares int64

	LastBuyPrice float64

	PlacedOrders []types.OrderExecutionState

	lotSize int64
}

func NewVault(lotSize int64, balance float64) *Vault {
	inst := &Vault{
		HoldingShares:            0,
		PendingBuyShares:         0,
		PendingSellShares:        0,
		NotConfirmedBlockedMoney: 0,
		LastBuyPrice:             0,
		lotSize: 0,
		LeftBalance: 0,
	}
	inst.lotSize = lotSize
	inst.LeftBalance = balance

	return inst
}

func (this *Vault) OnOrderSateChange(state types.OrderExecutionState) {
	l.Infof("Order state changed %v", state)

	if state.Status == types.ErrorPlacing {
		l.Error("Order placing error. State restored")
	}

	defer l.WithField("state", this).Info("State updated")

	if state.Status == types.New {
		this.PlacedOrders = append(this.PlacedOrders, state)
		l.Infof("Adding new order to placed list")
		return
	}
	if state.Status == types.Fill {
		filteredOrders := []types.OrderExecutionState{}

		for _, order := range this.PlacedOrders {
			if order.ID != state.ID {
				filteredOrders = append(filteredOrders, order)
			}
		}

		this.PlacedOrders = filteredOrders
	}

	if state.Status != types.PartiallyFill &&
		state.Status != types.Fill &&
		state.Status != types.ErrorPlacing &&
		state.Status != types.Cancelled {
		l.Warnf("Not processed order state change: %v", state)
		return
	}

	isBuyPlaceError := state.Direction == types.Buy && state.Status == types.ErrorPlacing
	isSellPlaceError := state.Direction == types.Sell && state.Status == types.ErrorPlacing
	isBuyCancel := state.Direction == types.Buy && state.Status == types.Cancelled
	isSellCancel := state.Direction == types.Sell && state.Status == types.Cancelled
	isSellOk := state.Direction == types.Sell && !isSellPlaceError && !isSellCancel
	isBuyOk := state.Direction == types.Buy && !isBuyPlaceError && !isBuyCancel

	if isBuyPlaceError {
		l.Trace("Updating state after buy order place error")
		this.LeftBalance += state.ExecutedOrderPrice
		this.PendingBuyShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.NotConfirmedBlockedMoney -= state.ExecutedOrderPrice
	} else if isSellPlaceError {
		this.PendingSellShares -= int64(state.LotsExecuted / int(this.lotSize))
	}

	if isSellOk || isBuyCancel {
		l.Trace("Updating state after sell order executed")
		this.PendingSellShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.LeftBalance += state.ExecutedOrderPrice
		this.HoldingShares -= int64(state.LotsExecuted / int(this.lotSize))
		l.WithField("orderId", state.ID).Infof(
			"Lots executed (cancelled %v, erroPlacing: %v) %v of %v; Executed sell price %v",
			isBuyCancel,
			isBuyPlaceError,
			state.LotsExecuted,
			state.LotsRequested,
			state.ExecutedOrderPrice,
		)
	} else if isBuyOk || isSellPlaceError || isSellCancel {
		l.Trace("Updating state after buy order executed")
		this.HoldingShares += int64(state.LotsExecuted / int(this.lotSize))
		this.PendingBuyShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.NotConfirmedBlockedMoney -= state.ExecutedOrderPrice
		this.LeftBalance -= state.ExecutedOrderPrice
		this.LastBuyPrice = state.ExecutedOrderPrice / float64(state.LotsExecuted)
		l.WithField("orderId", state.ID).Infof(
			"Lots executed (cancelled %v, erroPlacing: %v) %v of %v; Executed buy price %v",
			isSellCancel,
			isSellPlaceError,
			state.LotsExecuted,
			state.LotsRequested,
			state.ExecutedOrderPrice,
		)
	} else {
		l.Warnf("Order state change not handled: %v", state)
	}
}
