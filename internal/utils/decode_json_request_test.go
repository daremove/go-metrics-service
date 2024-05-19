package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSONRequest(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedResult interface{}
		expectedError  error
	}{
		{
			name:           "Should decode JSON with correct content type",
			contentType:    "application/json",
			body:           `{"name":"John","age":30}`,
			expectedResult: map[string]interface{}{"name": "John", "age": float64(30)},
			expectedError:  nil,
		},
		{
			name:           "Should return error with invalid content type",
			contentType:    "text/plain",
			body:           `{"name":"John","age":30}`,
			expectedResult: nil,
			expectedError:  errors.New(UnsupportedContentTypeCode),
		},
		{
			name:           "Should return error with invalid JSON format",
			contentType:    "application/json",
			body:           `{"name": "John", "age":}`,
			expectedResult: nil,
			expectedError:  &json.SyntaxError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			req.Header.Add("Content-Type", tc.contentType)

			result, err := DecodeJSONRequest[map[string]interface{}](req)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.IsType(t, tc.expectedError, err)
			}

			if err == nil {
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
