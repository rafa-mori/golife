package golife

import (
	"os/exec"
)

type ManagedProcess struct {
	Cmd  *exec.Cmd
	Name string
}

func (p *ManagedProcess) Start() error {
	return p.Cmd.Start()
}

func (p *ManagedProcess) Stop() error {
	if p.Cmd != nil {
		return p.Cmd.Process.Kill()
	}
	return nil
}

func (p *ManagedProcess) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}

func (p *ManagedProcess) IsRunning() bool {
	if p.Cmd != nil {
		return p.Cmd.ProcessState == nil
	}
	return false
}

func (p *ManagedProcess) Pid() int {
	if p.Cmd != nil {
		return p.Cmd.Process.Pid
	}
	return 0
}

//import (
//	"fmt"
//	"os/exec"
//	"runtime"
//	"strconv"
//	"strings"
//	"syscall"
//	"time"
//)
//
//type ManagedProcess struct {
//	Cmd        *exec.Cmd
//	Args       []string
//	Env        []string
//	Name       string                   `json:"name;omitempty" gorm:""`
//	Dir        string                   `json:"dir;omitempty" gorm:""`
//	Common     *ManagedCommonProperties `json:"common;omitempty" gorm:""`
//	Monit      *ManagedMonitProperties  `json:"monit;omitempty" gorm:""`
//	Events     *ManagedProcessEvents    `json:"events;omitempty" gorm:""`
//	Container  *ManagedContainer        `json:"container;omitempty" gorm:""`
//	Serverless *ManagedServerless       `json:"serverless;omitempty" gorm:""`
//	SaaS       *ManagedSaas             `json:"saaS;omitempty" gorm:""`
//	PaaS       *ManagedPaas             `json:"paaS;omitempty" gorm:""`
//	IaaS       *ManagedIaas             `json:"iaaS;omitempty" gorm:""`
//	Spawned    *ManagedSpawned          `json:"spawned;omitempty" gorm:""`
//	Managed    *ManagedManaged          `json:"managed;omitempty" gorm:""`
//	FunctionCA *ManagedFunctionCA       `json:"functionCA;omitempty" gorm:""`
//	Goroutine  *ManagedGoroutine        `json:"goroutine;omitempty" gorm:""`
//}
//
//func (p *ManagedProcess) Start() error {
//	p.Cmd = exec.Command("meu_binario_go", "--config=config.json")
//	return p.Cmd.Start()
//}
//func (p *ManagedProcess) Stop() error {
//	if p.Cmd != nil {
//		return p.Cmd.Process.Kill()
//	}
//	return nil
//}
//func (p *ManagedProcess) Wait() error {
//	if p.Cmd != nil {
//		return p.Cmd.Wait()
//	}
//	return nil
//}
//func (p *ManagedProcess) Restart() error {
//	if p.Cmd != nil {
//		if err := p.Stop(); err != nil {
//			return err
//		}
//		return p.Start()
//	}
//	return nil
//}
//func (p *ManagedProcess) Daemonize() error {
//	if p.Cmd != nil {
//		p.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
//	}
//	return nil
//}
//func (p *ManagedProcess) IsRunning() bool {
//	if p.Cmd != nil {
//		return p.Cmd.ProcessState == nil
//	}
//	return false
//}
//func (p *ManagedProcess) Pid() int {
//	if p.Cmd != nil {
//		return p.Cmd.Process.Pid
//	}
//	return 0
//}
//func (p *ManagedProcess) Signal(sig syscall.Signal) error {
//	if p.Cmd != nil {
//		return p.Cmd.Process.Signal(sig)
//	}
//	return nil
//}
//func (p *ManagedProcess) SetArgs(args []string) {
//	if p.Cmd != nil {
//		p.Cmd.Args = args
//	}
//}
//func (p *ManagedProcess) SetDir(dir string) {
//	if p.Cmd != nil {
//		p.Cmd.Dir = dir
//	}
//}
//func (p *ManagedProcess) SetEnv(env []string) {
//	if p.Cmd != nil {
//		p.Cmd.Env = env
//	}
//}
//func (p *ManagedProcess) Monitor() error {
//	if p.Cmd != nil {
//		return monitorProcess(p.Cmd.Process.Pid)
//	}
//	return nil
//}
//
//func monitorProcess(pid int) error {
//	for {
//		// Coletar informações do sistema
//		var memStats runtime.MemStats
//		runtime.ReadMemStats(&memStats)
//
//		// Usar o comando ps para coletar informações do processo
//		out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid,%cpu,%mem").Output()
//		if err != nil {
//			return err
//		}
//
//		// Parsear a saída do comando ps
//		output := strings.Split(string(out), "\n")
//		if len(output) > 1 {
//			info := strings.Fields(output[1])
//			if len(info) >= 3 {
//				fmt.Printf("PID: %s, CPU: %s%%, MEM: %s%%\n", info[0], info[1], info[2])
//			}
//		}
//
//		// Mostrar informações do uso de memória da aplicação
//		fmt.Printf("Alloc = %v MiB", bToMb(memStats.Alloc))
//		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(memStats.TotalAlloc))
//		fmt.Printf("\tSys = %v MiB", bToMb(memStats.Sys))
//		fmt.Printf("\tNumGC = %v\n", memStats.NumGC)
//
//		time.Sleep(5 * time.Second)
//	}
//}
//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
//
//type ManagedProcessFactory struct{}
//
//func (mpf *ManagedProcessFactory) New() *ManagedProcess { return &ManagedProcess{} }
//func (mpf *ManagedProcessFactory) CopyWith(values map[string]interface{}) *ManagedProcess {
//	//common := ManagedCommonProperties{
//	//	Args: []string{},
//	//	Dir:  "",
//	//	Env:  []string{},
//	//}
//	//monit := ManagedMonitProperties{
//	//	Monitoring: false,
//	//	Interval:   0,
//	//	Timeout:    0,
//	//	Delay:      0,
//	//	Timeouts:   []time.Duration{},
//	//	Delays:     []time.Duration{},
//	//	Running:    false,
//	//	Stopped:    false,
//	//	Failed:     false,
//	//	Success:    false,
//	//}
//	//events := ManagedProcessEvents{
//	//	EventFns:  map[string]func(interface{}){},
//	//	TriggerCh: make(chan interface{}),
//	//}
//	//container := ManagedContainer{
//	//	containerID:   "",
//	//	containerCmd:  "",
//	//	containerArgs: []string{},
//	//	containerEnv:  []string{},
//	//	containerDir:  "",
//	//	containerPort: 0,
//	//	containerHost: "",
//	//	containerUser: "",
//	//	containerPass: "",
//	//}
//	//serverless := ManagedServerless{
//	//	cloudFunctionName:     "",
//	//	cloudFunctionType:     "",
//	//	cloudFunctionCmd:      "",
//	//	cloudFunctionArgs:     []string{},
//	//	cloudFunctionEnv:      []string{},
//	//	cloudFunctionDir:      "",
//	//	cloudFunctionPort:     0,
//	//	cloudFunctionHost:     "",
//	//	cloudFunctionUser:     "",
//	//	cloudFunctionPass:     "",
//	//	cloudFunctionCert:     "",
//	//	cloudFunctionKey:      "",
//	//	cloudFunctionCA:       "",
//	//	cloudFunctionSSL:      false,
//	//	cloudFunctionTLS:      false,
//	//	cloudFunctionAuth:     false,
//	//	cloudFunctionAuthType: "",
//	//	cloudFunctionAuthUser: "",
//	//	cloudFunctionAuthPass: "",
//	//	cloudFunctionAuthCert: "",
//	//	cloudFunctionAuthKey:  "",
//	//	cloudFunctionAuthCA:   "",
//	//	cloudFunctionAuthSSL:  false,
//	//	cloudFunctionAuthTLS:  false,
//	//	cloudFunctionAuthAuth: false,
//	//}
//	//saaS := ManagedSaas{
//	//	serviceName: "",
//	//	serviceType: "",
//	//	serviceCmd:  "",
//	//	serviceArgs: []string{},
//	//	serviceEnv:  []string{},
//	//	serviceDir:  "",
//	//	servicePort: 0,
//	//	serviceHost: "",
//	//	serviceUser: "",
//	//	servicePass: "",
//	//	serviceCert: "",
//	//	serviceKey:  "",
//	//	serviceCA:   "",
//	//	serviceSSL:  false,
//	//	serviceTLS:  false,
//	//	serviceAuth: false,
//	//}
//
//	return &ManagedProcess{
//		Cmd:        nil,
//		Args:       nil,
//		Env:        nil,
//		Name:       "",
//		Dir:        "",
//		Common:     nil,
//		Monit:      nil,
//		Events:     nil,
//		Container:  nil,
//		Serverless: nil,
//		SaaS:       nil,
//		PaaS:       nil,
//		IaaS:       nil,
//		Spawned:    nil,
//		Managed:    nil,
//		FunctionCA: nil,
//		Goroutine:  nil,
//	}
//}
