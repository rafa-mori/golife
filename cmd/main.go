package main

import (
	"github.com/faelmori/golife/internal"
	l "github.com/faelmori/logz"
	"os"
)

func main() {
	if rootErr := RootCmd().Execute(); rootErr != nil {
		l.GetLogger("GoSpyder").Println(rootErr)
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
		nil,
	)
	if err != nil {
		l.GetLogger("GoSpyder").Printf("Erro ao registrar o broker: %v\n", err)
		return
	}
	if regEvErr := lifecycle.RegisterEvent("error", "running", func(dt interface{}) {
		l.GetLogger("GoSpyder").Printf("Erro ao executar o broker: %v\n", dt)
	}); regEvErr != nil {
		l.GetLogger("GoSpyder").Printf("Erro ao registrar o evento: %v\n", regEvErr)
		return
	}
	if regStErr := lifecycle.DefineStage("running"); regStErr != nil {
		l.GetLogger("GoSpyder").Printf("Erro ao definir o estágio: %v\n", regStErr)
		return
	}
	if err = lifecycle.Start(); err != nil {
		l.GetLogger("GoSpyder").Printf("Erro ao iniciar o lifecycle: %v\n", err)
		lifecycle.Send("error", err.Error())
		return
	}

	l.GetLogger("GoSpyder").Println("Broker registrado e em execução pelo GoLife!")
}
