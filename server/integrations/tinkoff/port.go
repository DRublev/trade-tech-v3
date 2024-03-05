package tinkoff

import (
	"context"
	"fmt"
	"main/types"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

const ENDPOINT = "invest-public-api.tinkoff.ru:443"

// https://github.com/RussianInvestments/invest-api-go-sdk

// TODO: Пора разделять методы по файлам

func (c *TinkoffBrokerPort) GetAccounts() ([]types.Account, error) {
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return []types.Account{}, err
	}

	us := sdk.NewUsersServiceClient()
	accountsRes, err := us.GetAccounts()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot get accounts ", err)
		return []types.Account{}, err
	}
	accounts := []types.Account{}

	for _, acc := range accountsRes.Accounts {
		isOpen := acc.Status == 2                                 //pb.AccountStatus_ACCOUNT_STATUS_OPEN
		hasAccess := acc.AccessLevel == 1 || acc.AccessLevel == 2 //AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS || AccessLevel_ACCOUNT_ACCESS_LEVEL_READ_ONLY
		isValidType := acc.Type == 1
		fmt.Println(acc) // pb.AccountType_ACCOUNT_TYPE_TINKOFF

		if isOpen && hasAccess && isValidType {
			accounts = append(accounts, types.Account{Id: acc.GetId(), Name: acc.GetName()})
		}
	}

	return accounts, nil
}

func (c *TinkoffBrokerPort) SetAccount(accountId string) error {
	return nil
}

func (c *TinkoffBrokerPort) GetCandles(instrumentId string, interval types.Interval, start time.Time, end time.Time) ([]types.OHLC, error) {
	// Инициализируем investgo sdk
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return []types.OHLC{}, err
	}

	// Сервис для работы с катировками
	candlesService := sdk.NewMarketDataServiceClient()

	// Получаем свечи по инструменту за определенный промежуток времени и интервал (переодичность)
	candlesRes, err := candlesService.GetCandles(instrumentId, investapi.CandleInterval(interval), start, end)
	if err != nil {
		fmt.Println("Cannot get candles", err)
		return []types.OHLC{}, err
	}

	candles := []types.OHLC{}

	// Конвертируем в нужный тип
	for _, candle := range candlesRes.Candles {
		candles = append(candles, toOHLC(candle))
	}
	return candles, nil
}

func (c *TinkoffBrokerPort) SubscribeCandles(ctx context.Context, ohlcCh *chan types.OHLC, instrumentId string, interval types.Interval) error {
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return err
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	// TODO: Эту штуку нужно переиспользовать в других эндпоинтах
	candlesStreamService := sdk.NewMarketDataStreamClient()

	candlesStream, err := candlesStreamService.MarketDataStream()
	if err != nil {
		fmt.Println("Cannot create stream ", err)
		return err
	}

	wg := &sync.WaitGroup{}

	// TODO: Докинуть обработку стакана и вообще вынести эту логику в некий Subscriber (глянуть паттерны)
	// Стрим не работает по выходным, см https://t.me/c/1436923108/53910/59213
	candlesCh, err := candlesStream.SubscribeCandle([]string{instrumentId}, investapi.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE, false)
	if err != nil {
		fmt.Println("Cannot subscribe ", err)
		return err
	}

	go func() {
		<-backCtx.Done()
		candlesStream.UnSubscribeCandle([]string{instrumentId}, investapi.SubscriptionInterval(interval), false)
		sdk.Stop()
	}()

	// Собирать свечи руками исходя из последних сделок, для выходных дней
	lastPriceCh := make(chan *investapi.LastPrice)
	// lastPriceCh, err := candlesStream.SubscribeLastPrice([]string{instrumentId})
	// if err != nil {
	// 	fmt.Println("Cannot subscribe ", err)
	// 	return err
	// }

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := candlesStream.Listen()

		fmt.Println("117 port", "listen end")

		if err != nil {
			fmt.Println("erorr in candles stream", err)
		}

	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		candles := make(map[int]types.OHLC)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("context closed for ", instrumentId)
				err := candlesStream.UnSubscribeCandle([]string{instrumentId}, investapi.SubscriptionInterval(interval), false)
				if err != nil {
					fmt.Println("Cannot unsubscribe ", instrumentId, err)
				}
				return
			case candle, ok := <-candlesCh:
				if !ok {
					fmt.Println("stream done for ", instrumentId)
					return
				}
				ohlc := toOHLC(candle)
				*ohlcCh <- ohlc
			// Врубать только для дебага графика в выходные!
			case lastPrice, ok := <-lastPriceCh:
				if !ok {
					fmt.Println("stream done for ", instrumentId)
					return
				}
				dealTime := lastPrice.Time.AsTime()
				if candle, exists := candles[dealTime.Minute()]; !exists {
					candles[dealTime.Minute()] = types.OHLC{
						Time:   dealTime,
						Open:   toQuant(lastPrice.Price),
						Close:  toQuant(lastPrice.Price),
						Low:    toQuant(lastPrice.Price),
						High:   toQuant(lastPrice.Price),
						Volume: 0,
					}
				} else {
					c := types.OHLC{
						Time:   dealTime,
						Open:   toQuant(lastPrice.Price),
						Close:  toQuant(lastPrice.Price),
						Low:    candle.Low,
						High:   candle.High,
						Volume: 0,
					}
					l := quantToNumber(candle.Low)
					h := quantToNumber(candle.High)
					if l > lastPrice.Price.ToFloat() {
						c.Low = toQuant(lastPrice.Price)
					}
					if h < lastPrice.Price.ToFloat() {
						c.High = toQuant(lastPrice.Price)
					}
					candles[dealTime.Minute()] = c
				}

				*ohlcCh <- candles[dealTime.Minute()]

				fmt.Println("164 port", lastPrice, candles[dealTime.Minute()])
			}
		}
	}(ctx)

	wg.Wait()

	return nil
}

func (c *TinkoffBrokerPort) GetShares(instrumentStatus types.InstrumentStatus) ([]types.Share, error) {
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return []types.Share{}, err
	}

	instrumentService := sdk.NewInstrumentsServiceClient()

	sharesRes, err := instrumentService.Shares(investapi.InstrumentStatus(instrumentStatus))
	if err != nil {
		fmt.Println("Cannot get shares", err)
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

	return shares, nil
}

func (c *TinkoffBrokerPort) SubscribeOrderbook(ctx context.Context, orderbookCh *chan *types.Orderbook, instrumentId string, depth int32) error {
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return err
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	streamService := sdk.NewMarketDataStreamClient()
	orderbookStream, err := streamService.MarketDataStream()
	if err != nil {
		fmt.Println("Cannot create stream ", err)
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := orderbookStream.Listen()

		if err != nil {
			fmt.Println("erorr in orderbooks stream", err)
		}

	}()

	unsubscribe := func() {
		err := orderbookStream.UnSubscribeOrderBook([]string{instrumentId}, depth)
		if err != nil {
			fmt.Println("Cannot unsubscribe ", instrumentId, err)
		}
		close(*orderbookCh)
	}

	go func() {
		<-backCtx.Done()
		fmt.Printf("290 port %v\n", backCtx.Err())
		unsubscribe()
		sdk.Stop()
	}()

	orderbookChan, err := orderbookStream.SubscribeOrderBook([]string{instrumentId}, depth)
	if err != nil {
		fmt.Println("Cannot subscribe ", err)
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
					fmt.Println("stream done for ", instrumentId)
					return
				}
				item := types.Orderbook{
					InstrumentId: instrumentId,
					Depth:        depth,
					Time:         orderbook.Time.AsTime(),
					LimitUp:      toQuant(orderbook.LimitUp),
					LimitDown:    toQuant(orderbook.LimitDown),
					Bids:         toBidAsk(orderbook.Bids),
					Asks:         toBidAsk(orderbook.Asks),
				}
				*orderbookCh <- &item
			}
		}
	}(ctx)

	wg.Wait()

	return nil
}

var accountId string

func (c *TinkoffBrokerPort) PlaceOrder(order *types.PlaceOrder) (types.OrderId, error) {
	// TODO: PlaceOrder -> TinkoffPlaceOrder
	fmt.Printf("336 port %v\n", order)
	sdk, err := c.GetSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return "", err
	}
	oc := sdk.NewOrdersServiceClient()

	direction := investapi.OrderDirection_ORDER_DIRECTION_BUY
	if order.Direction == types.Sell {
		direction = investapi.OrderDirection_ORDER_DIRECTION_SELL
	}
	
	price := toQuotation(float64(order.Price))
	if len(accountId) == 0 {
		accountIDRaw, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			return "", err
		}
		accountId = string(accountIDRaw)
	}

	orderResp, err := oc.PostOrder(&investgo.PostOrderRequest{
		InstrumentId: order.InstrumentID,
		Quantity: order.Quantity,
		Direction: direction,
		Price: &price,
		AccountId: accountId,
		OrderType: investapi.OrderType_ORDER_TYPE_LIMIT,
		OrderId: string(order.IdempodentID),
	})
	if err!=nil {
		fmt.Printf("362 port %v; accounId: %v\n", err, sdk.Config.AccountId)
		return "", err
	}
fmt.Printf("364 port %v\n", orderResp)
	return types.OrderId(orderResp.OrderId), err
}
