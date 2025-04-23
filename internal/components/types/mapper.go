package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
)

// IMapper is a generic interface for serializing and deserializing objects of type T.
type IMapper[T any] interface {
	// Serialize converts an object of type T to a byte array in the specified format.
	Serialize(data []byte, object *T, format string) ([]byte, error)
	// Deserialize converts a byte array in the specified format to an object of type T.
	Deserialize(data []byte, object *T, format string) error
}

// Mapper is a generic struct that implements the IMapper interface for serializing and deserializing objects.
type Mapper[T any] struct{}

// NewMapper creates a new instance of Mapper.
func NewMapper[T any]() IMapper[T] { return &Mapper[T]{} }

// Serialize converts an object of type T to a byte array in the specified format.
func (m *Mapper[T]) Serialize(data []byte, object *T, format string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("os dados estão vazios")
	}
	if object == nil {
		return nil, fmt.Errorf("o ponteiro do tipo está nil")
	}
	switch format {
	case "json":
		return json.Marshal(object)
	case "yaml":
		return yaml.Marshal(object)
	case "xml":
		return xml.Marshal(object)
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// Deserialize converts a byte array in the specified format to an object of type T.
func (m *Mapper[T]) Deserialize(data []byte, object *T, format string) error {
	if len(data) == 0 {
		return fmt.Errorf("os dados estão vazios")
	}
	if object == nil {
		return fmt.Errorf("o ponteiro do tipo está nil")
	}
	switch format {
	case "json":
		return json.Unmarshal(data, object)
	case "yaml":
		return yaml.Unmarshal(data, object)
	case "xml":
		return xml.Unmarshal(data, object)
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}
}
