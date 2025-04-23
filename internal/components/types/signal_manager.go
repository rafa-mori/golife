package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type SignalManager struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*Reference
	// SigChan is the channel for the signal.
	SigChan    chan os.Signal
	channelCtl ci.IChannelCtl[string]
}

// NewSignalManager creates a new SignalManager instance.
func newSignalManager(channelCtl ci.IChannelCtl[string], logger l.Logger) *SignalManager {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	return &SignalManager{
		Logger:     logger,
		Reference:  newReference("SignalManager"),
		SigChan:    make(chan os.Signal, 1),
		channelCtl: channelCtl,
	}
}

// NewSignalManager creates a new SignalManager instance.
func NewSignalManager(channelCtl ci.IChannelCtl[string], logger l.Logger) ci.ISignalManager {
	return newSignalManager(channelCtl, logger)
}

// ListenForSignals sets up the signal channel to listen for specific signals.
func (sm *SignalManager) ListenForSignals() error {
	signal.Notify(sm.SigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range sm.SigChan {
			fmt.Printf("Sinal recebido: %s\n", sig.String())
			if sm.channelCtl != nil {
				if chCtlWrp, chCtlWrpType, chCtlWrpOk := sm.channelCtl.GetSubChannelByName("ctl"); chCtlWrpOk && chCtlWrpType == reflect.TypeOf(new(string)) {
					chCtlWrp.GetChannel() <- fmt.Sprintf("{\"context\":\"%s\", \"message\":\"%s\"}", sm.GetName(), ""+sig.String())
				} else {
					gl.LogObjLogger(&sm, "error", "Erro ao enviar callback para control channel do Lifecycle")
				}
			} else {
				fmt.Println("Canal de controle não definido.")
			}
		}
	}()
	return nil
}

// StopListening stops listening for signals and closes the channel.
func (sm *SignalManager) StopListening() {
	signal.Stop(sm.SigChan) // 🔥 Para de escutar sinais
	close(sm.SigChan)       // 🔥 Fecha o canal para evitar vazamento de goroutines
	gl.LogObjLogger(sm, "info", "Parando escuta de sinais")
}
