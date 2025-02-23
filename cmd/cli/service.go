package cli

import (
	"fmt"
	. "github.com/faelmori/golife/internal"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"syscall"
)

var manager GWebLifeCycleManager

const availableStages = "all,prestart,poststart,prestop,poststop"

func ServiceCmdList() []*cobra.Command {
	return []*cobra.Command{
		startCommand(),
		stopCommand(),
		statusCommand(),
		restartCommand(),
		serviceCommand(),
	}
}

func startCommand() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait, restart bool
	var stages []string
	var triggers []string
	var processEvents map[string]func(interface{})

	var startCmd = &cobra.Command{
		Use:  "start",
		Long: Banner + `Start the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, mgrErr := createManager(processName, processCmd, stages, processEvents, triggers, processArgs, processWait, restart)
			if mgrErr != nil {
				return mgrErr
			} else {
				manager = mgr
				return nil
			}
		},
	}

	startCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	startCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	startCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	startCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "Wait for the process to finish before returning")
	startCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Restart the process if it is already running")
	startCmd.Flags().StringSliceVarP(&stages, "stages", "s", []string{}, "Stages to listen for and trigger")
	startCmd.Flags().StringSliceVarP(&triggers, "triggers", "t", []string{}, "Triggers to listen for and trigger")

	return startCmd
}
func stopCommand() *cobra.Command {
	var processName string
	var stopCmd = &cobra.Command{
		Use:  "stop",
		Long: Banner + `Stop the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if manager == nil {
				return fmt.Errorf("no manager found")
			} else {
				stopErr := manager.Stop()
				if stopErr != nil {
					return stopErr
				}
			}
			return nil
		},
	}

	stopCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")

	return stopCmd
}
func statusCommand() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:  "status",
		Long: Banner + `Check the status of the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if manager == nil {
				return fmt.Errorf("no manager found")
			} else {
				fmt.Println(manager.Status())
				return nil
			}
		},
	}

	return statusCmd
}
func restartCommand() *cobra.Command {
	var restartCmd = &cobra.Command{
		Use:  "restart",
		Long: Banner + `Restart the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if manager == nil {
				return fmt.Errorf("no manager found")
			} else {
				return manager.Restart()
			}
		},
	}

	return restartCmd
}
func serviceCommand() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait bool
	var processEvents map[string]string

	var serviceCmd = &cobra.Command{
		Use:  "service",
		Long: Banner + `Manage the application as a service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			usrEnvs := os.Environ()
			envPath := os.Getenv("PATH")
			usrEnvs = append(usrEnvs, fmt.Sprintf("PATH=%s", envPath))
			appBinPath, appBinPathErr := exec.LookPath(AppName)
			if appBinPathErr != nil {
				return appBinPathErr
			}
			cmdStartSpawner := fmt.Sprintf("%s service start -n %s -c %s -a %s -w %t -e %s", appBinPath, processName, processCmd, processArgs, processWait, processEvents)
			cmdStartErr := syscall.Exec("/bin/sh", []string{"-c", cmdStartSpawner}, os.Environ())
			if cmdStartErr != nil {
				return cmdStartErr
			} else {
				return nil
			}
		},
	}

	serviceCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	serviceCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	serviceCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	serviceCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "Wait for the process to finish before returning")
	serviceCmd.Flags().StringToStringVarP(&processEvents, "events", "e", map[string]string{}, "Events to listen for and trigger")

	return serviceCmd
}

func createManager(processName, processCmd string, stages []string, processEvents map[string]func(interface{}), triggers []string, processArgs []string, processWait, restart bool) (GWebLifeCycleManager, error) {
	if processName == "" {
		return nil, fmt.Errorf("no process name provided")
	}
	if processCmd == "" {
		return nil, fmt.Errorf("no command provided")
	}
	if len(stages) == 0 {
		stages = []string{"all"}
	}
	if len(processEvents) == 0 {
		processEvents = make(map[string]func(interface{}))
	} else {
		for _, trigger := range triggers {
			for _, stage := range stages {
				processEvents[trigger] = func(data interface{}) {
					manager.Trigger(stage, trigger, data)
				}
			}
		}
	}

	var events []IManagedProcessEvents
	var processes = make(map[string]IManagedProcess)
	var iStages = make(map[string]IStage)
	var sigChan = make(chan os.Signal, 1)
	var doneChan = make(chan struct{}, 1)
	var eventsChan = make(chan interface{}, 1)
	var eventsCh = make(chan IManagedProcessEvents, 1)

	for _, stage := range stages {
		iStage := NewStage(stage, stage, stage, "stage")
		iStages[stage] = iStage
	}

	for _, trigger := range triggers {
		iStage := NewStage(trigger, trigger, trigger, "trigger")
		iStages[trigger] = iStage
	}

	for _, stage := range iStages {
		for _, trigger := range triggers {
			stage.OnEvent(trigger, func(data interface{}) {
				eventsChan <- trigger
			})
		}
	}

	processes[processName] = NewManagedProcess(processName, processCmd, processArgs, processWait)
	if processEvents != nil {
		iEvent := NewManagedProcessEvents(processEvents, eventsChan)
		events = append(events, iEvent)
	}

	manager = NewLifecycleManager(
		processes,
		iStages,
		sigChan,
		doneChan,
		events,
		eventsCh,
	)

	regProcErr := manager.RegisterProcess(processName, processCmd, processArgs, restart)
	if regProcErr != nil {
		return nil, regProcErr
	}

	for _, stage := range iStages {
		manager.DefineStage(stage.Name())
	}

	go func() {
		startAllErr := manager.Start()
		if startAllErr != nil {
			fmt.Println("Erro ao iniciar processos:", startAllErr)
			return
		}

		listenErr := manager.ListenForSignals()
		if listenErr != nil {
			fmt.Println("Erro ao ouvir sinais:", listenErr)
			return
		}
		<-doneChan
	}()

	select {
	case ev := <-eventsChan:
		if ev != nil {
			for _, event := range events {
				event.RegisterEvent(event.Event(), func(data interface{}) {
					event.Trigger(data.(string), event.Event(), data)
				})
			}
		}
	case <-doneChan:
		if len(processes) > 0 {
			startAllErr := manager.StartAll()
			if startAllErr != nil {
				return nil, startAllErr
			}
		}
	case <-sigChan:
		sigChan <- syscall.SIGTERM
		stopAllErr := manager.StopAll()
		if stopAllErr != nil {
			return nil, stopAllErr
		}
	}

	return manager, nil
}
