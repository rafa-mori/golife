package internal

import (
	"context"
	"fmt"
	"github.com/faelmori/kbx/mods/conf"
	"github.com/faelmori/kbx/mods/logz"
	"github.com/faelmori/kbx/mods/task"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
	"time"
)

func taskzHandler(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	taskReq := "Taskz wrap cmd: " + cmd.Name()

	runner := task.NewTkzRunner(ctx, taskReq)

	runner.Then(func(ctx context.Context, taskReq *string) error {
		//_ = logz.Log("debug", fmt.Sprintf("Taskz: %s", *taskReq), "kbx", "quiet")
		return nil
	})

	_, err := runner.Result()
	if err != nil {
		return logz.ErrorLog(fmt.Sprintf("Taskz error: %s", err.Error()), "kbx")
	}

	return nil
}

func regTimeLog(start time.Time) {
	duration := time.Since(start).Round(time.Millisecond)
	// Só registra o tempo de execução se for maior que 1 segundo
	if duration > 1*time.Second {
		_ = logz.DebugLog("Tempo de execução ("+strings.Join(os.Args, " ")+"): "+duration.String(), "kbx", logz.QUIET)
	}
}

func wrapAllCommands(cmd *cobra.Command) {
	originalRunE := cmd.RunE
	start := time.Now()
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		defer regTimeLog(start)
		if err := taskzHandler(cmd, args); err != nil {
			return err
		}
		if originalRunE != nil {
			commandErr := originalRunE(cmd, args)
			regTimeLog(start)
			return commandErr
		}
		return nil
	}

	// Captura as informações detalhadas das flags do comando atual
	var flags []conf.FlagInfo
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flags = append(flags, conf.FlagInfo{
			Command:      cmd.Name(),
			Name:         flag.Name,
			Shorthand:    flag.Shorthand,
			Usage:        flag.Usage,
			DefaultValue: flag.DefValue,
			Required:     flag.NoOptDefVal == "",
			Active:       flag.Changed,
			Variable:     flag.Value.String(),
			ValueType:    flag.Value.Type(),
		})
	})
	conf.CommandFlagsMap[cmd.Name()] = conf.CommandFlags{
		Command: cmd.Name(),
		Flags:   flags,
	}

	for _, subCmd := range cmd.Commands() {
		wrapAllCommands(subCmd)
	}
}

func handleCommand(cmdRoot *cobra.Command) error {
	_ = filDefaultModules()

	// Percorre os módulos registrados e adiciona os comandos ao comando principal caso sejam ativos e válidos.
	manager := conf.NewManager()
	modsListMap := ModsRegistryMap
	for _, mod := range modsListMap {
		if mod.Active() {
			manager.RegisterModule(mod.Module(), mod)
			if mod.Command() != nil {
				cmdRoot.AddCommand(mod.Command())
			} else {
				_ = logz.Log("error", "no command found: "+mod.Module(), "kbx")
			}
		}
	}

	wrapAllCommands(cmdRoot)

	return nil
}
