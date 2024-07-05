package strategies

// Activity - любой майлстоун в логике стратегии, о котором стоит уведомить
// Это может быть определение цены входа в позицию или определение важного для стратегии уровня
type Activity struct {
	Kind  ActivityKind
	Value any
}

// StrategyActivityPubSub Реализация IStrategyActivityPubSub
type StrategyActivityPubSub struct {
	containers map[string]IStrategyActivityPubSub
}

var inst *StrategyActivityPubSub

// NewActivityPubSub создает синглтон StrategyActivityPubSub
func NewActivityPubSub() StrategyActivityPubSub {
	if inst == nil {
		inst = &StrategyActivityPubSub{}
	}

	return *inst
}

// Container Создает контейнер, обслуживающий все ивенты для containerID
func (sa StrategyActivityPubSub) Container(containerID string) IStrategyActivityPubSub {
	container, exists := sa.containers[containerID]
	if !exists {
		sa.containers[containerID] = ActivityContainer{
			activities: make(map[string]Activity),
		}
		return sa.containers[containerID]
	}

	return container
}

func (sa StrategyActivityPubSub) Subscribe(containerID string) *chan Activity {
	return sa.containers[containerID].GetNewSubscription()
}

// IStrategyActivityPubSub Служит для подписки на действия (Activity) стратегий
type IStrategyActivityPubSub interface {
	Track(id string, kind ActivityKind, value any)
	GetNewSubscription() *chan Activity
}

type ActivityContainer struct {
	IStrategyActivityPubSub

	activities map[string]Activity

	subscribers []*chan Activity
}

func (ac ActivityContainer) GetNewSubscription() *chan Activity {
	subscription := make(chan Activity)
	ac.subscribers = append(ac.subscribers, &subscription)
	return &subscription
}

func (ac ActivityContainer) Track(id string, kind ActivityKind, value any) {
	ac.activities[id] = Activity{
		Kind:  kind,
		Value: value,
	}

	go func() {
		for _, subscription := range ac.subscribers {
			if subscription != nil {
				*subscription <- ac.activities[id]
			}

		}
	}()
}

// ActivityKind вид активности: значение, уровень
type ActivityKind string

// Valid Валидация
func (candidate ActivityKind) Valid() bool {
	switch candidate {
	case "point", "level":
		return true
	default:
		return false
	}
}

type PointActivityValue[X any, Y any] struct {
	X X
	Y Y
}
type LevelActivityValue float64
type LineActivityValue[X any, Y any] struct {
	X X
	Y Y
}
