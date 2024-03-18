package utils

import (
	"strconv"
)

func ValidateLuhn(input string) bool {
	var sum int
	var alternate bool
	for i := len(input) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(input[i]))
		if err != nil {
			return false
		}
		if alternate {
			n *= 2
			if n > 9 {
				n = (n % 10) + 1
			}
		}
		sum += n
		alternate = !alternate
	}
	return sum%10 == 0
}
