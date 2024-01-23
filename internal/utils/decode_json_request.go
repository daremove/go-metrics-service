package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	UnsupportedContentTypeCode = "UnsupportedContentTypeCode"
	DecoderErrorCode           = "DecoderErrorCode"
	ReadBufferErrorCode        = "ReadBufferErrorCode"
)

func DecodeJSONRequest[Model interface{}](r *http.Request) (Model, error) {
	var emptyResult Model

	if r.Header.Get("Content-Type") != "application/json" {
		return emptyResult, errors.New(UnsupportedContentTypeCode)
	}

	var parsedData Model
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		return emptyResult, errors.New(UnsupportedContentTypeCode)
	}

	if err = json.Unmarshal(buf.Bytes(), &parsedData); err != nil {
		return emptyResult, errors.New(DecoderErrorCode)
	}

	return parsedData, nil
}
