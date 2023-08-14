package validator

import (
	"fmt"
	"regexp"
)

func IsDomaineNameValid(domainName string) bool {
	// A regular expression pattern for a typical domain name
	// Matches domains that start with alphanumeric characters, possibly containing dashes in the middle, and ending with alphanumeric characters
	// The domain name must contain at least one dot, and each section (separated by dots) must be between 1 and 63 characters long
	// The TLD (top-level domain) must be between 2 and 6 alphanumeric characters
	const pattern = `^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,6}$`

	// Compile the pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return false
	}

	// Check if the domain matches the pattern
	return re.MatchString(domainName)
}
