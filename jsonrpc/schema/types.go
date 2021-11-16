package schema

// Schema builder
type SchemaBuildError struct {
	info  string
	paths []string
}

type SchemaBuilder struct {
}

// Fix string map issue from yaml format
type NonStringMap struct {
	paths []string
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

// schema sucblasses
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
	Minimum *float64
	Maximum *float64
}

type IntegerSchema struct {
	SchemaMixin
	Minimum *int64
	Maximum *int64
}

type StringSchema struct {
	SchemaMixin
	MaxLength int
}

type AnyOfSchema struct {
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
	Params  []Schema
	Returns Schema
}
