// Пакет utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"encoding/json"
	"net/http"
)

// EncodeJSONRequest кодирует переданную модель данных в JSON и записывает её в HTTP ответ.
// При возникновении ошибки в процессе кодирования или записи, функция возвращает ошибку.
func EncodeJSONRequest[Model interface{}](w http.ResponseWriter, data Model) error {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = w.Write(resp)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
