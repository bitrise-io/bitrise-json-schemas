package schemas_test

import (
	"strings"
	"testing"

	schemas "github.com/bitrise-io/bitrise-json-schemas"
	"github.com/bitrise-io/bitrise-json-schemas/validator"
)

func TestStepSchema(t *testing.T) {
	for _, tt := range tests {
		validator, err := validator.NewJSONSchemaValidator(schemas.StepSchema)
		if err != nil {
			t.Fatalf("unexpected schema compile error: %v", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			_, errors, err := validator.Validate(tt.stepYML)
			if err != nil {
				t.Errorf("validation has failed: %s", err)
			}

			if tt.wantErr == "" && errors != nil {
				t.Errorf("unexpected error: %v", strings.Join(errors, "\n"))
			}
			if tt.wantErr != "" && errors == nil {
				t.Errorf("expected error: %s, got nil", tt.wantErr)
			}
			if tt.wantErr != "" && errors != nil && strings.Join(errors, "\n") != tt.wantErr {
				t.Errorf("expected error: %s, got: %s", tt.wantErr, strings.Join(errors, "\n"))
			}
		})
	}
}

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
    summary: Type your script here.
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
		wantErr: `I[#/website] S[#/definitions/URL/type] expected string, but got null`,
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
		wantErr: `I[#/support_url] S[#/definitions/URL/type] expected string, but got null`,
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
		wantErr: `I[#/source_code_url] S[#/definitions/URL/type] expected string, but got null`,
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
		wantErr: `I[#/website] S[#/definitions/URL/format] "git@github.com:bitrise-steplib/steps-script.git" is not valid "uri"
I[#/website] S[#/definitions/URL/pattern] does not match pattern "^https?://"`,
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
		wantErr: `I[#/source_code_url] S[#/definitions/URL/format] "git@github.com:bitrise-steplib/steps-script.git" is not valid "uri"
I[#/source_code_url] S[#/definitions/URL/pattern] does not match pattern "^https?://"`,
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
		wantErr: `I[#/support_url] S[#/definitions/URL/format] "git@github.com:bitrise-steplib/steps-script.git" is not valid "uri"
I[#/support_url] S[#/definitions/URL/pattern] does not match pattern "^https?://"`,
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
		wantErr: `I[#/is_always_run] S[#/then/properties/is_always_run/const] value must be true`,
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
		wantErr: `I[#/deps/brew/0/name] S[#/definitions/BrewDepModel/properties/name/type] expected string, but got null`,
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
		wantErr: `I[#/deps/brew/0/name] S[#/definitions/BrewDepModel/properties/name/not] not failed`,
	},
	{
		name: "deprecated dependencies property is not allowed",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
dependencies:
- manager: brew
  name: tee
`,
		wantErr: `I[#] S[#/additionalProperties] additionalProperties "dependencies" not allowed`,
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
		wantErr: `I[#/toolkit] S[#/definitions/StepToolkitModel/maxProperties] maximum 1 properties allowed, but found 2 properties`,
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
		wantErr: `I[#/toolkit/go/package_name] S[#/definitions/GoStepToolkitModel/properties/package_name/type] expected string, but got null`,
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
		wantErr: `I[#/toolkit/bash/entry_file] S[#/definitions/BashStepToolkitModel/properties/entry_file/type] expected string, but got null`,
	},
	{
		name: "deprecated host os tags property is not allowed",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
host_os_tags:
- linux-docker-android-20.04
`,
		wantErr: `I[#] S[#/additionalProperties] additionalProperties "host_os_tags" not allowed`,
	},
	{
		name: "deprecated is requires admin user property is not allowed",
		stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
is_requires_admin_user: true
`,
		wantErr: `I[#] S[#/additionalProperties] additionalProperties "is_requires_admin_user" not allowed`,
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
    summary: Script content
`,
		wantErr: `I[#/inputs/0/opts] S[#/definitions/EnvVarOpts/required] missing properties: "title"`,
	},
	{
		name: "input summary is required",
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
		wantErr: `I[#/inputs/0/opts] S[#/definitions/EnvVarOpts/required] missing properties: "summary"`,
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
    summary: Script content
    is_sensitive: true
    is_expand: false
`,
		wantErr: `I[#/inputs/0/opts/is_expand] S[#/definitions/EnvVarOpts/then/properties/is_expand/const] value must be true`,
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
    summary: Script content
    value_options:
    - "yes"
    - "no"
`,
		wantErr: `I[#/inputs/0/content] S[#/definitions/InputEnvVar/allOf/1/then/additionalProperties/type] expected string, but got null`,
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
    summary: Script content
    value_options:
    - true
    - false
`,
		wantErr: `I[#/inputs/0/opts/value_options/0] S[#/definitions/EnvVarOpts/properties/value_options/items/type] expected string, but got boolean
I[#/inputs/0/opts/value_options/1] S[#/definitions/EnvVarOpts/properties/value_options/items/type] expected string, but got boolean`,
	},
}
