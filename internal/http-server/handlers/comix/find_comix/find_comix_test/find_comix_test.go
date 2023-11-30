package find_comix_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/find_comix"
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

func (m *MockComixGetter) FindComixFromAllComix(name string, pageToDisplay int) ([]postgres.ComixFromAllComix, error) {
	return nil, nil
}

func TestComixFind_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := find_comix.New(mockLogger, &MockComixGetter{})

	requestBody := map[string]interface{}{
		"name":       "exampleName",
		"pageNumber": 1,
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/findcomix", bytes.NewBuffer(jsonBody))
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

func TestComixFind_EmptyValueParams(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{
			"pageNumber": 1,
		},
		{
			"name": "exampleName",
		},
	}

	for _, m := range requestsBody {
		handler := find_comix.New(mockLogger, &MockComixGetter{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/findcomix", bytes.NewBuffer(jsonBody))
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
