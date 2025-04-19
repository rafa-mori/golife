package main

import (
	g "github.com/faelmori/golife"
	l "github.com/faelmori/logz"

	"net/http"
	"os"
)

func mainB() {
	mux := http.NewServeMux()
	//log.RegisterSSEEndpoint(mux)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			l.Error("ErrorCtx starting web server: "+err.Error(), nil)
			os.Exit(1)
		}
	}()

	if rootErr := RegX().Execute(); rootErr != nil {
		l.Error("ErrorCtx executing command: "+rootErr.Error(), nil)
		os.Exit(1)
	}
}

func main() {
	logger := l.GetLogger("TestLogger")
	goLife := g.NewGoLife[any](logger, false)
	goLife.TestInitialization()
}
