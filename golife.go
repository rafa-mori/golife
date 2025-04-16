package golife

import (
	t "github.com/faelmori/golife/internal/types"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
	"sync"
)

type GoLife[T any] struct {
	id uuid.UUID

	mu   sync.RWMutex
	muL  sync.RWMutex
	wg   sync.WaitGroup
	cond *sync.Cond

	// properties is a map of string keys to DynamicProperty values.
	properties map[string]t.Property[T]

	// The size is implicitly defined with the new instance of the interface IChannel.
	chanCtl c.IChannel[t.IJob[any], int]

	// meta is a map of string keys to EventMetadata values.
	meta map[string]*t.EventMetadata
}

func NewGoLife[T any]() *GoLife[T] {
	return &GoLife[T]{
		id: uuid.New(),
		mu: sync.RWMutex{},
		wg: sync.WaitGroup{},

		properties: make(map[string]t.Property[T]),
	}
}

func (g *GoLife[T]) GetID() uuid.UUID {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.id
}
func (g *GoLife[T]) GetMeta() t.EventMetadata {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.meta
}
func (g *GoLife[T]) SetMeta(meta t.EventMetadata) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.meta) > 0 {
		for k, v := range g.meta {
			if _, ok := meta[k]; !ok {
				meta[k] = v
			}
		}
	}
}
