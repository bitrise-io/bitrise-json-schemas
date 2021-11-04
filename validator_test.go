package schemas

import (
	"reflect"
	"testing"
)

func TestJSONSchemaValidator_Validate(t *testing.T) {
	tests := []struct {
		name           string
		stepYML        string
		warningPattern string
		wantWarnings   []string
		wantErrors     []string
	}{
		{
			name: "deprecated properties (is requires admin user, host os tags, dependencies) are not allowed",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
is_requires_admin_user: true
`,
			warningPattern: AdditionalPropertiesNotAllowedPattern,
			wantWarnings:   []string{`I[#] S[#/additionalProperties] additionalProperties "is_requires_admin_user" not allowed`},
		},
		{
			name: "source code url and support url are not empty",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
`,
			warningPattern: MissingPropertiesPattern,
			wantWarnings: []string{
				`I[#] S[#/required] missing properties: "support_url", "source_code_url"`,
			},
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
			warningPattern: SummaryDoesNotMatchPattern,
			wantWarnings:   []string{`I[#/summary] S[#/properties/summary/pattern] does not match pattern "^.{1,100}$"`},
		},
		{
			name: "go listed as dependency",
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
			warningPattern: DepsNotFailedPattern,
			wantWarnings:   []string{`I[#/deps/brew/0/name] S[#/definitions/BrewDepModel/properties/name/not] not failed`},
		},
		{
			name: "step input summary missing",
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
			warningPattern: InputOutputMissingSummaryPattern,
			wantWarnings:   []string{`I[#/inputs/0/opts] S[#/definitions/EnvVarOpts/required] missing properties: "summary"`},
		},
		{
			name: "step input empty summary",
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
    summary: ""
`,
			warningPattern: InputOutputEmptySummaryPattern,
			wantWarnings:   []string{`I[#/inputs/0/opts/summary] S[#/definitions/EnvVarOpts/properties/summary/minLength] length must be >= 1, but got 0`},
		},
		{
			name: "step input value is not a string or null",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
inputs:
- content: 2
  opts:
    title: Script content
    summary: Type your script here.
`,
			warningPattern: InputValueOptionsDefaultValuePattern,
			wantWarnings:   []string{`I[#/inputs/0/content] S[#/definitions/InputEnvVar/additionalProperties/type] expected null or string, but got number`},
		},
		{
			name: "step input value options has at least 2 elements",
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
    value_options: []
`,
			warningPattern: InputValueOptionsMinItemsPattern,
			wantWarnings:   []string{`I[#/inputs/0/opts/value_options] S[#/definitions/EnvVarOpts/properties/value_options/minItems] minimum 2 items allowed, but found 0 items`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewJSONSchemaValidator(StepSchema)
			if err != nil {
				t.Fatalf("Failed to create validator: %s", err)
			}

			var warningPatterns []string
			if tt.warningPattern != "" {
				warningPatterns = append(warningPatterns, tt.warningPattern)
			}
			warnings, errors, err := v.Validate(tt.stepYML, warningPatterns...)
			if err != nil {
				t.Fatalf("Validate() error = %v", err)
			}

			if !reflect.DeepEqual(errors, tt.wantErrors) {
				t.Errorf("Validate() got errors = %v, want errors %v", errors, tt.wantErrors)
			}
			if !reflect.DeepEqual(warnings, tt.wantWarnings) {
				t.Errorf("Validate() got warnings = %v, want warnings %v", warnings, tt.wantWarnings)
			}
		})
	}
}
