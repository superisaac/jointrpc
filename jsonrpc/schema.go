package jsonrpc

import (
	"fmt"
	"strings"
	simplejson "github.com/bitly/go-simplejson"	
)

type SchemaValidator struct {
	paths []string
}
type Schema interface {
	// returns the generated
	Type() String
	IsValid(validator *SchemaValidator, data interface{}) bool
}

func (self *SchemaValidator) IsValid(schema Schema, path string, data interface{}) bool {
	if path != "" {
		self.paths = append(self.paths, path)
	}
	isv := schema.IsValid(data)
	if path != "" {
		self.paths = self.paths[:-1]
	}
	return isv
}

// type = "any"
type AnySchema struct {
}
func (self AnySchema) Type() string {
	return "any"
}
func (self *AnySchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	return true
}

// type = "none"
type NoneSchema struct {
}
func (self NoneSchema) Type() string {
	return "none"
}
func (self *NoneSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	if data != nil {
		return false
	}
	return true
}

// type= "bool"
type BoolSchema struct {
}
func (self BoolSchema) Type() string {
	return "bool"
}
func (self *BoolSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	if _, ok := data.(bool); ok {
		return true
	}
	return false
}

// type = "number"
type NumberSchema struct {
}
func (self NumberSchema) Type() string {
	return "number"
}
func (self *NumberSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	if _, ok := data.(json.Number); ok {
		return true
	}
	if _, ok := data.(int); ok {
		return true
	}
	if _, ok := data.(float); ok {
		return true
	}
	return false
}

// type = "string"
type StringSchema struct {
}
func (self StringSchema) Type() string {
	return "string"
}
func (self *StringSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	if _, ok := data.(string); ok {
		return true
	}
	return false
}

// type = "anyOf"
type UnionSchema struct {
	Choices []Schema
}
type NewUnionSchema() *UnionSchema {
	return &UnionScheme{Choices: make([]Schema, 0)}
}

func (self UnionSchema) Type() string {
	var arr []string
	for _, schema := self.Choices {
		arr = append(arr, schema.Type())
	}
	return fmt.Sprintf("any of %s", strings.Join(arr, ", "))
}
func (self *UnionSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	for _, schema := range self.Choices {
		if validator.IsValid(schema, schema.Type(), data) {
			return true
		}
	}
	return false
}

// type = "array", items is object
type ListSchema struct {
	Sub Schema
}

func NewListSchema() *ListSchema {
	return &ListSchema{}
}

func (self ListSchema) Type() string {
	return fmt.Sprintf("list of %s", self.Sub.Type())
}
func (self *ListSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	items, ok := data.([]interface{})
	if !ok {
		return false
	}
	for i, item := range items {
		if !validator.IsValid(self.Sub, fmt.Sprintf("[%d]", i), item) {
			return false
		}
	}
	return true
}

// type = "array", items is list
type TupleSchema struct {
	Children []Schema
}
type NewTupleSchema() *TupleSchema {
	return &TupleSchema{Children: make([]Schema, 0)}
}
func (self TupleSchema) Type() string {
	var arr []string
	for _, schema := range self.Children {
		arr = append(arr, schema.Type())
	}
	return fmt.Sprintf("tuple of %s", strings.Join(arr, ", "))
}

func (self *TupleSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	items, ok := data.([]interface{})
	if !ok {
		return false
	}
	if len(items) != len(self.Children) {
		return false
	}

	for i, item := range items {
		schema := self.Children[i]
		if !validator.IsValid(schema, fmt.Sprintf("[%d]", i), item) {
			return false
		}
	}
}

// type = "object"
type ObjectSchema struct {
	Properties map[string]Schema
	Required map[string]bool
}

func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		Properties: make(map[string]Schema),
		Required: make(map[string]bool),
	}		
}

func (self ObjectSchema) Type() string {
	return "object"
}

func (self *ObjectSchema) IsValid(validator *SchemaValidator, data interface{}) bool {
	obj, ok := data.(map[string]Schema)
	if !ok {
		return false
	}
	for prop, schema := range self.Properties {
		if v, found := data[prop]; found {
			if !validator.IsValid(schema, "." + prop, v) {
				return false
			}
			
		} else {
			if _, required := self.Required[prop]; required {
				// prop is required
				return false
			}
		}
	}
	return true
}
