package utils

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

func NewCompressWriter(writer http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:     writer,
		gzipWriter: gzip.NewWriter(writer),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	c.writer.Header().Set("Content-Encoding", "gzip")

	return c.gzipWriter.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.gzipWriter.Close()
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

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

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzipReader.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzipReader.Close()
}
