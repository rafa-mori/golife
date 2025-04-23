package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
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
	pm.processes[s] = t.GetValue()
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
