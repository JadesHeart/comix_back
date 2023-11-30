package get_number_of_comix_form_name_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_number_of_comix_form_name"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseMock struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type NumberComixGetter struct{}

func (n *NumberComixGetter) GetComixQuantityFromName(name string) (int, error) {
	return 0, nil
}

func TestGetNumberOfComix_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := get_number_of_comix_form_name.New(mockLogger, &NumberComixGetter{})

	requestBody := map[string]interface{}{
		"name": "someName",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/getquantityname", bytes.NewBuffer(jsonBody))
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

func TestGetNumberOfComix_EmptyValueParams(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{},
	}

	for _, m := range requestsBody {
		handler := get_number_of_comix_form_name.New(mockLogger, &NumberComixGetter{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/getquantityname", bytes.NewBuffer(jsonBody))
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
