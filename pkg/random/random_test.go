package random

import (
	"strings"
	"testing"
)

func TestRandomInt(t *testing.T) {
	min := int64(5)
	max := int64(10)
	for i := 0; i < 1000; i++ {
		result := RandomInt(min, max)
		if result < min || result > max {
			t.Errorf("RandomInt(%d, %d) = %d; want value in range [%d, %d]", min, max, result, min, max)
		}
	}
}

func TestRandomString(t *testing.T) {
	length := 10
	result := RandomString(length)
	if len(result) != length {
		t.Errorf("RandomString(%d) = %s; want length of %d", length, result, length)
	}
	for _, char := range result {
		if !strings.ContainsRune(alphabet, char) {
			t.Errorf("RandomString(%d) contains unexpected character '%c'", length, char)
		}
	}
}

func TestRandomDomainName(t *testing.T) {
	domain := RandomDomainName()
	if !strings.HasSuffix(domain, ".com") || len(domain) <= 4 {
		t.Errorf("RandomDomainName() = %s; want domain ending with .com and length greater than 4", domain)
	}
}

func TestRandomOwner(t *testing.T) {
	owner := RandomOwner()
	if len(owner) != 6 {
		t.Errorf("RandomOwner() = %s; want length of 6", owner)
	}
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()
	if !strings.HasSuffix(email, "@email.com") || len(email) <= 10 {
		t.Errorf("RandomEmail() = %s; want email ending with @email.com and length greater than 10", email)
	}
}
