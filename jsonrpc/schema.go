package jsonrpc

import (
	"fmt"
	//"errors"
	//"reflect"
	json "encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"strings"
)

func (self SchemaBuildError) Error() string {
	return fmt.Sprintf("SchemaBuildError %s", self.info)
}

func NewBuildError(info string) *SchemaBuildError {
	return &SchemaBuildError{info: info}
}

// Schema Validator

func (self ErrorPos) Path() string {
	return strings.Join(self.paths, "")
}

func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

func (self *SchemaValidator) NewErrorPos(hint string) *ErrorPos {
	var newPaths []string
	for _, path := range self.paths {
		newPaths = append(newPaths, path)
	}
	return &ErrorPos{paths: newPaths, hint: hint}
}

func (self *SchemaValidator) ValidateBytes(schema Schema, bytes []byte) *ErrorPos {
	data, err := simplejson.NewJson(bytes)
	if err != nil {
		panic(err)
	}
	return self.Scan(schema, "", data.Interface())
}

func (self *SchemaValidator) Validate(schema Schema, data interface{}) *ErrorPos {
	return self.Scan(schema, "", data)
}

func (self *SchemaValidator) pushPath(path string) {
	if path != "" {
		self.paths = append(self.paths, path)
	}
}

func (self *SchemaValidator) popPath(path string) {
	if path != "" {
		self.paths = self.paths[:len(self.paths)-1]
	}
}

func (self *SchemaValidator) Scan(schema Schema, path string, data interface{}) *ErrorPos {
	self.pushPath(path)
	errPos := schema.Scan(self, data)
	self.popPath(path)
	return errPos
}

// type = "any"
type AnySchema struct {
}

func (self AnySchema) Type() string {
	return "any"
}
func (self *AnySchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	return nil
}

// type = "null"
type NullSchema struct {
}

func (self NullSchema) Type() string {
	return "null"
}
func (self *NullSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if data != nil {
		return validator.NewErrorPos("data is not null")
	}
	return nil
}

// type= "bool"
type BoolSchema struct {
}

func (self BoolSchema) Type() string {
	return "bool"
}
func (self *BoolSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if _, ok := data.(bool); ok {
		return nil
	}
	return validator.NewErrorPos("data is not bool")
}

// type = "number"
type NumberSchema struct {
}

func NewNumberSchema() *NumberSchema {
	return &NumberSchema{}
}
func (self NumberSchema) Type() string {
	return "number"
}
func (self *NumberSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if _, ok := data.(json.Number); ok {
		return nil
	}
	if _, ok := data.(int); ok {
		return nil
	}

	if _, ok := data.(float64); ok {
		return nil
	}
	return validator.NewErrorPos("data is not number")
}

// type = "string"
type StringSchema struct {
}

func (self StringSchema) Type() string {
	return "string"
}
func (self *StringSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if _, ok := data.(string); ok {
		return nil
	}
	return validator.NewErrorPos("data is not string")
}

// type = "anyOf"
type UnionSchema struct {
	Choices []Schema
}

func NewUnionSchema() *UnionSchema {
	return &UnionSchema{Choices: make([]Schema, 0)}
}

func (self UnionSchema) Type() string {
	return "anyOf"
}
func (self *UnionSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	for _, schema := range self.Choices {
		if errPos := validator.Scan(schema, "", data); errPos == nil {
			return nil
		}
	}
	return validator.NewErrorPos("data is not any of the types")
}

// type = "array", items is object
type ListSchema struct {
	Item Schema
}

func NewListSchema() *ListSchema {
	return &ListSchema{}
}

func (self ListSchema) Type() string {
	return "list"
}
func (self *ListSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	items, ok := data.([]interface{})
	if !ok {
		return validator.NewErrorPos("data is not a list")
	}
	for i, item := range items {
		if errPos := validator.Scan(self.Item, fmt.Sprintf("[%d]", i), item); errPos != nil {
			return errPos
		}
	}
	return nil
}

// type = "array", items is list
type TupleSchema struct {
	Children []Schema
}

func NewTupleSchema() *TupleSchema {
	return &TupleSchema{Children: make([]Schema, 0)}
}
func (self TupleSchema) Type() string {
	return "list"
}

func (self *TupleSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	items, ok := data.([]interface{})
	if !ok {
		return validator.NewErrorPos("data is not a list")
	}
	if len(items) != len(self.Children) {
		return validator.NewErrorPos("tuple length mismatch")
	}

	for i, item := range items {
		schema := self.Children[i]
		if errPos := validator.Scan(schema, fmt.Sprintf("[%d]", i), item); errPos != nil {
			return errPos
		}
	}
	return nil
}

// type = "object"
type ObjectSchema struct {
	Properties map[string]Schema
	Requires   map[string]bool
}

func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		Properties: make(map[string]Schema),
		Requires:   make(map[string]bool),
	}
}

func (self ObjectSchema) Type() string {
	return "object"
}

func (self *ObjectSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return validator.NewErrorPos("data is not an object")
	}
	for prop, schema := range self.Properties {
		if v, found := obj[prop]; found {
			if errPos := validator.Scan(schema, "."+prop, v); errPos != nil {
				return errPos
			}

		} else {
			if _, required := self.Requires[prop]; required {
				validator.pushPath("." + prop)
				errPos := validator.NewErrorPos("required prop is not present")
				validator.popPath("." + prop)
				return errPos
			}
		}
	}
	return nil
}

// Builder
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{}
}

func (self *SchemaBuilder) BuildBytes(bytes []byte) (Schema, error) {
	js, err := simplejson.NewJson(bytes)
	if err != nil {
		return nil, err
	}
	return self.Build(js.Interface())
}

func (self *SchemaBuilder) Build(data interface{}) (Schema, error) {
	if node, ok := data.(map[string](interface{})); ok {
		return self.buildNode(node)
	} else {
		//fmt.Printf("data type is %+v", reflect.TypeOf(data))
		return nil, NewBuildError("data is not a map")
	}
}

func (self *SchemaBuilder) buildNode(node map[string](interface{})) (Schema, error) {
	nodeType, ok := node["type"]
	if !ok {
		return nil, NewBuildError("no type presented")
	}
	switch nodeType {
	case "number":
		return NewNumberSchema(), nil
	case "any":
		return &AnySchema{}, nil
	case "null":
		return &NullSchema{}, nil
	case "string":
		return &StringSchema{}, nil
	case "bool":
		return &BoolSchema{}, nil
	case "list":
		return self.buildListSchema(node)
	case "object":
		return self.buildObjectSchema(node)
	default:
		return nil, NewBuildError("unknown type")
	}
}

func (self *SchemaBuilder) buildListSchema(node map[string](interface{})) (Schema, error) {
	items, ok := node["items"]
	if !ok {
		return nil, NewBuildError("no items")
	}
	if itemsMap, ok := items.(map[string](interface{})); ok {
		itemSchema, err := self.buildNode(itemsMap)
		if err != nil {
			return nil, err
		}
		schema := NewListSchema()
		schema.Item = itemSchema
		return schema, nil
	}
	if itemsTuple, ok := items.([]interface{}); ok {
		schema := NewTupleSchema()
		for _, item := range itemsTuple {
			itemNode, ok := item.(map[string]interface{})
			if !ok {
				return nil, NewBuildError("tuple item not a map")
			}
			childSchema, err := self.buildNode(itemNode)
			if err != nil {
				return nil, err
			}
			schema.Children = append(schema.Children, childSchema)
		}
		return schema, nil
	}
	return nil, NewBuildError("fail to build list schema")
}

func (self *SchemaBuilder) buildObjectSchema(node map[string](interface{})) (Schema, error) {
	props, ok := node["properties"]
	if !ok {
		return nil, NewBuildError("no properties")
	}
	propNodes, ok := props.(map[string]interface{})
	if !ok {
		return nil, NewBuildError("wrong properties type")
	}

	schema := NewObjectSchema()
	for propName, pv := range propNodes {
		propNode, ok := pv.(map[string]interface{})
		if !ok {
			return nil, NewBuildError("wrong prop item type")
		}
		child, err := self.buildNode(propNode)
		if err != nil {
			return nil, err
		}
		schema.Properties[propName] = child
	}

	if req, ok := node["requires"]; ok {
		requireList, ok := req.([]interface{})
		if !ok {
			return nil, NewBuildError("requires is not a list")
		}
		for _, r := range requireList {
			reqProp, ok := r.(string)
			if !ok {
				return nil, NewBuildError("requires items is not string")
			}
			if _, found := schema.Properties[reqProp]; !found {
				return nil, NewBuildError("cannot find required prop")
			}

			schema.Requires[reqProp] = true
		}
	}
	return schema, nil
}
