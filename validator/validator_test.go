package validator

import (
	"reflect"
	"testing"

	schemas "github.com/bitrise-io/bitrise-json-schemas"
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
			name: "Valid step.yml",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
source_code_url: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
		},
		{
			name: "error: source code url missing",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
			wantErrors: []string{`I[#] S[#/required] missing properties: "source_code_url"`},
		},
		{
			name: "warning: source code url missing",
			stepYML: `
title: Script
summary: Run any custom script you want. The power is in your hands. Use it wisely!
website: https://github.com/bitrise-io/steps-script
support_url: https://github.com/bitrise-io/steps-script/issues
`,
			warningPattern: `I\[#\] S\[#/required\] missing properties: .+`,
			wantWarnings:   []string{`I[#] S[#/required] missing properties: "source_code_url"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewJSONSchemaValidator(schemas.StepSchema)
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
