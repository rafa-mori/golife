package proto

import (
	"fmt"
)

// Protobuf is a struct that represents a protobuf message.
type Protobuf struct {
	// Name is the name of the protobuf message.
	Name string `json:"name"`
	// Fields is a slice of Field that represents the fields of the protobuf message.
	Fields []Field `json:"fields"`
	// Imports is a slice of string that represents the imports of the protobuf message.
	Imports []string `json:"imports"`
	// Options is a slice of string that represents the options of the protobuf message.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the protobuf message.
	Comments []string `json:"comments"`
	// Package is the package name of the protobuf message.
	Package string `json:"package"`
	// Syntax is the syntax of the protobuf message.
	Syntax string `json:"syntax"`
	// SourceCodeInfo is the source code information of the protobuf message.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// OneofDecl is a slice of OneofDecl that represents the oneof declarations of the protobuf message.
	OneofDecl []OneofDecl `json:"oneof_decl"`
	// EnumType is a slice of EnumType that represents the enum types of the protobuf message.
	EnumType []EnumType `json:"enum_type"`
	// MessageType is a slice of MessageType that represents the message types of the protobuf message.
	MessageType []MessageType `json:"message_type"`
	// Extend is a slice of Extend that represents the extend of the protobuf message.
	Extend []Extend `json:"extend"`
	// ReservedRange is a slice of ReservedRange that represents the reserved range of the protobuf message.
	ReservedRange []ReservedRange `json:"reserved_range"`
	// ReservedName is a slice of string that represents the reserved names of the protobuf message.
	ReservedName []string `json:"reserved_name"`
	// OptionsMap is a map of string to string that represents the options of the protobuf message.
	OptionsMap map[string]string `json:"options_map"`
}

// Field is a struct that represents a field in a protobuf message.
type Field struct {
	// Name is the name of the field.
	Name string `json:"name"`
	// Number is the number of the field.
	Number int32 `json:"number"`
	// Type is the type of the field.
	Type string `json:"type"`
	// Label is the label of the field.
	Label string `json:"label"`
	// DefaultValue is the default value of the field.
	DefaultValue string `json:"default_value"`
	// JsonName is the JSON name of the field.
	JsonName string `json:"json_name"`
	// Options is a slice of string that represents the options of the field.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the field.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the field.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// OneofIndex is the index of the oneof field.
	OneofIndex int32 `json:"oneof_index"`
	// Proto3Optional is a boolean that indicates if the field is optional in proto3.
	Proto3Optional bool `json:"proto3_optional"`
	// Proto3Map is a boolean that indicates if the field is a map in proto3.
	Proto3Map bool `json:"proto3_map"`
	// Proto3Packed is a boolean that indicates if the field is packed in proto3.
	Proto3Packed bool `json:"proto3_packed"`
	// Proto3Singular is a boolean that indicates if the field is singular in proto3.
	Proto3Singular bool `json:"proto3_singular"`
	// Proto3Repeated is a boolean that indicates if the field is repeated in proto3.
	Proto3Repeated bool `json:"proto3_repeated"`
	// Proto3Required is a boolean that indicates if the field is required in proto3.
	Proto3Required bool `json:"proto3_required"`
	// Proto3Weak is a boolean that indicates if the field is weak in proto3.
	Proto3Weak bool `json:"proto3_weak"`
	// Proto3Deprecated is a boolean that indicates if the field is deprecated in proto3.
	Proto3Deprecated bool `json:"proto3_deprecated"`
	// Proto3JsonName is the JSON name of the field in proto3.
	Proto3JsonName string `json:"proto3_json_name"`
}

// SourceCodeInfo is a struct that represents the source code information of a protobuf message.
type SourceCodeInfo struct {
	// Location is a slice of Location that represents the location of the source code information.
	Location []Location `json:"location"`
	// Path is the path of the source code information.
	Path []int32 `json:"path"`
	// Span is a slice of int32 that represents the span of the source code information.
	Span []int32 `json:"span"`
	// LeadingComments is the leading comments of the source code information.
	LeadingComments string `json:"leading_comments"`
	// TrailingComments is the trailing comments of the source code information.
	TrailingComments string `json:"trailing_comments"`
}

// Location is a struct that represents the location of the source code information.
type Location struct {
	// Path is a slice of int32 that represents the path of the location.
	Path []int32 `json:"path"`
	// Span is a slice of int32 that represents the span of the location.
	Span []int32 `json:"span"`
	// LeadingComments is the leading comments of the location.
	LeadingComments string `json:"leading_comments"`
	// TrailingComments is the trailing comments of the location.
	TrailingComments string `json:"trailing_comments"`
}

// OneofDecl is a struct that represents the oneof declaration of a protobuf message.
type OneofDecl struct {
	// Name is the name of the oneof declaration.
	Name string `json:"name"`
	// Options is a slice of string that represents the options of the oneof declaration.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the oneof declaration.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the oneof declaration.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
}

// EnumType is a struct that represents the enum type of a protobuf message.
type EnumType struct {
	// Name is the name of the enum type.
	Name string `json:"name"`
	// EnumValue is a slice of EnumValue that represents the enum values of the enum type.
	EnumValue []EnumValue `json:"enum_value"`
	// Options is a slice of string that represents the options of the enum type.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the enum type.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the enum type.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// ReservedRange is a slice of ReservedRange that represents the reserved range of the enum type.
	ReservedRange []ReservedRange `json:"reserved_range"`
	// ReservedName is a slice of string that represents the reserved names of the enum type.
	ReservedName []string `json:"reserved_name"`
	// OptionsMap is a map of string to string that represents the options of the enum type.
	OptionsMap map[string]string `json:"options_map"`
}

// EnumValue is a struct that represents the enum value of an enum type.
type EnumValue struct {
	// Name is the name of the enum value.
	Name string `json:"name"`
	// Number is the number of the enum value.
	Number int32 `json:"number"`
	// Options is a slice of string that represents the options of the enum value.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the enum value.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the enum value.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// OptionsMap is a map of string to string that represents the options of the enum value.
	OptionsMap map[string]string `json:"options_map"`
}

// MessageType is a struct that represents the message type of a protobuf message.
type MessageType struct {
	// Name is the name of the message type.
	Name string `json:"name"`
	// Fields is a slice of Field that represents the fields of the message type.
	Fields []Field `json:"fields"`
	// OneofDecl is a slice of OneofDecl that represents the oneof declarations of the message type.
	OneofDecl []OneofDecl `json:"oneof_decl"`
	// EnumType is a slice of EnumType that represents the enum types of the message type.
	EnumType []EnumType `json:"enum_type"`
	// MessageType is a slice of MessageType that represents the message types of the message type.
	MessageType []MessageType `json:"message_type"`
	// Extend is a slice of Extend that represents the extend of the message type.
	Extend []Extend `json:"extend"`
	// Options is a slice of string that represents the options of the message type.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the message type.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the message type.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// ReservedRange is a slice of ReservedRange that represents the reserved range of the message type.
	ReservedRange []ReservedRange `json:"reserved_range"`
	// ReservedName is a slice of string that represents the reserved names of the message type.
	ReservedName []string `json:"reserved_name"`
	// OptionsMap is a map of string to string that represents the options of the message type.
	OptionsMap map[string]string `json:"options_map"`
}

// Extend is a struct that represents the extend of a protobuf message.
type Extend struct {
	// Name is the name of the extend.
	Name string `json:"name"`
	// Fields is a slice of Field that represents the fields of the extend.
	Fields []Field `json:"fields"`
	// OneofDecl is a slice of OneofDecl that represents the oneof declarations of the extend.
	OneofDecl []OneofDecl `json:"oneof_decl"`
	// EnumType is a slice of EnumType that represents the enum types of the extend.
	EnumType []EnumType `json:"enum_type"`
	// MessageType is a slice of MessageType that represents the message types of the extend.
	MessageType []MessageType `json:"message_type"`
	// Options is a slice of string that represents the options of the extend.
	Options []string `json:"options"`
	// Comments is a slice of string that represents the comments of the extend.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the extend.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// ReservedRange is a slice of ReservedRange that represents the reserved range of the extend.
	ReservedRange []ReservedRange `json:"reserved_range"`
	// ReservedName is a slice of string that represents the reserved names of the extend.
	ReservedName []string `json:"reserved_name"`
	// OptionsMap is a map of string to string that represents the options of the extend.
	OptionsMap map[string]string `json:"options_map"`
}

// ReservedRange is a struct that represents the reserved range of a protobuf message.
type ReservedRange struct {
	// Start is the start of the reserved range.
	Start int32 `json:"start"`
	// End is the end of the reserved range.
	End int32 `json:"end"`
	// Comments is a slice of string that represents the comments of the reserved range.
	Comments []string `json:"comments"`
	// SourceCodeInfo is the source code information of the reserved range.
	SourceCodeInfo *SourceCodeInfo `json:"source_code_info"`
	// Options is a slice of string that represents the options of the reserved range.
	Options []string `json:"options"`
	// OptionsMap is a map of string to string that represents the options of the reserved range.
	OptionsMap map[string]string `json:"options_map"`
}

// String returns the string representation of the Protobuf struct.
func (p *Protobuf) String() string {
	return fmt.Sprintf("Protobuf{Name: %s, Fields: %v, Imports: %v, Options: %v, Comments: %v, Package: %s, Syntax: %s, SourceCodeInfo: %v, OneofDecl: %v, EnumType: %v, MessageType: %v, Extend: %v, ReservedRange: %v, ReservedName: %v, OptionsMap: %v}",
		p.Name,
		p.Fields,
		p.Imports,
		p.Options,
		p.Comments,
		p.Package,
		p.Syntax,
		p.SourceCodeInfo,
		p.OneofDecl,
		p.EnumType,
		p.MessageType,
		p.Extend,
		p.ReservedRange,
		p.ReservedName,
		p.OptionsMap)
}

// String returns the string representation of the Field struct.
func (f *Field) String() string {
	return fmt.Sprintf("Field{Name: %s, Number: %d, Type: %s, Label: %s, DefaultValue: %s, JsonName: %s, Options: %v, Comments: %v, SourceCodeInfo: %v, OneofIndex: %d, Proto3Optional: %t, Proto3Map: %t, Proto3Packed: %t, Proto3Singular: %t, Proto3Repeated: %t, Proto3Required: %t, Proto3Weak: %t, Proto3Deprecated: %t, Proto3JsonName: %s}",
		f.Name,
		f.Number,
		f.Type,
		f.Label,
		f.DefaultValue,
		f.JsonName,
		f.Options,
		f.Comments,
		f.SourceCodeInfo,
		f.OneofIndex,
		f.Proto3Optional,
		f.Proto3Map,
		f.Proto3Packed,
		f.Proto3Singular,
		f.Proto3Repeated,
		f.Proto3Required,
		f.Proto3Weak,
		f.Proto3Deprecated,
		f.Proto3JsonName)
}

// String returns the string representation of the SourceCodeInfo struct.
func (s *SourceCodeInfo) String() string {
	return fmt.Sprintf("SourceCodeInfo{Location: %v, Path: %v, Span: %v, LeadingComments: %s, TrailingComments: %s}",
		s.Location,
		s.Path,
		s.Span,
		s.LeadingComments,
		s.TrailingComments)
}

// String returns the string representation of the Location struct.
func (l *Location) String() string {
	return fmt.Sprintf("Location{Path: %v, Span: %v, LeadingComments: %s, TrailingComments: %s}",
		l.Path,
		l.Span,
		l.LeadingComments,
		l.TrailingComments)
}

// String returns the string representation of the OneofDecl struct.
func (o *OneofDecl) String() string {
	return fmt.Sprintf("OneofDecl{Name: %s, Options: %v, Comments: %v, SourceCodeInfo: %v}",
		o.Name,
		o.Options,
		o.Comments,
		o.SourceCodeInfo)
}

// String returns the string representation of the EnumType struct.
func (e *EnumType) String() string {
	return fmt.Sprintf("EnumType{Name: %s, EnumValue: %v, Options: %v, Comments: %v, SourceCodeInfo: %v, ReservedRange: %v, ReservedName: %v, OptionsMap: %v}",
		e.Name,
		e.EnumValue,
		e.Options,
		e.Comments,
		e.SourceCodeInfo,
		e.ReservedRange,
		e.ReservedName,
		e.OptionsMap)
}

// String returns the string representation of the EnumValue struct.
func (e *EnumValue) String() string {
	return fmt.Sprintf("EnumValue{Name: %s, Number: %d, Options: %v, Comments: %v, SourceCodeInfo: %v, OptionsMap: %v}",
		e.Name,
		e.Number,
		e.Options,
		e.Comments,
		e.SourceCodeInfo,
		e.OptionsMap)
}

// String returns the string representation of the MessageType struct.
func (m *MessageType) String() string {
	return fmt.Sprintf("MessageType{Name: %s, Fields: %v, OneofDecl: %v, EnumType: %v, MessageType: %v, Extend: %v, Options: %v, Comments: %v, SourceCodeInfo: %v, ReservedRange: %v, ReservedName: %v, OptionsMap: %v}",
		m.Name,
		m.Fields,
		m.OneofDecl,
		m.EnumType,
		m.MessageType,
		m.Extend,
		m.Options,
		m.Comments,
		m.SourceCodeInfo,
		m.ReservedRange,
		m.ReservedName,
		m.OptionsMap)
}

// String returns the string representation of the Extend struct.
func (e *Extend) String() string {
	return fmt.Sprintf("Extend{Name: %s, Fields: %v, OneofDecl: %v, EnumType: %v, MessageType: %v, Options: %v, Comments: %v, SourceCodeInfo: %v, ReservedRange: %v, ReservedName: %v, OptionsMap: %v}",
		e.Name,
		e.Fields,
		e.OneofDecl,
		e.EnumType,
		e.MessageType,
		e.Options,
		e.Comments,
		e.SourceCodeInfo,
		e.ReservedRange,
		e.ReservedName,
		e.OptionsMap)
}

// String returns the string representation of the ReservedRange struct.
func (r *ReservedRange) String() string {
	return fmt.Sprintf("ReservedRange{Start: %d, End: %d, Comments: %v, SourceCodeInfo: %v, Options: %v, OptionsMap: %v}",
		r.Start,
		r.End,
		r.Comments,
		r.SourceCodeInfo,
		r.Options,
		r.OptionsMap)
}
