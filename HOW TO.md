# Гайды

## Как получить данные из брокера? (тинькофф)
Например, нужно получить данные о свечах за прошедший период из тинькофф. Для этого:

- В интерфейсе `Broker` описать метод получения свечей
  ```go
    // types/broker.go
    type IBroker interface {
        // ...
        GetCandles(string, Interval, time.Time, time.Time) ([]OHLC, error)
    }
  ```

- Если не хватает типов, нужно их описать и вынести в `types`
  ```go
    // types/candles.go
    package types

    import (
        "time"
    )

    type Interval byte

    type Quant struct {
        // Целая часть цены
        Units int
        // Дробная часть цены
        Nano int
    }

    type OHLC struct {
        // Цена открытия
        Open Quant
        // Максимальная цена за интервал
        High Quant
        // Минимальная цена за интервал
        Low Quant
        // Цена закрытия
        Close Quant
        Time  time.Time
    }
  ```

- Пойти в `TinkoffBrokerPort` и создать там метод, вызвав соответствующий метод из `investgo`
  ```go
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
            log.Errorf("Cannot init sdk! %v", err)
            return []types.OHLC{}, err
        }

        // Сервис для работы с катировками
        candlesService := sdk.NewMarketDataServiceClient()

        // Получаем свечи по инструменту за определенный промежуток времени и интервал (переодичность)
        candlesRes, err := candlesService.GetCandles(instrumentId, investapi.CandleInterval(interval), start, end)
        if err != nil {
			log.Warnff("Cannot get candles %v", err)

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
  ```

  

- Далее создаем gRPC эндпоинт, если его еще нет
  Начнет с создания `.proto` файла в папке `protobuf` в корне проекта
  Или изменения имеющегося файла, если есть подходящий по логике (в нашем кейсе его нет)
    ```proto
    // marketData.proto
    syntax = "proto3";

    package marketData;

    // Название пакета для go
    option go_package = "grpcGW.marketdata";

    import "google/protobuf/timestamp.proto";

    service MarketData {
        // Название нашего эндпоинта
        rpc GetCandles(GetCandlesRequest) returns (GetCandlesResponse);
    }


    // Что ждем от клиента
    message GetCandlesRequest {
        string instrumentId = 1;
        int32 interval = 2;
        google.protobuf.Timestamp start = 3; 
        google.protobuf.Timestamp end = 4; 
        
    }

    // Что хотим вернуть
    message GetCandlesResponse {
        message Quant { 
            int32 units = 1;
            int32 nano = 2;
        }

        message OHLC {
            Quant open = 1;
            Quant high = 2;
            Quant low = 3;
            Quant close = 4;
            google.protobuf.Timestamp time = 5;
        }

        repeated OHLC candles = 1;
    }
    ```

- Запускаем команду генерации сервисов из proto для go (см README.md)
- Идем в `server/server/init.go` и указываем что мы теперь обрабатываем контракт для `marketdata`
  ```go
    // server/server/init.go
    import (
        // ...
	    marketdata "main/grpcGW/grpcGW.marketdata"
        // ...
    )

    type Server struct {
        // ...
        marketdata.UnimplementedMarketDataServer
    }

    func Start(ctx context.Context, port int) {
    //...
    srv := &Server{}

	marketdata.RegisterMarketDataServer(s, srv)
    // ...
  ```
- Создать в папке `server/server` файл для обработки gRPC запроса от клиента
  ```go
  // server/server/marketData.go
    package server

    import (
        "context"
        "main/bot"
        marketdata "main/grpcGW/grpcGW.marketdata"
        "main/types"

        "google.golang.org/protobuf/types/known/timestamppb"
	    log "github.com/sirupsen/logrus"
    )

    // Обьявляем нвоый обработчик эндпоинта GetCandles
    func (s *Server) GetCandles(ctx context.Context, in *marketdata.GetCandlesRequest) (*marketdata.GetCandlesResponse, error) {
        err := bot.Init(ctx, types.Tinkoff)
        if err != nil {
            lot.Warnf("marketdata GetCandles request err %v", err)
            return &marketdata.GetCandlesResponse{Candles: []*marketdata.GetCandlesResponse_OHLC{}}, err
        }

        var res []*marketdata.GetCandlesResponse_OHLC

        // Вызываем созданный ранее сервис
        candles, err := bot.Broker.GetCandles(
            in.InstrumentId,
            types.Interval(in.Interval),
            in.Start.AsTime(),
            in.End.AsTime())

        if err != nil {
            return &marketdata.GetCandlesResponse{Candles: res}, err
        }

        // Мапим в нужный формат
        for _, candle := range candles {
            o := marketdata.GetCandlesResponse_Quant{
                Units: int32(candle.Open.Units),
                Nano:  int32(candle.Open.Nano),
            }
            h := marketdata.GetCandlesResponse_Quant{
                Units: int32(candle.High.Units),
                Nano:  int32(candle.High.Nano),
            }
            l := marketdata.GetCandlesResponse_Quant{
                Units: int32(candle.Low.Units),
                Nano:  int32(candle.Low.Nano),
            }
            c := marketdata.GetCandlesResponse_Quant{
                Units: int32(candle.Close.Units),
                Nano:  int32(candle.Close.Nano),
            }
            res = append(res, &marketdata.GetCandlesResponse_OHLC{
                Open:  &o,
                High:  &h,
                Low:   &l,
                Close: &c,
                Time:  timestamppb.New(candle.Time),
            })
        }

        log.Trace("marketdata GetCandles request")
        return &marketdata.GetCandlesResponse{Candles: res}, nil
    }

  ```

- Готово, теперь сервер обрабатывает новый эндпоинт и ходит за данными в сдк Тинькофф

