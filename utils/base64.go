package utils

import (
	"encoding/base64"
	"golang.org/x/sys/unix"
)

func B64UrlEncode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func B64UrlDecode(data string) (string, error) {
	decByte, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	return unix.ByteSliceToString(decByte), nil
}
