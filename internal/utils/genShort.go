package utils

func GenShort(count int) string {
	var result string
	for count > 0 {
		count--
		result = string(rune('A'+count%26)) + result
		count /= 26
	}
	return result

}
