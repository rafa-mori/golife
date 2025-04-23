package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/signal"
	"syscall"
)

type SignalManager[T chan string] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*Reference
	// SigChan is the channel for the signal.
	SigChan    chan os.Signal
	channelCtl T
}

// NewSignalManager creates a new SignalManager instance.
func newSignalManager[T chan string](channelCtl T, logger l.Logger) *SignalManager[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	return &SignalManager[T]{
		Logger:     logger,
		Reference:  newReference("SignalManager"),
		SigChan:    make(chan os.Signal, 1),
		channelCtl: channelCtl,
	}
}

// NewSignalManager creates a new SignalManager instance.
func NewSignalManager[T chan string](channelCtl chan string, logger l.Logger) ci.ISignalManager[T] {
	return newSignalManager[T](channelCtl, logger)
}

// ListenForSignals sets up the signal channel to listen for specific signals.
func (sm *SignalManager[T]) ListenForSignals() error {
	signal.Notify(sm.SigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range sm.SigChan {
			fmt.Printf("Sinal recebido: %s\n", sig.String())
			if sm.channelCtl != nil {
				sm.channelCtl <- fmt.Sprintf("{\"context\":\"%s\", \"message\":\"%s\"}", sm.GetName(), ""+sig.String())
			} else {
				fmt.Println("Canal de controle não definido.")
			}
		}
	}()
	return nil
}

// StopListening stops listening for signals and closes the channel.
func (sm *SignalManager[T]) StopListening() {
	signal.Stop(sm.SigChan) // 🔥 Para de escutar sinais
	close(sm.SigChan)       // 🔥 Fecha o canal para evitar vazamento de goroutines
	gl.LogObjLogger(sm, "info", "Parando escuta de sinais")
}
