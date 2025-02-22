package cmd

import "github.com/spf13/cobra"

func EventsCmdList() []*cobra.Command {

	return []*cobra.Command{
		triggerCmd(),
		registerEventCmd(),
		removeEventCmd(),
		stopEventsCmd(),
	}
}

func triggerCmd() *cobra.Command {
	// triggerCmd represents the trigger command
	// Here we define vars to be used in the command through flags or arguments

	var cmdTrigger = &cobra.Command{
		Use:  "trigger [stage] [event] [data]",
		Long: Banner + `Trigger an event in a stage`,
		Args: cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			//stage := args[0]
			//event := args[1]
			//data := args[2]
			//manager.Trigger(stage, event, data)
		},
	}

	// Here we define flags to be used in the command
	//triggerCmd.Flags().StringVarP(&stage, "stage", "s", "", "The stage to trigger the event in") //EXAMPLE

	// Here we define arguments to be used in the command
	//triggerCmd.Args = cobra.MinimumNArgs(3) //EXAMPLE

	return cmdTrigger
}
func registerEventCmd() *cobra.Command {
	// registerEventCmd represents the registerEvent command
	// Here we define vars to be used in the command through flags or arguments

	var cmdRegisterEvent = &cobra.Command{
		Use:  "regEvent",
		Long: Banner + `Register an event in a stage`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//stage := args[0]
			//event := args[1]
			//manager.RegisterEvent(stage, event)
		},
	}

	// Here we define flags to be used in the command
	//registerEventCmd.Flags().StringVarP(&stage, "stage", "s", "", "The stage to register the event in") //EXAMPLE

	// Here we define arguments to be used in the command
	//registerEventCmd.Args = cobra.MinimumNArgs(2) //EXAMPLE

	return cmdRegisterEvent
}
func removeEventCmd() *cobra.Command {
	// removeEventCmd represents the removeEvent command
	// Here we define vars to be used in the command through flags or arguments

	var cmdRemoveEvent = &cobra.Command{
		Use:  "removeEvent [stage] [event]",
		Long: Banner + `Remove an event from a stage`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//stage := args[0]
			//event := args[1]
			//manager.RemoveEvent(stage, event)
		},
	}

	// Here we define flags to be used in the command
	//removeEventCmd.Flags().StringVarP(&stage, "stage", "s", "", "The stage to remove the event from") //EXAMPLE

	// Here we define arguments to be used in the command
	//removeEventCmd.Args = cobra.MinimumNArgs(2) //EXAMPLE

	return cmdRemoveEvent
}
func stopEventsCmd() *cobra.Command {
	// stopEventsCmd represents the stopEvents command
	// Here we define vars to be used in the command through flags or arguments

	var cmdStopEvents = &cobra.Command{
		Use:  "stopEvents",
		Long: Banner + `Stop all events`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			//manager.StopEvents()
		},
	}

	// Here we define flags to be used in the command
	//stopEventsCmd.Flags().StringVarP(&stage, "stage", "s", "", "The stage to stop events in") //EXAMPLE

	// Here we define arguments to be used in the command
	//stopEventsCmd.Args = cobra.NoArgs //EXAMPLE

	return cmdStopEvents
}
