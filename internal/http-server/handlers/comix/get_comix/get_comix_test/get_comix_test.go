package get_comix_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_comix"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"jadesheart/comix_back/internal/storage/postgres"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseMock struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type MockComixGetter struct{}

func (m *MockComixGetter) GetComixByName(tagName string, name string) (postgres.Comix, error) {
	return postgres.Comix{}, nil
}
func (m *MockComixGetter) CheckComixExists(tagName string, name string) (bool, error) {
	if name == "comixIsNotExist" {
		return true, nil
	} else if name == "comixIsNotExist" {
		return false, nil
	}
	return false, nil
}
func (m *MockComixGetter) TagExist(tagName string) (bool, error) {
	if tagName == "tagExist" {
		return true, nil
	} else if tagName == "tagIsNotExist" {
		return false, nil
	}
	return false, nil
}

func TestGetComix_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := get_comix.New(mockLogger, &MockComixGetter{})

	requestBody := map[string]interface{}{
		"tagName": "tagExist",
		"name":    "comixIsNotExist",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/getcomix", bytes.NewBuffer(jsonBody))
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

func TestGetComix_NotExistParam(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{
			"tagName": "tagExist",
			"name":    "comixIsNotExist",
		},
		{
			"tagName": "tagIsNotExist",
			"name":    "comixIsNotExist",
		},
	}

	for _, m := range requestsBody {
		handler := get_comix.New(mockLogger, &MockComixGetter{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/getcomix", bytes.NewBuffer(jsonBody))
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

func TestGetComix_EmptyValueParams(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{
			"name": "comixIsNotExist",
		},
		{
			"tagName": "tagExist",
		},
	}

	for _, m := range requestsBody {
		handler := get_comix.New(mockLogger, &MockComixGetter{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/getcomix", bytes.NewBuffer(jsonBody))
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
