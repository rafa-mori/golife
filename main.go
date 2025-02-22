package main

import (
	"fmt"
	cli2 "github.com/faelmori/golife/cli"
	"os"
)

func main() {
	cli := cli2.RootCmd()

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
