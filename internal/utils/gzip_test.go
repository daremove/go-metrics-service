package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressWriter(t *testing.T) {
	t.Run("Should write and compress data", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		compressWriter := NewCompressWriter(recorder)

		data := []byte("test data")
		_, err := compressWriter.Write(data)
		assert.NoError(t, err)

		err = compressWriter.Close()
		assert.NoError(t, err)

		assert.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))

		reader, err := gzip.NewReader(recorder.Body)
		assert.NoError(t, err)

		decompressedData, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, data, decompressedData)
	})

	t.Run("Should set headers correctly", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		compressWriter := NewCompressWriter(recorder)

		headers := compressWriter.Header()
		headers.Set("Content-Type", "application/json")

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("Should write status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		compressWriter := NewCompressWriter(recorder)

		compressWriter.WriteHeader(http.StatusNotFound)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}

func TestCompressReader(t *testing.T) {
	t.Run("Should read and decompress data", func(t *testing.T) {
		data := []byte("test data")
		var compressedData bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedData)

		_, err := gzipWriter.Write(data)
		assert.NoError(t, err)

		err = gzipWriter.Close()
		assert.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(compressedData.Bytes()))
		compressReader, err := NewCompressReader(reader)
		assert.NoError(t, err)

		decompressedData, err := io.ReadAll(compressReader)
		assert.NoError(t, err)
		assert.Equal(t, data, decompressedData)

		err = compressReader.Close()
		assert.NoError(t, err)
	})

	t.Run("Should return error for invalid gzip data", func(t *testing.T) {
		invalidData := []byte("invalid gzip data")
		reader := io.NopCloser(bytes.NewReader(invalidData))

		_, err := NewCompressReader(reader)
		assert.Error(t, err)
	})
}
