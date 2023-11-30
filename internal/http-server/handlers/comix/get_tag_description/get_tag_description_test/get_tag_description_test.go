package get_tag_description_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_tag_description"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseMock struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type TagDescriptionGetterMock struct{}

func (t *TagDescriptionGetterMock) GetTagDescription(tagName string) (string, error) {
	return "", nil
}

func TestGetTagDescription_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := get_tag_description.New(mockLogger, &TagDescriptionGetterMock{})

	requestBody := map[string]interface{}{
		"tagName": "someTag",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/gettagdescription", bytes.NewBuffer(jsonBody))
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

func TestGetTagDescription_EmptyValueParams(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger() // инициализируйте ваш mock логгер здесь

	requestsBody := []map[string]interface{}{
		{},
	}

	for _, m := range requestsBody {
		handler := get_tag_description.New(mockLogger, &TagDescriptionGetterMock{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/gettagdescription", bytes.NewBuffer(jsonBody))
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
