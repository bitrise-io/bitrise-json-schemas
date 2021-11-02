package schemas_test

import (
	"fmt"
	"strings"
	"testing"

	schemas "github.com/bitrise-io/bitrise-json-schemas"
	"github.com/santhosh-tekuri/jsonschema/v3"
	"gopkg.in/yaml.v2"
)

var tests = []struct {
	name    string
	stepYML string
	wantErr string
}{
	// Passing test cases
	{
		name: "Minimal valid step",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
	},
	// Failing test cases - bitrise run
	{
		name: "deps name can not be empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
deps:
  brew:
  - name: 
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/deps] S[#/properties/deps/$ref] doesn't validate with "#/definitions/DepsModel"`,
	},
	{
		name: "dependencies name can not be empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
dependencies:
- manager: brew
  name: 
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/dependencies/0] S[#/properties/dependencies/items/$ref] doesn't validate with "#/definitions/DependencyModel"`,
	},
	{
		name: "dependencies manager is either brew or _",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
dependencies:
- manager: brew
  name: tee
- manager: _
  name: xcode
- manager: aptget
  name: zip
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/dependencies/2] S[#/properties/dependencies/items/$ref] doesn't validate with "#/definitions/DependencyModel"`,
	},
	{
		name: "go toolkit package name is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
toolkit:
  go:
    package_name: 
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/toolkit] S[#/properties/toolkit/$ref] doesn't validate with "#/definitions/StepToolkitModel"`,
	},
	{
		name: "bash toolkit entry file is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
toolkit:
  bash:
    entry_file: 
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/toolkit] S[#/properties/toolkit/$ref] doesn't validate with "#/definitions/StepToolkitModel"`,
	},
	{
		name: "timtout > 0",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
timeout: 0
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/timeout] S[#/properties/timeout/exclusiveMinimum] must be > 0/1 but found 0`,
	},
	{
		name: "title is not empty",
		stepYML: `
title: 
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/title] S[#/properties/title/type] expected string, but got null`,
	},
	{
		name: "support url is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: 
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/support_url] S[#/properties/support_url/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "source code url is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: 
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: "Script content"
outputs:
- RUNNER_BIN: value
  opts:
    title: Runner binary
`,
		wantErr: `I[#/source_code_url] S[#/properties/source_code_url/$ref] doesn't validate with "#/definitions/URL"`,
	},
}

func TestStepSchema(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validate(tt.stepYML, schemas.StepSchema)
			if tt.wantErr == "" && gotErr != nil {
				t.Errorf("unexpected error: %v", gotErr)
			}
			if tt.wantErr != "" && gotErr == nil {
				t.Errorf("expected error: %s, got nil", tt.wantErr)
			}
			if tt.wantErr != "" && gotErr != nil && gotErr.Error() != tt.wantErr {
				t.Errorf("expected error: %s, got: %s", tt.wantErr, gotErr)
			}
		})
	}
}

func validate(stepYML, schemaStr string) error {
	var m interface{}
	err := yaml.Unmarshal([]byte(stepYML), &m)
	if err != nil {
		return err
	}
	m, err = recursiveJSONMarshallable(m)
	if err != nil {
		return err
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaStr)); err != nil {
		return err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return err
	}
	if err := schema.ValidateInterface(m); err != nil {
		return err
	}
	fmt.Println("validation successfull")
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
