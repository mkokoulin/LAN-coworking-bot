package encoder

import (
	"encoding/base64"
	"fmt"
)

func Encode(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

func Decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode error: %s", err)
	}

	return string(decoded), nil
}