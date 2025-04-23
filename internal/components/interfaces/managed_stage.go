package interfaces

import (
	"github.com/google/uuid"
)

type IStage[T any] interface {
	GetID() uuid.UUID
	GetName() string
	GetStageType() string
	GetDescription() string

	GetData() *T
	WithData(data *T) IStage[T]

	GetChannelCtl() chan any
	WithChannelCtl(channelCtl IChannelCtl[any]) IStage[T]

	EventExists(event string) bool
	Dispatch(task func()) error

	GetEvents() map[string]func(...any) any
	GetEvent(event string) func(...any) any

	On(...any) IStage[T]
	Off(...any) IStage[T]

	CheckTransition(fromStage string, toStage string) (bool, error)
	RegisterTransition(fromStage string, toStage string) error

	WithAutoScale(enable bool, limit int, f func(...any) error) IStage[T]
	GetWorkerCount() int

	GetWorkerPool() IWorkerPool
	WithWorkerPool(pool IWorkerPool) IStage[T]

	GetTags() []string
	SetTags(tags []string) IStage[T]

	GetMeta(key string) (any, bool)
	SetMeta(key string, value any) IStage[T]

	GetEventFns() map[string]func(interface{})
	CanTransitionTo(stageID string) bool

	Initialize() error
}
