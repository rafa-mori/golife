package interfaces

import "time"

type IMutexes interface {
	MuLock()
	MuUnlock()
	MuRLock()
	MuRUnlock()
	MuTryLock() bool
	MuTryRLock() bool

	MuWaitCond()
	MuSignalCond()
	MuBroadcastCond()

	GetMuSharedCtx() any
	SetMuSharedCtx(ctx any)
	GetMuSharedCtxValidate() func(any) (bool, error)
	SetMuSharedCtxValidate(validate func(any) (bool, error))
	MuWaitCondWithTimeout(timeout time.Duration) bool

	MuAdd(delta int)
	MuDone()
	MuWait()
}
