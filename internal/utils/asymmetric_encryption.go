package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/daremove/go-metrics-service/internal/logger"
	"go.uber.org/zap"
)

const (
	RSAKeySize = 2048
	ChunkSize  = (RSAKeySize / 8) - 2*sha256.Size - 2
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pemFile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemFile)

	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}

func EncryptWithPublicKey(data []byte, key *rsa.PublicKey) ([]byte, error) {
	numChunks := (len(data) + ChunkSize - 1) / ChunkSize
	encryptedChunks := make([]byte, 0, numChunks*RSAKeySize/8)

	for start := 0; start < len(data); start += ChunkSize {
		end := start + ChunkSize

		if end > len(data) {
			end = len(data)
		}

		chunk := data[start:end]
		encryptedChunk, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, chunk, nil)

		if err != nil {
			return nil, err
		}

		encryptedChunks = append(encryptedChunks, encryptedChunk...)
	}

	return encryptedChunks, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemFile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemFile)

	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func DecryptWithPrivateKey(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	numChunks := len(data) / ChunkSize
	decryptedMessage := make([]byte, 0, numChunks*ChunkSize)

	chunkSize := RSAKeySize / 8

	for start := 0; start < len(data); start += chunkSize {
		end := start + chunkSize

		if end > len(data) {
			end = len(data)
		}

		chunk := data[start:end]
		decryptedChunk, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, chunk, nil)

		if err != nil {
			return nil, err
		}

		decryptedMessage = append(decryptedMessage, decryptedChunk...)
	}

	return decryptedMessage, nil
}

func DecryptMiddleware(privateKey *rsa.PrivateKey) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)

			if err != nil {
				logger.Log.Error("error read body data", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			decryptedMessage, err := DecryptWithPrivateKey(body, privateKey)

			if err != nil {
				http.Error(w, "Failed to decrypt message", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decryptedMessage))

			next.ServeHTTP(w, r)
		}
	}
}
