package delete_comix_test

import (
	"bytes"
	"encoding/json"
	"jadesheart/comix_back/internal/http-server/handlers/comix/delete_comix"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert" // Импорт библиотеки testify для написания тестов
)

// Определение структуры-заглушки для интерфейса ComixDeleter
type mockComixDeleter struct{}

func (m *mockComixDeleter) DeleteComixFromTagTable(tag string, name string) error {
	return nil
}

func (m *mockComixDeleter) DeleteComixFromAllComixTable(name string) error {
	return nil
}

func (m *mockComixDeleter) CheckPass(inputPass string) (bool, error) {
	if inputPass == "password" {
		return true, nil
	}
	return false, nil
}

type MockResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Тест успешного выполнения обработчика
func TestDeleteComixHandler_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	// Создание обработчика, передача логгера и объекта-заглушки
	handler := delete_comix.New(mockLogger, &mockComixDeleter{})

	// Создание тела запроса в формате JSON
	requestBody := map[string]interface{}{
		"password": "password",
		"tagName":  "exampleTag",
		"name":     "exampleName",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var responseBody MockResponse

	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Errorf("Ошибка при распоковке JSON: %s", err)
		return
	}

	assert.Equal(t, http.StatusOK, responseBody.Status)

}

// Тест обработки недопустимого запроса
func TestDeleteComixHandler_InvalidRequest(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	handler := delete_comix.New(mockLogger, &mockComixDeleter{})

	// Создание недопустимого тела запроса (отсутствует обязательное поле "password")
	invalidRequestBody := map[string]interface{}{
		"tagName": "exampleTag",
		"name":    "exampleName",
	}

	jsonBody, _ := json.Marshal(invalidRequestBody)
	req, err := http.NewRequest("POST", "localhost:8080/delete", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err) // Проверка отсутствия ошибок при создании запроса

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var responseBody MockResponse

	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Errorf("Ошибка при распоковке JSON: %s", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, responseBody.Status) // Проверка кода статуса ответа
}

// Тест обработки запроса с неверным паролем
func TestDeleteComixHandler_FailedPasswordVerification(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := delete_comix.New(mockLogger, &mockComixDeleter{})

	// Создание запроса с неверным паролем
	requestBody := map[string]interface{}{
		"password": "wrong_password",
		"tagName":  "exampleTag",
		"name":     "exampleName",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err) // Проверка отсутствия ошибок при создании запроса

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var responseBody MockResponse

	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Errorf("Ошибка в распоковке JSON: %s", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, responseBody.Status) // Проверка кода статуса ответа
}

// Дополнительные тесты могут быть добавлены для покрытия других сценариев ошибок в вашем обработчике
