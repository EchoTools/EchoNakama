package login

import (
	"testing"
)

func TestFilterDisplayName(t *testing.T) {
	tests := []struct {
		displayName    string
		expectedResult string
	}{
		{"John_Doe", "John_Doe"},                         // Test with valid display name
		{"Jane-Smith!@#", "Jane-Smith"},                  // Test with valid display name
		{"Jane-S!@#mith", "Jane-Smith"},                  // Test with valid display name
		{"Jane-Smith", "Jane-Smith"},                     // Test with valid display name
		{"Jane-Smith [DEMO]", "Jane-Smith[DEMO]"},        // Test with valid display name
		{"12345678901234567890", "12345678901234567890"}, // Test with display name exceeding 20 characters
		{"a", ""}, // Test with display name exceeding 20 characters
		{"!@#$%^&*()_+{}|:\"<>?`-=[]\\;',./", ""}, // Test with display name containing special characters only
		{"", ""}, // Test with empty display name
	}

	for _, tt := range tests {
		result := FilterDisplayName(tt.displayName)
		if result != tt.expectedResult {
			t.Errorf("FilterDisplayName(%s) = %s, want %s", tt.displayName, result, tt.expectedResult)
		}
	}
}
