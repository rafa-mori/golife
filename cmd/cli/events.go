package cli

import (
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
		Long: Banner + `Trigger an event in a stage`,
		Args: cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			data = args[2]
			manager.Trigger(stage, event, data)
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
		Long: Banner + `Register an event in a stage`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			regEvErr := manager.RegisterEvent(stage, event, func(data interface{}) {
				manager.Trigger(stage, event, data)
			})
			if regEvErr != nil {
				return
			}
		},
	}

	cmdRegisterEvent.Flags().StringVarP(&stage, "stage", "s", "", "The stage to register the event in")
	cmdRegisterEvent.Flags().StringVarP(&event, "event", "e", "", "The event to register")

	return cmdRegisterEvent
}

func removeEventCmd() *cobra.Command {
	var stage, event string

	var cmdRemoveEvent = &cobra.Command{
		Use:  "removeEvent [stage] [event]",
		Long: Banner + `Remove an event from a stage`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stage = args[0]
			event = args[1]
			manager.RemoveEvent(stage, event)
		},
	}

	cmdRemoveEvent.Flags().StringVarP(&stage, "stage", "s", "", "The stage to remove the event from")
	cmdRemoveEvent.Flags().StringVarP(&event, "event", "e", "", "The event to remove")

	return cmdRemoveEvent
}

func stopEventsCmd() *cobra.Command {
	var cmdStopEvents = &cobra.Command{
		Use:  "stopEvents",
		Long: Banner + `Stop all events`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			manager.StopEvents()
		},
	}

	return cmdStopEvents
}
