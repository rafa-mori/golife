package golife

import (
	"fmt"
	"github.com/faelmori/golife/components/process_input"
	p "github.com/faelmori/golife/components/types"
	i "github.com/faelmori/golife/internal"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
)

// GoLife is a generic struct that implements the IGoLife inte-rface.
type GoLife[T any] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*p.Reference
	// Mutexes is the mutexes for this GoLife instance.
	*p.Mutexes
	// Properties is a map of properties for this GoLife instance.
	properties map[string]interface{}
	// metadata is a map of metadata for this GoLife instance.
	metadata map[string]interface{}
	// layers is a map of layers for this GoLife instance.
	layers map[string]interface{}
}

// NewGoLife creates a new GoLife instance with the provided Logger.
func NewGoLife[T any](logger l.Logger, debug bool) *GoLife[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	gl.SetDebug(debug)
	ggl := GoLife[T]{
		Logger:    logger,
		Mutexes:   p.NewMutexes(),
		Reference: p.NewReference("GoLife"),
	}
	ggl.initialize()
	return &ggl
}

// WithProperty sets the property for the GoLife instance.
func (g *GoLife[T]) WithProperty(name string, property *T, withMetrics bool, cb func(any) (bool, error)) *GoLife[T] {
	if g.properties == nil {
		g.properties = make(map[string]interface{})
	}
	if _, ok := g.properties[name]; ok {
		gl.LogObjLogger[GoLife[T]](g, "error", "Property already exists")
		return g
	}
	prop := WithProperty[T](name, property, withMetrics, cb)
	if prop == nil {
		gl.LogObjLogger[GoLife[T]](g, "error", "Property is nil")
		return g
	}
	g.properties[name] = prop
	return g
}

// WithLayer sets the layer for the GoLife instance.
func (g *GoLife[T]) WithLayer(scope string, layerObject any) *GoLife[T] {
	if g.layers == nil {
		g.layers = make(map[string]interface{})
	}
	if _, ok := g.layers[scope]; ok {
		gl.LogObjLogger[GoLife[T]](g, "error", "Layer already exists")
		return g
	}
	layer := WithLayer[any](scope, &layerObject)
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

	// Falta iniciar esse objeto, aí conseguimos começar a TESTAR a inicialização.. rsrs
	//procInput := p.NewProcessInput[T]()

	procInput := process_input.NewSystemProcessInput[T](
		"defaultProcess", // Nome do processo
		"ls",             // Comando padrão
		[]string{"-l"},   // Argumentos padrão
		false,            // Não esperar pelo término do comando
		false,            // Sem reinício automático
		nil,              // Sem função customizada
		g.Logger,         // Logger compartilhado do GoLife
	)
	if procInput == nil {
		gl.LogObjLogger[GoLife[T]](g, "error", "Failed to initialize ProcessInput")
		return
	}

	lcm := i.NewLifeCycle[T](procInput)

	if g.properties == nil {
		g.properties = make(map[string]interface{})
	}
	g.properties["lifeCycle"] = WithProperty[i.ILifeCycle[T]]("lifeCycle", &lcm, true, func(any) (bool, error) {
		// Will create the callback function for the lifecycle manager
		// This is a placeholder for the actual callback logic
		gl.LogObjLogger[GoLife[T]](g, "debug", "Lifecycle manager callback executed")
		return true, nil
	})
	if g.layers == nil {
		g.layers = make(map[string]interface{})
	}
	g.layers["lifeCycle"] = p.NewLayer[p.Property[i.ILifeCycle[T]]]("lifeCycle", g.properties["lifeCycle"].(*p.Property[i.ILifeCycle[T]]))
}

func (g *GoLife[T]) TestInitialization() {
	gl.LogObjLogger[GoLife[T]](g, "info", "Testing GoLife initialization...")
	if propInterface, exists := g.properties["lifeCycle"]; exists {
		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle property exists!")

		prop := propInterface.(*p.Property[i.ILifeCycle[T]])

		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle property object:", prop.Prop.ID.String())
		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle property object:", prop.Prop.Name)
		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle property object:", prop.Prop.Type().String())
		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle property object:", fmt.Sprintf("%v", prop.Prop.Get(false)))
	} else {
		gl.LogObjLogger[GoLife[T]](g, "error", "Lifecycle property does not exist!")
	}

	if lcmInterface, exists := g.layers["lifeCycle"]; exists {
		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle layer exists!")

		lcm := lcmInterface.(*p.Layer[p.Property[i.ILifeCycle[T]]])

		gl.LogObjLogger[GoLife[T]](g, "info", "Lifecycle layer object:", lcm.Snapshot())
	} else {
		gl.LogObjLogger[GoLife[T]](g, "error", "Lifecycle layer does not exist!")
	}
}

// WithProperty creates a new property with the given name, property, and callback function.
func WithProperty[P any](name string, property *P, withMetrics bool, cb func(any) (bool, error)) *p.Property[P] {
	if property == nil {
		property = interface{}(nil).(*P)
	}
	return p.NewProperty[P](name, property, withMetrics, cb)
}

// WithLayer creates a new layer with the given scope and layer object.
func WithLayer[L any](scope string, layerObject *L) *p.Layer[L] {
	if layerObject == nil {
		return nil
	}
	return p.NewLayer[L](scope, layerObject)
}
