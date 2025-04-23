package interfaces

import (
	"github.com/google/uuid"
)

type IStage[T any] interface {
	GetID() uuid.UUID
	GetName() string
	GetType() string
	GetDescription() string

	GetData() *T
	SetData(data *T) IStage[T]

	GetChannelCtl() IChannelCtl[T]
	SetChannelCtl(channelCtl IChannelCtl[T]) IStage[T]

	EventExists(event string) bool
	Dispatch(task func()) error
	GetEvents() map[string]func(...any) any
	GetEvent(event string) func(...any) any

	On(...any) IStage[T]
	Off(...any) IStage[T]

	CheckTransition(fromStage string, toStage string) (bool, error)
	RegisterTransition(fromStage string, toStage string) error

	AutoScale(bool, int, func(...any) error) IStage[T]
	GetWorkerCount() int

	GetWorkerPool() IWorkerPool
	SetWorkerPool(pool IWorkerPool) IStage[T]

	GetTags() []string
	SetTags(tags []string) IStage[T]

	GetMeta(key string) (any, bool)
	SetMeta(key string, value any) IStage[T]
}
