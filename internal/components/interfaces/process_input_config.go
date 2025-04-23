package interfaces

type IProcessInputConfig interface {
	GetWaitFor() bool
	GetRestart() bool
	GetProcessType() string
	GetMetadata(key string) (any, bool)
	SetMetadata(key string, value any)
}
