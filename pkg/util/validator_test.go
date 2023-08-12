package util

import "testing"

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		domain string
		want   bool
	}{
		{"example.com", true},
		{"www.example.co.uk", true},
		{"test-domain.com", true},
		{"-invalid.com", false},
		{"invalid-.com", false},
		{"invalid.", false},
		{"", false},
		{"invalid", false},
		{"ex@ample.com", false},
		{"example..com", false},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			if got := IsDomaineNameValid(tt.domain); got != tt.want {
				t.Errorf("isValidDomain(%q) = %v, want %v", tt.domain, got, tt.want)
			}
		})
	}
}
