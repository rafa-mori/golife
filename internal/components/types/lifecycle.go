package types

import (
	"fmt"
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
	l "github.com/rafa-mori/logz"
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
	if input == nil {
		gl.LogObjLogger(input, "error", "ProcessInput não pode ser nulo")
		return nil
	}

	if logger == nil {
		logger = l.GetLogger("GoLife")
	}

	channelCtl := NewChannelCtl[T]("LifecycleCtl", logger)

	lc := &LifeCycle[T, P]{
		Logger:     logger,
		Mutexes:    NewMutexesType(),
		Reference:  NewReference("LifeCycle").GetReference(),
		Object:     input,
		Metadata:   make(map[string]any),
		channelCtl: channelCtl,
	}

	if rawCtl, rawCtlType, rawCtlOk := channelCtl.GetSubChannelByName("ctl"); !rawCtlOk {
		gl.LogObjLogger(lc, "fatal", fmt.Sprintf("Control channel does not exist: %s", rawCtlType))
		return nil
	} else {
		if rawCtlType != reflect.TypeFor[string]() {
			gl.LogObjLogger(lc, "fatal", fmt.Sprintf("Control channel type is not string: %s", rawCtlType))
			return nil
		} else {
			chCtlR := reflect.ValueOf(rawCtl).Interface().(ci.IChannelBase[string])
			chCtlT, chCtlType := chCtlR.GetChannel()
			if reflect.ValueOf(chCtlT).Kind() != reflect.Chan {
				gl.LogObjLogger(lc, "fatal", fmt.Sprintf("Control channel type is not chan: %s", chCtlType))
				return nil
			} else {
				if chCtl, ok := reflect.ValueOf(chCtlT).Interface().(chan string); ok {
					lc.Components = newComponents[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](chCtl, logger)
				} else {
					gl.LogObjLogger(lc, "fatal", fmt.Sprintf("Control channel type is not string: %s", chCtlType))
					return nil
				}
			}
		}
	}
	return lc
}

func NewLifeCycle[T any, P ci.IProperty[ci.IProcessInput[T]]](input *P, logger l.Logger) ci.ILifeCycle[T, P] {
	return newLifeCycle[T, P](input, logger)
}

// GetConfig returns the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) GetConfig() *P {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return nil
	}

	lc.Mutexes.MuRLock()
	defer lc.Mutexes.MuRUnlock()

	obj := *lc.Object
	if err := obj.LoadFromFile("config.json", "json"); err != nil {
		gl.LogObjLogger(lc, "error", "Erro ao carregar configuração: %s", err.Error())
	}

	return lc.Object
}

// SetConfig sets the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) SetConfig(processInput P) {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return
	}
	lc.Mutexes.MuLock()
	defer lc.Mutexes.MuUnlock()

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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return nil, false
	}

	lc.Mutexes.MuRLock()
	defer lc.Mutexes.MuRUnlock()

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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return nil
	}
	lc.Mutexes.MuRLock()
	defer lc.Mutexes.MuRUnlock()
	return lc.Logger
}

// SetLogger sets the logger of the lifecycle.
func (lc *LifeCycle[T, P]) SetLogger(logger l.Logger) {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return
	}
	lc.Mutexes.MuLock()
	defer lc.Mutexes.MuUnlock()
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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

	if !reflect.ValueOf(lc.Object).IsValid() {
		return fmt.Errorf("erro: O ProcessInput está vazio")
	}

	if vResult := lc.ValidateLifecycle(); !vResult.GetIsValid() {
		gl.LogObjLogger(&lc, "error", "Lifecycle não é válido. %s", vResult.GetMessage())
		return vResult.GetError()
	}

	lc.Mutexes.MuLock()
	defer lc.Mutexes.MuUnlock()

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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

	gl.LogObjLogger(&lc, "info", "Desligando Lifecycle...")

	// Os mutexes serão usados dentro dos responsáveis, aqui seria lock só de leitura
	//lc.Mutexes.MuRLock()
	//defer lc.Mutexes.MuRUnlock()

	// Depois vou separar essas etapas designando cada uma ao seu responsável (Manager)

	// Só de exemplo, remove o processo "processName". Pra encher linguiça.. hahaha
	if rawProp, ok := lc.GetComponent("processManager"); !ok {
		gl.LogObjLogger(&lc, "error", "Erro ao encerrar processos: %s", "ProcessManager não encontrado")
		return fmt.Errorf("erro: ProcessManager não encontrado")
	} else {
		if propPtr, ok := rawProp.(ci.IProperty[ci.IProcessInput[T]]); ok {
			gl.LogObjLogger(&lc, "Encerrando processos...")
			propWrpr := propPtr.GetValue()
			if prop, ok := propWrpr.(ci.IProcessInput[T]); ok {
				if err := prop.Send("stop", func(msg T) {
					if chCtlWrp, chCtlWrpType, chCtlWrpOk := lc.channelCtl.GetSubChannelByName("ctl"); chCtlWrpOk && chCtlWrpType == reflect.TypeFor[string]() {
						if chCtlWrpObj, chCtlWrpObjOk := reflect.ValueOf(chCtlWrp).Interface().(ci.IChannelBase[string]); !chCtlWrpObjOk {
							gl.LogObjLogger(&lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
							return
						} else {
							chCtlWrpChan, _ := chCtlWrpObj.GetChannel()
							if reflect.ValueOf(chCtlWrpChan).Kind() != reflect.Chan {
								gl.LogObjLogger(&lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
								return
							} else {
								chCtlWrpChan.(chan T) <- msg
							}
						}
					} else {
						gl.LogObjLogger(&lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
					}
				}); err != nil {
					gl.LogObjLogger(&lc, "error", "Erro ao enviar mensagem para o processo: %s", err.Error())
				}
			}
		}
	}
	if rawProp, ok := lc.GetComponent("stageManager"); !ok {
		gl.LogObjLogger(&lc, "error", "Erro ao encerrar stages: %s", "StageManager não encontrado")
		return fmt.Errorf("erro: StageManager não encontrado")
	} else {
		if prop, ok := rawProp.(ci.IStageManager); ok {
			gl.LogObjLogger(&lc, "Encerrando stages...")
			// Enchendo linguiça de novo, remove o stage "stageName"
			_ = prop.RemoveStage("stageName")
		}
	}
	if rawProp, ok := lc.GetComponent("eventManager"); !ok {
		gl.LogObjLogger(&lc, "error", "Erro ao encerrar eventos: %s", "EventManager não encontrado")
		return fmt.Errorf("erro: EventManager não encontrado")
	} else {
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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return nil
	}

	lc.Mutexes.MuRLock()
	defer lc.Mutexes.MuRUnlock()

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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

	gl.LogObjLogger(lc, "success", "Iniciando Lifecycle...")

	if err := lc.Initialize(); err != nil {
		return err
	}

	if rawProp, ok := lc.GetComponent("processManager"); ok {
		if propPtr, ok := rawProp.(ci.IProperty[ci.IProcessInput[T]]); ok {
			gl.LogObjLogger(lc, "info", "Iniciando processos...")
			propWrpr := propPtr.GetValue()
			if prop, ok := propWrpr.(ci.IProcessInput[T]); ok {
				if err := prop.Send("start", func(msg T) {
					if chCtlWrp, chCtlWrpType, chCtlWrpOk := lc.channelCtl.GetSubChannelByName("ctl"); chCtlWrpOk && chCtlWrpType == reflect.TypeFor[string]() {
						if chCtlWrpObj, chCtlWrpObjOk := reflect.ValueOf(chCtlWrp).Interface().(ci.IChannelBase[string]); !chCtlWrpObjOk {
							gl.LogObjLogger(lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
							return
						} else {

							chCtlWrpChan, _ := chCtlWrpObj.GetChannel()
							vlCtlWrpChan := reflect.ValueOf(chCtlWrpChan)
							if vlCtlWrpChan.Kind() != reflect.Chan || !vlCtlWrpChan.IsValid() {
								gl.LogObjLogger(lc, "error", "Erro ao enviar callback para control channel do Lifecycle")
								return
							} else {
								chCtlWrpChan.(chan T) <- msg
							}
						}
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
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

	gl.LogObjLogger(lc, "warning", "Parando Lifecycle...")

	return lc.Shutdown()
}

// RestartLifecycle restarts the lifecycle.
func (lc *LifeCycle[T, P]) RestartLifecycle() error {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

	gl.LogObjLogger(lc, "info", "Reiniciando Lifecycle...")

	if err := lc.StopLifecycle(); err != nil {
		return err
	}

	return lc.StartLifecycle()
}

// StatusLifecycle returns the status of the lifecycle.
func (lc *LifeCycle[T, P]) StatusLifecycle() string {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return ""
	}

	if lc.Object == nil {
		return "Lifecycle não está inicializado"
	}

	return "Lifecycle em execução"
}

// ValidateConfig validates the configuration of the lifecycle.
func (lc *LifeCycle[T, P]) ValidateConfig() error {
	if lc == nil {
		gl.LogObjLogger(lc, "fatal", "Lifecycle não pode ser nulo")
		return fmt.Errorf("erro: Lifecycle não pode ser nulo")
	}

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
//	lc.Mutexes.MuLock()
//	defer lc.Mutexes.MuUnlock()
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
