package edit_comix_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/edit_comix"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseMock struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type MockComixEditor struct{}

func (m *MockComixEditor) EditComixFromAllComixTable(name string, param string, newValue string) error {
	return nil
}
func (m *MockComixEditor) EditComixFromTagTable(tag string, name string, param string, newValue string) error {
	return nil
}
func (m *MockComixEditor) EditComixTag(tag string, name string, newValue string) error {
	return nil
}

func (m *MockComixEditor) CheckPass(inputPass string) (bool, error) {
	if inputPass == "password" {
		return true, nil
	}
	return false, nil
}

func TestEdit_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	handler := edit_comix.New(mockLogger, &MockComixEditor{})

	requestBody := map[string]interface{}{
		"password": "password",
		"tagName":  "exampleTag",
		"name":     "exampleName",
		"param":    "exampleParam",
		"newValue": "exampleNewValue",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/editcomix", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var responseBody ResponseMock

	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Errorf("Ошибка при распоковке JSON: %s", err)
		return
	}

	assert.Equal(t, http.StatusOK, responseBody.Status)
}

func TestEdit_IncorrectPassword(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	handler := edit_comix.New(mockLogger, &MockComixEditor{})

	requestBody := map[string]interface{}{
		"password": "wrongPass",
		"tagName":  "exampleTag",
		"name":     "exampleName",
		"param":    "exampleParam",
		"newValue": "exampleNewValue",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/editcomix", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var responseBody ResponseMock

	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Errorf("Ошибка при распоковке JSON: %s", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, responseBody.Status)
}

func TestEdit_EmptyValueParams(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{
			"tagName": "exampleTag",
			"name":    "exampleName",
			"param":   "exampleParam",
		}, {
			"password": "password",
			"tagName":  "exampleTag",
			"param":    "exampleParam",
			"newValue": "exampleNewValue",
		}, {
			"password": "password",
			"name":     "exampleName",
			"param":    "exampleParam",
			"newValue": "exampleNewValue",
		},
	}

	for _, m := range requestsBody {
		handler := edit_comix.New(mockLogger, &MockComixEditor{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/editcomix", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var responseBody ResponseMock

		if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
			t.Errorf("Ошибка при распоковке JSON: %s", err)
			return
		}

		assert.Equal(t, http.StatusBadRequest, responseBody.Status)
	}

}
