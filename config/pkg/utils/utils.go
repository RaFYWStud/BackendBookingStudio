package utils

// SafeCompareString securely compares two strings to prevent timing attacks.
// This function ensures that the comparison time remains constant regardless of differences in characters,
// making it useful for comparing sensitive data such as passwords or tokens.
func SafeCompareString(a, b string) bool {
	// If the lengths are different, return false immediately
	if len(a) != len(b) {
		return false
	}

	// Variable to store the result of bitwise comparisons
	var result byte = 0

	// Iterate through each character in the string and perform a bitwise XOR operation
	for i := range a {
		result |= a[i] ^ b[i] // If there is a difference, result will become non-zero
	}

	// If result remains 0, all characters are identical, return true
	return result == 0
}
