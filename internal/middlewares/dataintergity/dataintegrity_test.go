package dataintergity

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataIntegrityMiddleware(t *testing.T) {
	signingKey := "test-signing-key"
	config := DataIntegrityMiddlewareConfig{
		SigningKey: signingKey,
	}

	middleware := NewMiddleware(config)

	t.Run("Should pass valid data through middleware with valid signature", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("response data"))
		})

		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString("request data"))
		w := httptest.NewRecorder()

		hash := hmac.New(sha256.New, []byte(signingKey))
		hash.Write([]byte("request data"))
		signature := hex.EncodeToString(hash.Sum(nil))
		r.Header.Set(HeaderKeyHash, signature)

		middleware(handler).ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")
		assert.Equal(t, "response data", w.Body.String(), "Should return correct response data")

		responseHash := w.Header().Get(HeaderKeyHash)
		require.NotEmpty(t, responseHash, "Response should have a hash header")

		expectedHash, err := utils.SignData([]byte("response data"), signingKey)
		require.NoError(t, err, "Should not error on generating expected hash")

		decodedResponseHash, err := hex.DecodeString(responseHash)
		require.NoError(t, err, "Should decode response hash without error")
		assert.True(t, hmac.Equal(decodedResponseHash, expectedHash), "Hashes should match")
	})

	t.Run("Should return error on unauthenticated data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("response data"))
		})

		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString("request data"))
		w := httptest.NewRecorder()

		r.Header.Set(HeaderKeyHash, "invalid-signature")

		middleware(handler).ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Should return 500 Internal Server Error on invalid signature")
	})

	t.Run("Should handle no header provided", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("response data"))
		})

		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString("request data"))
		w := httptest.NewRecorder()

		middleware(handler).ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 when no header is provided due to error in CI tests")
	})
}
