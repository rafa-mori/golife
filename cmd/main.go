package main

import (
	"fmt"
	g "github.com/faelmori/golife"
	pi "github.com/faelmori/golife/components/process_input"
	i "github.com/faelmori/golife/internal"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var logger l.Logger

// main initializes the logger and creates a new GoLife instance.
func main() {
	//mainProd()
	logger = l.GetLogger("TestLogger")
	gl.SetDebug(false)
	if len(os.Args) > 1 {
		if os.Args[1] == "_supervisor" {
			gl.Log("info", "Launching GoLife process in parallel...")
			newArgs := os.Args[2:]
			willWait := "-w"
			for _, arg := range os.Args {
				if arg == "--wait" || arg == "-w" {
					willWait = ""
					break
				}
			}
			os.Args = append(os.Args, willWait)
			newArgs = append(newArgs, willWait)
			go func() {
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
				go mainA(newArgs...)
				defer func() {
					if r := recover(); r != nil {
						gl.Log("error", fmt.Sprintf("Recovered from panic: %v", r))
					} else {
						// release the routine process
						gl.Log("info", "Releasing GoLife process...")
						mainRoutineRuntime := runtime.NumGoroutine()
						gl.Log("info", fmt.Sprintf("GoLife process released, goroutines: %d", mainRoutineRuntime))
						if mainRoutineRuntime > 1 {
							gl.Log("info", "Waiting for GoLife process to finish...")
							for i := 0; i < mainRoutineRuntime; i++ {
								select {
								case <-ch:
									gl.Log("info", "Received interrupt signal, stopping GoLife process...")
									gl.Log("info", "Stopping GoLife process...")
									return
								default:
									gl.Log("info", fmt.Sprintf("GoLife process %d is still running...", i))
									continue
								}
							}
						} else {
							gl.Log("info", "GoLife process finished.")
						}
					}
				}()
				select {
				case <-ch:
					gl.Log("info", "Received interrupt signal, stopping GoLife process...")
					gl.Log("info", "Stopping GoLife process...")
					return
				}
			}()
		} else {
			gl.Log("info", "Test A - Start")
			mainA("A")
			gl.Log("info", "Test A - End")
		}
	} else {
		gl.Log("info", "Test A - Start")
		mainA("A")
		gl.Log("info", "Test A - End")
	}
}

// main initializes the logger and creates a new GoLife instance.
func mainA(args ...string) {
	goLife := g.NewGoLifeTest[i.ILifeCycle[pi.ProcessInput[any]]](args[0], logger, false)
	if goLife == nil {
		l.Error("ErrorCtx creating GoLife instance", nil)
		return
	}
	if err := goLife.Object.Initialize(); err != nil {
		l.Error("ErrorCtx initializing GoLife instance: "+err.Error(), nil)
		return
	}
	g.TestInitialization[i.ILifeCycle[pi.ProcessInput[any]]](goLife)
}

// main initializes the logger and creates a new GoLife instance.
func mainB() {
	goLife := g.NewGoLifeTest[i.ILifeCycle[pi.ProcessInput[any]]]("B", logger, false)
	if goLife == nil {
		l.Error("ErrorCtx creating GoLife instance", nil)
		return
	}
	if err := goLife.Object.Initialize(); err != nil {
		l.Error("ErrorCtx initializing GoLife instance: "+err.Error(), nil)
		return
	}
	g.TestInitialization[i.ILifeCycle[pi.ProcessInput[any]]](goLife)
}
