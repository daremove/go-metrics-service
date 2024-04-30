// Пакет utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// UnsupportedContentTypeCode константа для обозначения ошибки несоответствия типа содержимого.
const (
	UnsupportedContentTypeCode = "UnsupportedContentTypeCode"
)

// ModelParameter тип-ограничение для обобщенной функции DecodeJSONRequest, позволяющее передавать любую модель.
type ModelParameter interface {
	interface{} | []interface{}
}

// DecodeJSONRequest декодирует JSON из HTTP-запроса в указанную модель.
// Возвращает декодированную модель или ошибку, если содержимое запроса не соответствует ожидаемому.
func DecodeJSONRequest[Model ModelParameter](r *http.Request) (Model, error) {
	var emptyResult Model

	if r.Header.Get("Content-Type") != "application/json" {
		return emptyResult, errors.New(UnsupportedContentTypeCode)
	}

	var parsedData Model
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		return emptyResult, err
	}

	if err = json.Unmarshal(buf.Bytes(), &parsedData); err != nil {
		return emptyResult, err
	}

	return parsedData, nil
}
