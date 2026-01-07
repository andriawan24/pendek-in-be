package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

const (
	shortCodeAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultShortCodeN = 8
)

func GenerateShortCode() string {
	const n = defaultShortCodeN

	out := make([]byte, n)
	max := big.NewInt(int64(len(shortCodeAlphabet)))

	for i := range n {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			return fallbackShortCode(n)
		}

		out[i] = shortCodeAlphabet[v.Int64()]
	}

	return string(out)
}

func fallbackShortCode(n int) string {
	if n <= 0 {
		return ""
	}

	x := uint64(time.Now().UnixNano())
	out := make([]byte, 0, n)
	for len(out) < n {
		out = append(out, shortCodeAlphabet[x%uint64(len(shortCodeAlphabet))])
		x = x/uint64(len(shortCodeAlphabet)) + uint64(time.Now().UnixNano())
	}
	return string(out[:n])
}
