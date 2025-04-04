package tinkoff

import (
	"context"
	"errors"
	"main/types"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	log "github.com/sirupsen/logrus"
)


// TODO: Пора разделять методы по файлам

func (c *TinkoffBrokerPort) GetAccounts() ([]types.Account, error) {
	sdkL.Info("GetAccounts")
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return []types.Account{}, err
	}

	sdkL.Trace("Creating new users service client")
	us := sdk.NewUsersServiceClient()
	sdkL.Trace("Requesting accounts")
	accountsRes, err := us.GetAccounts()
	if err != nil {
		sdkL.Errorf("Failed getting accounts: %v", err)
		return []types.Account{}, err
	}
	accounts := []types.Account{}

	for _, acc := range accountsRes.Accounts {
		//pb.AccountStatus_ACCOUNT_STATUS_OPEN
		isOpen := acc.Status == 2
		//AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS || AccessLevel_ACCOUNT_ACCESS_LEVEL_READ_ONLY
		hasAccess := acc.AccessLevel == 1 || acc.AccessLevel == 2
		// pb.AccountType_ACCOUNT_TYPE_TINKOFF
		isValidType := acc.Type == 1

		if isOpen && hasAccess && isValidType {
			accounts = append(accounts, types.Account{Id: acc.GetId(), Name: acc.GetName()})
		}
	}
	sdkL.Infof("Found %v accounts", len(accounts))
	return accounts, nil
}

func (c *TinkoffBrokerPort) SetAccount(accountId string) error {
	return nil
}

func (c *TinkoffBrokerPort) GetCandles(instrumentID string, interval types.Interval, start time.Time, end time.Time) ([]types.OHLC, error) {
	sdkL.Infof("Getting candles for %v %v %v", instrumentID, start, end)
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return []types.OHLC{}, err
	}

	sdkL.Trace("Creating new marketdata service client")
	candlesService := sdk.NewMarketDataServiceClient()

	sdkL.Tracef("Requesting candles, instrument: %v; from: %v; to: %v; interval: %v", instrumentID, start, end, interval)
	candlesRes, err := candlesService.GetCandles(instrumentID, investapi.CandleInterval(interval), start, end)
	if err != nil {
		sdkL.Errorf("Failed getting candles: %v", err)
		return []types.OHLC{}, err
	}

	candles := []types.OHLC{}

	sdkL.Tracef("Mapping %v candles", len(candlesRes.Candles))
	for _, candle := range candlesRes.Candles {
		candles = append(candles, toOHLC(candle))
	}

	return candles, nil
}

func (c *TinkoffBrokerPort) SubscribeCandles(ctx context.Context, ohlcCh *chan types.OHLC, instrumentID string, interval types.Interval, waitForCandleClose bool) error {
	sdkL.Infof("Subscribe candles for %v", instrumentID)

	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return err
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	// TODO: Эту штуку нужно переиспользовать в других эндпоинтах
	candlesStreamService := sdk.NewMarketDataStreamClient()

	sdkL.Trace("Creating new candles stream")
	candlesStream, err := candlesStreamService.MarketDataStream()
	if err != nil {
		sdkL.Errorf("Failed creating new marketdata stream: %v", err)
		return err
	}

	wg := &sync.WaitGroup{}

	sdkL.Tracef("Subscribing for candles, instrument: %v", instrumentID)

	// Стрим не работает по выходным, см https://t.me/c/1436923108/53910/59213
	candlesCh, err := candlesStream.SubscribeCandle([]string{instrumentID}, investapi.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE, waitForCandleClose)
	if err != nil {
		sdkL.Errorf("Failed to subscribe for candles for %v: %v", instrumentID, err)
		return err
	}

	go func() {
		<-backCtx.Done()

		sdkL.Infof("Unsubscribing from candles for %v", instrumentID)
		err := candlesStream.UnSubscribeCandle([]string{instrumentID}, investapi.SubscriptionInterval(interval), waitForCandleClose)
		if err != nil {
			sdkL.Errorf("Failed to unsubscribe from candles: %v", err)
		}
	}()

	wg.Add(1)
	go func(ctx context.Context, ohlcCh *chan types.OHLC, candlesCh <-chan *investapi.Candle) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				sdkL.Infof("Unsubscribing from candles for %v", instrumentID)
				err := candlesStream.UnSubscribeCandle([]string{instrumentID}, investapi.SubscriptionInterval(interval), waitForCandleClose)
				if err != nil {
					sdkL.Warnf("Failed unsubscribing: %v", err)
				}
				return
			case candle, ok := <-candlesCh:
				if !ok {
					sdkL.Infof("Candles stream is done, %v", instrumentID)
					return
				}

				ohlc := toOHLCWithTrade(candle)
				sdkL.Tracef("Notifying about new candle. Candle time: %v", ohlc.Time)
				*ohlcCh <- ohlc
			}
		}
	}(backCtx, ohlcCh, candlesCh)

	wg.Add(1)
	go func() {
		defer wg.Done()

		sdkL.Trace("Start listening candles stream")
		err := candlesStream.Listen()
		if err != nil {
			sdkL.Errorf("Failed to listen candles stream: %v", err)
		}
	}()

	// wg.Wait()

	return nil
}

func (c *TinkoffBrokerPort) GetShares(instrumentStatus types.InstrumentStatus) ([]types.Share, error) {
	sdkL.Info("Get shares")
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return []types.Share{}, err
	}

	sdkL.Trace("Creating new instrument service client")
	instrumentService := sdk.NewInstrumentsServiceClient()

	sdkL.Trace("Sending get shares request")
	sharesRes, err := instrumentService.Shares(investapi.InstrumentStatus(instrumentStatus))
	if err != nil {
		sdkL.Errorf("Failed getting shares: %v", err)
		return []types.Share{}, err
	}

	etfsRes, err := instrumentService.Etfs(investapi.InstrumentStatus(instrumentStatus))
	if err != nil {
		sdkL.Errorf("Failed getting etf: %v", err)
		return []types.Share{}, err
	}

	shares := []types.Share{}

	for _, share := range sharesRes.Instruments {
		if share.ShareType == investapi.ShareType_SHARE_TYPE_COMMON &&
			!share.ForQualInvestorFlag &&
			share.ApiTradeAvailableFlag &&
			share.BuyAvailableFlag &&
			share.SellAvailableFlag {
			shares = append(shares, types.Share{
				Name:                share.Name,
				Figi:                share.Figi,
				Exchange:            share.Exchange,
				Ticker:              share.Ticker,
				Lot:                 share.Lot,
				IpoDate:             share.IpoDate.AsTime(),
				TradingStatus:       types.TradingStatus(share.TradingStatus),
				MinPriceIncrement:   toQuant(share.MinPriceIncrement),
				Uid:                 share.Uid,
				First1minCandleDate: share.First_1MinCandleDate.AsTime(),
				First1dayCandleDate: share.First_1DayCandleDate.AsTime(),
			})
		}
	}

	for _, etf := range etfsRes.Instruments {
		if !etf.ForQualInvestorFlag &&
			etf.ApiTradeAvailableFlag &&
			etf.BuyAvailableFlag &&
			etf.SellAvailableFlag {
			shares = append(shares, types.Share{
				Name:                etf.Name,
				Figi:                etf.Figi,
				Exchange:            etf.Exchange,
				Ticker:              etf.Ticker,
				Lot:                 etf.Lot,
				IpoDate:             etf.ReleasedDate.AsTime(),
				TradingStatus:       types.TradingStatus(etf.TradingStatus),
				MinPriceIncrement:   toQuant(etf.MinPriceIncrement),
				Uid:                 etf.Uid,
				First1minCandleDate: etf.First_1MinCandleDate.AsTime(),
				First1dayCandleDate: etf.First_1DayCandleDate.AsTime(),
			})
		}
	}

	sdkL.Infof("Got %v shares", len(shares))
	return shares, nil
}

func (c *TinkoffBrokerPort) SubscribeOrderbook(ctx context.Context, orderbookCh *chan *types.Orderbook, instrumentID string, depth int32) error {
	sdkL.WithFields(log.Fields{
		"instrument": instrumentID,
		"depth":      depth,
	}).Infof("Subscribe for orderbook")

	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return err
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	streamService := sdk.NewMarketDataStreamClient()

	sdkL.Trace("Creating new marketdata stream")

	orderbookStream, err := streamService.MarketDataStream()
	if err != nil {
		sdkL.Errorf("Failed creating marketdata stream: %v", err)
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sdkL.Trace("Start listening orderbook stream")
		err := orderbookStream.Listen()
		if err != nil {
			sdkL.Errorf("Failed to listen orderbook stream: %v", err)
		}
	}()

	unsubscribe := func() {
		sdkL.Infof("Unsubscribing from orderbook for %v", instrumentID)
		err := orderbookStream.UnSubscribeOrderBook([]string{instrumentID}, depth)
		if err != nil {
			sdkL.Errorf("Failed to unsubscribe from orderbook: %v", err)
		}
		close(*orderbookCh)
	}

	go func() {
		<-backCtx.Done()

		sdkL.Trace("SubscribeOrderbook context is closed")
		unsubscribe()
	}()

	sdkL.Tracef("Subscribing for orderbook for %v", instrumentID)
	orderbookChan, err := orderbookStream.SubscribeOrderBook([]string{instrumentID}, depth)
	if err != nil {
		sdkL.Errorf("Failed to subscribe for orderbook: %v", err)
		return err
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				unsubscribe()
				return
			case orderbook, ok := <-orderbookChan:
				if !ok {
					sdkL.Errorf("Orderbook stream is done for %v", instrumentID)
					return
				}
				item := types.Orderbook{
					InstrumentId: instrumentID,
					Depth:        depth,
					Time:         orderbook.Time.AsTime(),
					LimitUp:      toQuant(orderbook.LimitUp),
					LimitDown:    toQuant(orderbook.LimitDown),
					Bids:         toBidAsk(orderbook.Bids),
					Asks:         toBidAsk(orderbook.Asks),
				}
				sdkL.Tracef("New orderbook item. Time: %v", orderbook.Time)
				*orderbookCh <- &item
			}
		}
	}(ctx)

	wg.Wait()

	return nil
}

var accountID string

func (c *TinkoffBrokerPort) getAccountId() string {
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return ""
	}

	// TODO: Вынести как отдельную сущность с провайдером/репозиторием
	var accId = accountID
	accountIdOkCh := make(chan struct{}, 1)
	if len(accId) == 0 {
		sdkL.Trace("No accountID")
		accountIDRaw, err := dbInstance.Get([]string{"accounts"})
		if err == nil {
			sdkL.Trace("Got accountID from db")
			accId = string(accountIDRaw)
			accountIdOkCh <- struct{}{}
		} else {
			sdkL.Trace("Got accountID from sdk")
			accId = sdk.Config.AccountId
			accountIdOkCh <- struct{}{}
		}
	} else {
		accountIdOkCh <- struct{}{}
	}
	<-accountIdOkCh

	return accId
}

// TODO: Как это кешировать? На каком уровне? На уровне брокера?
var cachedInstruments = make(map[string]*investapi.Instrument)

func (c *TinkoffBrokerPort) PlaceOrder(order *types.PlaceOrder) (types.OrderID, error) {
	sdkL.WithFields(log.Fields{
		"instrument": order.InstrumentID,
		"direction":  order.Direction,
	}).Infof("Placing order")

	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return "", err
	}
	oc := sdk.NewOrdersServiceClient()

	if len(order.CancelOrder) > 0 {
		err = c.CancelOrder(order.CancelOrder)
		if err != nil {
			sdkL.Warnf("Error closing order %v: %v", order.CancelOrder, err)
			return "", errors.New("error closing order")
		}
		return "", errors.New("order cancelled")
	}

	direction := investapi.OrderDirection_ORDER_DIRECTION_BUY
	if order.Direction == types.Sell {
		direction = investapi.OrderDirection_ORDER_DIRECTION_SELL
	}

	cachedInstrument, exists := cachedInstruments[order.InstrumentID]
	if !exists {
		ic := sdk.NewInstrumentsServiceClient()
		instrument, err := ic.InstrumentByUid(order.InstrumentID)
		if err != nil {
			sdkL.Errorf("Cannot get instrument info: %v", err)
		} else {
			cachedInstrument = instrument.Instrument
		}
	}

	// Дефолтное значение, на случай если мы не смогли получить инфу об инструменте
	// Для большиства инструментов отработает ок, для остальных вернется ошибка при выставлении ордера
	minPriceIncrement := &investapi.Quotation{
		Units: 0,
		Nano:  10000000, // Ok
		// Nano: 10000, // VTBR
		// Nano: 10000000, // TMOS
	}
	if cachedInstrument != nil {
		minPriceIncrement = cachedInstrument.MinPriceIncrement
	}

	price := FloatToQuotation(float64(order.Price), minPriceIncrement)

	o := &investgo.PostOrderRequest{
		InstrumentId: order.InstrumentID,
		Quantity:     order.Quantity,
		Direction:    direction,
		Price:        &price,
		AccountId:    c.getAccountId(),
		OrderType:    investapi.OrderType_ORDER_TYPE_LIMIT,
		OrderId:      string(order.IdempodentID),
	}

	sdkL.Tracef("Placing order %v", order)
	orderResp, err := oc.PostOrder(o)
	if err != nil {
		sdkL.Errorf("Failed to place order: %v", err)
		return "", err
	}

	sdkL.Infof("Order placed, id: %v", orderResp.OrderId)
	return types.OrderID(orderResp.OrderId), err
}


func (c *TinkoffBrokerPort) SubscribeOrders(cb func(types.OrderExecutionState)) error {
	sdkL.Info("Subscribing for order states")
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return err
	}

	ordersStreamClient := sdk.NewOrdersStreamClient()

	sdkL.Trace("Creating new trades stream")
	tradesStream, err := ordersStreamClient.TradesStream([]string{
		c.getAccountId(),
	})
	if err != nil {
		sdkL.Errorf("Failed to create tradees stream: %v", err)
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sdkL.Trace("Start listening trades stream")

		err := tradesStream.Listen()
		if err != nil {
			sdkL.Errorf("Failed to listen trades stream: %v", err)
		}
	}()

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		<-backCtx.Done()

		sdkL.Info("Unsubscribing from trades stream")
		tradesStream.Stop()
	}()

	wg.Add(1)
	go func(ctx context.Context, ts *investgo.TradesStream, cb func(types.OrderExecutionState)) {
		defer wg.Done()

		for {
			select {
			case tradeState, ok := <-ts.Trades():
				// TODO: there are some nil trades.
				if !ok {
					return
				}

				sdkL.Infof("New state of order: %v; direction: %v; time %v", tradeState.OrderId, tradeState.Direction, tradeState.CreatedAt)

				lotsExecuted := 0
				var executedPrice float64 = 0
				for _, t := range tradeState.Trades {
					sdkL.Tracef("Trade of order: quantity %v; price: %v", int(t.Quantity), t.Price.ToFloat())
					lotsExecuted += int(t.Quantity)
					executedPrice += t.Price.ToFloat() * float64(t.Quantity)
				}

				changeEvent := types.OrderExecutionState{
					ID:                 types.OrderID(tradeState.OrderId),
					Direction:          types.OperationType(tradeState.Direction),
					InstrumentID:       tradeState.InstrumentUid,
					LotsExecuted:       lotsExecuted,
					Status:             0, // TODO: Научиться определять статус заявки
					ExecutedOrderPrice: executedPrice,
					ExecutionTime:      tradeState.Trades[len(tradeState.Trades)-1].DateTime.AsTime(),
					// TODO: Научиться считать вот это все (из tradeState.Trades видимо)
					// LotsRequested      int
					// InitialOrderPrice  types.Money
					// ExecutedOrderPrice types.Money
					// InitialComission   types.Money
					// ExecutedComission  types.Money
				}

				state, err := c.GetOrderState(types.OrderID(tradeState.OrderId))
				if err == nil {
					changeEvent.Status = state.Status
					changeEvent.LotsRequested = state.LotsRequested
				}

				sdkL.Infof("Order state changed, notifying: %v", changeEvent)
				go cb(changeEvent)
			}
		}
	}(backCtx, tradesStream, cb)

	wg.Wait()

	return nil
}

func (c *TinkoffBrokerPort) GetTradingSchedules(exchange string, from time.Time, to time.Time) ([]types.TradingSchedule, error) {
	sdkL.Info("Getting trading schedules")

	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return []types.TradingSchedule{}, err
	}

	instrumentService := sdk.NewInstrumentsServiceClient()

	tradingSchedulesRes, err := instrumentService.TradingSchedules(exchange, from, to)
	if err != nil {
		sdkL.Errorf("Cannot get trading schedules: %v", err)
		return []types.TradingSchedule{}, err
	}

	exchanges := []types.TradingSchedule{}

	for _, ex := range tradingSchedulesRes.Exchanges {
		days := []types.TradingDay{}

		for _, day := range ex.Days {
			days = append(days, types.TradingDay{
				Date:         day.Date.AsTime(),
				IsTradingDay: day.IsTradingDay,
				StartTime:    day.StartTime.AsTime(),
				EndTime:      day.EndTime.AsTime(),
			})
		}
		exchanges = append(exchanges, types.TradingSchedule{
			Exchange: ex.Exchange,
			Days:     days,
		})
	}

	return exchanges, nil
}
func (c *TinkoffBrokerPort) GetOrderState(orderID types.OrderID) (types.OrderExecutionState, error) {
	sdkL.Infof("Getting state of order %v", orderID)
	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Failed to create tradees stream: %v", err)
		return types.OrderExecutionState{}, err
	}

	oc := sdk.NewOrdersServiceClient()

	sdkL.Trace("Sending get order state request")
	state, err := oc.GetOrderState(c.getAccountId(), string(orderID))
	if err != nil {
		sdkL.Errorf("Failed to get order state: %v", err)
		return types.OrderExecutionState{}, err
	}
	var status types.ExecutionStatus = types.Unspecified
	if state.LotsExecuted == state.LotsRequested || state.ExecutionReportStatus == 1 {
		status = types.Fill
	} else if state.ExecutionReportStatus == 4 {
		status = types.New
	} else if state.ExecutionReportStatus == 2 {
		status = types.ErrorPlacing
	}
	exectuionTime := time.Now()
	if len(state.Stages) > 0 {
		exectuionTime = state.Stages[len(state.Stages)-1].ExecutionTime.AsTime()
	}

	orderState := types.OrderExecutionState{
		ID:                 types.OrderID(state.OrderId),
		IdempodentID:       types.IdempodentID(state.OrderRequestId),
		Direction:          types.OperationType(state.Direction),
		InstrumentID:       state.InstrumentUid,
		LotsExecuted:       int(state.LotsExecuted),
		LotsRequested:      int(state.LotsRequested),
		Status:             status, // TODO: Научиться определять статус заявки
		ExecutedOrderPrice: state.ExecutedOrderPrice.ToFloat(),
		ExecutionTime:      exectuionTime,
		InitialOrderPrice:  state.InitialOrderPrice.ToFloat(),
	}

	sdkL.Infof("Got order state %v", orderState)
	return orderState, nil
}

func (c *TinkoffBrokerPort) CancelOrder(orderID types.OrderID) error {
	sdkL.WithFields(log.Fields{
		"orderID": orderID,
	}).Infof("Cancelling order")

	sdk, err := c.GetSdk()
	if err != nil {
		sdkL.Errorf("Cannot init sdk: %v", err)
		return err
	}
	oc := sdk.NewOrdersServiceClient()

	accID := c.getAccountId()

	res, err := oc.CancelOrder(accID, string(orderID))

	sdkL.Infof("Cancelling order %v %v", accID, string(orderID))
	if err != nil {
		sdkL.Infof("Cancelling order error: %v", err)
	} else {
		sdkL.Infof("Cancelling order response: %v", res.CancelOrderResponse)
	}

	return err
}
