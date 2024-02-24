package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func SignData(data []byte, signingKey string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(signingKey))

	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
