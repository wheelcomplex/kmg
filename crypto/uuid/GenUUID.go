package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var GenUUIDErrors = errors.New("gen uuid fail")

func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", GenUUIDErrors
	}
	return hex.EncodeToString(uuid), nil
}

func MustGenUUID() string {
	val, err := GenUUID()
	if err != nil {
		panic(err)
	}
	return val
}
