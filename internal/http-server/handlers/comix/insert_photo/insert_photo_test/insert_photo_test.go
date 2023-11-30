package insert_photo_test

import (
	"bytes"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogdiscard"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"jadesheart/comix_back/internal/http-server/handlers/comix/insert_photo"
)

type ComixSaverMock struct {
}

func (c *ComixSaverMock) CheckPass(inputPass string) (bool, error) {
	return false, nil
}

func TestInsertPhotoHandler(t *testing.T) {
	mockPasswordVerifier := &ComixSaverMock{}
	mockLogger := slogdiscard.NewDiscardLogger()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("password", "test_password")
	_ = writer.WriteField("tag", "tag_name")
	_ = writer.WriteField("name", "comix_name")
	part, _ := writer.CreateFormFile("photo", "test_image.jpg")
	part.Write([]byte("file content"))
	writer.Close()
	req, err := http.NewRequest("POST", "/insertphoto", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	handler := insert_photo.New(mockLogger, mockPasswordVerifier)
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code, "статус должен быть 200")
	expectedContentType := "application/json"
	assert.Equal(t, expectedContentType, recorder.Header().Get("Content-Type"), "")
}
