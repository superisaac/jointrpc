package schema

// Schema builder
type SchemaBuildError struct {
	info string
}

type SchemaBuilder struct {
}

// Schema validator
type SchemaValidator struct {
	paths     []string
	hint      string
	errorPath string
}

type ErrorPos struct {
	paths []string
	hint  string
}

type Schema interface {
	// returns the generated
	Type() string
	RebuildType() map[string]interface{}
	Scan(validator *SchemaValidator, data interface{}) *ErrorPos
	SetName(name string)
	GetName() string
	SetDescription(desc string)
	GetDescription() string
}

type SchemaMixin struct {
	name        string
	description string
}

// schema subclasses
type AnySchema struct {
	SchemaMixin
}

type NullSchema struct {
	SchemaMixin
}
type BoolSchema struct {
	SchemaMixin
}

type NumberSchema struct {
	SchemaMixin
}
type StringSchema struct {
	SchemaMixin
}

type UnionSchema struct {
	SchemaMixin
	Choices []Schema
}
type ListSchema struct {
	SchemaMixin
	Item Schema
}

type TupleSchema struct {
	SchemaMixin
	Children         []Schema
	AdditionalSchema Schema
}

type ObjectSchema struct {
	SchemaMixin
	Properties map[string]Schema
	Requires   map[string]bool
}

type MethodSchema struct {
	SchemaMixin
	Params []Schema
	Result Schema
}
