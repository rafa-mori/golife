package interfaces

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

	MuAdd(delta int)
	MuDone()
	MuWait()
}
