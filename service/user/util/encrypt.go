package util

import (
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"log"
)

func ScryptPw(password string) (string, error) {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 6, 66, 22, 222, 11}
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw, nil
}
