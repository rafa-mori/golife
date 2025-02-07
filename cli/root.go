package cli

import (
	"fmt"
	lcmd "github.com/faelmori/golife/cli/cmd"
	"github.com/spf13/cobra"
)

const Version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "life",
	Aliases: []string{"lfc, goLife"},
	Example: "life start",
	Version: Version,
	Short:   "Programatically can be used to manage the life cycle of an application, service or module",
	Long: `
  #####                #                          
 #     #   ####        #        #  ######  ###### 
 #        #    #       #        #  #       #      
 #  ####  #    #       #        #  #####   #####  
 #     #  #    #       #        #  #       #      
 #     #  #    #       #        #  #       #      
  #####    ####        #######  #  #       ###### 

Life is a CLI tool that can be used to manage the life cycle of an application, service or module.
It can be used to start, stop, restart, pause, resume, trigger events and more.
With the capability to attach to a running process and listen for events, it can be used to orchestrate the life cycle of almost any application.
Hope you enjoy using it as much as I enjoyed creating it. For more information, visit: https://github.com/faelmori/goLife
Happy coding! Happy Life!`,
}

func init() {
	setUsageDefinition(rootCmd)

	rootCmd.AddCommand(lcmd.ServiceCmdList()...)
	rootCmd.AddCommand(lcmd.EventsCmdList()...)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
