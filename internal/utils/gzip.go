// Package utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter обертывает http.ResponseWriter, добавляя функциональность сжатия gzip.
type compressWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

// NewCompressWriter создает новый экземпляр compressWriter.
func NewCompressWriter(writer http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:     writer,
		gzipWriter: gzip.NewWriter(writer),
	}
}

// Header возвращает заголовки HTTP ответа.
func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

// Write записывает данные в gzip writer, автоматически сжимая их перед отправкой клиенту.
func (c *compressWriter) Write(p []byte) (int, error) {
	c.writer.Header().Set("Content-Encoding", "gzip")

	return c.gzipWriter.Write(p)
}

// WriteHeader отправляет HTTP статус код.
func (c *compressWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

// Close закрывает gzip writer и освобождает все связанные с ним ресурсы.
func (c *compressWriter) Close() error {
	return c.gzipWriter.Close()
}

// compressReader обертывает io.ReadCloser, добавляя функциональность декомпрессии gzip.
type compressReader struct {
	reader     io.ReadCloser // Оригинальный reader
	gzipReader *gzip.Reader  // Reader gzip для декомпрессии данных
}

// NewCompressReader создает новый экземпляр compressReader.
func NewCompressReader(reader io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(reader)

	if err != nil {
		return nil, err
	}

	return &compressReader{
		reader,
		gzipReader,
	}, nil
}

// Read читает и декомпрессирует данные из gzip stream.
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzipReader.Read(p)
}

// Close закрывает gzip reader и освобождает все связанные с ним ресурсы.
func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzipReader.Close()
}
