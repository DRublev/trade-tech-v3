package tinkoff

import (
	"fmt"
	"main/types"
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
			Time:  candle.Time.AsTime(),
			Open:  toQuant(candle.Open),
			Close: toQuant(candle.Close),
			Low:   toQuant(candle.Low),
			High:  toQuant(candle.High),
		})
	}
	return candles, nil
}
