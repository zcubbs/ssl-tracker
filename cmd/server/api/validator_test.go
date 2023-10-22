package api

import "testing"

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		domain string
		noErr  bool
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
			if err := ValidateDomainName(tt.domain); err != nil && tt.noErr {
				t.Errorf("ValidateDomainName() error = %v, wantErr %v", err, tt.noErr)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		noErr    bool
	}{
		{"username", true},
		{"user_name", true},
		{"user-name", true},
		{"user.name", true},
		{"user_name123", true},
		{"user-name123", true},
		{"user.name123", true},
		{"user_name-123", true},
		{"user-name-123", true},
		{"user.name-123", true},
		{"user_name-123", true},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if err := ValidateUsername(tt.username); err != nil && tt.noErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.noErr)
			}
		})
	}
}

func TestIsValidFullName(t *testing.T) {
	tests := []struct {
		fullName string
		noErr    bool
	}{
		{"John Doe", true},
		{"John-Doe", true},
		{"John.Doe", true},
		{"John_Doe", true},
		{"John Doe123", true},
		{"John-Doe123", true},
		{"John.Doe123", true},
		{"John_Doe123", true},
	}

	for _, tt := range tests {
		t.Run(tt.fullName, func(t *testing.T) {
			if err := ValidateFullName(tt.fullName); err != nil && tt.noErr {
				t.Errorf("ValidateFullName() error = %v, wantErr %v", err, tt.noErr)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		noErr    bool
	}{
		{"password", true},
		{"password123", true},
		{"password-123", true},
		{"password_123", true},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			if err := ValidatePassword(tt.password); err != nil && tt.noErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.noErr)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		noErr bool
	}{
		{"test@example.com", true},
		{"test.com", false},
		{"test@example", false},
		{"test", false},
		{"test@", false},
		{"@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if err := ValidateEmail(tt.email); err != nil && tt.noErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.noErr)
			}
		})
	}
}

func TestIsValidEmailID(t *testing.T) {
	// TODO: Implement this test.
}

func TestIsValidEmailVerificationCode(t *testing.T) {
	// TODO: Implement this test.
}

func TestIsValidUUID(t *testing.T) {
	// TODO: Implement this test.
}

func TestValidateEmailVerificationCode(t *testing.T) {
	// TODO: Implement this test.
}
