package strategies

import (
	"sync"
	"time"
)

// StrategyActivityPubSub Реализация IStrategyActivityPubSub
type StrategyActivityPubSub struct {
	containers map[string]IStrategyActivityPubSub
}

var inst *StrategyActivityPubSub

// NewActivityPubSub создает синглтон StrategyActivityPubSub
func NewActivityPubSub() StrategyActivityPubSub {
	if inst == nil {
		inst = &StrategyActivityPubSub{
			containers: make(map[string]IStrategyActivityPubSub),
		}
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
	Track(id string, kind ActivityKind, value interface{})
	GetNewSubscription() *chan Activity
}

type ActivityContainer struct {
	IStrategyActivityPubSub

	mx sync.Mutex
	activities map[string]Activity

	subscribers []*chan Activity
}

func (ac ActivityContainer) GetNewSubscription() *chan Activity {
	subscription := make(chan Activity)
	ac.subscribers = append(ac.subscribers, &subscription)

	go func() {
		for {
			<-time.After(time.Second * 5)
			for _, a := range ac.activities {
				subscription <- a
			}
		}
	}()

	return &subscription
}

func (ac ActivityContainer) Track(id string, kind ActivityKind, value interface{}) {
	act := Activity{
		ID:   id,
		Kind: kind,
	}
	if point, isPoint := value.(PointActivityValue[time.Time, float64]); isPoint {
		act.Value = point
	}
	if level, isLevel := value.(LevelActivityValue); isLevel {
		act.Value = level
	}
	if line, isLine := value.(LineActivityValue[time.Time, float64]); isLine {
		act.Value = line
	}

	ac.mx.Lock()
	ac.activities[id] = act
	ac.mx.Unlock()

	for _, subscription := range ac.subscribers {
		if subscription != nil {
			*subscription <- ac.activities[id]
		}

	}
}

// ActivityKind вид активности: значение, уровень
type ActivityKind string

// Activity - любой майлстоун в логике стратегии, о котором стоит уведомить
// Это может быть определение цены входа в позицию или определение важного для стратегии уровня
type Activity struct {
	ID    string
	Kind  ActivityKind
	Value any
}

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
	X    X
	Y    Y
	Text string
	DeleteFlag bool
}
type LevelActivityValue struct {
	Level float64
	Text  string
	DeleteFlag bool
}
type LineActivityValue[X any, Y any] struct {
	X1   X
	Y1   Y
	X2   X
	Y2   Y
	Text string
	DeleteFlag bool
}
