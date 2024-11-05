package helpers

import "unicode"

func IsValidLuhn(number string) bool {
	var sum int
	double := false

	for i := len(number) - 1; i >= 0; i-- {
		r := rune(number[i])

		if !unicode.IsDigit(r) {
			return false
		}

		n := int(r - '0')
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}

	return sum%10 == 0
}
