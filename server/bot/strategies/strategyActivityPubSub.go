package strategies

// IStrategyActivityPubSub Служит для подписки на действия (activity) стратегий
// Activity - любой майлстоун в логике стратегии, о котором стоит уведомить
// Это может быть определение цены входа в позицию или определение важного для стратегии уровня
type IStrategyActivityPubSub interface {
}

// StrategyActivityPubSub Реализация IStrategyActivityPubSub
type StrategyActivityPubSub struct {
}

var inst *StrategyActivityPubSub

// NewActivityPubSub создает синглтон StrategyActivityPubSub
func NewActivityPubSub() StrategyActivityPubSub {
	if inst == nil {
		inst = &StrategyActivityPubSub{}
	}

	return *inst
}
