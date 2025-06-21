package base62

import "math/rand/v2"

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

// Generates a random base62 string with the length in the range [min, max]
func RandomSeqRange(min int, max int) string {
	return RandSeq(min + rand.IntN(max-min+1))
}
