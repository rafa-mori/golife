package interfaces

import (
	"github.com/google/uuid"
)

type IStage[T any] interface {
	EventExists(event string) bool
	GetEvent(event string) func(interface{})
	GetEventFns() map[string]func(interface{})
	GetData() interface{}
	CanTransitionTo(stageID string) bool
	OnEnter(fn func()) IStage[T]
	OnExit(fn func()) IStage[T]
	OnEvent(event string, fn func(interface{})) IStage[T]
	AutoScale(size int) IStage[T]
	SetMeta(key string, value any) IStage[T]
	GetMeta(key string) (any, bool)
	SetTags(tags []string) IStage[T]
	GetTags() []string
	SetWorkerPool(pool IWorkerPool) IStage[T]
	GetWorkerPool() IWorkerPool
	GetStageID() uuid.UUID
	GetStageName() string
	GetStageType() string
	GetStageDesc() string
	GetStageData() *T
	SetPossibleNext(stages []string) IStage[T]
	GetPossibleNext() []string
	SetPossiblePrev(stages []string) IStage[T]
	GetPossiblePrev() []string
	GetWorkerCount() int
	Dispatch(task func()) error
	Description() string
	Name() string
	ID() uuid.UUID
}
