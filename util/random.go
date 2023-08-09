package util

import "math/rand"

// Funcs used for testing.

// RandomString generates a random string of length n.
func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[RandomInt(len(letters))]
	}
	return string(result)
}

// RandomInt generates a random integer between 0 and max.
func RandomInt(max int) int {
	return rand.Intn(max)
}
