package schemas

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v3"
	"gopkg.in/yaml.v2"
)

var WarningPatters = []string{
	`#: additionalProperties .+ not allowed`, // is_requires_admin_user, host_os_tags, dependencies are deprecated
	`#: missing properties: .+`,              // support_url, source_code_url are required
	`#/summary: does not match pattern "\^\.\{1,100\}\$"`,
	`#/deps/(brew|apt_get)+/\d+/(name|bin_name)+: not failed`, // go listed as dependency
	`#/(inputs|outputs)+/\d+/opts: missing properties: "summary"`,
	`#/(inputs|outputs)+/\d+/opts/summary: length must be >= 1, but got 0`,
	`#/inputs/\d+/.+: expected .+, but got .+`, // input value is not a string or null
	`#/inputs/\d+/opts/value_options: minimum 2 items allowed, but found \d+ items`,
}

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

func (v JSONSchemaValidator) Validate(ymlStr string, warningPatterns ...string) ([]string, []string, error) {
	var m interface{}
	err := yaml.Unmarshal([]byte(ymlStr), &m)
	if err != nil {
		return nil, nil, err
	}
	m, err = recursiveJSONMarshallable(m)
	if err != nil {
		return nil, nil, err
	}

	if err := v.schema.ValidateInterface(m); err != nil {
		validationErr := &jsonschema.ValidationError{}
		if errors.As(err, &validationErr) {
			warnings, errors := collectIssues(*validationErr, warningPatterns)
			return warnings, errors, nil
		}
		return nil, nil, err
	}

	return nil, nil, nil
}

func collectIssues(err jsonschema.ValidationError, warningPatterns []string) (warnings []string, errors []string) {
	var issues []string
	issues = recursivelyCollectIssues(err, issues)

	for _, issue := range issues {
		isWarning := false
		for _, pattern := range warningPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(issue) {
				isWarning = true
				warnings = append(warnings, issue)
				break
			}
		}
		if !isWarning {
			errors = append(errors, issue)
		}
	}

	return warnings, errors
}

func recursivelyCollectIssues(err jsonschema.ValidationError, issues []string) []string {
	if len(err.Causes) == 0 {
		issues = append(issues, fmt.Sprintf("I[%s] S[%s] %s", err.InstancePtr, err.SchemaPtr, err.Message))
		return issues
	}

	for _, cause := range err.Causes {
		issues = recursivelyCollectIssues(*cause, issues)
	}

	return issues
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
