package util

// IntPtr returns a pointer to the given int value.
func IntPtr(i int) *int {
	return &i
}
