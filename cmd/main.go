package main

import (
	"github.com/faelmori/golife/internal"
	l "github.com/faelmori/golife/internal/log"
	"os"
)

func main() {
	if rootErr := RootCmd().Execute(); rootErr != nil {
		l.Println(rootErr)
		os.Exit(1)
	}
}

func startBroker() {
	// Inicializando o GoLife
	lifecycle := internal.NewLifecycleManager(
		make(map[string]internal.IManagedProcess),
		make(map[string]internal.IStage),
		make(chan os.Signal, 1),
		make(chan struct{}),
		nil,
		make(chan internal.IManagedProcessEvents),
	)

	// Registrando o gkbxsrv como subprocesso
	err := lifecycle.RegisterProcess(
		"gkbxsrv",                     // Nome do processo
		"../coreflux/gkbxsrv/gkbxsrv", // Caminho do binário do broker
		[]string{"broker,start"},      // Parâmetros passados ao broker
		true,                          // Reiniciar em caso de falha
	)
	if err != nil {
		l.Printf("Erro ao registrar o broker: %v\n", err)
		return
	}
	if regEvErr := lifecycle.RegisterEvent("error", "running"); regEvErr != nil {
		l.Printf("Erro ao registrar o evento: %v\n", regEvErr)
		return
	}
	if regStErr := lifecycle.DefineStage("running"); regStErr != nil {
		l.Printf("Erro ao definir o estágio: %v\n", regStErr)
		return
	}
	if err = lifecycle.Start(); err != nil {
		l.Printf("Erro ao iniciar o lifecycle: %v\n", err)
		lifecycle.Send("error", err)
		return
	}

	l.Println("Broker registrado e em execução pelo GoLife!")
}
