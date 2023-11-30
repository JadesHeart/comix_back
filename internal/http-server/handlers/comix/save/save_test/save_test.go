package save_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/save"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ResponseMock struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ComixSaverMock struct{}

func (c *ComixSaverMock) CreateNewTag(tagName string) error {
	return nil
}
func (c *ComixSaverMock) TagExist(tagName string) (bool, error) {
	if tagName == "tagExist" {
		return true, nil
	} else if tagName == "tagIsNotExist" {
		return false, nil
	}
	return false, nil
}
func (c *ComixSaverMock) CheckPass(inputPass string) (bool, error) {
	if inputPass == "password" {
		return true, nil
	}
	return false, nil
}
func (c *ComixSaverMock) AddComixTagToAllTags(tagName string) error {
	return nil
}
func (c *ComixSaverMock) AddTagDescription(tag string, description string) error {
	return nil
}

func TestGetTagDescription_Success(t *testing.T) {
	mockLogger := setupLogger("local")

	handler := save.New(mockLogger, &ComixSaverMock{})

	requestBody := map[string]interface{}{
		"tagName":     "someTag",
		"password":    "password",
		"description": "someText",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/newtag", bytes.NewBuffer(jsonBody))
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
		{
			"password":    "password",
			"description": "someText",
		}, {
			"tagName":     "someTag",
			"description": "someText",
		}, {
			"tagName":  "someTag",
			"password": "password",
		},
	}

	for _, m := range requestsBody {
		handler := save.New(mockLogger, &ComixSaverMock{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/newtag", bytes.NewBuffer(jsonBody))
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

func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case "local":
		logger = setupPrettySlog()
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
