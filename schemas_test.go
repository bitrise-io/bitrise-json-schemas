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
`,
	},
	{
		name: "Minimal valid step with inputs",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: Script content
`,
	},
	// Failing test cases
	{
		name: "title is not empty",
		stepYML: `
title: 
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/title] S[#/properties/title/type] expected string, but got null`,
	},
	{
		name: "summary is not empty",
		stepYML: `
title: Script
summary: 
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/summary] S[#/properties/summary/type] expected string, but got null`,
	},
	{
		name: "summary is a single line text",
		stepYML: `
title: Script
summary: |-
  Run any custom script you want.  
  The power is in your hands.  
  Use it wisely!  
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/summary] S[#/properties/summary/pattern] does not match pattern "^.{1,100}$"`,
	},
	{
		name: "summary has 100 chars at max",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely! Too long line! Too long line!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/summary] S[#/properties/summary/pattern] does not match pattern "^.{1,100}$"`,
	},
	{
		name: "website is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: 
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/website] S[#/properties/website/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "support url is not empty",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url:
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
`,
		wantErr: `I[#/source_code_url] S[#/properties/source_code_url/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "website url is in http format",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: git@github.com:bitrise-steplib/steps-script.git
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/website] S[#/properties/website/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "source code url is in http format",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: git@github.com:bitrise-steplib/steps-script.git
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		wantErr: `I[#/source_code_url] S[#/properties/source_code_url/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "support url is in http format",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: git@github.com:bitrise-steplib/steps-script.git
`,
		wantErr: `I[#/support_url] S[#/properties/support_url/$ref] doesn't validate with "#/definitions/URL"`,
	},
	{
		name: "unsupported type tag",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
type_tags:
- utility
- invalid
`,
		wantErr: `I[#/type_tags/1] S[#/properties/type_tags/items/enum] value must be one of "access-control", "artifact-info", "installer", "deploy", "utility", "dependency", "code-sign", "build", "test", "notification"`,
	},
	{
		name: "is always run is set if notification type tag",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
is_always_run: false
type_tags:
- notification
`,
		wantErr: `I[#] S[#/then] if-then failed`,
	},
	{
		name: "unsupported project type tag",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
project_type_tags:
- ios
- unsupported
`,
		wantErr: `I[#/project_type_tags/1] S[#/properties/project_type_tags/items/enum] value must be one of "ios", "macos", "android", "xamarin", "react-native", "cordova", "ionic", "flutter"`,
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
`,
		wantErr: `I[#/timeout] S[#/properties/timeout/exclusiveMinimum] must be > 0/1 but found 0`,
	},
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
`,
		wantErr: `I[#/deps] S[#/properties/deps/$ref] doesn't validate with "#/definitions/DepsModel"`,
	},
	{
		name: "go is not listed as deps",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
deps:
  brew:
  - name: go
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
`,
		wantErr: `I[#/dependencies/2] S[#/properties/dependencies/items/$ref] doesn't validate with "#/definitions/DependencyModel"`,
	},
	{
		name: "go is not listed as dependencies",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
dependencies:
- manager: brew
  name: go
`,
		wantErr: `I[#/dependencies/0] S[#/properties/dependencies/items/$ref] doesn't validate with "#/definitions/DependencyModel"`,
	},
	{
		name: "toolkit is either bas, go or undefined",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
toolkit:
  go:
    package_name: github.com/bitrise-io/steps-script
  bash:
    entry_file: step.sh

`,
		wantErr: `I[#/toolkit] S[#/properties/toolkit/$ref] doesn't validate with "#/definitions/StepToolkitModel"`,
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
`,
		wantErr: `I[#/toolkit] S[#/properties/toolkit/$ref] doesn't validate with "#/definitions/StepToolkitModel"`,
	},
	// Input tests
	{
		name: "input title is required",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: 
`,
		wantErr: `I[#/inputs/0] S[#/properties/inputs/items/$ref] doesn't validate with "#/definitions/InputEnvVar"`,
	},
	{
		name: "is expand is set if is sensitive is set",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: ""
  opts:
    title: Script content
    is_sensitive: true
    is_expand: false
`,
		wantErr: `I[#/inputs/0] S[#/properties/inputs/items/$ref] doesn't validate with "#/definitions/InputEnvVar"`,
	},
	{
		name: "input has default value is value options defined",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: 
  opts:
    title: Script content
    value_options:
    - "yes"
    - "no"
`,
		wantErr: `I[#/inputs/0] S[#/properties/inputs/items/$ref] doesn't validate with "#/definitions/InputEnvVar"`,
	},
	{
		name: "input value option elements are strings",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: "true"
  opts:
    title: Script content
    value_options:
    - true
    - false
`,
		wantErr: `I[#/inputs/0] S[#/properties/inputs/items/$ref] doesn't validate with "#/definitions/InputEnvVar"`,
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
