package bot

import "main/types"

// IMIddleware Общий интерфейс для мидлварей
// TODO: Общий тип. Перенести на уровень проекта
type IMIddleware[T any] interface {
	Do(arg T) error
}

// IPlaceOrderMiddleware Мидлвари, вызываемые перед выставлением обдеда
type IPlaceOrderMiddleware interface {
	IMIddleware[types.Order]
}
