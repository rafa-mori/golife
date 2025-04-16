package utils

import (
	"fmt"
	c "github.com/faelmori/golife/internal/routines/chan"
	j "github.com/faelmori/golife/internal/routines/taskz/jobs"
	t "github.com/faelmori/golife/internal/types"
)

// ValidateWorkerLimit valida o limite de workers
func ValidateWorkerLimit(value any) error {
	if limit, ok := value.(int); ok {
		if limit < 0 {
			return fmt.Errorf("worker limit cannot be negative")
		}
	} else {
		return fmt.Errorf("invalid type for worker limit")
	}
	return nil
}

// validateWorkerPool valida o pool de workers
func validateWorkerPool(value any) error {
	if pool, ok := value.(*t.IWorkerPool); ok {
		if pool == nil {
			return fmt.Errorf("worker pool cannot be nil")
		}
	} else {
		return fmt.Errorf("invalid type for worker pool")
	}
	return nil
}

// validateWorkerChannel valida o canal de trabalho
func validateWorkerChannel(value any) error {
	if channel, ok := value.(c.IChannel[j.IJob, int]); ok {
		if channel == nil {
			return fmt.Errorf("worker channel cannot be nil")
		}
	} else {
		return fmt.Errorf("invalid type for worker channel")
	}
	return nil
}
