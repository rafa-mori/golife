package golife

import (
	"fmt"
	pi "github.com/faelmori/golife/components/process_input"
	p "github.com/faelmori/golife/components/types"
	i "github.com/faelmori/golife/internal"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"reflect"
	"strings"
	"time"
)

// GoLife is a generic struct that implements the IGoLife interface.
type GoLife[T i.ILifeCycle[pi.ProcessInput[any]]] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*p.Reference
	// Mutexes is the mutexes for this GoLife instance.
	*p.Mutexes
	// Object is the object to pass to the command.
	Object T
	// Properties is a map of properties for this GoLife instance.
	properties map[string]interface{}
	// metadata is a map of metadata for this GoLife instance.
	metadata map[string]interface{}
	// layers is a map of layers for this GoLife instance.
	layers map[string]interface{}
}

// NewGoLifeTest creates a new GoLife instance with the provided Logger.
func NewGoLifeTest[T i.ILifeCycle[pi.ProcessInput[any]]](test string, logger l.Logger, debug bool) *GoLife[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	gl.SetDebug(debug)

	pFunc := p.NewValidation[pi.ProcessInput[any]]()

	commandNameArg := "defaultProcess"
	commandArg := ""
	argsArg := make([]string, 0)
	waitForArg := false
	restartArg := false
	debugArg := debug
	args := os.Args[1:]
	if len(args) > 0 {
		gl.Log("info", "Arguments provided:")
		for index, arg := range args {
			switch arg {
			case "--name", "-n":
				if index+1 < len(args) {
					commandNameArg = args[index+1]
				}
			case "--command", "-c":
				if index+1 < len(args) {
					commandArg = args[index+1]
				}
			case "--args", "-a":
				if index+1 < len(args) {
					argArg := args[index+1]
					if strings.Contains(argArg, ",") {
						aargsArg := strings.Split(argArg, ",")
						for indexArg := 0; indexArg < len(aargsArg); indexArg++ {
							aArg := aargsArg[indexArg]
							if strings.Contains(aArg, "=") {
								argsArg[indexArg] = strings.Split(aArg, "=")[1]
							} else {
								argsArg = append(argsArg, aArg)
							}
						}
					} else {
						argsArg = append(argsArg, argArg)
					}
				}
			case "--wait", "-w":
				waitForArg = true
			case "--restart", "-r":
				restartArg = true
			case "--debug", "-d":
				debugArg = true
			}
		}
	} else {
		commandArg = "tail"
		argsArg = []string{"-f", "/dev/null"}
		waitForArg = true
		restartArg = false
		debugArg = debug
	}

	var lcm i.ILifeCycle[pi.ProcessInput[any]]

	fn := p.ValidationFunc[pi.ProcessInput[any]]{
		Priority: 0,
		Func: func(obj *pi.ProcessInput[any], args ...any) *p.ValidationResult {
			objT := reflect.ValueOf(obj).Interface()
			if len(args) > 0 {
				for _, arg := range args {
					if reflect.TypeOf(arg) == reflect.TypeFor[func(*any, ...any) *p.ValidationResult]() {
						return arg.(func(*any, ...any) *p.ValidationResult)(&objT, args...)
					}
				}
				return nil
			}
			return nil
		},
		Result: nil,
	}
	if addPfnErr := pFunc.AddValidator(fn); addPfnErr != nil {
		l.FatalC(fmt.Sprintf("Error adding validation function: %s", addPfnErr.Error()), nil)
		return nil
	}

	lcm = i.NewLifeCycle[pi.ProcessInput[any]](pi.NewSystemProcessInput[any](
		commandNameArg,
		commandArg,
		argsArg,
		waitForArg,
		restartArg,
		&fn,
		logger,
		debugArg,
	), strings.ToUpper(test) == "A")

	if _, ok := lcm.(i.ILifeCycle[pi.ProcessInput[any]]); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type ILifeCycle[ProcessInput[any]] for test %s", test), nil)
		return nil
	}
	if lcmT, ok := lcm.(T); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type %s for test %s", reflect.TypeOf(lcm).String(), test), nil)
		return nil
	} else {
		ggl := GoLife[T]{
			Logger:    logger,
			Mutexes:   p.NewMutexes(),
			Reference: p.NewReference("GoLife"),
			Object:    lcmT,
		}

		ggl.initializeTest()
		return &ggl
	}
}

// WithProperty sets the property for the GoLife instance.
func (g *GoLife[T]) WithProperty(name string, property i.ILifeCycle[pi.ProcessInput[any]], withMetrics bool, cb func(any) (bool, error)) *GoLife[T] {
	if g.properties == nil {
		g.properties = make(map[string]interface{})
	}
	if _, ok := g.properties[name]; ok {
		gl.LogObjLogger[GoLife[T]](g, "error", "Property already exists")
		return g
	}
	prop := WithProperty[i.ILifeCycle[pi.ProcessInput[any]]](name, &property, withMetrics, cb)
	if prop == nil {
		gl.LogObjLogger[GoLife[T]](g, "error", "Property is nil")
		return g
	}
	g.properties[name] = prop
	return g
}

// WithLayer sets the layer for the GoLife instance.
func (g *GoLife[T]) WithLayer(scope string, layerObject *p.Property[i.ILifeCycle[pi.ProcessInput[any]]]) *GoLife[T] {
	if g.layers == nil {
		g.layers = make(map[string]interface{})
	}
	if _, ok := g.layers[scope]; ok {
		gl.LogObjLogger[GoLife[T]](g, "error", "Layer already exists")
		return g
	}
	layer := WithLayer[p.Property[i.ILifeCycle[pi.ProcessInput[any]]]](scope, layerObject)
	if layer == nil {
		gl.LogObjLogger[GoLife[T]](g, "error", "Layer is nil")
		return g
	}
	g.layers[scope] = layer
	return g
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
	arrMap := []map[string]interface{}{g.layers, g.properties, g.metadata}
	for key, m := range arrMap {
		if m == nil {
			m = make(map[string]interface{})
			arrMap[key] = m
		}
	}

	if g.properties == nil {
		g.properties = make(map[string]interface{})
	}
	g.properties["lifeCycle"] = WithProperty[T]("lifeCycle", &g.Object, true, func(any) (bool, error) {
		// Will create the callback function for the lifecycle manager
		// This is a placeholder for the actual callback logic
		gl.LogObjLogger[GoLife[T]](g, "debug", "Lifecycle manager callback executed")
		return true, nil
	})
	if g.layers == nil {
		g.layers = make(map[string]interface{})
	}
	g.layers["lifeCycle"] = p.NewLayer[p.Property[T]]("lifeCycle", g.properties["lifeCycle"].(*p.Property[T]))
}

// initialize initializes the GoLife properties, mutexes, metadata and other fields correctly.
func (g *GoLife[T]) initializeTest() {
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
	arrMap := []map[string]interface{}{g.layers, g.properties, g.metadata}
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
	g.properties["lifeCycle"] = WithProperty[T]("lifeCycle", &g.Object, true, func(any) (bool, error) {
		// Will create the callback function for the lifecycle manager
		// This is a placeholder for the actual callback logic
		gl.LogObjLogger[GoLife[T]](g, "debug", "Lifecycle manager callback executed")
		return true, nil
	})
	if g.layers == nil {
		g.layers = make(map[string]interface{})
	}
	g.layers["lifeCycle"] = p.NewLayer[p.Property[T]]("lifeCycle", g.properties["lifeCycle"].(*p.Property[T]))

}

// WithProperty creates a new property with the given name, property, and callback function.
func WithProperty[T i.ILifeCycle[pi.ProcessInput[any]]](name string, property *T, withMetrics bool, cb func(any) (bool, error)) *p.Property[T] {
	if property == nil {
		property = interface{}(nil).(*T)
	}
	return p.NewProperty[T](name, property, withMetrics, cb)
}

// WithLayer creates a new layer with the given scope and layer object.
func WithLayer[T p.Property[i.ILifeCycle[pi.ProcessInput[any]]]](scope string, layerObject *T) *p.Layer[T] {
	if layerObject == nil {
		return nil
	}
	return p.NewLayer[T](scope, layerObject)
}

// TestInitialization tests the initialization of the GoLife instance.
func TestInitialization[T i.ILifeCycle[pi.ProcessInput[any]]](g *GoLife[i.ILifeCycle[pi.ProcessInput[any]]]) {

	gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "info", "Testing GoLife initialization...")

	if propInterface, exists := g.properties["lifeCycle"]; !exists {
		gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "error", "Lifecycle property does not exist!")
		return
	} else {
		gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "success", "Lifecycle property instance created successfully!")

		if prop := propInterface.(*p.Property[T]); prop == nil {
			gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "error", "Lifecycle property is nil!")
			return
		} else {
			lcProp := prop.Prop.Get(false)

			if pp := reflect.ValueOf(lcProp).Interface().(*i.ILifeCycle[pi.ProcessInput[any]]); !reflect.ValueOf(pp).IsValid() {
				gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "warn", "Lifecycle expect property type:", reflect.TypeFor[*i.ILifeCycle[pi.ProcessInput[any]]]().String())
				gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "warn", "Lifecycle have property type:  ", reflect.TypeOf(lcProp).String())
				gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "error", "Lifecycle property does not implement ILifeCycle")
				return
			} else {
				ppp := *pp

				if err := ppp.StartLifecycle(); err != nil {
					gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "error", "Error starting lifecycle:", err.Error())
					return
				} else {
					gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "success", "Lifecycle started successfully!")
					time.Sleep(5 * time.Second)
					//if err := ppp.StopLifecycle(); err != nil {
					//	gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "error", "Error stopping lifecycle:", err.Error())
					//	return
					//} else {
					//	gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "success", "Lifecycle stopped successfully!")
					//}
					gl.LogObjLogger[GoLife[i.ILifeCycle[pi.ProcessInput[any]]]](g, "info", "Quitting test...")
				}
			}
		}
	}
}

// NewGoLife creates a new GoLife instance with the provided Logger.
func NewGoLife[T i.ILifeCycle[pi.ProcessInput[any]]](logger l.Logger, debug bool) *GoLife[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	gl.SetDebug(debug)
	pFunc := p.NewValidation[pi.ProcessInput[any]]()
	fn := p.ValidationFunc[pi.ProcessInput[any]]{
		Priority: 0,
		Func: func(obj *pi.ProcessInput[any], args ...any) *p.ValidationResult {
			objT := reflect.ValueOf(obj).Interface()
			if len(args) > 0 {
				for _, arg := range args {
					if reflect.TypeOf(arg) == reflect.TypeFor[func(*any, ...any) *p.ValidationResult]() {
						return arg.(func(*any, ...any) *p.ValidationResult)(&objT, args...)
					}
				}
				return nil
			}
			return nil
		},
		Result: nil,
	}
	if addPfnErr := pFunc.AddValidator(fn); addPfnErr != nil {
		l.FatalC(fmt.Sprintf("Error adding validation function: %s", addPfnErr.Error()), nil)
		return nil
	}
	pInput := pi.NewSystemProcessInput[any](
		"defaultProcess",
		"htop",
		[]string{"--tree"},
		true,
		false,
		&fn,
		logger,
		false,
	)
	// Example of creating a new LifeCycle for a SystemProcessInput
	lcm := i.NewLifeCycle[pi.ProcessInput[any]](pInput, false)

	if _, ok := lcm.(i.ILifeCycle[pi.ProcessInput[any]]); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type ILifeCycle[ProcessInput[any]]"), nil)
		return nil
	}
	if lcmT, ok := lcm.(T); !ok {
		l.FatalC(fmt.Sprintf("Lifecycle manager is not of type %s", reflect.TypeOf(lcm).String()), nil)
		return nil
	} else {
		ggl := GoLife[T]{
			Logger:    logger,
			Mutexes:   p.NewMutexes(),
			Reference: p.NewReference("GoLife"),
			Object:    lcmT,
		}

		ggl.initializeTest()
		return &ggl
	}
}
