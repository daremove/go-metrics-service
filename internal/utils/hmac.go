// Package utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

// SignData генерирует подпись для данных, используя HMAC с алгоритмом SHA-256.
// Эта функция принимает байтовый массив данных и строку секретного ключа для подписи.
// Возвращает подписанные данные или ошибку в случае неудачи.
func SignData(data []byte, signingKey string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(signingKey))

	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
