package crypto

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type HashKind uint

const (
	Invalid HashKind = iota
	BCRYPT
)

func PasswordHash(password []byte, algo HashKind) (string, error) {
	if algo != BCRYPT {
		return "", errors.New("invalid password hash algo")
	}
	result, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), err
}

func PasswordVerify(password []byte, hash string) (bool, error) {
	if len(hash) <= 4 {
		return false, errors.New("invalid password hash format [" + hash + "]")
	}
	hashAlgo := hash[1:3]
	if hashAlgo != "2a" && hashAlgo != "2y" {
		return false, errors.New("invalid password hash algo [" + hashAlgo + "]")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), password)
	return err == nil, nil
}
