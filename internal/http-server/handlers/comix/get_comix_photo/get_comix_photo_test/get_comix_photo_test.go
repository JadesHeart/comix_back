package get_comix_photo_test

import (
	"jadesheart/comix_back/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"jadesheart/comix_back/internal/http-server/handlers/comix/get_comix_photo"
)

type MockViewsAdder struct{}

func (m *MockViewsAdder) AddViews(tag string, name string) error {
	return nil
}

func TestGetComixPhotoHandler(t *testing.T) {
	// Создание фейкового ViewsAdder
	logger := setupLogger("local")

	viewsAdder := &MockViewsAdder{}

	// Создание фейкового запроса с URL параметрами
	reqURL := "/comix/tag_name/comix_name//" // URL с параметрами tag_name и comix_name
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Получение объекта Query параметров из URL
	reqURLWithParams, _ := url.Parse(reqURL)
	queryParams := reqURLWithParams.Query()

	// Добавление параметров tag и name в запрос
	queryParams.Set("tag", "tag_name")
	queryParams.Set("name", "comix_name")

	// Установка обновленных параметров в URL запроса
	reqURLWithParams.RawQuery = queryParams.Encode()

	// Установка обновленного URL обратно в запрос
	req.URL = reqURLWithParams

	// Создание фейкового ResponseWriter
	recorder := httptest.NewRecorder()

	// Создание хэндлера с передачей фейкового ViewsAdder и логгера
	handler := get_comix_photo.New(logger, viewsAdder)

	// Выполнение запроса
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа
	assert.Equal(t, http.StatusOK, recorder.Code, "status code should be 200")

	// Проверка ожидаемого формата JSON в ответе
	expectedContentType := "application/json"
	assert.Equal(t, expectedContentType, recorder.Header().Get("Content-Type"), "content type should be JSON")

	// Дополнительные проверки ожидаемых данных в ответе можно добавить с учетом предполагаемой логики
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
