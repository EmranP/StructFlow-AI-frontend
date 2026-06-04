package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost    uint32 = 1
	memoryCost  uint32 = 64 * 1024
	parallelism uint8  = 4
	keyLength   uint32 = 32
	saltLength         = 16
)

func Hash(password string) (string, error) {

	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memoryCost,
		parallelism,
		keyLength,
	)

	return base64.RawStdEncoding.EncodeToString(salt) +
			"." +
			base64.RawStdEncoding.EncodeToString(hash),
		nil
}

func Compare(password string, encoded string) bool {

	parts := make([]string, 0)

	for _, p := range split(encoded, '.') {
		parts = append(parts, p)
	}

	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memoryCost,
		parallelism,
		keyLength,
	)

	return subtle.ConstantTimeCompare(
		hash,
		expectedHash,
	) == 1
}

func split(s string, sep byte) []string {

	var result []string

	start := 0

	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}

	result = append(result, s[start:])

	return result
}
