package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeJSONRequest(t *testing.T) {
	tests := []struct {
		name          string
		data          interface{}
		expectedBody  string
		expectedError error
	}{
		{
			name:          "Should encode valid data into JSON",
			data:          map[string]interface{}{"name": "John", "age": 30},
			expectedBody:  `{"age":30,"name":"John"}`,
			expectedError: nil,
		},
		{
			name:          "Should handle nil data gracefully",
			data:          nil,
			expectedBody:  `null`,
			expectedError: nil,
		},
		{
			name:          "Should return error with non-encodable data",
			data:          make(chan int),
			expectedBody:  "",
			expectedError: &json.UnsupportedTypeError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := EncodeJSONRequest(w, tc.data)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.IsType(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedBody, w.Body.String())
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}
