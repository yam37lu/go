package utils

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm4"
)

const (
	SM4 = "9ijn5TGB8uhb7UJM"
)

func ConfigEncrypt(content string) (string, error) {
	dst, err := sm4.Sm4Ecb([]byte(SM4), []byte(content), true)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(dst), nil
}

func ConfigDecrypt(content string) (string, error) {
	encode, err := hex.DecodeString(content)
	if err != nil {
		return "", err
	}
	dst, err := sm4.Sm4Ecb([]byte(SM4), encode, false)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}
