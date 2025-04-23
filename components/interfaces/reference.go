package interfaces

import "github.com/google/uuid"

type IReference interface {
	GetID() uuid.UUID
	GetName() string
	SetName(name string)
	String() string
}
