package gzipm

import (
	"net/http"
	"strings"

	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/utils"
	"go.uber.org/zap"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodingHeader := r.Header.Get("Content-Encoding")
		isRequestCompressed := strings.Contains(contentEncodingHeader, "gzip")

		if isRequestCompressed {
			cr, err := utils.NewCompressReader(r.Body)

			if err != nil {
				logger.Log.Error("error decompress data", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = cr

			defer cr.Close()
		}

		h.ServeHTTP(w, r)
	})
}
