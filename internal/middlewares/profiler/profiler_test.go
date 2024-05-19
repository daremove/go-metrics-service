package profiler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfiler(t *testing.T) {
	handler := Profiler()

	tests := []struct {
		path         string
		expectedCode int
	}{
		{"/", http.StatusMovedPermanently},
		{"/pprof", http.StatusMovedPermanently},
		{"/pprof/", http.StatusOK},
		{"/pprof/cmdline", http.StatusOK},
		// профилирование занимает какое-то время
		//{"/pprof/profile", http.StatusOK},
		{"/pprof/symbol", http.StatusOK},
		{"/pprof/trace", http.StatusOK},
		{"/pprof/goroutine", http.StatusOK},
		{"/pprof/threadcreate", http.StatusOK},
		{"/pprof/mutex", http.StatusOK},
		{"/pprof/heap", http.StatusOK},
		{"/pprof/block", http.StatusOK},
		{"/pprof/allocs", http.StatusOK},
		{"/vars", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}
