package errors

import  (
	e "errors"
)

var (
	// UnknownStrategy Сигнализирует о невалидном ключе стратегии
	UnknownStrategy = e.New("unknown strategy key")
)