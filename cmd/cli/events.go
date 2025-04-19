package cli

import (
	"fmt"
	l "github.com/faelmori/logz"
	"github.com/spf13/cobra"
)

func EventsCmdList() []*cobra.Command {
	return []*cobra.Command{
		triggerCmd(),
		registerEventCmd(),
		removeEventCmd(),
		stopEventsCmd(),
	}
}

func triggerCmd() *cobra.Command {
	var stage, event, data string

	var cmdTrigger = &cobra.Command{
		Use:  "trigger [stage] [event] [data]",
		Args: cobra.MinimumNArgs(3),
		Annotations: GetDescriptions([]string{
			"Trigger an event in a stage",
			"Dispatch an event in a stage",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			data = args[2]
			//manager.Trigger(stage, event, data)
			l.Info(fmt.Sprintf("Event %s triggered in stage %s with data: %s", event, stage, data), map[string]interface{}{})
		},
	}

	cmdTrigger.Flags().StringVarP(&stage, "stage", "s", "", "The stage to trigger the event in")
	cmdTrigger.Flags().StringVarP(&event, "event", "e", "", "The event to trigger")
	cmdTrigger.Flags().StringVarP(&data, "data", "d", "", "The data to pass to the event")

	return cmdTrigger
}

func registerEventCmd() *cobra.Command {
	var stage, event string

	var cmdRegisterEvent = &cobra.Command{
		Use:  "regEvent [stage] [event]",
		Args: cobra.MinimumNArgs(2),
		Annotations: GetDescriptions([]string{
			"Register an event",
			"Register an event in a stage",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			regEvErr := manager.RegisterEvent(stage, event, func(data interface{}) {
				//manager.Trigger(stage, event, data)
			})
			if regEvErr != nil {
				l.Error(fmt.Sprintf("ErrorCtx registering event: %s", regEvErr.Error()), map[string]interface{}{})
				return
			}
			l.Info("Event registered successfully", map[string]interface{}{})
		},
	}

	cmdRegisterEvent.Flags().StringVarP(&stage, "stage", "s", "", "The stage to register the event in")
	cmdRegisterEvent.Flags().StringVarP(&event, "event", "e", "", "The event to register")

	return cmdRegisterEvent
}

func removeEventCmd() *cobra.Command {
	var stage, event string

	var cmdRemoveEvent = &cobra.Command{
		Use: "removeEvent [stage] [event]",
		Annotations: GetDescriptions([]string{
			"Remove an event",
			"Remove an event from a stage",
		}, false),
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			err := manager.RemoveEvent( /*stage, */ event)
			if err != nil {
				l.Error(fmt.Sprintf("ErrorCtx removing event: %s", err.Error()), map[string]interface{}{})
				return
			}
			l.Info("Event removed successfully", map[string]interface{}{})
		},
	}

	cmdRemoveEvent.Flags().StringVarP(&stage, "stage", "s", "", "The stage to remove the event from")
	cmdRemoveEvent.Flags().StringVarP(&event, "event", "e", "", "The event to remove")

	return cmdRemoveEvent
}

func stopEventsCmd() *cobra.Command {
	var cmdStopEvents = &cobra.Command{
		Use:  "stopEvents",
		Args: cobra.NoArgs,
		Annotations: GetDescriptions([]string{
			"Stop events",
			"Stop all events",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.StopEvents()
			if err != nil {
				l.Error(fmt.Sprintf("ErrorCtx stopping events: %s", err.Error()), map[string]interface{}{})
				return
			}
			l.Info("Events stopped successfully", map[string]interface{}{})
		},
	}

	return cmdStopEvents
}
