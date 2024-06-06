package indicators

// Indicator Описательный тип для всех индикаторов
// T - результат работы индикатора
// S - данные для работы индикатора
type Indicator[T interface{}, S interface{}] interface {
	Latest() (T, error)
	Get() []T
	Update(S)
}
