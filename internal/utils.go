package internal

import (
	"bufio"
	"fmt"
	pr "github.com/faelmori/golife/components/process"
	p "github.com/faelmori/golife/components/types"
	ev "github.com/faelmori/golife/internal/routines/taskz/events"
	st "github.com/faelmori/golife/internal/routines/taskz/stage"
	gl "github.com/faelmori/golife/logger"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	mutex = p.NewMutexes()
	ref   *p.Reference
)

func initializeChannels[T p.ProcessInput[any]](lm *LifeCycle[T]) error {
	if lm.controllers["channels"] == nil {
		gl.LogObjLogger(lm, "error", "Channels are nil")
		return fmt.Errorf("channels are nil")
	}

	bufSm, bufMd, bufLg := GetDefaultBufferSizes()
	channels := lm.controllers["channels"].(map[string]interface{})
	channels["chanCtl"] = p.NewChannelCtl[any]("chanCtl", &bufSm, lm.Logger)
	channels["chanProcess"] = p.NewChannelCtl[any]("chanProcess", &bufLg, lm.Logger)
	channels["chanStage"] = p.NewChannelCtl[any]("chanStage", &bufMd, lm.Logger)
	channels["chanEvent"] = p.NewChannelCtl[any]("chanEvent", &bufLg, lm.Logger)
	channels["chanSignal"] = p.NewChannelCtl[os.Signal]("chanSignal", &bufSm, lm.Logger)
	channels["chanDone"] = p.NewChannelCtl[bool]("chanDone", &bufSm, lm.Logger)
	channels["chanExit"] = p.NewChannelCtl[bool]("chanExit", &bufSm, lm.Logger)
	channels["chanError"] = p.NewChannelCtl[error]("chanError", &bufSm, lm.Logger)
	channels["chanMessage"] = p.NewChannelCtl[any]("chanMessage", &bufLg, lm.Logger)
	lm.controllers["channels"] = channels
	return nil
}
func initializeProcess[T p.ProcessInput[any]](lm *LifeCycle[T]) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	if lm.controllers["processes"] == nil {
		gl.LogObjLogger(lm, "fatal", "Processes are nil")
		return fmt.Errorf("processes are nil")
	}

	if mainProcessInput, ok := reflect.ValueOf(lm.ProcessInput).Interface().(*p.ProcessInput[any]); !ok {
		gl.LogObjLogger(lm, "fatal", fmt.Sprintf("Error casting main process input. type: %s", reflect.TypeOf(lm.ProcessInput)))
		return fmt.Errorf(fmt.Sprintf("Error casting main process input type: %s", reflect.TypeOf(lm.ProcessInput)))
	} else {
		if mainProcessInput == nil {
			gl.LogObjLogger(lm, "fatal", "Main process input is nil")
			return fmt.Errorf("main process input is nil")
		} else {
			if mainProcessInput.IsRunning {
				gl.LogObjLogger(lm, "error", "Main process is already running")
				return fmt.Errorf("main process is already running")
			}

			processes := lm.controllers["processes"].(map[string]interface{})

			if _, exists := processes[mainProcessInput.Name]; !exists {
				gl.LogObjLogger(lm, "debug", fmt.Sprintf("Instantiating main process '%s'", mainProcessInput.Name))
				newProcess := pr.NewManagedProcess(
					mainProcessInput.Name,
					mainProcessInput.GetCommand(),
					mainProcessInput.GetArgs(),
					mainProcessInput.GetWaitFor(),
					mainProcessInput.GetFunction(),
				)
				processes[mainProcessInput.Name] = newProcess

				lm.controllers["currentProcess"] = newProcess

				lm.controllers["processes"] = processes
			} else {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Main process '%s' already exists!", mainProcessInput.Name))
				return fmt.Errorf("main process '%s' already exists", mainProcessInput.Name)
			}
		}
	}
	return nil
}
func initializeStages[T p.ProcessInput[any]](lm *LifeCycle[T]) error {
	if lm.controllers["stages"] == nil {
		gl.LogObjLogger(lm, "fatal", "Stages are nil")
		return fmt.Errorf("stages are nil")
	}

	stages := lm.controllers["stages"].(map[string]any)

	baseStages := setBaseStages(lm)

	if baseStages == nil {
		gl.LogObjLogger(lm, "error", "Base stages are nil")
		return fmt.Errorf("base stages are nil")
	}

	for _, stage := range baseStages {
		if _, exists := stages[stage.Name()]; exists {
			gl.LogObjLogger(lm, "warn", fmt.Sprintf("Stage '%s' already registered: %v", stage.Name(), stage.GetStageID()))
			continue
		}
		stages[stage.Name()] = stage
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Added stage '%s'", stage.Name()))
	}

	lm.controllers["stages"] = stages

	gl.LogObjLogger(lm, "success", "Lifecycle stages initialized successfully!")
	return nil
}
func initializeEvents[T p.ProcessInput[any]](lm *LifeCycle[T]) error {
	if lm.controllers["events"] == nil {
		gl.LogObjLogger(lm, "fatal", "Events are nil")
		return fmt.Errorf("events are nil")
	}
	events := getBaseEvents()
	for name, event := range events {
		if err := lm.AddEvent(name, event); err != nil {
			gl.LogObjLogger(lm, "error", err.Error())
		}
	}
	return nil
}
func GetDefaultBufferSizes() (sm, md, lg int) {
	return 2, 5, 10
}
func setBaseStages[T p.ProcessInput[any]](lm *LifeCycle[T]) map[string]st.IStage[any] {
	stageMap := map[string]st.IStage[any]{
		"init":    st.NewStage[any]("init", "Initialization stage", "base", nil),
		"execute": st.NewStage[any]("execute", "Execution stage", "base", nil),
		"end":     st.NewStage[any]("end", "End stage", "base", nil),
	}

	// End stage can transition back to init (for restart)
	stageMap["end"] = stageMap["end"].
		SetPossibleNext([]string{"init"}).
		SetPossiblePrev([]string{"execute"})

	stageMap["execute"] = stageMap["execute"].
		SetPossibleNext([]string{"end"}).
		SetPossiblePrev([]string{"init"}).
		OnEnter(func() {
			gl.Log("debug", fmt.Sprintf("Stage '%s' is executing...", stageMap["execute"].Name()))
			gl.Log("debug", fmt.Sprintf("Stage '%s' is executing processes...", stageMap["execute"].Name()))
		}).
		OnEvent("start", func(msg interface{}) {
			gl.LogObjLogger(lm, "debug", fmt.Sprintf("Stage '%s' started!", stageMap["execute"].Name()))
			gl.LogObjLogger(lm, "debug", fmt.Sprintf("Executing processes..."))

			// Assuming msg is a string with the process name, or will use the current process
			processName := ""
			ok := false
			if msg != nil {
				if processName, ok = msg.(string); ok {
					processName = msg.(string)
				}
			} else {
				processName = "currentProcess"
			}
			var process pr.IManagedProcess[p.ProcessInput[any]]
			var processT any
			if processT = lm.controllers[processName]; processT == nil {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Process '%s' not found!", processName))
				return
			} else {
				if process, ok = processT.(process.IManagedProcess[p.ProcessInput[any]]); ok && process != nil {
					gl.LogObjLogger(lm, "info", fmt.Sprintf("Starting process '%s' in stage '%s'", processName, stageMap["execute"].Name()))
					if err := process.Start(); err != nil {
						gl.LogObjLogger(lm, "error", fmt.Sprintf("Error starting process '%s': %v", processName, err))
					} else {
						gl.LogObjLogger(lm, "info", fmt.Sprintf("Process '%s' started successfully!", processName))
					}
				} else {
					gl.LogObjLogger(lm, "error", fmt.Sprintf("Error casting process '%s' to IManagedProcess: %v", processName, processT))
					return
				}
			}
		}).
		OnEvent("stop", func(msg interface{}) {
			gl.Log("debug", fmt.Sprintf("Stage '%s' stopped!", stageMap["execute"].Name()))
			endEvent := stageMap["end"].GetEvent("start")
			if endEvent != nil {
				gl.Log("debug", fmt.Sprintf("Triggering 'start' event in stage: '%s'", stageMap["end"].Name()))
				endEvent(nil)
			}
		})

	// Transition mapping between stages
	stageMap["init"] = stageMap["init"].
		SetPossibleNext([]string{"execute"}).
		OnEnter(func() {
			gl.Log("debug", fmt.Sprintf("Stage '%s' started!", stageMap["init"].Name()))
			gl.Log("debug", fmt.Sprintf("Stage '%s' is initializing...", stageMap["init"].Name()))
		}).
		OnEvent("start", func(msg interface{}) {
			gl.Log("debug", fmt.Sprintf("Stage '%s' is starting...", stageMap["init"].Name()))
			gl.Log("debug", fmt.Sprintf("Stage '%s' is initializing processes...", stageMap["init"].Name()))
		}).
		OnEvent("stop", func(msg interface{}) {
			gl.Log("debug", fmt.Sprintf("Stage '%s' is stopping...", stageMap["init"].Name()))
			startExecute := stageMap["execute"].GetEvent("start")
			if startExecute != nil {
				gl.Log("debug", fmt.Sprintf("Triggering 'start' event in stage: '%s'", stageMap["execute"].Name()))
				startExecute(msg) // Passa mensagem ao próximo estágio
			} else {
				gl.Log("debug", fmt.Sprintf("No 'start' event found in stage: '%s'", stageMap["execute"].Name()))
			}
		})

	return stageMap
}
func getBaseEvents() map[string]ev.IManagedProcessEvents[any] {
	return map[string]ev.IManagedProcessEvents[any]{
		"start": ev.NewManagedProcessEvents[ev.IManagedProcessEvents[any]](),
		"stop":  ev.NewManagedProcessEvents[ev.IManagedProcessEvents[any]](),
	}
}

func listenStdin(ppp ILifeCycle[p.ProcessInput[any]], chErr chan error, chDone chan bool, chSignal chan os.Signal) {
	gl.LogObjLogger(&ppp, "success", "Lifecycle started successfully!")

	if ref == nil {
		ref = p.NewReference("golife")
	} else {
		return
	}

	go func() {
		gl.LogObjLogger(&ppp, "notice", "Lifecycle started listening for stdin input...")
		chMessage := readStdin(ppp, chErr, chDone)
		for {
			select {
			case input := <-chMessage:
				if input == "" {
					continue
				}
				if input == "quit" || input == "q" {
					gl.LogObjLogger(&ppp, "info", "Stopping stdin listener...")
					chErr <- fmt.Errorf("stdin listener stopped")
					break
				}
			case sig := <-chSignal:
				if sig == nil {
					continue
				}
				handleSignal(ppp, chDone, chErr)
			case <-time.After(600 * time.Millisecond):
				processChannels(ppp, chErr, chDone, chSignal)
			case <-time.After(1800 * time.Millisecond):
				checkChannelErrors(ppp, chErr, chDone, chSignal)
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	<-chDone
	gl.LogObjLogger(&ppp, "notice", "Lifecycle stopped listening for stdin input.")
}

// 🔹 Função para tratar sinais do sistema (Ctrl+C, kill, etc.)
func handleSignal(ppp ILifeCycle[p.ProcessInput[any]], chDone chan bool, chErr chan error) {
	if err := ppp.StopLifecycle(); err != nil {
		gl.LogObjLogger(&ppp, "error", "Error stopping lifecycle:", err.Error())
		chErr <- err
	}
	chDone <- true
}

// 🔹 Função para ler a entrada do usuário via stdin
func readStdin(ppp ILifeCycle[p.ProcessInput[any]], chErr chan error, chDone chan bool) <-chan string {
	chMessage := handleInput(ppp, chErr, chDone)

	go func() {
		stdin := os.Stdin
		if stdin == nil {
			gl.LogObjLogger(&ppp, "error", "Error opening stdin")
			chErr <- fmt.Errorf("error opening stdin")
			return
		}
		reader := bufio.NewReader(stdin)
		if reader == nil {
			gl.LogObjLogger(&ppp, "error", "Error creating buffered reader")
			chErr <- fmt.Errorf("error creating buffered reader")
			return
		}
		for {
			select {
			case <-time.After(600 * time.Millisecond):
				inputT, err := reader.ReadString('\n')
				if err != nil {
					gl.LogObjLogger(&ppp, "error", "Error reading stdin:", err.Error())
					chErr <- err
					return
				}
				if input := strings.TrimSpace(inputT); input == "" {
					continue
				} else {
					chMessage <- input
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return chMessage
}

// 🔹 Função para processar comandos do usuário via stdin
func handleInput(ppp ILifeCycle[p.ProcessInput[any]], chErr chan error, chDone chan bool) chan string {
	chMessage := make(chan string, 2)
	go func() {
		defer close(chMessage)

		for {
			select {
			case input := <-chMessage:
				switch input {
				case "stop", "quit", "exit", "q":
					gl.LogObjLogger(&ppp, "notice", "Received stop command, stopping lifecycle...")
					chDone <- true
					break
				case "start":
					gl.LogObjLogger(&ppp, "notice", "Received start command, starting lifecycle...")
					if err := ppp.StartLifecycle(); err != nil {
						gl.LogObjLogger(&ppp, "error", "Error starting lifecycle:", err.Error())
						chErr <- err
					}
					break
				case "status":
					status := ppp.StatusLifecycle()
					gl.LogObjLogger(&ppp, "info", fmt.Sprintf("Lifecycle status: %s", status))
					break
				case "restart":
					gl.LogObjLogger(&ppp, "info", "Received restart command, restarting lifecycle...")
					if err := ppp.RestartLifecycle(); err != nil {
						gl.LogObjLogger(&ppp, "error", "Error restarting lifecycle:", err.Error())
						chErr <- err
					}
					break
				default:
					continue
				}
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	return chMessage
}

// 🔹 Função para verificar se os canais estão funcionando corretamente
func processChannels(ppp ILifeCycle[p.ProcessInput[any]], chErr chan error, chDone chan bool, chSignal chan os.Signal) {
	if chErr == nil || chDone == nil || chSignal == nil {
		gl.LogObjLogger(&ppp, "error", "Failed to get channels")
		chErr <- fmt.Errorf("failed to get channels")
	}
}

// 🔹 Função para detectar erros nos canais e evitar bloqueios
func checkChannelErrors(ppp ILifeCycle[p.ProcessInput[any]], chErr chan error, chDone chan bool, chSignal chan os.Signal) {
	if len(chErr) > 0 || len(chDone) > 0 || len(chSignal) > 0 {
		gl.LogObjLogger(&ppp, "error", "Error detected in channels")
		chErr <- fmt.Errorf("error in channels")
	}
}

// TestInitialization tests the initialization of the GoLife instance.
func TestInitialization(glm ILifeCycle[p.ProcessInput[any]]) {
	gl.LogObjLogger(&glm, "info", "Testing GoLife initialization...")
	if glm == nil {
		gl.LogObjLogger(&glm, "error", "Failed to get lifecycle manager")
	} else {
		go func(prc ILifeCycle[p.ProcessInput[any]]) {
			if err := prc.StartLifecycle(); err != nil {
				gl.LogObjLogger(&glm, "error", "Error starting lifecycle:", err.Error())
			}
		}(glm)
		time.Sleep(500 * time.Millisecond)
		gl.LogObjLogger(&glm, "info", "Lifecycle manager started successfully!")
		gl.LogObjLogger(&glm, "notice", "Lifecycle manager listening for stdin input...")
		if err := glm.ListenForTerminalInput(); err != nil {
			gl.LogObjLogger(&glm, "error", "Error listening for stdin input:", err.Error())
		}
		gl.LogObjLogger(&glm, "success", "Lifecycle manager work like a charm!")
		gl.LogObjLogger(&glm, "info", "Quitting main process...")
	}
}
