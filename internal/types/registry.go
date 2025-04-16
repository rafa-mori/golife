package types

import (
	"fmt"
	f "github.com/faelmori/golife/internal/property"
	"reflect"
	"sync/atomic"
)

type RegistryEntry[RU any] struct {
	// name é o nome do registro.
	name string
	// metadata armazena informações adicionais sobre o registro.
	metadata f.Metadata
	// valueType é o tipo do valor armazenado no registro.
	valueType reflect.Type
	// value é o valor atual do registro, armazenado em um ponteiro atômico para
	// garantir acesso seguro em ambientes concorrentes.
	value atomic.Pointer[RU]
	// validators é uma lista de funções de validação para o valor do registro.
	validators []func(RU) error
}

// NewRegistryEntry cria um novo registro com o nome e valor especificados.
// Se o valor for nil, o registro é inicializado com o valor zero do tipo E.
func NewRegistryEntry[E any](name string, value *E) *RegistryEntry[E] {
	entry := &RegistryEntry[E]{
		name:       name,
		valueType:  reflect.TypeFor[E](),
		value:      atomic.Pointer[E]{},
		metadata:   make(f.Metadata),
		validators: make([]func(E) error, 0),
	}
	if value == nil {
		entry.value.Store(new(E))
	} else {
		// Se o valor não for nil, armazena o valor no ponteiro atômico.
		entry.value.Store(value)
	}
	return entry
}

// Registry é um tipo que atua como o "map" dinâmico para armazenar os registros.
type Registry[R any, E any] struct {
	// store é um mapa que armazena os registros em ponteiros atômicos para garantir
	// acesso seguro em ambientes concorrentes e evitar cópias desnecessárias.
	store map[string]interface{}
}

// NewRegistry cria uma instância vazia de Registry.
func NewRegistry[R any, E any]() *Registry[R, E] {
	return &Registry[R, E]{
		store: make(map[string]interface{}),
	}
}

// Add adiciona um novo registro ao Registry.
func (r *Registry[R, E]) Add(name string, config E) error {
	// Verifica se o nome já existe no registro
	if _, exists := r.store[name]; exists {
		// Se o registro já existe, retorna um erro
		return fmt.Errorf("registro já existe: %s", name)
	} else {
		vCfg := reflect.ValueOf(config)
		if !vCfg.IsValid() {
			return fmt.Errorf("valor inválido para o registro: %s", name)
		}
		if vCfg.Kind() != reflect.Ptr {
			return fmt.Errorf("valor deve ser um ponteiro: %s", name)
		}
		if vCfg.Type().Elem().AssignableTo(reflect.TypeFor[E]()) {
			// Se o valor for um ponteiro, armazena o valor no ponteiro atômico.
			r.store[name] = NewRegistryEntry[E](name, &config)
		} else {
			return fmt.Errorf("tipo de valor incompatível para o registro: %s", name)
		}
	}

	return nil
}

// GetByName busca um registro pelo nome.
func (r *Registry[R, E]) GetByName(name string) (interface{}, error) {
	config, exists := r.store[name]
	if !exists {
		return nil, fmt.Errorf("registro não encontrado: %s", name)
	}
	return config, nil
}

// GetByType busca um registro pelo tipo.
func (r *Registry[R, E]) GetByType(targetType reflect.Type) ([]interface{}, error) {
	var matches []interface{}
	for _, config := range r.store {
		if entry, ok := config.(*RegistryEntry[E]); ok {
			if entry.valueType == targetType {
				matches = append(matches, entry)
			}
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("nenhum registro encontrado para o tipo: %s", targetType)
	}
	return matches, nil
}

// List retorna todos os registros armazenados.
func (r *Registry[R, E]) List() map[string]interface{} {
	return r.store
}
