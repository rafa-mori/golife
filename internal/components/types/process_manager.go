package types

import (
	"context"
	"fmt"
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
	"reflect"
)

type ProcessManager[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]] struct {
	processes map[string]ci.IProcessInput[ci.IManagedProcess[any]]
}

func newProcessManager[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]() *ProcessManager[T] {
	return &ProcessManager[T]{processes: make(map[string]ci.IProcessInput[ci.IManagedProcess[any]])}
}

func NewProcessManager[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]() ci.IProcessManager[T] {
	return newProcessManager[T]()
}

func (pm *ProcessManager[T]) AddProcess(s string, t T) error {
	if !reflect.ValueOf(t).IsValid() {
		return fmt.Errorf("erro: ProcessInput não pode ser nulo ao adicionar um processo")
	}
	vl := t.GetValue()
	if vl == nil {
		return fmt.Errorf("erro: ProcessInput não pode ser nulo ao adicionar um processo")
	}
	pm.processes[s] = vl
	return nil
}

func (pm *ProcessManager[T]) GetProcess(name string) (ci.IProcessInput[ci.IManagedProcess[any]], error) {
	if process, exists := pm.processes[name]; exists {
		return process, nil
	}
	return nil, fmt.Errorf("process %s not found", name)
}

func (pm *ProcessManager[T]) CurrentProcess() ci.IProcessInput[ci.IManagedProcess[any]] {
	for _, process := range pm.processes {
		return process // Retorna um qualquer para simplificação
	}
	return nil
}

func (pm *ProcessManager[T]) RemoveProcess(name string) error {
	delete(pm.processes, name)
	return nil
}

func (pm *ProcessManager[T]) StartProcess(name string) error {
	process, err := pm.GetProcess(name)
	if err != nil {
		return err
	}
	if process == nil {
		return fmt.Errorf("process %s not found", name)
	}
	if sendErr := process.Send("start", nil); sendErr != nil {
		gl.Log("error", "Error starting process:", name, sendErr.Error())
	} else {
		ctx := context.Background()

		if obj, receiveErr := process.Receive(ctx, func(msg string) {
			gl.Log("info", "Received message from process:", name, msg)
		}); receiveErr != nil {
			gl.Log("error", "Error receiving process:", name, receiveErr.Error())
		} else if obj != nil {
			gl.Log("info", "Received object from process:", name, obj.(string))
		} else {
			gl.Log("info", "No object received from process:", name)
		}
	}
	return nil
}
