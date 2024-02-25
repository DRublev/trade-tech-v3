package tinkoff

import (
	"context"
	"fmt"
	"main/types"
	"os"
	"os/signal"
	"sync"
	"time"

	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

const ENDPOINT = "sandbox-invest-public-api.tinkoff.ru:443"

// https://github.com/RussianInvestments/invest-api-go-sdk
// TODO: Хорошо бы явно наследовать types.Broker (чтоб были подсказки при имплементации метода)
type TinkoffBrokerPort struct{}

func (c *TinkoffBrokerPort) GetAccounts() ([]types.Account, error) {
	sdk, err := c.getSdk()
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

func toQuant(iq *investapi.Quotation) types.Quant {
	return types.Quant{
		Units: int(iq.Units),
		Nano:  int(iq.Nano),
	}
}

func (c *TinkoffBrokerPort) GetCandles(instrumentId string, interval types.Interval, start time.Time, end time.Time) ([]types.OHLC, error) {
	// Инициализируем investgo sdk
	sdk, err := c.getSdk()
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
		candles = append(candles, types.OHLC{
			Time:   candle.Time.AsTime(),
			Open:   toQuant(candle.Open),
			Close:  toQuant(candle.Close),
			Low:    toQuant(candle.Low),
			High:   toQuant(candle.High),
			Volume: candle.Volume,
		})
	}
	return candles, nil
}

const nanoPrecision = 1_000_000_000

func quantToNumber(q types.Quant) float64 {
	return float64(q.Units) + (float64(q.Nano) / nanoPrecision)
}

func (c *TinkoffBrokerPort) SubscribeCandles(ctx context.Context, ohlcCh *chan types.OHLC, instrumentId string, interval types.Interval) error {
	sdk, err := c.getSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return err
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-backCtx.Done()
		sdk.Stop()
	}()

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
				err := candlesStream.UnSubscribeAll()
				if err != nil {
					fmt.Println("Cannot unsubscribe ", instrumentId, err)
				}
				return
			case candle, ok := <-candlesCh:
				if !ok {
					fmt.Println("stream done for ", instrumentId)
					return
				}
				ohlc := types.OHLC{
					Time: candle.Time.AsTime(),
					Volume: candle.Volume,
					Open: toQuant(candle.Open),
					High: toQuant(candle.High),
					Low: toQuant(candle.Low),
					Close: toQuant(candle.Close),
				}
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

	// wg.Wait()

	return nil
}
