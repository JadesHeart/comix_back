package delete_comix

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"reflect"
)

type Request struct {
	Password string `json:"password" validator:"required"`
	TagName  string `json:"tagName" validator:"required"`
	Name     string `json:"name" validator:"required"`
}
type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ComixDeleter interface {
	DeleteComixFromTagTable(tag string, name string) error
	DeleteComixFromAllComixTable(name string) error
	CheckPass(inputPass string) (bool, error)
}

func New(log *slog.Logger, comixDeleter ComixDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.delete_comix.New"

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

		res, err := comixDeleter.CheckPass(req.Password)
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

		err = comixDeleter.DeleteComixFromTagTable(req.TagName, req.Name)
		if err != nil {
			log.Error("Cannot delete comix from bd", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot delete comix from bd"))

			return
		}

		err = comixDeleter.DeleteComixFromAllComixTable(req.Name)
		if err != nil {
			log.Error("Cannot delete comix from bd", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot delete comix from bd"))

			return
		}

		responseOK(w, r)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: resp.StatusOK,
	})
}
