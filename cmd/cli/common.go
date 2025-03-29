package cli

import (
	"os"
	"strings"
)

const Banner = `
  #####                #                          
 #     #   ####        #        #  ######  ###### 
 #        #    #       #        #  #       #      
 #  ####  #    #       #        #  #####   #####  
 #     #  #    #       #        #  #       #      
 #     #  #    #       #        #  #       #      
  #####    ####        #######  #  #       ###### 
`

const Version = "1.0.0"

func GetDescriptions(descriptionArg []string, _ bool) map[string]string {
	var description string
	if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
		description = descriptionArg[0]
	} else {
		description = descriptionArg[1]
	}

	return map[string]string{"banner": Banner, "description": description}
}
