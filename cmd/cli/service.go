package cli

import (
	"fmt"
	//ci "github.com/faelmori/golife/internal/components/interfaces"
	//. "github.com/faelmori/golife/internal/components/process"
	//"github.com/faelmori/golife/internal/components/types"
	//pi "github.com/faelmori/golife/internal/components/types"
	l "github.com/faelmori/logz"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"syscall"
)

//var manager ci.ILifeCycle[ci.IProcessInput[ci.IManagedProcess[any]]]

func ServiceCmdList() []*cobra.Command {
	return []*cobra.Command{
		lifeCycleManagerCmd(),
		startCommand(),
		stopCommand(),
		statusCommand(),
		restartCommand(),
		serviceCommand(),
	}
}

func lifeCycleManagerCmd() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait, restart bool
	var stages []string
	var triggers []string
	//var processEvents map[string]func(interface{})

	var lCMCmd = &cobra.Command{
		Use:    "lfm",
		Hidden: true,
		Annotations: GetDescriptions([]string{
			"Create a life cycle manager",
			"Create a life cycle manager for the application/process",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			//mgr, mgrErr := createManager(processName, processCmd, stages, processEvents, triggers, processArgs, processWait, restart)
			//if mgrErr != nil {
			//	l.Error(fmt.Sprintf("Fail to create manager: %s", mgrErr), map[string]interface{}{})
			//} else {
			//	manager = mgr
			//}
		},
	}

	lCMCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	lCMCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	lCMCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	lCMCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "MuWait for the process to finish before returning")
	lCMCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Restart the process if it is already running")
	lCMCmd.Flags().StringSliceVarP(&stages, "stages", "s", []string{}, "Stages to listen for and trigger")
	lCMCmd.Flags().StringSliceVarP(&triggers, "triggers", "t", []string{}, "Triggers to listen for and trigger")

	return lCMCmd
}
func startCommand() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait, restart bool
	var stages []string
	var triggers []string
	//var processEvents map[string]func(interface{})

	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"Start a process with a life cycle manager",
			"Start a process with a life cycle manager, optionally waiting for it to finish",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			if processWait {
				//mgr, mgrErr := createManager(processName, processCmd, stages, processEvents, triggers, processArgs, processWait, restart)
				//if mgrErr != nil {
				//	l.Error(fmt.Sprintf("Fail to create manager: %s", mgrErr), map[string]interface{}{})
				//	return
				//} else {
				//	manager = mgr
				//	return
				//}
			} else {
				appFullPath, appFullPathErr := exec.LookPath("golife")
				if appFullPathErr != nil {
					l.Error(fmt.Sprintf("Fail to find golife binary: %s", appFullPathErr), map[string]interface{}{})
					return
				}
				argsStr, waitFlag, restartFlag, stagesStr, triggersStr := getFlagsAsSliceStr(processWait, restart, processArgs, stages, triggers)
				mgrCmdStr := fmt.Sprintf("%s lfm -n %s -c %s %s %s %s -s %s -t %s", appFullPath, processName, processCmd, argsStr, waitFlag, restartFlag, stagesStr, triggersStr)
				mgrCmd := exec.Command("/bin/sh", "-c", mgrCmdStr)
				mgrCmd.Stdout = os.Stdout
				mgrCmd.Stderr = os.Stderr
				mgrCmd.Stdin = os.Stdin
				mgrCmdErr := mgrCmd.Start()
				if mgrCmdErr != nil {
					l.Error(fmt.Sprintf("Fail to start manager command: %s", mgrCmdErr), map[string]interface{}{})
					return
				} else {
					releaseErr := mgrCmd.Process.Release()
					if releaseErr != nil {
						l.Error(fmt.Sprintf("Fail to release process: %s", releaseErr), map[string]interface{}{})
						return
					}
					return
				}
			}
		},
	}

	startCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	startCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	startCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	startCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "MuWait for the process to finish before returning")
	startCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Restart the process if it is already running")
	startCmd.Flags().StringSliceVarP(&stages, "stages", "s", []string{}, "Stages to listen for and trigger")
	startCmd.Flags().StringSliceVarP(&triggers, "triggers", "t", []string{}, "Triggers to listen for and trigger")

	return startCmd
}
func stopCommand() *cobra.Command {
	var processName string
	var stopCmd = &cobra.Command{
		Use: "stop",
		Annotations: GetDescriptions([]string{
			"Stop a process with a life cycle manager",
			"Stop a process with a life cycle manager",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Help()
		},
	}

	stopCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")

	return stopCmd
}
func statusCommand() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use: "status",
		Annotations: GetDescriptions([]string{
			"Get the status of a process with a life cycle manager",
			"Get the status of a process with a life cycle manager. Shows PID, status, and other information.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			//if manager == nil {
			//	l.Error("no manager found", map[string]interface{}{})
			//} else {
			//	l.Info(manager.StatusLifecycle(), map[string]interface{}{})
			//}
		},
	}

	return statusCmd
}
func restartCommand() *cobra.Command {
	var restartCmd = &cobra.Command{
		Use: "restart",
		Annotations: GetDescriptions([]string{
			"Restart a process with a life cycle manager",
			"Restart a process with a life cycle manager",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			//if manager == nil {
			//	l.Error("no manager found", map[string]interface{}{})
			//} else {
			//	if err := manager.RestartLifecycle(); err != nil {
			//		l.Error(fmt.Sprintf("Fail to restart process: %s", err), map[string]interface{}{})
			//	} else {
			//		l.Info("Process restarted successfully", map[string]interface{}{})
			//	}
			//}
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
		Use: "service",
		Annotations: GetDescriptions([]string{
			"Start a process with a life cycle manager (in background)",
			"Start a process with a life cycle manager (in background). This command is used to start a process in the background and detach it from the terminal.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			usrEnvs := os.Environ()
			envPath := os.Getenv("PATH")
			usrEnvs = append(usrEnvs, fmt.Sprintf("PATH=%s", envPath))
			appBinPath, appBinPathErr := exec.LookPath("golife")
			if appBinPathErr != nil {
				l.Error(fmt.Sprintf("Fail to find golife binary: %s", appBinPathErr), map[string]interface{}{})
				return
			}
			cmdStartSpawner := fmt.Sprintf("%s service start -n %s -c %s -a %s -w %t -e %s", appBinPath, processName, processCmd, processArgs, processWait, processEvents)
			cmdStartErr := syscall.Exec("/bin/sh", []string{"-c", cmdStartSpawner}, os.Environ())
			if cmdStartErr != nil {
				l.Error(fmt.Sprintf("Fail to start service: %s", cmdStartErr), map[string]interface{}{})
			}
		},
	}

	serviceCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	serviceCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	serviceCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	serviceCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "MuWait for the process to finish before returning")
	serviceCmd.Flags().StringToStringVarP(&processEvents, "events", "e", map[string]string{}, "Events to listen for and trigger")

	return serviceCmd
}

//	func createManager(processName, processCmd string, stages []string, processEvents map[string]func(interface{}), triggers []string, processArgs []string, processWait, restart bool) (ci.ILifeCycle[ci.IProcessInput[ci.IManagedProcess[any]]], error) {
//		if processName == "" {
//			return nil, fmt.Errorf("no process name provided")
//		}
//		if processCmd == "" {
//			return nil, fmt.Errorf("no command provided")
//		}
//		if len(stages) == 0 {
//			stages = []string{"all"}
//		}
//		if len(processEvents) == 0 {
//			processEvents = make(map[string]func(interface{}))
//		} else {
//			//for _, trigger := range triggers {
//			//	for _, stage := range stages {
//			//		processEvents[trigger] = func(data interface{}) {
//			//			//manager.Trigger(stage, trigger, data)
//			//		}
//			//	}
//			//}
//		}
//
//		var events []IManagedProcessEvents[any]
//		var processes = make(map[string]IManagedProcess[types.ProcessInput[any]])
//		var iStages = make(map[string]IStage[any])
//		//var sigChan = make(chan os.Signal, 1)
//		//var doneChan = make(chan struct{}, 1)
//		var eventsChan = make(chan interface{}, 1)
//		//var eventsCh = make(chan IManagedProcessEvents[any], 1)
//
//		for _, stage := range stages {
//			iStage := NewStage[any](stage, stage, "stage", nil)
//			iStages[stage] = iStage
//		}
//
//		for _, trigger := range triggers {
//			iStage := NewStage[any](trigger, trigger, "trigger", nil)
//			iStages[trigger] = iStage
//		}
//
//		for _, stage := range iStages {
//			for _, trigger := range triggers {
//				stage.OnEvent(trigger, func(data interface{}) {
//					eventsChan <- trigger
//				})
//			}
//		}
//
//		processes[processName] = NewManagedProcess(processName, processCmd, processArgs, processWait, nil)
//		if processEvents != nil {
//			iEvent := NewManagedProcessEvents[any]() //(processEvents, eventsChan)
//			events = append(events, iEvent)
//		}
//
//		regProcErr := manager.AddProcess(processName, pi.NewSystemProcessInput[any](
//			processName,
//			processCmd,
//			processArgs,
//			processWait,
//			restart,
//			nil,
//			nil,
//			false,
//		))
//		if regProcErr != nil {
//			return nil, regProcErr
//		}
//
//		//for _, stage := range iStages {
//		//	defStageErr := manager.DefineStage(stage.Name())
//		//	if defStageErr != nil {
//		//		return nil, defStageErr
//		//	}
//		//}
//
//		startAllErr := manager.StartLifecycle()
//		if startAllErr != nil {
//			return nil, startAllErr
//		}
//
//		return manager, manager.ListenForSignals()
//	}
func getFlagsAsSliceStr(processWait, restart bool, processArgs, stages, triggers []string) (string, string, string, string, string) {
	waitFlag := ""
	if processWait {
		waitFlag = "-w"
	}
	restartFlag := ""
	if restart {
		restartFlag = "-r"
	}
	argsFlags := make([]string, 0)
	for _, arg := range processArgs {
		argsFlags = append(argsFlags, fmt.Sprintf("-a %s", arg))
	}
	argsStr := ""
	if len(argsFlags) > 0 {
		argsStr = fmt.Sprintf("%s", argsFlags)
	}
	stagesFlag := make([]string, 0)
	for _, stage := range stages {
		stagesFlag = append(stagesFlag, fmt.Sprintf("-s %s", stage))
	}
	stagesStr := ""
	if len(stagesFlag) > 0 {
		stagesStr = fmt.Sprintf("%s", stagesFlag)
	}
	triggersFlag := make([]string, 0)
	for _, trigger := range triggers {
		triggersFlag = append(triggersFlag, fmt.Sprintf("-t %s", trigger))
	}
	triggersStr := ""
	if len(triggersFlag) > 0 {
		triggersStr = fmt.Sprintf("%s", triggersFlag)
	}
	return waitFlag, restartFlag, argsStr, stagesStr, triggersStr
}
