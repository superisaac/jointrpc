package schema

import (
	"fmt"
	//"reflect"
	json "encoding/json"
	simplejson "github.com/bitly/go-simplejson"
)

// SchemaMixin
func (self *SchemaMixin) SetName(name string) {
	self.name = name
}
func (self SchemaMixin) GetName() string {
	return self.name
}

func (self *SchemaMixin) SetDescription(desc string) {
	self.description = desc
}

func (self SchemaMixin) GetDescription() string {
	return self.description
}

func (self SchemaMixin) rebuildType(nType string) map[string]interface{} {
	tp := map[string]interface{}{
		"type": nType,
	}
	if self.name != "" {
		tp["name"] = self.name
	}
	if self.description != "" {
		tp["description"] = self.description
	}
	return tp
}

// type = "any"

func (self AnySchema) Type() string {
	return "any"
}
func (self AnySchema) RebuildType() map[string]interface{} {
	return self.rebuildType(self.Type())
}

func (self *AnySchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	return nil
}

// type = "null"
func (self NullSchema) Type() string {
	return "null"
}
func (self NullSchema) RebuildType() map[string]interface{} {
	return self.rebuildType(self.Type())
}

func (self *NullSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if data != nil {
		return validator.NewErrorPos("data is not null")
	}
	return nil
}

// type= "bool"
func (self BoolSchema) Type() string {
	return "bool"
}
func (self BoolSchema) RebuildType() map[string]interface{} {
	return self.rebuildType(self.Type())
}

func (self *BoolSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if _, ok := data.(bool); ok {
		return nil
	}
	return validator.NewErrorPos("data is not bool")
}

// type = "number"

func NewNumberSchema() *NumberSchema {
	return &NumberSchema{}
}
func (self NumberSchema) Type() string {
	return "number"
}
func (self NumberSchema) RebuildType() map[string]interface{} {
	return self.rebuildType(self.Type())
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
func (self StringSchema) Type() string {
	return "string"
}
func (self StringSchema) RebuildType() map[string]interface{} {
	return self.rebuildType(self.Type())
}

func (self *StringSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	if _, ok := data.(string); ok {
		return nil
	}
	return validator.NewErrorPos("data is not string")
}

// type = "anyOf"
func NewUnionSchema() *UnionSchema {
	return &UnionSchema{Choices: make([]Schema, 0)}
}

func (self UnionSchema) RebuildType() map[string]interface{} {
	tp := self.rebuildType(self.Type())
	arr := make([](map[string]interface{}), 0)
	for _, choice := range self.Choices {
		arr = append(arr, choice.RebuildType())
	}
	tp["anyOf"] = arr
	return tp
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

func NewListSchema() *ListSchema {
	return &ListSchema{}
}

func (self ListSchema) Type() string {
	return "list"
}
func (self ListSchema) RebuildType() map[string]interface{} {
	tp := self.rebuildType(self.Type())
	tp["items"] = self.Item.RebuildType()
	return tp
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

func NewTupleSchema() *TupleSchema {
	return &TupleSchema{Children: make([]Schema, 0)}
}
func (self TupleSchema) Type() string {
	return "list"
}

func (self TupleSchema) RebuildType() map[string]interface{} {
	tp := self.rebuildType(self.Type())
	arr := make([](map[string]interface{}), 0)
	for _, child := range self.Children {
		arr = append(arr, child.RebuildType())
	}
	tp["items"] = arr
	return tp
}

func (self *TupleSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	items, ok := data.([]interface{})
	if !ok {
		return validator.NewErrorPos("data is not a list")
	}
	if self.AdditionalSchema == nil {
		if len(items) != len(self.Children) {
			return validator.NewErrorPos("tuple items length mismatch")
		}
	} else {
		if len(items) < len(self.Children) {
			return validator.NewErrorPos("data items length smaller than expected")
		}
	}

	for i, schema := range self.Children {
		item := items[i]
		if errPos := validator.Scan(schema, fmt.Sprintf("[%d]", i), item); errPos != nil {
			return errPos
		}
	}
	if self.AdditionalSchema != nil {
		for i, item := range items[len(self.Children):] {
			pos := fmt.Sprintf("[%d]", i+len(self.Children))
			if errPos := validator.Scan(self.AdditionalSchema, pos, item); errPos != nil {
				return errPos
			}
		}
	}
	return nil
}

// type = "method"
func NewMethodSchema() *MethodSchema {
	return &MethodSchema{Params: make([]Schema, 0), Result: nil}
}
func (self MethodSchema) Type() string {
	return "method"
}

func (self MethodSchema) RebuildType() map[string]interface{} {
	tp := self.rebuildType(self.Type())
	arr := make([](map[string]interface{}), 0)
	for _, p := range self.Params {
		arr = append(arr, p.RebuildType())
	}
	tp["params"] = arr
	tp["result"] = self.Result
	return tp
}

func (self *MethodSchema) Scan(validator *SchemaValidator, data interface{}) *ErrorPos {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return validator.NewErrorPos("data is not object")
	}

	if params, ok := convertList(dataMap, "params", false); ok {
		errPos := self.ScanParams(validator, params)
		if errPos != nil {
			return errPos
		}
		return nil
	}

	if result, ok := dataMap["result"]; ok {
		errPos := self.ScanResult(validator, result)
		return errPos
	}

	return validator.NewErrorPos("data is not a JSONRPC message")
}

func (self *MethodSchema) ScanParams(validator *SchemaValidator, params []interface{}) *ErrorPos {
	validator.pushPath(".params")
	defer validator.popPath(".params")

	if len(params) != len(self.Params) {
		return validator.NewErrorPos("length of params mismatch")
	}
	for i, paramSchema := range self.Params {
		errPos := validator.Scan(paramSchema, fmt.Sprintf("[%d]", i), params[i])
		if errPos != nil {
			return errPos
		}
	}
	return nil
}

func (self *MethodSchema) ScanResult(validator *SchemaValidator, result interface{}) *ErrorPos {
	if self.Result != nil {
		return validator.Scan(self.Result, ".result", result)
	}
	return nil
}

// type = "object"
func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		Properties: make(map[string]Schema),
		Requires:   make(map[string]bool),
	}
}

func (self ObjectSchema) Type() string {
	return "object"
}

func (self ObjectSchema) RebuildType() map[string]interface{} {
	tp := self.rebuildType(self.Type())
	props := make(map[string]interface{})
	for name, p := range self.Properties {
		props[name] = p.RebuildType()
	}
	arr := make([]string, 0)
	for name, _ := range self.Requires {
		arr = append(arr, name)
	}
	tp["requires"] = arr
	return tp
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
				// prop is required but not present
				validator.pushPath("." + prop)
				errPos := validator.NewErrorPos("required prop is not present")
				validator.popPath("." + prop)
				return errPos
			}
		}
	}
	return nil
}

func SchemaToString(schema Schema) string {
	s := schema.RebuildType()
	schemaJson := simplejson.New()
	schemaJson.SetPath(nil, s)
	schemaBytes, err := schemaJson.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(schemaBytes)
}
