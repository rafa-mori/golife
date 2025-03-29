package main

import (
	l "github.com/faelmori/logz"
	"os"
)

func main() {
	if rootErr := RegX().Execute(); rootErr != nil {
		l.Error("Error executing command: "+rootErr.Error(), nil)
		os.Exit(1)
	}
}
