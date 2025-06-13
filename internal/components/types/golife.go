package types

import (
	"fmt"
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	//iu "github.com/rafa-mori/golife/internal/utils"
	gl "github.com/rafa-mori/golife/logger"
	l "github.com/rafa-mori/logz"
	"reflect"
)

// GoLife is a generic struct that implements the IGoLife interface.
type GoLife[T any, P ci.ILifeCycle[T, ci.IProperty[ci.IProcessInput[T]]]] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	ci.IReference
	*Reference
	// Mutexes is the mutexes for this GoLife instance.
	ci.IMutexes
	*Mutexes
	// Object is the object to pass to the command.
	Object *T
	// Properties is a map of properties for this GoLife instance.
	properties map[string]interface{}
	// metadata is a map of metadata for this GoLife instance.
	metadata map[string]interface{}
}

// NewGoLife creates a new GoLife instance with the provided Logger.
func NewGoLife[T any, P ci.ILifeCycle[T, ci.IProperty[ci.IProcessInput[T]]]](input *ci.IProperty[ci.IProcessInput[T]], logger l.Logger, debug bool) *GoLife[T, P] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	gl.SetDebug(debug)
	var lcm ci.ILifeCycle[T, ci.IProperty[ci.IProcessInput[T]]]

	lcm = NewLifeCycle[T, ci.IProperty[ci.IProcessInput[T]]](input, logger)

	if inp, ok := lcm.(ci.ILifeCycle[T, ci.IProperty[ci.IProcessInput[T]]]); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type ILifeCycle[ProcessInput[any]] for test %s", inp.GetName()), nil)
		return nil
	}
	if lcmT, ok := lcm.(T); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type %s for test %s", reflect.TypeOf(lcm).String(), lcm.GetName()), nil)
		return nil
	} else {
		chCtl := NewChannelCtl[string]("goLife", logger)
		chMon := NewChannelCtl[string]("goLifeMonitor", logger)
		ggl := GoLife[T, P]{
			Logger:    logger,
			Mutexes:   NewMutexesType(),
			Reference: NewReference("GoLife").GetReference(),
			Object:    &lcmT,
			properties: map[string]interface{}{
				"lifeCycle": WithProperty(lcm.GetName(), &lcm, true, func(any) (bool, error) {
					// Will create the callback function for the lifecycle manager
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Lifecycle manager callback executed")
					return true, nil
				}),
				"chCtl": WithProperty[ci.IChannelCtl[string]]("channel", &chCtl, true, func(any) (bool, error) {
					// Will create the callback function for the channel
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Channel callback executed")
					return true, nil
				}),
				"telemetry": WithProperty("telemetry", NewTelemetry(), true, func(any) (bool, error) {
					// Will create the callback function for the telemetry
					// This is a placeholder for the actual callback logic
					gl.Log("debug", "Telemetry callback executed")
					return true, nil
				}),
				"chMon": WithProperty("monitor", &chMon, true, func(any) (bool, error) {
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
func (g *GoLife[T, P]) GetProperty(name string) (*P, bool) {
	if prop, ok := g.properties[name]; ok {
		if propObj, ok := prop.(*P); ok {
			return propObj, true
		} else {
			gl.LogObjLogger(g, "error", "Property is not of type Property")
		}
	}
	return nil, false
}

// initialize initializes the GoLife properties, mutexes, metadata and other fields correctly.
func (g *GoLife[T, P]) initialize() {
	if g.Logger == nil {
		g.Logger = l.GetLogger("GoLife")
	}
	gl.LogObjLogger(g, "notice", "Initializing GoLife instance")
	if g.Reference == nil {
		g.Reference = NewReference("GoLife").GetReference()
	}
	if g.Mutexes == nil {
		g.Mutexes = NewMutexesType()
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
	g.properties["lifeCycle"] = WithProperty[T]("lifeCycle", g.Object, true, func(any) (bool, error) {
		// Will create the callback function for the lifecycle manager
		// This is a placeholder for the actual callback logic
		gl.LogObjLogger(g, "debug", "Lifecycle manager callback executed")
		return true, nil
	})
}

// WithProperty creates a new property with the given name, property, and callback function.
func WithProperty[T any](name string, property *T, withMetrics bool, cb func(any) (bool, error)) ci.IProperty[T] {
	if property == nil {
		property = interface{}(nil).(*T)
	}
	return NewProperty[T](name, property, withMetrics, cb)
}
