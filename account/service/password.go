package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

func hashPassword(password string) (string, error) {
	// example for making salt -- https://play.golang.org/p/_Aw6WeWC42I
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from -- https://godoc.org/golang.org/x/crypto/scrypt
	sHash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(sHash), hex.EncodeToString(salt))
	return hashedPW, nil
}

func comparePasswords(storedPassword string, suppliedPassword string) (bool, error) {
	pwSalt := strings.Split(storedPassword, ".")

	salt, err := hex.DecodeString(pwSalt[1])
	if err != nil {
		return false, fmt.Errorf("Unable to verify user password\n")
	}

	sHash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, err
	}
	return hex.EncodeToString(sHash) == pwSalt[0], nil
}
