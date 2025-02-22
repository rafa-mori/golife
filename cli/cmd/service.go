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

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			manager = internal.NewLifecycleManager()

		},
	}
	return startCmd
}
