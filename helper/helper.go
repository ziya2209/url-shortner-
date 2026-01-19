package helper

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func generateSalt() string {
	// ato z AtoZ 0to9 it will return 5 digit char
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 5
	solt := make([]byte, length)
	for i := range length {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		solt[i] = charset[v.Int64()]

	}

	return string(solt)
}

func GenerateShortUrl(user_id, Long_url string) string {
	hash := md5.Sum([]byte(user_id + Long_url + generateSalt()))
	shortURL := hex.EncodeToString(hash[:])[:5] // Take first 5 characters of hex hash
	return shortURL
}
