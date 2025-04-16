package utils

import (
	"fmt"
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
