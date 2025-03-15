package api

import i "github.com/faelmori/golife/internal"

type Stage = i.IStage

func NewStage(id, name, desc, stageType string) Stage {
	return i.NewStage(id, name, desc, stageType)
}
