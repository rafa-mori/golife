package golife

import (
	"fmt"
	p "github.com/faelmori/golife/components/types"
	i "github.com/faelmori/golife/internal"
	iu "github.com/faelmori/golife/internal/utils"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
)

// GoLife is a generic struct that implements the IGoLife interface.
type GoLife[T i.ILifeCycle[p.ProcessInput[any]]] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*p.Reference
	// Mutexes is the mutexes for this GoLife instance.
	*p.Mutexes
	// Object is the object to pass to the command.
	Object *T
	// Properties is a map of properties for this GoLife instance.
	properties map[string]interface{}
	// metadata is a map of metadata for this GoLife instance.
	metadata map[string]interface{}
}

// NewGoLife creates a new GoLife instance with the provided Logger.
func NewGoLife[T i.ILifeCycle[p.ProcessInput[any]]](input *p.ProcessInput[any], logger l.Logger, debug bool) *GoLife[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	gl.SetDebug(debug)
	var lcm i.ILifeCycle[p.ProcessInput[any]]

	lcm = i.NewLifeCycle[p.ProcessInput[any]](input)

	if _, ok := lcm.(i.ILifeCycle[p.ProcessInput[any]]); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type ILifeCycle[ProcessInput[any]] for test %s", input.Name), nil)
		return nil
	}
	if lcmT, ok := lcm.(T); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type %s for test %s", reflect.TypeOf(lcm).String(), input.Name), nil)
		return nil
	} else {
		ggl := GoLife[T]{
			Logger:    logger,
			Mutexes:   p.NewMutexes(),
			Reference: p.NewReference("GoLife"),
			Object:    &lcmT,
			properties: map[string]interface{}{
				"lifeCycle": iu.WithProperty(input.Name, &lcm, true, func(any) (bool, error) {
					// Will create the callback function for the lifecycle manager
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Lifecycle manager callback executed")
					return true, nil
				}),
				"chCtl": iu.WithProperty("channel", p.NewChannelCtl[string]("goLife", func(b int) *int { return &b }(10), logger), true, func(any) (bool, error) {
					// Will create the callback function for the channel
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Channel callback executed")
					return true, nil
				}),
				"telemetry": iu.WithProperty("telemetry", p.NewTelemetry(), true, func(any) (bool, error) {
					// Will create the callback function for the telemetry
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Telemetry callback executed")
					return true, nil
				}),
				"chMon": iu.WithProperty("monitor", p.NewChannelCtl[string]("goLifeMonitor", func(b int) *int { return &b }(10), logger), true, func(any) (bool, error) {
					// Will create the callback function for the monitor
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Monitor callback executed")
					return true, nil
				}),
			},
		}
		// Initialize the GoLife instance
		ggl.initialize()

		// return the GoLife instance
		return &ggl
	}
}

// GetProperty returns the property for the GoLife instance.
func (g *GoLife[T]) GetProperty(name string) (*p.Property[i.ILifeCycle[p.ProcessInput[any]]], bool) {
	if prop, ok := g.properties[name]; ok {
		if propObj, ok := prop.(*p.Property[i.ILifeCycle[p.ProcessInput[any]]]); ok {
			return propObj, true
		} else {
			gl.LogObjLogger[GoLife[T]](g, "error", "Property is not of type Property")
		}
	}
	return nil, false
}

// initialize initializes the GoLife properties, mutexes, metadata and other fields correctly.
func (g *GoLife[T]) initialize() {
	if g.Logger == nil {
		g.Logger = l.GetLogger("GoLife")
	}
	gl.LogObjLogger[GoLife[T]](g, "notice", "Initializing GoLife instance")
	if g.Reference == nil {
		g.Reference = p.NewReference("GoLife")
	}
	if g.Mutexes == nil {
		g.Mutexes = p.NewMutexes()
	}
	arrMap := []map[string]interface{}{g.properties, g.metadata}
	for key, m := range arrMap {
		if m == nil {
			m = make(map[string]interface{})
			arrMap[key] = m
		}
	}
	if g.properties == nil {
		g.properties = make(map[string]interface{})
	}
	//obj := reflect.ValueOf(g.Object).Interface().(i.ILifeCycle[pi.ProcessInput[any]])
	g.properties["lifeCycle"] = iu.WithProperty[T]("lifeCycle", g.Object, true, func(any) (bool, error) {
		// Will create the callback function for the lifecycle manager
		// This is a placeholder for the actual callback logic
		gl.LogObjLogger[GoLife[T]](g, "debug", "Lifecycle manager callback executed")
		return true, nil
	})
}
