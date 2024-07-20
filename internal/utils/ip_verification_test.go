package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsTrustedIP(t *testing.T) {
	tests := []struct {
		name           string
		IP             string
		trustedSubnet  string
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "Should return true for IP in trusted subnet",
			IP:             "192.168.1.10",
			trustedSubnet:  "192.168.1.0/24",
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Should return false for IP not in trusted subnet",
			IP:             "192.168.2.10",
			trustedSubnet:  "192.168.1.0/24",
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Should return error for invalid IP",
			IP:             "invalid-ip",
			trustedSubnet:  "192.168.1.0/24",
			expectedResult: false,
			expectError:    true,
		},
		{
			name:           "Should return error for invalid subnet",
			IP:             "192.168.1.10",
			trustedSubnet:  "invalid-subnet",
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := isTrustedIP(tt.IP, tt.trustedSubnet)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestVerifyIPMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		headerIP      string
		expectedCode  int
	}{
		{
			name:          "Should allow request for trusted IP",
			trustedSubnet: "192.168.1.0/24",
			headerIP:      "192.168.1.10",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Should deny request for untrusted IP",
			trustedSubnet: "192.168.1.0/24",
			headerIP:      "192.168.2.10",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "Should deny request for missing IP header",
			trustedSubnet: "192.168.1.0/24",
			headerIP:      "",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "Should allow request if no trusted subnet is provided",
			trustedSubnet: "",
			headerIP:      "192.168.1.10",
			expectedCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)

			if tt.headerIP != "" {
				req.Header.Set("X-Real-IP", tt.headerIP)
			}

			rr := httptest.NewRecorder()

			handler := VerifyIPMiddleware(tt.trustedSubnet)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
