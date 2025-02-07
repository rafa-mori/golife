package cmd

import "github.com/spf13/cobra"

func EventsCmdList() []*cobra.Command {

	return []*cobra.Command{
		TriggerCmd(),
	}
}

func TriggerCmd() *cobra.Command {
	// triggerCmd represents the trigger command
	// Here we define vars to be used in the command through flags or arguments

	var cmdTrigger = &cobra.Command{
		Use:   "trigger [stage] [event] [data]",
		Short: "Trigger an event in a stage",
		Args:  cobra.MinimumNArgs(3),
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
