package cli

import "strconv"

func EnsureInt(t string, l rune) bool {
	_, err := strconv.Atoi(t)
	return err == nil
}
