package strategies

import (
	"fmt"
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

	PlacedBuyOrders  []types.OrderExecutionState
	PlacedSellOrders []types.OrderExecutionState

	lotSize int64
}

func NewVault(lotSize int64, balance float64) *Vault {
	inst := &Vault{
		HoldingShares:            0,
		PendingBuyShares:         0,
		PendingSellShares:        0,
		NotConfirmedBlockedMoney: 0,
		LastBuyPrice:             0,
		lotSize:                  0,
		LeftBalance:              0,
	}
	inst.lotSize = lotSize
	inst.LeftBalance = balance
	inst.PlacedBuyOrders = []types.OrderExecutionState{}
	inst.PlacedSellOrders = []types.OrderExecutionState{}

	return inst
}

func (s *Vault) String() string {
	return fmt.Sprintf(
		"Holding %v\nLeft balance %v; Blocked money %v\nPending buy %v, sell %v\nLast buy price %v",
		s.HoldingShares,
		s.LeftBalance,
		s.NotConfirmedBlockedMoney,
		s.PendingBuyShares,
		s.PendingSellShares,
		s.LastBuyPrice,
	)
}

func (this *Vault) updateBuyOrders(state types.OrderExecutionState) {
	if state.Direction != types.Buy {
		return
	}
	if state.Status == types.New {
		this.PlacedBuyOrders = append(this.PlacedBuyOrders, state)
		l.Infof("Adding new buy order to placed list (ID: %v)", state.ID)
		return
	}
	if state.Status == types.Fill || state.Status == types.ErrorPlacing {
		filteredOrders := []types.OrderExecutionState{}
		for _, order := range this.PlacedBuyOrders {
			if order.ID == state.ID || order.IdempodentID == state.IdempodentID {
				l.Infof("Removing cancelled buy order from pending list: %v", state.ID)
				continue
			}
			filteredOrders = append(filteredOrders, order)
		}

		this.PlacedBuyOrders = filteredOrders
	}
}
func (this *Vault) updateSellOrders(state types.OrderExecutionState) {
	if state.Direction != types.Sell {
		return
	}
	if state.Status == types.New {
		this.PlacedSellOrders = append(this.PlacedSellOrders, state)
		l.Infof("Adding new sell order to placed list")
		return
	}
	if state.Status == types.Fill || state.Status == types.ErrorPlacing {
		filteredOrders := []types.OrderExecutionState{}

		for _, order := range this.PlacedSellOrders {
			if order.ID == state.ID || order.IdempodentID == state.IdempodentID {
				l.Infof("Removing cancelled sell order from pending list: %v", state.ID)
				continue
			}
			filteredOrders = append(filteredOrders, order)
		}

		this.PlacedSellOrders = filteredOrders
	}
}

func (this *Vault) OnOrderSateChange(state types.OrderExecutionState) {
	l.Infof("Order state changed %v", state)

	if state.Status == types.ErrorPlacing {
		l.Error("Order placing error. State restored")
	}

	defer l.WithField("state", this).Info("State updated")

	this.updateBuyOrders(state)
	this.updateSellOrders(state)

	if state.Status != types.PartiallyFill &&
		state.Status != types.Fill &&
		state.Status != types.ErrorPlacing &&
		state.Status != types.Cancelled &&
		state.Status != types.New {
		l.Warnf("Not processed order state change: %v", state)
		return
	}

	isBuyPlaceError := state.Direction == types.Buy && state.Status == types.ErrorPlacing
	isSellPlaceError := state.Direction == types.Sell && state.Status == types.ErrorPlacing
	isBuyCancel := state.Direction == types.Buy && state.Status == types.Cancelled
	isSellCancel := state.Direction == types.Sell && state.Status == types.Cancelled
	isSellOk := state.Direction == types.Sell && !isSellPlaceError && !isSellCancel
	isBuyOk := state.Direction == types.Buy && !isBuyPlaceError && !isBuyCancel

	if isBuyPlaceError || isBuyCancel {
		l.Info("Updating state after buy order place error")
		this.LeftBalance += state.ExecutedOrderPrice
		this.PendingBuyShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.NotConfirmedBlockedMoney -= state.ExecutedOrderPrice
		l.Infof("NotConfirmedBlockedMoney %v; ExecutedOrderPrice %v", this.NotConfirmedBlockedMoney, state.ExecutedOrderPrice)
		return
	} else if isSellPlaceError || isSellCancel {
		l.Info("Updating state after sell order place error")
		this.PendingSellShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.HoldingShares += int64(state.LotsExecuted / int(this.lotSize))
		return
	}

	if isSellOk {
		l.Trace("Updating state after sell order executed")
		this.PendingSellShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.HoldingShares -= int64(state.LotsExecuted / int(this.lotSize))

		if state.Status != types.New {
			this.LeftBalance += state.ExecutedOrderPrice
		}

		l.WithField("orderId", state.ID).Infof(
			"Lots executed (cancelled %v, erroPlacing: %v) %v of %v; Executed sell price %v",
			isBuyCancel,
			isBuyPlaceError,
			state.LotsExecuted,
			state.LotsRequested,
			state.ExecutedOrderPrice,
		)
	} else if isBuyOk {
		l.Tracef("Updating state after buy order executed, oder: %v; self state %v", state.String(), this.String())
		this.HoldingShares += int64(state.LotsExecuted / int(this.lotSize))
		this.PendingBuyShares -= int64(state.LotsExecuted / int(this.lotSize))
		this.LastBuyPrice = state.ExecutedOrderPrice

		if state.Status != types.New {
			this.NotConfirmedBlockedMoney -= state.ExecutedOrderPrice / float64(state.LotsExecuted) * float64(this.lotSize)

			this.LeftBalance -= state.ExecutedOrderPrice
		}

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
