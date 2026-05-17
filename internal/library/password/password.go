package password

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const saltLength = 16

func GenerateSalt() (string, error) {
	bytes := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashPassword(password string, salt string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", salt, password)))
	return hex.EncodeToString(h.Sum(nil))
}

func CreatePassword(password string) (hash string, salt string, err error) {
	salt, err = GenerateSalt()
	if err != nil {
		return "", "", err
	}
	hash = HashPassword(password, salt)
	return hash, salt, nil
}

func VerifyPassword(password, salt, hash string) bool {
	return HashPassword(password, salt) == hash
}
