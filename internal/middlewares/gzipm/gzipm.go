package gzipm

import (
	"github.com/daremove/go-metrics-service/internal/utils"
	"net/http"
	"strings"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodingHeader := r.Header.Get("Content-Encoding")
		isRequestCompressed := strings.Contains(contentEncodingHeader, "gzip")

		if isRequestCompressed {
			cr, err := utils.NewCompressReader(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr

			defer cr.Close()
		}

		h.ServeHTTP(w, r)
	})
}
