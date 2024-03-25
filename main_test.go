package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenBodyNotEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code
	status_200 := http.StatusOK

	body := responseRecorder.Body.String()

	require.Equal(t, status_200, statusCode)
	require.NotEmpty(t, body)
}

// Город, который передаётся в параметре city, не поддерживается.
// Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenWrongCityValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=WrongCity", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code
	status_400 := http.StatusBadRequest

	body := responseRecorder.Body.String()
	err := "wrong city value"

	require.Equal(t, status_400, statusCode)
	require.Equal(t, err, body)
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := len(cafeList["moscow"])
	count := strconv.Itoa(totalCount + 1)

	req := httptest.NewRequest("GET", "/cafe?count="+count+"&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	statusCode := responseRecorder.Code
	status_200 := http.StatusOK

	splitBody := strings.Split(responseRecorder.Body.String(), ",")

	require.Equal(t, status_200, statusCode)
	assert.Len(t, splitBody, totalCount)
}
