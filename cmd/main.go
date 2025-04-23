package main

import (
	"fmt"
	g "github.com/faelmori/golife"
	i "github.com/faelmori/golife/internal"
	pi "github.com/faelmori/golife/internal/components/process_input"
	p "github.com/faelmori/golife/internal/components/types"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var logger l.Logger

// main initializes the logger and creates a new GoLife instance.
func main() {
	logger = l.GetLogger("GoLife")
	gl.SetDebug(false)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP)
	var glife *g.GoLife[i.ILifeCycle[p.ProcessInput[any]]]
	newArgs := []string{"A"}

	gl.Log("info", "Test A - Start")

	if len(os.Args) > 1 {
		if os.Args[1] == "_supervisor" {
			newArgs = os.Args[2:]
			glife = mainA(newArgs...)
		}
	}

	if glife == nil {
		glife = mainA("A")
	}
	if glife.Object == nil {
		gl.Log("error", "Error creating Lifecycle instance")
		os.Exit(1)
	} else {
		obj := *glife.Object

		// Routine to startup supervisor process without blocking the main thread
		go func(obj i.ILifeCycle[p.ProcessInput[any]]) {
			gl.Log("info", "Starting GoLife process...")
			if err := obj.StartLifecycle(); err != nil {
				gl.Log("error", "starting GoLife instance: "+err.Error())
				os.Exit(1)
			}
		}(obj)

		// Routine to monitor the GoLife process
		for {
			select {
			case <-ch:
				gl.Log("info", "Received interrupt signal, stopping GoLife process...")
				gl.Log("info", "Stopping GoLife process...")
				if err := obj.StopLifecycle(); err != nil {
					gl.Log("error", "stopping GoLife instance: "+err.Error())
				}
				gl.Log("notice", "GoLife process stopped successfully")
				gl.Log("info", "Exiting GoLife process...")
				os.Exit(0)
				return
			case <-time.After(10 * time.Second):
				gl.Log("notice", "Current process: "+obj.StatusLifecycle())
				gl.Log("info", "GoLife process is still running...")
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
		return
	}
}

// main initializes the logger and creates a new GoLife instance.
func mainA(args ...string) *g.GoLife[i.ILifeCycle[p.ProcessInput[any]]] {
	pFunc := p.NewValidation[p.ProcessInput[any]]()

	commandNameArg := "defaultProcess"
	commandArg := ""
	argsArg := make([]string, 0)
	waitForArg := false
	restartArg := false
	debugArg := false
	if len(os.Args[1:]) > 0 {
		gl.Log("info", "Arguments provided:")
		for index, arg := range os.Args[1:] {
			switch arg {
			case "--name", "-n":
				if index+1 < len(args) {
					commandNameArg = args[index+1]
				}
			case "--command", "-c":
				if index+1 < len(args) {
					commandArg = args[index+1]
				}
			case "--args", "-a":
				if index+1 < len(args) {
					argArg := args[index+1]
					if strings.Contains(argArg, ",") {
						aargsArg := strings.Split(argArg, ",")
						for indexArg := 0; indexArg < len(aargsArg); indexArg++ {
							aArg := aargsArg[indexArg]
							if strings.Contains(aArg, "=") {
								argsArg[indexArg] = strings.Split(aArg, "=")[1]
							} else {
								argsArg = append(argsArg, aArg)
							}
						}
					} else {
						argsArg = append(argsArg, argArg)
					}
				}
			case "--wait", "-w":
				waitForArg = true
			case "--restart", "-r":
				restartArg = true
			case "--debug", "-d":
				debugArg = true
			}
		}
	} else {
		gl.Log("error", "No arguments provided, using default values")
		os.Exit(1)
	}

	fn := p.ValidationFunc[p.ProcessInput[any]]{
		Priority: 0,
		Func: func(obj *p.ProcessInput[any], args ...any) *p.ValidationResult {
			objT := reflect.ValueOf(obj).Interface()
			if len(args) > 0 {
				for _, arg := range args {
					if reflect.TypeOf(arg) == reflect.TypeFor[func(*any, ...any) *p.ValidationResult]() {
						return arg.(func(*any, ...any) *p.ValidationResult)(&objT, args...)
					}
				}
				return nil
			}
			return nil
		},
		Result: nil,
	}
	if addPfnErr := pFunc.AddValidator(fn); addPfnErr != nil {
		l.FatalC(fmt.Sprintf("Error adding validation function: %s", addPfnErr.Error()), nil)
		return nil
	}

	input := pi.NewSystemProcessInput[any](
		commandNameArg,
		commandArg,
		argsArg,
		waitForArg,
		restartArg,
		&fn,
		logger,
		debugArg,
	)

	goLife := g.NewGoLife[i.ILifeCycle[p.ProcessInput[any]]](input, logger, false)
	if goLife == nil {
		l.Error("ErrorCtx creating GoLife instance", nil)
		os.Exit(1)
	}
	return goLife
}
