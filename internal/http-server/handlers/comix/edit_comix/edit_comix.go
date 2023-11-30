package edit_comix

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"io"
	"io/ioutil"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

type Request struct {
	Password string `json:"password" validator:"required"`
	TagName  string `json:"tagName" validator:"required"`
	Name     string `json:"name" validator:"required"`
	Param    string `json:"param" validator:"required"`
	NewValue string `json:"newValue" validator:"required"`
}
type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ComixEditor interface {
	EditComixFromAllComixTable(name string, param string, newValue string) error
	EditComixFromTagTable(tag string, name string, param string, newValue string) error
	EditComixTag(tag string, name string, newValue string) error
	CheckPass(inputPass string) (bool, error)
}

func New(log *slog.Logger, comixEditor ComixEditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.edit_comix.New"

		log.With(
			slog.String("op", op),
			slog.String("request", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed decode request json", sl.Err(err))

			render.JSON(w, r, "failed decode request json")

			return
		}

		reqType := reflect.TypeOf(req)
		for i := 0; i < reqType.NumField(); i++ {
			field := reqType.Field(i)
			fieldValue := reflect.ValueOf(req).FieldByName(field.Name)
			if fieldValue.IsZero() {
				errorMsg := fmt.Sprintf("zero point value: %s", fieldValue)
				render.JSON(w, r, resp.Error(errorMsg))
				return
			}
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed validate", sl.Err(err))

			render.JSON(w, r, resp.ValidateErrors(validateErr))

			return
		}

		res, err := comixEditor.CheckPass(req.Password)
		if err != nil {
			log.Error("failed password verified", sl.Err(err))

			render.JSON(w, r, resp.Error("failed password verified"))

			return
		}

		if !res {
			log.Info("INCORRECT PASSWORD", slog.Any("Trying input password: ", req.Password))

			render.JSON(w, r, resp.Error("incorrect password"))

			return
		}

		err = comixEditor.EditComixFromAllComixTable(req.Name, req.Param, req.NewValue)
		if err != nil {
			log.Error("failed edit comix", sl.Err(err))

			render.JSON(w, r, resp.Error("failed edit comix"))

			return
		}
		if req.Param == "comix_tag" {
			err = comixEditor.EditComixTag(req.TagName, req.Name, req.NewValue)
			if err != nil {
				log.Error("failed edit comix tag", sl.Err(err))

				render.JSON(w, r, resp.Error("failed edit comix tag"))

				return
			}
			err = editComixDir(req.TagName, req.Name[:len(req.Name)-1], req.NewValue)
			if err != nil {
				log.Error("failed edit comix photo dir", sl.Err(err))

				render.JSON(w, r, resp.Error("failed edit comix photo dir"))

				return
			}

		} else {
			err = comixEditor.EditComixFromTagTable(req.TagName, req.Name, req.Param, req.NewValue)
			if err != nil {
				log.Error("failed edit comix", sl.Err(err))

				render.JSON(w, r, resp.Error("failed edit comix"))

				return
			}
		}
		responseOK(w, r)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: resp.StatusOK,
	})
}

func editComixDir(tag string, name string, newTag string) error {
	sourceDir := fmt.Sprintf("internal/storage/web/photos/%s/%s", tag, name)
	destDir := fmt.Sprintf("internal/storage/web/photos/%s/%s", newTag, name)

	// Создание новой директории
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// Получение файлов из исходной директории
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	// Копирование файлов в новую директорию
	for _, file := range files {
		sourceFile := filepath.Join(sourceDir, file.Name())
		destFile := filepath.Join(destDir, file.Name())

		err = copyFile(sourceFile, destFile)
		if err != nil {
			return err
		}
	}

	// Удаление старой директории
	err = os.RemoveAll(sourceDir)
	if err != nil {
		return err
	}
	return nil
}

// Функция для копирования файла
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
