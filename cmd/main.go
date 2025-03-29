package main

import (
	"net/http"
	"os"

	l "github.com/faelmori/logz"
	"github.com/faelmori/golife/api"
)

func main() {
	mux := http.NewServeMux()
	api.RegisterSSEEndpoint(mux)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			l.Error("Error starting web server: "+err.Error(), nil)
			os.Exit(1)
		}
	}()

	if rootErr := RegX().Execute(); rootErr != nil {
		l.Error("Error executing command: "+rootErr.Error(), nil)
		os.Exit(1)
	}
}
