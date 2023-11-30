package get_number_of_comix_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_number_of_comics"
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

func (n *NumberComixGetter) GetComixQuantity() (int, error) {
	return 0, nil
}

func TestGetComix_Success(t *testing.T) {
	mockLogger := slogdiscard.NewDiscardLogger()

	handler := get_number_of_comics.New(mockLogger, &NumberComixGetter{})

	req, err := http.NewRequest("GET", "/getmainpagecomix", nil)
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
