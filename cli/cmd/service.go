package cmd

import (
	"github.com/faelmori/golife/internal"
	"github.com/spf13/cobra"
)

var manager *internal.GWebLifeCycleManager

func ServiceCmdList() []*cobra.Command {

	return []*cobra.Command{
		StartCmd(),
	}
}

func StartCmd() *cobra.Command {
	var processName, processCmd string
	var processArgs []string
	var processWait bool
	var processEvents map[string]string

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
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
