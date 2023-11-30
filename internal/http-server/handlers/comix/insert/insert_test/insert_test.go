package insert_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/insert"
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

type ComixAdderMock struct{}

func (c *ComixAdderMock) AddComixByTagName(tagName string, name string, description string, currentDate string) error {
	return nil
}
func (c *ComixAdderMock) AddComixToAllComixTable(tagName string, name string, description string, currentDate string) error {
	return nil
}
func (c *ComixAdderMock) CheckComixExists(tagName string, name string) (bool, error) {
	if name == "comixExist" {
		return true, nil
	} else if name == "comixIsNotExist" {
		return false, nil
	}
	return false, nil
}
func (c *ComixAdderMock) TagExist(tagName string) (bool, error) {
	if tagName == "tagExist" {
		return true, nil
	} else if tagName == "tagIsNotExist" {
		return false, nil
	}
	return false, nil
}
func (c *ComixAdderMock) CheckPass(inputPass string) (bool, error) {
	if inputPass == "password" {
		return true, nil
	}
	return false, nil
}

func TestGetTagDescription_Success(t *testing.T) {
	mockLogger := setupLogger("local")

	handler := insert.New(mockLogger, &ComixAdderMock{})

	requestBody := map[string]interface{}{
		"password":    "password",
		"tagName":     "tagExist",
		"name":        "comixIsNotExist",
		"description": "someText",
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/newcomix", bytes.NewBuffer(jsonBody))
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
			"tagName":     "tagExist",
			"name":        "comixIsNotExist",
			"description": "someText",
		}, {
			"password":    "password",
			"name":        "comixIsNotExist",
			"description": "someText",
		}, {
			"password":    "password",
			"tagName":     "tagExist",
			"description": "someText",
		}, {
			"password": "password",
			"tagName":  "tagExist",
			"name":     "comixIsNotExist",
		},
	}

	for _, m := range requestsBody {
		handler := insert.New(mockLogger, &ComixAdderMock{})

		jsonBody, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", "/newcomix", bytes.NewBuffer(jsonBody))
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
