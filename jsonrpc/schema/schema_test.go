package schema

import (
	//json "encoding/json"
	//"fmt"
	"testing"
	//"reflect"
	"github.com/stretchr/testify/assert"
)

func TestBuildBasicSchema(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{"type": "number"}`)
	builder := NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("number", s.Type())

	s1 = []byte(`"string"`)
	builder = NewSchemaBuilder()
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("string", s.Type())

	s1 = []byte(`{"type": "bad"}`)
	s, err = builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError unknown type", err.Error())

	s1 = []byte(`"bad2"`)
	s, err = builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError data is not an object", err.Error())

	s1 = []byte(`{"aa": 4}`)
	s, err = builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError no type presented", err.Error())

	s1 = []byte(`{"type": "string"}`)
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("string", s.Type())

	s1 = []byte(`{"type": "null"}`)
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("null", s.Type())

	s1 = []byte(`{"type": "any"}`)
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("any", s.Type())

}

func TestBuildListSchema(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{"type": "list"}`)
	builder := NewSchemaBuilder()
	_, err := builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError no items", err.Error())

	s1 = []byte(`{
"type": "list",
"items": "number"
}`)
	builder = NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("list", s.Type())

	s1 = []byte(`{
"type": "list",
"items": [
"number",
{"type": "string"},
"bool"
]}`)
	builder = NewSchemaBuilder()
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("list", s.Type())
	tupleSchema, ok := s.(*TupleSchema)
	assert.True(ok)
	assert.Equal(3, len(tupleSchema.Children))
	assert.Equal("number", tupleSchema.Children[0].Type())
	assert.Equal("string", tupleSchema.Children[1].Type())
	assert.Equal("bool", tupleSchema.Children[2].Type())
}

func TestBuildObjectSchema(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{"type": "object"}`)
	builder := NewSchemaBuilder()
	_, err := builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError properties is not a map of objects", err.Error())

	s1 = []byte(`{
"type": "object",
"properties": {
  "aaa": "string",
  "bbb": {"type": "number"}
}
}`)
	builder = NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("object", s.Type())

	obj, ok := s.(*ObjectSchema)
	assert.True(ok)
	assert.Equal("string", obj.Properties["aaa"].Type())
	assert.Equal("number", obj.Properties["bbb"].Type())

	s1 = []byte(`{
"type": "object",
"properties": {
  "aaa": {"type": "string"},
  "bbb": {"type": "number"}
},
"requires": ["nosuch", "aaa"]
}`)
	builder = NewSchemaBuilder()
	s, err = builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError cannot find required prop", err.Error())

	s1 = []byte(`{
"type": "object",
"properties": {
  "aaa": {"type": "string"},
  "bbb": {"type": "number"},
  "ccc": {"type": "list", "items": {"type": "string"}}
},
"requires": ["aaa", "bbb"]
}`)
	builder = NewSchemaBuilder()
	s, err = builder.BuildBytes(s1)
	assert.Nil(err)

	obj, ok = s.(*ObjectSchema)
	assert.True(ok)

	assert.Equal(2, len(obj.Requires))
}

func TestBasicValidator(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{"type": "number"}`)
	builder := NewSchemaBuilder()
	schema, err := builder.BuildBytes(s1)
	assert.Nil(err)
	numberSchema, ok := schema.(*NumberSchema)
	assert.True(ok)

	validator := NewSchemaValidator()
	errPos := validator.ValidateBytes(numberSchema, []byte(`6`))
	assert.Nil(errPos)

	validator = NewSchemaValidator()
	errPos = validator.ValidateBytes(numberSchema, []byte(`"a string"`))
	assert.NotNil(errPos)
	assert.Equal("data is not number", errPos.hint)
	assert.Equal("", errPos.Path())
}

func TestUnionValidator(t *testing.T) {
	assert := assert.New(t)
	s1 := []byte(`{"type": "union"}`)
	builder := NewSchemaBuilder()
	_, err := builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError no valid anyOf attribute", err.Error())

	s1 = []byte(`{
"type": "union",
"anyOf": [
  {"type": "number"},
  {"type": "string"}
]
}`)
	builder = NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)

	uschema, ok := s.(*UnionSchema)
	assert.True(ok)

	validator := NewSchemaValidator()
	data := []byte(`true`)
	errPos := validator.ValidateBytes(uschema, data)
	assert.NotNil(errPos)
	assert.Equal("", errPos.Path())
	assert.Equal("data is not any of the types", errPos.hint)

	validator = NewSchemaValidator()
	data = []byte(`{}`)
	errPos = validator.ValidateBytes(uschema, data)
	assert.NotNil(errPos)
	assert.Equal("", errPos.Path())
	assert.Equal("data is not any of the types", errPos.hint)

	validator = NewSchemaValidator()
	data = []byte(`-3.88`)
	errPos = validator.ValidateBytes(uschema, data)
	assert.Nil(errPos)

	validator = NewSchemaValidator()
	data = []byte(`"a string"`)
	errPos = validator.ValidateBytes(uschema, data)
	assert.Nil(errPos)

}

func TestComplexValidator(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{
"type": "list",
"items": [
   {"type": "string"},
   {"type": "object",
    "properties": {
       "abc": {"type": "string"},
       "def": {"type": "number"}
    },
    "requires": ["abc"]
   }
],
"additionalItems": {"type": "string"}
}`)
	builder := NewSchemaBuilder()
	schema, err := builder.BuildBytes(s1)
	assert.Nil(err)
	s, ok := schema.(*TupleSchema)
	assert.True(ok)

	validator := NewSchemaValidator()
	data := []byte(`["hello", {"abc": "world"}]`)
	errPos := validator.ValidateBytes(s, data)
	assert.Nil(errPos)

	validator = NewSchemaValidator()
	data = []byte(`["hello", {"abc": "world"}, "hello1", "hello2"]`)
	errPos = validator.ValidateBytes(s, data)
	assert.Nil(errPos)

	validator = NewSchemaValidator()
	data = []byte(`["hello", {"abc": 8}]`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not string", errPos.hint)
	assert.Equal("[1].abc", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`["hello", {"def": 7}]`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("required prop is not present", errPos.hint)
	assert.Equal("[1].abc", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`["hello", {"abc": "world"}, "hello1", 123]`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not string", errPos.hint)
	assert.Equal("[3]", errPos.Path())

}

func TestMethodValidator(t *testing.T) {
	assert := assert.New(t)

	s1 := []byte(`{
"type": "method"
}`)
	builder := NewSchemaBuilder()
	_, err := builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError params is not a list of objects", err.Error())

	s1 = []byte(`{
"type": "method",
"params": [
  {"type": "number", "name": "a"},
  {"type": "string", "name": "b"},
  {
    "type": "object",
    "name": "options",
    "description": "calc options",
    "properties": {"aaa": {"type": "string"}, "bbb": {"type": "number"}},
    "requires": ["aaa"]
  }
],
"returns": {"type": "string"}
}`)
	builder = NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)

	assert.Equal("method", s.Type())
	methodSchema, ok := s.(*MethodSchema)
	assert.True(ok)
	assert.Equal("calc options", methodSchema.Params[2].GetDescription())

	validator := NewSchemaValidator()
	data := []byte(`["hello", 5, {"abc": 8}]`)
	errPos := validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not object", errPos.hint)
	assert.Equal("", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`{"params": ["hello", 5, {"abc": 8}]}`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not number", errPos.hint)
	assert.Equal(".params[0]", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`{"params": [5, "hello", {"abc": 8}]}`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("required prop is not present", errPos.hint)
	assert.Equal(".params[2].aaa", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`{"params": [5, "hello", {"aaa": 8}]}`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not string", errPos.hint)
	assert.Equal(".params[2].aaa", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`{"params": [5, "hello", {"aaa": "a string"}]}`)
	errPos = validator.ValidateBytes(s, data)
	assert.Nil(errPos)

	validator = NewSchemaValidator()
	data = []byte(`{"result": 8}`)
	errPos = validator.ValidateBytes(s, data)
	assert.NotNil(errPos)
	assert.Equal("data is not string", errPos.hint)
	assert.Equal(".result", errPos.Path())

	validator = NewSchemaValidator()
	data = []byte(`{"result": "a string"}`)
	errPos = validator.ValidateBytes(s, data)
	assert.Nil(errPos)

}
