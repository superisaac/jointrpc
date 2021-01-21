package schema

import (
	"fmt"
	//"reflect"
	simplejson "github.com/bitly/go-simplejson"
)

// Schema build error
func (self SchemaBuildError) Error() string {
	return fmt.Sprintf("SchemaBuildError %s", self.info)
}

func NewBuildError(info string) *SchemaBuildError {
	return &SchemaBuildError{info: info}
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
	return self.buildNode(data)
}

func (self *SchemaBuilder) buildNode(data interface{}) (Schema, error) {
	if typeMap, ok := convertTypeMap(data); ok {
		return self.buildNodeMap(typeMap)
	} else {
		return nil, NewBuildError("data is not an object")
	}
}

func (self *SchemaBuilder) buildNodeMap(node map[string](interface{})) (Schema, error) {
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

	// build tuple
	if itemsTuple, ok := items.([]interface{}); ok {
		schema := NewTupleSchema()
		for _, item := range itemsTuple {
			childSchema, err := self.buildNode(item)
			if err != nil {
				return nil, err
			}
			schema.Children = append(schema.Children, childSchema)
		}
		// additional items
		if additional, ok := node["additionalItems"]; ok {
			addSchema, err := self.buildNode(additional)
			if err != nil {
				return nil, err
			}
			schema.AdditionalSchema = addSchema
		}

		return schema, nil
	}

	// build list
	//if itemsMap, ok := items.(map[string]interface{}); ok {

	itemSchema, err := self.buildNode(items)
	if err != nil {
		return nil, err
	}
	schema := NewListSchema()
	schema.Item = itemSchema
	return schema, nil
	//}

	return nil, NewBuildError("fail to build list schema")
}

func (self *SchemaBuilder) buildUnionSchema(node map[string](interface{})) (*UnionSchema, error) {
	schema := NewUnionSchema()
	if choices, ok := convertAttrListOfMap(node, "anyOf", false); ok {
		for _, choiceNode := range choices {
			c, err := self.buildNodeMap(choiceNode)
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
	if paramsNodes, ok := convertAttrListOfMap(node, "params", false); ok {
		for _, paramNode := range paramsNodes {
			c, err := self.buildNodeMap(paramNode)
			if err != nil {
				return nil, err
			}
			schema.Params = append(schema.Params, c)
		}
	} else {
		return nil, NewBuildError("params is not a list of objects")
	}

	if resultNode, ok := convertAttrMap(node, "result", true); ok {
		if _, ok := resultNode["type"]; !ok {
			resultNode["type"] = "any"
		}
		c, err := self.buildNodeMap(resultNode)
		if err != nil {
			return nil, err
		}
		schema.Result = c
	}
	return schema, nil
}

func (self *SchemaBuilder) buildObjectSchema(node map[string](interface{})) (*ObjectSchema, error) {
	schema := NewObjectSchema()
	if propNodes, ok := convertAttrMapOfMap(node, "properties", false); ok {
		for propName, propNode := range propNodes {
			child, err := self.buildNodeMap(propNode)
			if err != nil {
				return nil, err
			}
			schema.Properties[propName] = child
		}
	} else {
		return nil, NewBuildError("properties is not a map of objects")
	}

	if requireList, ok := convertAttrListOfString(node, "requires", true); ok {
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
