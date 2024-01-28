package gzipm

import (
	"github.com/daremove/go-metrics-service/internal/utils"
	"net/http"
	"strings"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w

		acceptEncodingHeader := r.Header.Get("Accept-Encoding")
		contentTypeHeader := r.Header.Get("Content-Type")
		isGzipSupported := false

		if contentTypeHeader == "" {
			contentTypeHeader = "text/html"
		}

		if (contentTypeHeader == "application/json" || contentTypeHeader == "text/html") && strings.Contains(acceptEncodingHeader, "gzip") {
			isGzipSupported = true
		}

		if isGzipSupported {
			cw := utils.NewCompressWriter(w)
			originalWriter = cw

			defer cw.Close()
		}

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

		h.ServeHTTP(originalWriter, r)
	})
}
