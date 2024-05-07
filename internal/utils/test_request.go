// Package utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRequest выполняет HTTP-запрос к тестовому серверу и проверяет наличие ошибок в процессе выполнения.
// Эта функция упрощает написание тестов за счет автоматизации создания запросов и обработки ответов.
// Принимает следующие параметры:
// - t: указатель на объект тестирования *testing.T для регистрации результатов тестирования.
// - ts: тестовый сервер, к которому будет выполнен запрос.
// - method: HTTP метод запроса (например, "GET", "POST").
// - path: путь URL, к которому добавляется базовый URL сервера.
// - headers: словарь заголовков, которые нужно добавить в HTTP-запрос.
// - body: тело запроса в виде io.Reader (может быть nil, если тело не требуется).
// Функция возвращает объект ответа *http.Response и строку с телом ответа.
func TestRequest(t *testing.T, ts *httptest.Server, method, path string, headers map[string]string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	req.Header.Set("Accept-Encoding", "identity")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
