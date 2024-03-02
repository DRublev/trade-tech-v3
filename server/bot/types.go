package bot

import "main/types"

type IMIddleware[T any] interface {
	Do(arg T) error
}

type IPlaceOrderMiddleware interface {
	IMIddleware[types.Order]
}
