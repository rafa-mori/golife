package api

import i "github.com/faelmori/golife/internal"

type ManagedProcess = i.IManagedProcess

func NewManagedProcess(name string, command string, args []string, waitFor bool) ManagedProcess {
	return i.NewManagedProcess(name, command, args, waitFor)
}
