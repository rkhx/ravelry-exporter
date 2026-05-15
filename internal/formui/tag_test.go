package formui

import (
	"reflect"
	"testing"
)

func TestParseFyneFormTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected map[string]string
	}{
		{
			name: "key-value pairs",
			tag:  "label=Username,placeholder=Enter name",
			expected: map[string]string{
				"label":       "Username",
				"placeholder": "Enter name",
			},
		},
		{
			name: "with flags (booleans)",
			tag:  "label=Password,required",
			expected: map[string]string{
				"label":    "Password",
				"required": "true",
			},
		},
		{
			name: "spaces around keys and values",
			tag:  "  label = Email , required , placeholder = user@example.com ",
			expected: map[string]string{
				"label":       "Email",
				"required":    "true",
				"placeholder": "user@example.com",
			},
		},
		{
			name:     "empty string",
			tag:      "",
			expected: nil,
		},
		{
			name: "only flags",
			tag:  "required,readonly",
			expected: map[string]string{
				"required": "true",
				"readonly": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseFyneFormTag(tt.tag)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("got %+v, want %+v", actual, tt.expected)
			}
		})
	}
}
