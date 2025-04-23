package process

import (
	"os"
	"time"
)

// ManagedCommonProperties holds common properties for managed processes.
type ManagedCommonProperties struct {
	Args []string // Arguments
	Dir  string   // Directory
	Env  []string // Environment variables
}

// ManagedContainer represents a managed container process.
type ManagedContainer struct {
	containerID       string   // Container ID
	containerCmd      string   // Container command
	containerArgs     []string // Container arguments
	containerEnv      []string // Container environment variables
	containerDir      string   // Container directory
	containerPort     int      // Container port
	containerHost     string   // Container host
	containerUser     string   // Container user
	containerPass     string   // Container password
	containerCert     string   // Container certificate
	containerKey      string   // Container key
	containerCA       string   // Container CA
	containerSSL      bool     // Container SSL enabled
	containerTLS      bool     // Container TLS enabled
	containerAuth     bool     // Container authentication enabled
	containerAuthType string   // Container authentication type
	containerAuthUser string   // Container authentication user
	containerAuthPass string   // Container authentication password
	containerAuthCert string   // Container authentication certificate
	containerAuthKey  string   // Container authentication key
	containerAuthCA   string   // Container authentication CA
	containerAuthSSL  bool     // Container authentication SSL enabled
	containerAuthTLS  bool     // Container authentication TLS enabled
	containerAuthAuth bool     // Container authentication enabled
}

// ManagedServerless represents a managed serverless process.
type ManagedServerless struct {
	cloudFunctionName     string   // Cloud function name
	cloudFunctionType     string   // Cloud function type
	cloudFunctionCmd      string   // Cloud function command
	cloudFunctionArgs     []string // Cloud function arguments
	cloudFunctionEnv      []string // Cloud function environment variables
	cloudFunctionDir      string   // Cloud function directory
	cloudFunctionPort     int      // Cloud function port
	cloudFunctionHost     string   // Cloud function host
	cloudFunctionUser     string   // Cloud function user
	cloudFunctionPass     string   // Cloud function password
	cloudFunctionCert     string   // Cloud function certificate
	cloudFunctionKey      string   // Cloud function key
	cloudFunctionCA       string   // Cloud function CA
	cloudFunctionSSL      bool     // Cloud function SSL enabled
	cloudFunctionTLS      bool     // Cloud function TLS enabled
	cloudFunctionAuth     bool     // Cloud function authentication enabled
	cloudFunctionAuthType string   // Cloud function authentication type
	cloudFunctionAuthUser string   // Cloud function authentication user
	cloudFunctionAuthPass string   // Cloud function authentication password
	cloudFunctionAuthCert string   // Cloud function authentication certificate
	cloudFunctionAuthKey  string   // Cloud function authentication key
	cloudFunctionAuthCA   string   // Cloud function authentication CA
	cloudFunctionAuthSSL  bool     // Cloud function authentication SSL enabled
	cloudFunctionAuthTLS  bool     // Cloud function authentication TLS enabled
	cloudFunctionAuthAuth bool     // Cloud function authentication enabled
}

// ManagedSaas represents a managed SaaS process.
type ManagedSaas struct {
	serviceName string   // Service name
	serviceType string   // Service type
	serviceCmd  string   // Service command
	serviceArgs []string // Service arguments
	serviceEnv  []string // Service environment variables
	serviceDir  string   // Service directory
	servicePort int      // Service port
	serviceHost string   // Service host
	serviceUser string   // Service user
	servicePass string   // Service password
	serviceCert string   // Service certificate
	serviceKey  string   // Service key
	serviceCA   string   // Service CA
	serviceSSL  bool     // Service SSL enabled
	serviceTLS  bool     // Service TLS enabled
	serviceAuth bool     // Service authentication enabled
}

// ManagedPaas represents a managed PaaS process.
type ManagedPaas struct {
	appName string   // Application name
	appType string   // Application type
	appCmd  string   // Application command
	appArgs []string // Application arguments
	appEnv  []string // Application environment variables
	appDir  string   // Application directory
	appPort int      // Application port
	appHost string   // Application host
	appUser string   // Application user
	appPass string   // Application password
	appCert string   // Application certificate
	appKey  string   // Application key
	appCA   string   // Application CA
	appSSL  bool     // Application SSL enabled
	appTLS  bool     // Application TLS enabled
}

// ManagedIaas represents a managed IaaS process.
type ManagedIaas struct {
	instanceID string        // Instance ID
	provider   string        // Provider
	region     string        // Region
	zone       string        // Zone
	image      string        // Image
	size       string        // Size
	sshKey     string        // SSH key
	sshUser    string        // SSH user
	sshPort    int           // SSH port
	sshPass    string        // SSH password
	sshCmd     string        // SSH command
	sshArgs    []string      // SSH arguments
	sshEnv     []string      // SSH environment variables
	sshDir     string        // SSH directory
	sshTimeout time.Duration // SSH timeout
	sshRetries int           // SSH retries
	sshDelay   time.Duration // SSH delay
}

// ManagedSpawned represents a managed spawned process.
type ManagedSpawned struct {
	spawnedProcess  *os.Process // Spawned process
	spawnedPid      int         // Spawned process PID
	spawnedCmd      string      // Spawned process command
	spawnedArgs     []string    // Spawned process arguments
	spawnedEnv      []string    // Spawned process environment variables
	spawnedDir      string      // Spawned process directory
	spawnedPort     int         // Spawned process port
	spawnedHost     string      // Spawned process host
	spawnedUser     string      // Spawned process user
	spawnedPass     string      // Spawned process password
	spawnedCert     string      // Spawned process certificate
	spawnedKey      string      // Spawned process key
	spawnedCA       string      // Spawned process CA
	spawnedSSL      bool        // Spawned process SSL enabled
	spawnedTLS      bool        // Spawned process TLS enabled
	spawnedAuth     bool        // Spawned process authentication enabled
	spawnedAuthType string      // Spawned process authentication type
	spawnedAuthUser string      // Spawned process authentication user
	spawnedAuthPass string      // Spawned process authentication password
	spawnedAuthCert string      // Spawned process authentication certificate
	spawnedAuthKey  string      // Spawned process authentication key
	spawnedAuthCA   string      // Spawned process authentication CA
	spawnedAuthSSL  bool        // Spawned process authentication SSL enabled
	spawnedAuthTLS  bool        // Spawned process authentication TLS enabled
	spawnedAuthAuth bool        // Spawned process authentication enabled
}

// ManagedManaged represents a managed process.
type ManagedManaged struct {
	managedName     string          // Managed process name
	managedType     string          // Managed process type
	managedCmd      string          // Managed process command
	managedArgs     []string        // Managed process arguments
	managedEnv      []string        // Managed process environment variables
	managedDir      string          // Managed process directory
	managedPort     int             // Managed process port
	managedHost     string          // Managed process host
	managedUser     string          // Managed process user
	managedPass     string          // Managed process password
	managedCert     string          // Managed process certificate
	managedKey      string          // Managed process key
	managedCA       string          // Managed process CA
	managedSSL      bool            // Managed process SSL enabled
	managedTLS      bool            // Managed process TLS enabled
	managedAuth     bool            // Managed process authentication enabled
	managedAuthType string          // Managed process authentication type
	managedAuthUser string          // Managed process authentication user
	managedAuthPass string          // Managed process authentication password
	managedAuthCert string          // Managed process authentication certificate
	managedAuthKey  string          // Managed process authentication key
	managedAuthCA   string          // Managed process authentication CA
	managedAuthSSL  bool            // Managed process authentication SSL enabled
	managedAuthTLS  bool            // Managed process authentication TLS enabled
	managedAuthAuth bool            // Managed process authentication enabled
	managedRestart  bool            // Managed process restart enabled
	managedRetries  int             // Managed process retries
	managedInterval time.Duration   // Managed process interval
	managedTimeout  time.Duration   // Managed process timeout
	managedDelay    time.Duration   // Managed process delay
	managedTimeouts []time.Duration // Managed process timeouts
	managedDelays   []time.Duration // Managed process delays
}

// ManagedFunctionCA represents a managed FaaS process.
type ManagedFunctionCA struct {
	functionName string   // Function name
	functionType string   // Function type
	functionCmd  string   // Function command
	functionArgs []string // Function arguments
	functionEnv  []string // Function environment variables
	functionDir  string   // Function directory
	functionPort int      // Function port
	functionHost string   // Function host
	functionUser string   // Function user
	functionPass string   // Function password
	functionCert string   // Function certificate
	functionKey  string   // Function key
	functionCA   string   // Function CA
}
