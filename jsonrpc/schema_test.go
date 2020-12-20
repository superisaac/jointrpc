package jsonrpc

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

	s1 = []byte(`{"type": "bad"}`)
	s, err = builder.BuildBytes(s1)
	assert.NotNil(err)
	assert.Equal("SchemaBuildError unknown type", err.Error())

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
"items": {"type": "number"}
}`)
	builder = NewSchemaBuilder()
	s, err := builder.BuildBytes(s1)
	assert.Nil(err)
	assert.Equal("list", s.Type())

	s1 = []byte(`{
"type": "list",
"items": [
{"type": "number"},
{"type": "string"},
{"type": "bool"}
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
	assert.Equal("SchemaBuildError no properties", err.Error())

	s1 = []byte(`{
"type": "object",
"properties": {
  "aaa": {"type": "string"},
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
]}`)
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
}
