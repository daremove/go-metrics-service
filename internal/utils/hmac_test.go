package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignData(t *testing.T) {
	t.Run("Should generate correct HMAC signature", func(t *testing.T) {
		data := []byte("test data")
		signingKey := "secret"
		expectedMAC := hmac.New(sha256.New, []byte(signingKey))
		expectedMAC.Write(data)
		expectedSignature := expectedMAC.Sum(nil)

		signature, err := SignData(data, signingKey)
		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, signature)
	})

	t.Run("Should return error when signing data", func(t *testing.T) {
		signingKey := "secret"

		signature, err := SignData(nil, signingKey)
		assert.NoError(t, err)
		assert.NotNil(t, signature)
	})
}
