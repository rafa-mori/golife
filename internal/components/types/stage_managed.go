package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
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
	channelCtl ci.IChannelCtl[any]

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

func (st *Stage[T]) GetType() string {
	if st == nil {
		gl.Log("error", "GetType: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return ""
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	return st.Type
}

func (st *Stage[T]) GetDescription() string {
	if st == nil {
		gl.Log("error", "GetDescription: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return ""
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	return st.Description
}

func (st *Stage[T]) GetData() *T {
	if st == nil {
		gl.Log("error", "GetData: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	if st.data == nil {
		gl.Log("error", "GetData: data does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}
	value := st.data.GetValue()
	if !reflect.ValueOf(value).IsValid() {
		gl.Log("error", "GetData: data is nil (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}
	return &value
}

func (st *Stage[T]) GetID() uuid.UUID {
	if st == nil {
		gl.Log("error", "GetID: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return uuid.Nil
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	return st.Reference.GetID()
}

func (st *Stage[T]) GetName() string {
	if st == nil {
		gl.Log("error", "GetName: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return ""
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	return st.Name
}

func (st *Stage[T]) GetStageType() string {
	if st == nil {
		gl.Log("error", "GetStageType: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return ""
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	return st.Type
}

func (st *Stage[T]) WithData(data *T) ci.IStage[T] {
	if st == nil {
		gl.Log("error", "WithData: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}

	st.Mutexes.MuLock()
	defer st.Mutexes.MuUnlock()

	if st.data == nil {
		gl.Log("error", "WithData: data does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}
	st.data.SetValue(data)
	return st
}

func (st *Stage[T]) GetChannelCtl() chan any {
	if st == nil {
		gl.Log("error", "GetChannelCtl: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}

	st.Mutexes.MuRLock()
	defer st.Mutexes.MuRUnlock()

	if st.channelCtl == nil {
		gl.Log("error", "GetChannelCtl: channelCtl does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}
	if rawChan := st.channelCtl.GetMainChannel(); rawChan != nil {
		if ch, ok := rawChan.(chan any); ok {
			return ch
		}
	}
	gl.Log("error", "GetChannelCtl: channelCtl is not a channel (", reflect.TypeFor[Stage[T]]().String(), ")")
	return nil
}

func (st *Stage[T]) WithChannelCtl(channelCtl ci.IChannelCtl[any]) ci.IStage[T] {
	if st == nil {
		gl.Log("error", "WithChannelCtl: stage does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}

	st.Mutexes.MuLock()
	defer st.Mutexes.MuUnlock()

	if channelCtl == nil {
		gl.Log("error", "WithChannelCtl: channelCtl does not exist (", reflect.TypeFor[Stage[T]]().String(), ")")
		return nil
	}
	st.channelCtl = channelCtl
	return st
}

func (st *Stage[T]) EventExists(event string) bool {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) Dispatch(task func()) error {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetEvents() map[string]func(...any) any {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetEvent(event string) func(...any) any {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) On(args ...any) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) Off(args ...any) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) CheckTransition(fromStage string, toStage string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) RegisterTransition(fromStage string, toStage string) error {
	if st.transitionRegistry == nil {
		st.transitionRegistry = make(map[string]string)
	}
	st.transitionRegistry[fromStage] = toStage
	return nil
}

func (st *Stage[T]) WithAutoScale(enable bool, limit int, f func(...any) error) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetWorkerCount() int {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetWorkerPool() ci.IWorkerPool {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) WithWorkerPool(pool ci.IWorkerPool) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetTags() []string {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) SetTags(tags []string) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetMeta(key string) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) SetMeta(key string, value any) ci.IStage[T] {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) GetEventFns() map[string]func(interface{}) {
	//TODO implement me
	panic("implement me")
}

func (st *Stage[T]) CanTransitionTo(stageID string) bool {
	//TODO implement me
	panic("implement me")
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

	channelCtl := NewChannelCtl[any](name, logger)

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

func NewStage[T any](name, stageType string, data *T, logger l.Logger) ci.IStage[T] {
	return newStage(name, stageType, &data, logger)
}

func (st *Stage[T]) Initialize() error {
	if st.data == nil {
		return fmt.Errorf("data is nil")
	}
	if st.channelCtl == nil {
		return fmt.Errorf("channelCtl is nil")
	}
	if st.workerPool == nil {
		return fmt.Errorf("workerPool is nil")
	}
	return nil
}

func getDefaultEventMap[T any](stage *Stage[T]) (*Stage[T], map[string]any, error) {
	stage.events = make(map[string]any)

	stage.events["onStart"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage iniciado:", stage.Name)
		return NewValidationResult(true, "Stage iniciado corretamente", nil, nil)
	})
	stage.events["onStop"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage encerrado:", stage.Name)
		return NewValidationResult(true, "Stage encerrado corretamente", nil, nil)
	})
	stage.events["onComplete"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Stage completo:", stage.Name)
		return NewValidationResult(true, "Stage completo corretamente", nil, nil)
	})
	stage.events["onError"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Erro capturado:", args)
		return NewValidationResult(false, "Erro identificado", nil, fmt.Errorf("erro desconhecido"))
	})

	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Auto-escalonamento acionado")
		return NewValidationResult(true, "Escalonamento automático acionado", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Pool de trabalhadores acionado")
		return NewValidationResult(true, "Pool de trabalhadores acionado", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Meta dados acionados")
		return NewValidationResult(true, "Meta dados acionados", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Tags acionadas")
		return NewValidationResult(true, "Tags acionadas", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Canal de controle acionado")
		return NewValidationResult(true, "Canal de controle acionado", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Contagem de trabalhadores acionada")
		return NewValidationResult(true, "Contagem de trabalhadores acionada", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Pausa acionada")
		return NewValidationResult(true, "Pausa acionada", nil, nil)
	})
	stage.events["ctl"] = NewValidationFunc[T](1, func(value *T, args ...any) ci.IValidationResult {
		fmt.Println("Retomada acionada")
		return NewValidationResult(true, "Retomada acionada", nil, nil)
	})

	return stage, stage.events, nil
}
