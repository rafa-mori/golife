package golife

import (
	"context"
	"os"
	"sync"
	"time"
)

type ManagedProcessEvents struct {
	// Callbacks
	EventFns map[string]func(interface{})
	// Triggers
	TriggerCh chan interface{}
}
type ManagedMonitProperties struct {
	// Monitoramento
	Monitoring bool
	Interval   time.Duration
	Timeout    time.Duration
	Delay      time.Duration
	Timeouts   []time.Duration
	Delays     []time.Duration

	// Status
	Running bool
	Stopped bool
	Failed  bool
	Success bool
}
type ManagedCommonProperties struct {
	// Configurações
	Args []string
	Dir  string
	Env  []string
}
type ManagedContainer struct {
	// Processo gerenciado (quarta opção - container)
	containerID       string
	containerCmd      string
	containerArgs     []string
	containerEnv      []string
	containerDir      string
	containerPort     int
	containerHost     string
	containerUser     string
	containerPass     string
	containerCert     string
	containerKey      string
	containerCA       string
	containerSSL      bool
	containerTLS      bool
	containerAuth     bool
	containerAuthType string
	containerAuthUser string
	containerAuthPass string
	containerAuthCert string
	containerAuthKey  string
	containerAuthCA   string
	containerAuthSSL  bool
	containerAuthTLS  bool
	containerAuthAuth bool
}
type ManagedServerless struct {
	// Processo gerenciado (nona opção - serverless)
	cloudFunctionName     string
	cloudFunctionType     string
	cloudFunctionCmd      string
	cloudFunctionArgs     []string
	cloudFunctionEnv      []string
	cloudFunctionDir      string
	cloudFunctionPort     int
	cloudFunctionHost     string
	cloudFunctionUser     string
	cloudFunctionPass     string
	cloudFunctionCert     string
	cloudFunctionKey      string
	cloudFunctionCA       string
	cloudFunctionSSL      bool
	cloudFunctionTLS      bool
	cloudFunctionAuth     bool
	cloudFunctionAuthType string
	cloudFunctionAuthUser string
	cloudFunctionAuthPass string
	cloudFunctionAuthCert string
	cloudFunctionAuthKey  string
	cloudFunctionAuthCA   string
	cloudFunctionAuthSSL  bool
	cloudFunctionAuthTLS  bool
	cloudFunctionAuthAuth bool
}
type ManagedSaas struct {
	// Processo gerenciado (oitava opção - SaaS)
	serviceName string
	serviceType string
	serviceCmd  string
	serviceArgs []string
	serviceEnv  []string
	serviceDir  string
	servicePort int
	serviceHost string
	serviceUser string
	servicePass string
	serviceCert string
	serviceKey  string
	serviceCA   string
	serviceSSL  bool
	serviceTLS  bool
	serviceAuth bool
}
type ManagedPaas struct {
	// Processo gerenciado (sexta opção - PaaS)
	appName string
	appType string
	appCmd  string
	appArgs []string
	appEnv  []string
	appDir  string
	appPort int
	appHost string
	appUser string
	appPass string
	appCert string
	appKey  string
	appCA   string
	appSSL  bool
	appTLS  bool
}
type ManagedIaas struct {
	// Processo gerenciado (quinta opção -IaaS)
	instanceID string
	provider   string
	region     string
	zone       string
	image      string
	size       string
	sshKey     string
	sshUser    string
	sshPort    int
	sshPass    string
	sshCmd     string
	sshArgs    []string
	sshEnv     []string
	sshDir     string
	sshTimeout time.Duration
	sshRetries int
	sshDelay   time.Duration
}
type ManagedSpawned struct {
	// Processo gerenciado (segunda opção - spawn process)
	spawnedProcess  *os.Process
	spawnedPid      int
	spawnedCmd      string
	spawnedArgs     []string
	spawnedEnv      []string
	spawnedDir      string
	spawnedPort     int
	spawnedHost     string
	spawnedUser     string
	spawnedPass     string
	spawnedCert     string
	spawnedKey      string
	spawnedCA       string
	spawnedSSL      bool
	spawnedTLS      bool
	spawnedAuth     bool
	spawnedAuthType string
	spawnedAuthUser string
	spawnedAuthPass string
	spawnedAuthCert string
	spawnedAuthKey  string
	spawnedAuthCA   string
	spawnedAuthSSL  bool
	spawnedAuthTLS  bool
	spawnedAuthAuth bool
}
type ManagedManaged struct {
	// Processo gerenciado (décima opção - managed)
	managedName     string
	managedType     string
	managedCmd      string
	managedArgs     []string
	managedEnv      []string
	managedDir      string
	managedPort     int
	managedHost     string
	managedUser     string
	managedPass     string
	managedCert     string
	managedKey      string
	managedCA       string
	managedSSL      bool
	managedTLS      bool
	managedAuth     bool
	managedAuthType string
	managedAuthUser string
	managedAuthPass string
	managedAuthCert string
	managedAuthKey  string
	managedAuthCA   string
	managedAuthSSL  bool
	managedAuthTLS  bool
	managedAuthAuth bool
	managedRestart  bool
	managedRetries  int
	managedInterval time.Duration
	managedTimeout  time.Duration
	managedDelay    time.Duration
	managedTimeouts []time.Duration
	managedDelays   []time.Duration
}
type ManagedFunctionCA struct {
	// Processo gerenciado (sétima opção - FaaS)
	functionName string
	functionType string
	functionCmd  string
	functionArgs []string
	functionEnv  []string
	functionDir  string
	functionPort int
	functionHost string
	functionUser string
	functionPass string
	functionCert string
	functionKey  string
	functionCA   string
}
type ManagedGoroutine struct {
	// Processo gerenciado (terceira opção - goroutine)
	goroutineFn          func()
	goroutineCh          chan struct{}
	goroutineErr         error
	goroutineDone        bool
	goroutineWG          sync.WaitGroup
	goroutineMu          sync.Mutex
	goroutineOnce        sync.Once
	goroutineCond        *sync.Cond
	goroutineLock        sync.Mutex
	goroutineDoneCh      chan struct{}
	goroutineErrCh       chan error
	goroutineCancel      func()
	goroutineCtx         context.Context
	goroutineCancelFn    context.CancelFunc
	goroutineTimeout     time.Duration
	goroutineDeadline    time.Time
	goroutineDeadlineSet bool
}
