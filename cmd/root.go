package cmd

import (
	"github.com/faelmori/golife/cmd/cli"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "golife",
		Aliases: []string{"lfc, goLife"},
		Example: "life start",
		Version: cli.Version,
		Long: cli.Banner + `
Life is a CLI tool that can be used to manage the life cycle of an application, service or module.
It can be used to start, stop, restart, pause, resume, trigger events and more.
With the capability to attach to a running process and listen for events, it can be used to orchestrate the life cycle of almost any application.
Hope you enjoy using it as much as I enjoyed creating it. For more information, visit: https://github.com/faelmori/goLife
Happy coding! Happy Life!`,
	}

	rootCmd.AddCommand(cli.ServiceCmdList()...)
	rootCmd.AddCommand(cli.EventsCmdList()...)

	setUsageDefinition(rootCmd)

	for _, cmd := range rootCmd.Commands() {
		setUsageDefinition(cmd)
	}

	return rootCmd
}
