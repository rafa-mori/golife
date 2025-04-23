package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
)

type StageManager struct {
	stages map[string]ci.IStage[any]
}

func newStageManager() *StageManager {
	return &StageManager{stages: make(map[string]ci.IStage[any])}
}

func NewStageManager() ci.IStageManager { return newStageManager() }

func (sm *StageManager) GetStage(name string) (ci.IStage[any], error) {
	if stage, exists := sm.stages[name]; exists {
		return stage, nil
	}
	return nil, fmt.Errorf("stage %s not found", name)
}

func (sm *StageManager) AddStage(name string, stageObj ci.IStage[any]) error {
	sm.stages[name] = stageObj
	return nil
}

func (sm *StageManager) CurrentStage() ci.IStage[any] {
	for _, stage := range sm.stages {
		return stage // Retorna um qualquer para simplificação
	}
	return nil
}

func (sm *StageManager) SetCurrentStage(name string) (ci.IStage[any], error) {
	return sm.GetStage(name)
}

func (sm *StageManager) RemoveStage(name string) error {
	delete(sm.stages, name)
	return nil
}
