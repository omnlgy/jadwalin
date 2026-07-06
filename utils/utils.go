package utils

import "strings"

func PascalToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 {
			prev := rune(s[i-1])
			// lowercase → uppercase boundary  (e.g. FullName → full_name)
			if prev >= 'a' && prev <= 'z' && r >= 'A' && r <= 'Z' {
				result = append(result, '_')
			}
			// uppercase → lowercase boundary within an acronym  (e.g. OTPCode → otp_code)
			if i < len(s)-1 {
				next := rune(s[i+1])
				if prev >= 'A' && prev <= 'Z' && r >= 'A' && r <= 'Z' && next >= 'a' && next <= 'z' {
					result = append(result, '_')
				}
			}
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
