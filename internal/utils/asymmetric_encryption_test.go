package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSAKeySize)

	require.NoError(t, err)

	return privateKey, &privateKey.PublicKey
}

func writePEMFile(t *testing.T, path string, block *pem.Block) {
	file, err := os.Create(path)
	require.NoError(t, err)

	defer file.Close()

	err = pem.Encode(file, block)

	require.NoError(t, err)
}

func TestLoadPublicKey(t *testing.T) {
	t.Run("Should load valid RSA public key from PEM file", func(t *testing.T) {
		_, publicKey := generateTestKeys(t)
		pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)

		require.NoError(t, err)

		pubPEM := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}
		pubPath := "test_public_key.pem"

		writePEMFile(t, pubPath, pubPEM)
		defer os.Remove(pubPath)

		loadedPubKey, err := LoadPublicKey(pubPath)

		assert.NoError(t, err)
		assert.Equal(t, publicKey, loadedPubKey)
	})

	t.Run("Should return error for invalid PEM file", func(t *testing.T) {
		invalidPEM := []byte("invalid pem content")
		invalidPath := "invalid_public_key.pem"

		err := os.WriteFile(invalidPath, invalidPEM, 0644)
		defer os.Remove(invalidPath)

		require.NoError(t, err)

		_, err = LoadPublicKey(invalidPath)

		assert.Error(t, err)
	})
}

func TestEncryptWithPublicKey(t *testing.T) {
	t.Run("Should encrypt data with RSA public key", func(t *testing.T) {
		_, publicKey := generateTestKeys(t)
		data := []byte("test data to encrypt")

		encryptedData, err := EncryptWithPublicKey(data, publicKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedData)
	})

	t.Run("Should handle empty data encryption", func(t *testing.T) {
		_, publicKey := generateTestKeys(t)
		data := []byte("")

		encryptedData, err := EncryptWithPublicKey(data, publicKey)

		assert.NoError(t, err)
		assert.Equal(t, encryptedData, []byte{})
	})
}

func TestLoadPrivateKey(t *testing.T) {
	t.Run("Should load valid RSA private key from PEM file", func(t *testing.T) {
		privateKey, _ := generateTestKeys(t)
		privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
		privPath := "test_private_key.pem"

		writePEMFile(t, privPath, privPEM)
		defer os.Remove(privPath)

		loadedPrivKey, err := LoadPrivateKey(privPath)

		assert.NoError(t, err)
		assert.Equal(t, privateKey, loadedPrivKey)
	})

	t.Run("Should return error for invalid PEM file", func(t *testing.T) {
		invalidPEM := []byte("invalid pem content")
		invalidPath := "invalid_private_key.pem"

		err := os.WriteFile(invalidPath, invalidPEM, 0644)
		defer os.Remove(invalidPath)

		require.NoError(t, err)

		_, err = LoadPrivateKey(invalidPath)

		assert.Error(t, err)
	})
}

func TestDecryptWithPrivateKey(t *testing.T) {
	t.Run("Should decrypt data with RSA private key", func(t *testing.T) {
		privateKey, publicKey := generateTestKeys(t)
		data := []byte("test data to encrypt and decrypt")

		encryptedData, err := EncryptWithPublicKey(data, publicKey)

		require.NoError(t, err)

		decryptedData, err := DecryptWithPrivateKey(encryptedData, privateKey)

		assert.NoError(t, err)
		assert.Equal(t, data, decryptedData)
	})

	t.Run("Should return error for invalid encrypted data", func(t *testing.T) {
		privateKey, _ := generateTestKeys(t)
		invalidData := []byte("invalid encrypted data")

		_, err := DecryptWithPrivateKey(invalidData, privateKey)

		assert.Error(t, err)
	})
}

func TestDecryptMiddleware(t *testing.T) {
	t.Run("Should decrypt request body and call next handler", func(t *testing.T) {
		privateKey, publicKey := generateTestKeys(t)
		data := []byte("test data to encrypt and decrypt")

		encryptedData, err := EncryptWithPublicKey(data, publicKey)

		require.NoError(t, err)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, data, body)
			w.WriteHeader(http.StatusOK)
		})

		middleware := DecryptMiddleware(privateKey)(nextHandler)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encryptedData))
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return error for invalid encrypted request body", func(t *testing.T) {
		privateKey, _ := generateTestKeys(t)
		invalidData := []byte("invalid encrypted data")

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fail()
		})

		middleware := DecryptMiddleware(privateKey)(nextHandler)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidData))
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
