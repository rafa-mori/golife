package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
)

type Stage[T any] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger

	// Mutexes is the mutexes for this GoLife instance.
	*Mutexes

	// Reference is the reference ID and name.
	*Reference

	// Description and Type are just for documentation/informational purposes.
	Description string
	Type        string

	// data is the object who is being managed by this stage.
	data ci.IProperty[T]

	// channelCtl is the control channel for async operations and communication.
	channelCtl ci.IChannelCtl[string]

	// events is a map of events and their corresponding functions. This map will be treated as map[string]ci.IValidationFunc[T] internally.
	events map[string]any

	// transitionRegistry is a map of possible transitions to other stages.
	transitionRegistry map[string]string

	// possibleNext is a list of possible next stages.
	workerPool ci.IWorkerPool

	// meta is a map of metadata for the stage.
	meta map[string]any

	// tags is a list of tags for the stage.
	tags []string
}

func newStage[T any](name, stageType string, data *T, logger l.Logger) *Stage[T] {
	if name == "" || stageType == "" {
		gl.Log("error", "Stage name and type cannot be empty")
		return nil
	}

	if logger == nil {
		logger = l.NewLogger("Stage")
	}

	var dataProperty ci.IProperty[T]
	if data == nil {
		data = reflect.New(reflect.TypeFor[T]()).Interface().(*T)
	}

	dataProperty = NewProperty[T](name, data, true, func(data any) (bool, error) {
		if data == nil {
			return false, fmt.Errorf("data is nil")
		}
		if _, ok := data.(*T); !ok {
			return false, fmt.Errorf("data is not of type %T", data)
		}
		return true, nil
	})

	channelCtl := NewChannelCtl[string](name, logger)

	return &Stage[T]{
		Logger:    l.NewLogger("Stage"),
		Mutexes:   NewMutexesType(),
		Reference: newReference(name),

		Type:       stageType,
		data:       dataProperty,
		channelCtl: channelCtl,
		events:     make(map[string]any),

		transitionRegistry: make(map[string]string),

		meta: make(map[string]any),
		tags: []string{},
	}
}

func NewStage[T any](name, stageType, desc string, data *T, logger l.Logger) ci.IStage[T] {
	return newStage(name, stageType, desc, data, logger)
}

func (s *Stage[T]) EventExists(event string) bool {
	_, exists := s.events[event]
	return exists
}

func (s *Stage[T]) GetEvent(event string) func(interface{}) {
	return s.events[event]
}

func (s *Stage[T]) GetEventFns() map[string]func(interface{}) {
	return s.events
}

func (s *Stage[T]) CanTransitionTo(stageID string) bool {
	for _, possible := range s.possibleNext {
		if possible == stageID {
			return true
		}
	}
	return false
}

func (s *Stage[T]) Dispatch(task func()) error {
	if s.workerPool == nil {
		return fmt.Errorf("worker pool não foi definido para o Stage %s", s.Name)
	}
	return s.workerPool.Submit(task)
}

func (s *Stage[T]) RegisterTransition(fromStage string, toStage string) error {
	if s.transitionRegistry == nil {
		s.transitionRegistry = make(map[string]string)
	}
	s.transitionRegistry[fromStage] = toStage
	return nil
}

func (s *Stage[T]) Initialize() error {
	if s.data == nil {
		return fmt.Errorf("data is nil")
	}
	if s.channelCtl == nil {
		return fmt.Errorf("channelCtl is nil")
	}
	if s.workerPool == nil {
		return fmt.Errorf("workerPool is nil")
	}
	return nil
}

func getDefaultEventMap[T any](stage *Stage[T]) (*Stage[T], map[string]any, error) {
	stage.events = make(map[string]any)
	stage.events["onStart"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage iniciado:", stage.Name)
		return NewValidationResult(true, "Stage iniciado corretamente", nil)
	})
	stage.events["onStop"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage encerrado:", stage.Name)
		return NewValidationResult(true, "Stage encerrado corretamente", nil)
	})
	stage.events["onComplete"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage completo:", stage.Name)
		return NewValidationResult(true, "Stage completo corretamente", nil)
	})
	stage.events["onError"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Erro capturado:", args)
		return NewValidationResult(false, "Erro identificado", fmt.Errorf("Erro desconhecido"))
	})

	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Auto-escalonamento acionado")
		return NewValidationResult(true, "Escalonamento automático acionado", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Pool de trabalhadores acionado")
		return NewValidationResult(true, "Pool de trabalhadores acionado", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Meta dados acionados")
		return NewValidationResult(true, "Meta dados acionados", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Tags acionadas")
		return NewValidationResult(true, "Tags acionadas", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Canal de controle acionado")
		return NewValidationResult(true, "Canal de controle acionado", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Contagem de trabalhadores acionada")
		return NewValidationResult(true, "Contagem de trabalhadores acionada", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Pausa acionada")
		return NewValidationResult(true, "Pausa acionada", nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Retomada acionada")
		return NewValidationResult(true, "Retomada acionada", nil)
	})

}
