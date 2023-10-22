package api

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_@.-]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString
)

func ValidateString(s string, min, max int) error {
	n := len(s)
	if n < min || n > max {
		return fmt.Errorf("string length must be between %d and %d", min, max)
	}
	return nil
}

func ValidateUsername(s string) error {
	if err := ValidateString(s, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(s) {
		return fmt.Errorf("username must be lowercase alphanumeric or underscore")
	}

	return nil
}

func ValidateFullName(s string) error {
	if err := ValidateString(s, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(s) {
		return fmt.Errorf("full name must be alphanumeric")
	}

	return nil
}

func ValidatePassword(s string) error {
	return ValidateString(s, 8, 64)
}

func ValidateEmail(s string) error {
	if err := ValidateString(s, 3, 64); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

func ValidateEmailID(s string) error {
	if err := ValidateString(s, 36, 36); err != nil {
		return fmt.Errorf("invalid email ID, must be UUID")
	}

	return nil
}

func ValidateEmailVerificationCode(s string) error {
	if err := ValidateString(s, 32, 32); err != nil {
		return fmt.Errorf("invalid email verification code")
	}

	return nil
}

func ValidateDomainName(domainName string) error {
	// A regular expression pattern for a typical domain name
	// Matches domains that start with alphanumeric characters, possibly containing dashes in the middle, and ending with alphanumeric characters
	// The domain name must contain at least one dot, and each section (separated by dots) must be between 1 and 63 characters long
	// The TLD (top-level domain) must be between 2 and 6 alphanumeric characters
	const pattern = `^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,6}$`

	// Compile the pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		wErr := fmt.Errorf("error compiling regex: %s", err)
		return fmt.Errorf("invalid domain name. example of valid domain name: example.com err=%w", wErr)
	}

	// Check if the domain matches the pattern
	valid := re.MatchString(domainName)

	if !valid {
		return fmt.Errorf("invalid domain name. example of valid domain name: example.com")
	}

	return nil
}
