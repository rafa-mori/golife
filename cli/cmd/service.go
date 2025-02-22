package cmd

import (
	"fmt"
	"github.com/faelmori/golife/internal"
	"github.com/spf13/cobra"
	"os"
	"syscall"
)

var manager *internal.GWebLifeCycleManager

func ServiceCmdList() []*cobra.Command {
	return []*cobra.Command{
		startCmd(),
		stopCmd(),
		statusCmd(),
		restartCmd(),
		serviceCmd(),
	}
}

func startCmd() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait bool
	var processEvents map[string]string

	var startCmd = &cobra.Command{
		Use:  "start",
		Long: Banner + `Start the application`,
		Run: func(cmd *cobra.Command, args []string) {
			iManager := internal.NewLifecycleManager()
			iManager.RegisterProcess(processName, processCmd, processArgs, processWait)

			for ev, stage := range processEvents {
				iManager.RegisterEvent(ev, stage)
			}

			iManager.StartAll()

			manager = &iManager
		},
	}

	startCmd.Flags().StringVarP(&processName, "name", "n", "", "Name of the process")
	startCmd.Flags().StringVarP(&processCmd, "cmd", "c", "", "Command to execute")
	startCmd.Flags().StringSliceVarP(&processArgs, "args", "a", []string{}, "Arguments to pass to the command")
	startCmd.Flags().BoolVarP(&processWait, "wait", "w", false, "Wait for the process to finish before returning")
	startCmd.Flags().StringToStringVarP(&processEvents, "events", "e", map[string]string{}, "Events to listen for and trigger")

	return startCmd
}
func stopCmd() *cobra.Command {
	var stopCmd = &cobra.Command{
		Use:  "stop",
		Long: Banner + `Stop the application`,
		Run: func(cmd *cobra.Command, args []string) {
			if manager == nil {
				return
			} else {
				mgr := *manager
				mgr.StopAll()
			}
		},
	}

	return stopCmd
}
func statusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:  "status",
		Long: Banner + `Check the status of the application`,
		Run: func(cmd *cobra.Command, args []string) {
			if manager == nil {
				return
			} else {
				mgr := *manager
				fmt.Println(mgr.Status())
			}
		},
	}

	return statusCmd
}
func restartCmd() *cobra.Command {
	var restartCmd = &cobra.Command{
		Use:  "restart",
		Long: Banner + `Restart the application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if manager == nil {
				return fmt.Errorf("No manager found")
			} else {
				mgr := *manager
				return mgr.Restart()
			}
		},
	}

	return restartCmd
}
func serviceCmd() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait bool
	var processEvents map[string]string

	var serviceCmd = &cobra.Command{
		Use:  "service",
		Long: Banner + `Manage the application as a service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdStartSpawner := fmt.Sprintf("golife service start -n %s -c %s -a %s -w %t -e %s", processName, processCmd, processArgs, processWait, processEvents)
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
