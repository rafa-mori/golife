package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
	"strings"
)

type LifeCycle[T any, P ci.IProperty[ci.IProcessInput[T]]] struct {
	// Public fields

	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Mutexes is the mutexes for this GoLife instance.
	*Mutexes
	// Reference is the reference ID and name.
	*Reference `json:"reference" yaml:"reference" xml:"reference" toml:"reference" gorm:"reference"`

	// Object is the object to pass to the command.
	Object *P `json:"process_input" yaml:"process_input" xml:"process_input" toml:"process_input" gorm:"process_input"`

	// Components is a map of properties for this GoLife instance.
	Components *Components[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]] `json:"components,omitempty" yaml:"components,omitempty" xml:"components,omitempty" toml:"components,omitempty" gorm:"components,omitempty"`

	// metadata is a map of metadata for this GoLife instance.
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty" xml:"metadata,omitempty" toml:"metadata,omitempty" gorm:"metadata,omitempty"`

	// Private fields

	// channelCtl is the channel control for this GoLife instance.
	channelCtl ci.IChannelCtl[T]
}

func newLifeCycle[T any, P ci.IProperty[ci.IProcessInput[T]]](input *P, logger l.Logger) *LifeCycle[T, P] {
	return &LifeCycle[T, P]{
		Logger:     logger,
		Mutexes:    NewMutexesType(),
		Reference:  NewReference("LifeCycle").GetReference(),
		Object:     input,
		Components: newComponents[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](logger),
		Metadata:   make(map[string]any),
		channelCtl: NewChannelCtl[T]("LifecycleCtl", logger),
	}
}

func NewLifeCycle[T any, P ci.IProperty[ci.IProcessInput[T]]](input *P, logger l.Logger) ci.ILifeCycle[T, P] {
	return newLifeCycle[T, P](input, logger)
}

// GetConfig returns the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) GetConfig() *P {
	lc.MuRLock()
	defer lc.MuRUnlock()

	obj := *lc.Object
	if err := obj.LoadFromFile("config.json", "json"); err != nil {
		gl.LogObjLogger(lc, "error", "Erro ao carregar configuração: %s", err.Error())
	}

	return lc.Object
}

// SetConfig sets the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) SetConfig(processInput P) {
	lc.MuLock()
	defer lc.MuUnlock()

	// Use reflection to check if processInput is nil or something else invalid. Check directly may cause panic with first class types.
	if !reflect.ValueOf(processInput).IsValid() {
		gl.LogObjLogger(lc, "error", "ProcessInput não pode ser nulo")
		return
	}

	lc.Object = &processInput

	obj := *lc.Object
	if err := obj.SaveToFile("config.json", "json"); err != nil {
		gl.LogObjLogger(lc, "error", "Erro ao salvar configuração: %s", err.Error())
	}
}

// GetComponent retrieves a component from the lifecycle.
func (lc *LifeCycle[T, P]) GetComponent(name string) (any, bool) {
	lc.MuRLock()
	defer lc.MuRUnlock()

	if lc.Components == nil {
		return nil, false
	}

	if component, ok := lc.Components.GetComponent(name); ok {
		return component, true
	}

	return nil, false
}

// GetLogger returns the logger of the lifecycle.
func (lc *LifeCycle[T, P]) GetLogger() l.Logger {
	lc.MuRLock()
	defer lc.MuRUnlock()
	return lc.Logger
}

// SetLogger sets the logger of the lifecycle.
func (lc *LifeCycle[T, P]) SetLogger(logger l.Logger) {
	lc.MuLock()
	defer lc.MuUnlock()
	if logger == nil {
		lc.Logger = l.GetLogger("GoLife")
		return
	}
	lc.Logger = logger
	lgrConfig := logger.GetConfig()
	gl.LogObjLogger(&lc, "info", "New Logger config:")
	gl.LogObjLogger(&lc, "info", "Output: %s", lgrConfig.Output())
	gl.LogObjLogger(&lc, "info", "Level: %s", lgrConfig.Level())
	gl.LogObjLogger(&lc, "info", "Format: %s", lgrConfig.Format())
}

// Initialize initializes the lifecycle.
func (lc *LifeCycle[T, P]) Initialize() error {
	if lc.Object == nil {
		return fmt.Errorf("erro: O ProcessInput está vazio")
	}

	if vResult := lc.ValidateLifecycle(); !vResult.GetIsValid() {
		gl.LogObjLogger(&lc, "error", "Lifecycle não é válido. %s", vResult.GetMessage())
		return vResult.GetError()
	}

	lc.MuLock()
	defer lc.MuUnlock()

	if lc.Logger == nil {
		lc.Logger = l.GetLogger("GoLife")
	}

	if lc.Metadata == nil {
		lc.Metadata = make(map[string]any)
	}

	gl.LogObjLogger(&lc, "success", "Lifecycle inicializado com sucesso")

	return nil
}

// Shutdown shuts down the lifecycle.
func (lc *LifeCycle[T, P]) Shutdown() error {
	gl.LogObjLogger(&lc, "info", "Desligando Lifecycle...")

	// Os mutexes serão usados dentro dos responsáveis, aqui seria lock só de leitura
	//lc.MuRLock()
	//defer lc.MuRUnlock()

	// Depois vou separar essas etapas designando cada uma ao seu responsável (Manager)

	// Só de exemplo, remove o processo "processName". Pra encher linguiça.. hahaha
	if rawProp, ok := lc.GetComponent("processManager"); ok {
		if propPtr, ok := rawProp.(ci.IProperty[ci.IProcessInput[T]]); ok {
			gl.LogObjLogger(&lc, "Encerrando processos...")
			propWrpr := propPtr.GetValue()
			if prop, ok := propWrpr.(ci.IProcessInput[T]); ok {
				if err := prop.Send("stop", func(msg string) {
					if chCtlWrp, chCtlWrpType, chCtlWrpOk := lc.channelCtl.GetSubChannelByName("ctl"); chCtlWrpOk && chCtlWrpType == reflect.TypeOf(new(string)) {
						chCtlWrp.GetChannel() <- fmt.Sprintf("{\"context\":\"%s\", \"message\":\"%s\"}", lc.GetName(), msg)
					} else {
						gl.LogObjLogger(&lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
					}
				}); err != nil {
					gl.LogObjLogger(&lc, "error", "Erro ao enviar mensagem para o processo: %s", err.Error())
				}
			}
		}
	}
	if rawProp, ok := lc.GetComponent("stageManager"); ok {
		if prop, ok := rawProp.(ci.IStageManager); ok {
			gl.LogObjLogger(&lc, "Encerrando stages...")
			// Enchendo linguiça de novo, remove o stage "stageName"
			_ = prop.RemoveStage("stageName")
		}
	}
	if rawProp, ok := lc.GetComponent("eventManager"); ok {
		if prop, ok := rawProp.(ci.IEventManager); ok {
			gl.LogObjLogger(&lc, "Encerrando eventos...")
			// Deee novo, enchendo linguiça, remove o evento "eventName"
			_ = prop.RemoveEvent("eventName")
		}
	}

	return nil
}

// ValidateLifecycle validates the lifecycle.
func (lc *LifeCycle[T, P]) ValidateLifecycle() ci.IValidationResult {
	lc.MuRLock()
	defer lc.MuRUnlock()

	validations := map[string]func() bool{
		// Critical objects
		"ProcessInput is nil":    func() bool { return lc.Object != nil },
		"Control channel is nil": func() bool { return lc.channelCtl != nil },

		"Components empty": func() bool { return lc.Components != nil },

		"ProcessManager is nil": func() bool { return lc.Components.ProcessManager != nil },
		"EventManager is nil":   func() bool { return lc.Components.EventManager != nil },
		"StageManager is nil":   func() bool { return lc.Components.StageManager != nil },
		"SignalManager is nil":  func() bool { return lc.Components.SignalManager != nil },

		"Lifecycle without mutexes":   func() bool { return lc.Mutexes != nil },
		"Lifecycle without reference": func() bool { return lc.Reference != nil },

		// Complementary objects
		//"Logger not set": func() bool { return lc.Logger != nil },
		//"Metadata empty": func() bool { return len(lc.Metadata) > 0 },
	}

	messageStringBuilder := strings.Builder{}
	var isInvalid bool
	for message, isValid := range validations {
		if !isValid() {
			isInvalid = true
			if messageStringBuilder.Len() > 0 {
				messageStringBuilder.WriteString(message + "\n")
			} else {
				messageStringBuilder.WriteString("Lifecycle is invalid:\n")
				messageStringBuilder.WriteString(message + "\n")
			}
		}
	}

	if isInvalid {
		return &ValidationResult{
			IsValid: false,
			Message: messageStringBuilder.String(),
			Error:   fmt.Errorf("lifecycle validation failed. Error: %s", messageStringBuilder.String()),
		}
	}

	return &ValidationResult{
		IsValid: true,
		Message: "Lifecycle is valid",
		Error:   nil,
	}
}

// StartLifecycle starts the lifecycle.
func (lc *LifeCycle[T, P]) StartLifecycle() error {
	gl.LogObjLogger(lc, "success", "Iniciando Lifecycle...")

	if err := lc.Initialize(); err != nil {
		return err
	}

	if rawProp, ok := lc.GetComponent("processManager"); ok {
		if propPtr, ok := rawProp.(ci.IProperty[ci.IProcessInput[T]]); ok {
			gl.LogObjLogger(lc, "info", "Iniciando processos...")
			propWrpr := propPtr.GetValue()
			if prop, ok := propWrpr.(ci.IProcessInput[T]); ok {
				if err := prop.Send("start", func(msg string) {
					if chCtlWrp, chCtlWrpType, chCtlWrpOk := lc.channelCtl.GetSubChannelByName("ctl"); chCtlWrpOk && chCtlWrpType == reflect.TypeOf(new(string)) {
						chCtlWrp.GetChannel() <- fmt.Sprintf("{\"context\":\"%s\", \"message\":\"%s\"}", lc.GetName(), msg)
					} else {
						gl.LogObjLogger(lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
					}
				}); err != nil {
					gl.LogObjLogger(lc, "error", "Erro ao enviar mensagem para o processo: %s", err.Error())
				}
			}
		}
	}

	return nil
}

// StopLifecycle stops the lifecycle.
func (lc *LifeCycle[T, P]) StopLifecycle() error {
	gl.LogObjLogger(lc, "warning", "Parando Lifecycle...")

	return lc.Shutdown()
}

// RestartLifecycle restarts the lifecycle.
func (lc *LifeCycle[T, P]) RestartLifecycle() error {
	gl.LogObjLogger(lc, "info", "Reiniciando Lifecycle...")

	if err := lc.StopLifecycle(); err != nil {
		return err
	}

	return lc.StartLifecycle()
}

// StatusLifecycle returns the status of the lifecycle.
func (lc *LifeCycle[T, P]) StatusLifecycle() string {
	if lc.Object == nil {
		return "Lifecycle não está inicializado"
	}

	return "Lifecycle em execução"
}

// ValidateConfig validates the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) ValidateConfig() error {
	if lc.Object == nil {
		return fmt.Errorf("erro: configuração do Lifecycle não pode estar vazia")
	}
	return nil
}

//// AddComponent adds a component to the lifecycle.
//func (lc *LifeCycle[T, P]) AddComponent(name string, component any) {
//	lc.Components[name] = component
//}

//// RemoveComponent removes a component from the lifecycle.
//func (lc *LifeCycle[T, P]) RemoveComponent(name string) error {
//	lc.MuLock()
//	defer lc.MuUnlock()
//
//	if _, exists := lc.Components[name]; exists {
//		delete(lc.Components, name)
//	} else {
//		gl.LogObjLogger(lc, "error", "Component %s não encontrado", name)
//		return fmt.Errorf("component %s não encontrado", name)
//	}
//
//	return nil
//}
