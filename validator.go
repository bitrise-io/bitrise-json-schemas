package schemas

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v3"
	"gopkg.in/yaml.v2"
)

type JSONSchemaValidator struct {
	schema *jsonschema.Schema
}

func NewJSONSchemaValidator(schemaStr string) (*JSONSchemaValidator, error) {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaStr)); err != nil {
		return nil, err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}

	return &JSONSchemaValidator{
		schema: schema,
	}, nil
}

func (v JSONSchemaValidator) Validate(ymlStr string) error {
	var m interface{}
	err := yaml.Unmarshal([]byte(ymlStr), &m)
	if err != nil {
		return err
	}
	m, err = recursiveJSONMarshallable(m)
	if err != nil {
		return err
	}

	if err := v.schema.ValidateInterface(m); err != nil {
		return err
	}
	return nil
}

func recursiveJSONMarshallable(source interface{}) (interface{}, error) {
	if array, ok := source.([]interface{}); ok {
		var convertedArray []interface{}
		for _, element := range array {
			convertedValue, err := recursiveJSONMarshallable(element)
			if err != nil {
				return nil, err
			}
			convertedArray = append(convertedArray, convertedValue)
		}
		return convertedArray, nil
	}

	if interfaceToInterfaceMap, ok := source.(map[interface{}]interface{}); ok {
		target := map[string]interface{}{}
		for key, value := range interfaceToInterfaceMap {
			strKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("failed to convert map key from type interface{} to string")
			}

			convertedValue, err := recursiveJSONMarshallable(value)
			if err != nil {
				return nil, err
			}
			target[strKey] = convertedValue
		}
		return target, nil
	}

	if stringToInterfaceMap, ok := source.(map[string]interface{}); ok {
		target := map[string]interface{}{}
		for key, value := range stringToInterfaceMap {
			convertedValue, err := recursiveJSONMarshallable(value)
			if err != nil {
				return nil, err
			}
			target[key] = convertedValue
		}
		return target, nil
	}

	return source, nil
}
