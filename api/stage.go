package api

import i "github.com/faelmori/golife/internal"

type Stage = i.IStage

func NewStage(name, desc, stageType string) Stage {
	return i.NewStage(name, desc, stageType)
}
