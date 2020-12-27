package schema

import (
	"errors"
	"fmt"
	//"reflect"
	json "encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	"strings"
)

// util functions
func convertMap(node map[string]interface{}, attrName string, optional bool) (map[string]interface{}, bool) {
	if v, ok := node[attrName]; ok {
		if aMap, ok := v.(map[string]interface{}); ok {
			return aMap, ok
		}
	} else if optional {
		return make(map[string]interface{}), true
	}
	return nil, false
}

func convertList(node map[string]interface{}, attrName string, optional bool) ([]interface{}, bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			return aList, ok
		}
	} else if optional {
		return [](interface{}){}, true
	}
	return nil, false
}

func convertMapOfMap(node map[string](interface{}), attrName string, optional bool) (map[string](map[string]interface{}), bool) {
	if mm, ok := convertMap(node, attrName, optional); ok {
		resMap := make(map[string](map[string]interface{}))
		for name, value := range mm {
			mv, ok := value.(map[string]interface{})
			if !ok {
				return nil, false
			}
			resMap[name] = mv
		}
		return resMap, true
	}
	return nil, false
}

func convertListOfMap(node map[string]interface{}, attrName string, optional bool) ([](map[string]interface{}), bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			arr := make([](map[string]interface{}), 0)
			for _, item := range aList {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					return nil, false
				}
				arr = append(arr, itemMap)
			}
			return arr, true
		}
	} else if optional {
		return [](map[string]interface{}){}, true
	}
	return nil, false
}

func convertListOfString(node map[string]interface{}, attrName string, optional bool) ([]string, bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			arr := make([]string, 0)
			for _, item := range aList {
				strItem, ok := item.(string)
				if !ok {
					return nil, false
				}
				arr = append(arr, strItem)
			}
			return arr, true
		}
	} else if optional {
		return []string{}, true
	}
	return nil, false
}

// schema buillder error
func (self SchemaBuildError) Error() string {
	return fmt.Sprintf("SchemaBuildError %s", self.info)
}

func NewBuildError(info string) *SchemaBuildError {
	return &SchemaBuildError{info: info}
}

// schema validator
func (self ErrorPos) Path() string {
	return strings.Join(self.paths, "")
}

func (self ErrorPos) Error() string {
	return fmt.Sprintf("Validation Error: %s %s", self.Path(), self.hint)
}

func (self ErrorPos) ToMessage(id interface{}) *jsonrpc.ErrorMessage {
	err := &jsonrpc.RPCError{10401, self.Error(), false}
	return err.ToMessage(id)
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
		if len(self.paths) < 1 || self.paths[len(self.paths)-1] != path {
			panic(errors.New(fmt.Sprintf("pop path %s is different from stack top %s", path, self.paths[len(self.paths)-1])))
		}
		self.paths = self.paths[:len(self.paths)-1]
	}
}

func (self *SchemaValidator) Scan(schema Schema, path string, data interface{}) *ErrorPos {
	self.pushPath(path)
	errPos := schema.Scan(self, data)
	self.popPath(path)
	return errPos
}

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
		return nil, NewBuildError("data is not a object")
	}
}

func (self *SchemaBuilder) buildNode(node map[string](interface{})) (Schema, error) {
	nodeType, ok := node["type"]
	if !ok {
		return nil, NewBuildError("no type presented")
	}
	var schema Schema = nil
	var err error = nil

	switch nodeType {
	case "number":
		schema = NewNumberSchema()
	case "any":
		schema = &AnySchema{}
	case "null":
		schema = &NullSchema{}
	case "string":
		schema = &StringSchema{}
	case "bool":
		schema = &BoolSchema{}
	case "union":
		schema, err = self.buildUnionSchema(node)
	case "list":
		schema, err = self.buildListSchema(node)
	case "object":
		schema, err = self.buildObjectSchema(node)
	case "method":
		schema, err = self.buildMethodSchema(node)
	default:
		err = NewBuildError("unknown type")
	}

	if err != nil {
		return nil, err
	}

	err = self.buildMixin(schema, node)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (self *SchemaBuilder) buildMixin(schema Schema, node map[string]interface{}) error {
	if name, ok := node["name"]; ok {
		if strName, ok := name.(string); ok {
			schema.SetName(strName)
		} else {
			return NewBuildError("name must be string")
		}
	}

	if desc, ok := node["description"]; ok {
		if strDesc, ok := desc.(string); ok {
			schema.SetDescription(strDesc)
		} else {
			return NewBuildError("decsription must be string")
		}
	}
	return nil
}

func (self *SchemaBuilder) buildListSchema(node map[string](interface{})) (Schema, error) {
	items, ok := node["items"]
	if !ok {
		return nil, NewBuildError("no items")
	}

	if itemsMap, ok := items.(map[string]interface{}); ok {
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

func (self *SchemaBuilder) buildUnionSchema(node map[string](interface{})) (*UnionSchema, error) {
	schema := NewUnionSchema()
	if choices, ok := convertListOfMap(node, "anyOf", false); ok {
		for _, choiceNode := range choices {
			c, err := self.buildNode(choiceNode)
			if err != nil {
				return nil, err
			}
			schema.Choices = append(schema.Choices, c)
		}
	} else {
		return nil, NewBuildError("no valid anyOf attribute")
	}
	return schema, nil
}

func (self *SchemaBuilder) buildMethodSchema(node map[string](interface{})) (*MethodSchema, error) {
	schema := NewMethodSchema()
	if paramsNodes, ok := convertListOfMap(node, "params", false); ok {
		for _, paramNode := range paramsNodes {
			c, err := self.buildNode(paramNode)
			if err != nil {
				return nil, err
			}
			schema.Params = append(schema.Params, c)
		}
	} else {
		return nil, NewBuildError("params is not a list of objects")
	}

	if resultNode, ok := convertMap(node, "result", true); ok {
		if _, ok := resultNode["type"]; !ok {
			resultNode["type"] = "any"
		}
		c, err := self.buildNode(resultNode)
		if err != nil {
			return nil, err
		}
		schema.Result = c
	}
	return schema, nil
}

func (self *SchemaBuilder) buildObjectSchema(node map[string](interface{})) (*ObjectSchema, error) {
	schema := NewObjectSchema()
	if propNodes, ok := convertMapOfMap(node, "properties", false); ok {
		for propName, propNode := range propNodes {
			child, err := self.buildNode(propNode)
			if err != nil {
				return nil, err
			}
			schema.Properties[propName] = child
		}
	} else {
		return nil, NewBuildError("properties is not a map of objects")
	}

	if requireList, ok := convertListOfString(node, "requires", true); ok {
		for _, reqProp := range requireList {
			if _, found := schema.Properties[reqProp]; !found {
				return nil, NewBuildError("cannot find required prop")
			}
			schema.Requires[reqProp] = true
		}
	} else {
		return nil, NewBuildError("requires is not a list of strings")
	}
	return schema, nil
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
