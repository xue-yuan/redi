package utils

func Base62Encode(b [32]byte) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	result := ""

	for _, v := range b {
		result += string(charset[v%62])
	}

	return result
}
