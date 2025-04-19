package main

import (
	cc "github.com/faelmori/golife/cmd/cli"
	vs "github.com/faelmori/golife/version"
	l "github.com/faelmori/logz"

	"github.com/spf13/cobra"

	"os"
	"strings"
)

type GoLife struct {
	parentCmdName string
	printBanner   bool
}

func (m *GoLife) Alias() string { return "" }
func (m *GoLife) ShortDescription() string {
	return "GoLife is a tool to manage the life cycle of an application, service or module."
}
func (m *GoLife) LongDescription() string {
	return `GoLife: a tool that can be used to manage the life cycle of an application, service or module.
It can be used to start, stop, restart, pause, resume, trigger events and more.
With the capability to attach to a running process and listen for events, it can be used to orchestrate the life cycle of almost any application.
Hope you enjoy using it as much as I enjoyed creating it. For more information, visit: https://github.com/faelmori/goLife
Happy coding! Happy Life!`
}
func (m *GoLife) Usage() string {
	return "golife [command] [args]"
}
func (m *GoLife) Examples() []string {
	return []string{"golife [command] [args]", "golife start -n myProcessName -c 'tail' -a 'f,/dev/null'"}
}
func (m *GoLife) Active() bool {
	return true
}
func (m *GoLife) Module() string {
	return "golife"
}
func (m *GoLife) Execute() error { return m.Command().Execute() }
func (m *GoLife) Command() *cobra.Command {
	l.GetLogger("GoLife")
	l.Info("Starting GoLife...", map[string]interface{}{
		"context":  "main",
		"action":   "start",
		"version":  vs.GetVersion(),
		"showData": false,
	})

	var rtCmd = &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Version: cc.Version,
		Annotations: cc.GetDescriptions([]string{
			m.LongDescription(),
			m.ShortDescription(),
		}, m.printBanner),
	}

	rtCmd.AddCommand(cc.ServiceCmdList()...)
	rtCmd.AddCommand(cc.EventsCmdList()...)

	rtCmd.AddCommand(vs.CliCommand())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(rtCmd)
	for _, c := range rtCmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return rtCmd
}
func (m *GoLife) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *GoLife) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}

func RegX() *GoLife {
	var printBannerV = os.Getenv("GOLIFE_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &GoLife{
		printBanner: strings.ToLower(printBannerV) == "true",
	}
}
