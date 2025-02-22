package cmd

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
	var processEvents map[string]string

	var startCmd = &cobra.Command{
		Use:  "start",
		Long: Banner + `Start the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager = NewLifecycleManager()
			//if len(stages) > 0 {
			//	processEvents = make(map[string]string)
			//	for _, stage := range stages {
			//		if stage == "all" {
			//			for _, stg := range strings.Split(availableStages, ",") {
			//				processEvents[stg] = stg
			//			}
			//			break
			//		} else {
			//			if strings.Contains(availableStages, stage) {
			//				processEvents[stage] = stage
			//			} else {
			//				return fmt.Errorf("invalid stage: %s", stage)
			//			}
			//		}
			//		processEvents[stage] = stage
			//	}
			//}
			//if processName == "" {
			//	return fmt.Errorf("process name is required")
			//}
			//if processCmd == "" {
			//	return fmt.Errorf("process command is required")
			//}
			//if len(processEvents) == 0 {
			//	processEvents["all"] = "all"
			//} else {
			//	for ev, _ := range processEvents {
			//		if !strings.Contains(availableStages, ev) {
			//			return fmt.Errorf("invalid event: %s", ev)
			//		} else {
			//			processEvents[ev] = ev
			//		}
			//	}
			//}

			regProcErr := manager.RegisterProcess(processName, processCmd, processArgs, restart, processWait)
			if regProcErr != nil {
				return regProcErr
			}

			startAllErr := manager.Start()
			if startAllErr != nil {
				return startAllErr
			}

			return nil
		},
	}

	startCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	startCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	startCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	startCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "Wait for the process to finish before returning")
	startCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Restart the process if it is already running")
	startCmd.Flags().StringSliceVarP(&stages, "stages", "s", []string{}, "Stages to listen for and trigger")
	startCmd.Flags().StringToStringVarP(&processEvents, "events", "e", map[string]string{}, "Events to listen for and trigger")

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
