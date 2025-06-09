package ascii

import (
	"strings"
	"testing"
)

func TestCenterText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "short text",
			input:    "test",
			expected: "                        test                         \n",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "                                                     \n",
		},
		{
			name:     "long text",
			input:    "this is a very long text that exceeds the grid width",
			expected: "this is a very long text that exceeds the grid width\n",
		},
		{
			name:     "exact width text",
			input:    strings.Repeat("x", GridWidth),
			expected: strings.Repeat("x", GridWidth) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := centerText(tt.input)
			if result != tt.expected {
				t.Errorf("centerText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
