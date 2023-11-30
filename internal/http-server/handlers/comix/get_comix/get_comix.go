package get_comix

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"jadesheart/comix_back/internal/storage/postgres"
	"log/slog"
	"net/http"
	"reflect"
)

type Request struct {
	Tag  string `json:"tagName" validator:"required"`
	Name string `json:"name" validator:"required"`
}
type Response struct {
	Status      int    `json:"status,omitempty"`
	Error       string `json:"error,omitempty"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UploadDate  string `json:"upload_date"`
	Views       int    `json:"views"`
}

type ComixGetter interface {
	GetComixByName(tagName string, name string) (postgres.Comix, error)
	CheckComixExists(tagName string, name string) (bool, error)
	TagExist(tagName string) (bool, error)
}

func New(log *slog.Logger, comixGetter ComixGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.get_comix.New"

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

		result, err := comixGetter.TagExist(req.Tag)
		if err != nil {
			log.Error("failed get table by tag", sl.Err(err))

			render.JSON(w, r, resp.Error("failed get table by tag, pleas check table exists"))

			return
		}
		if !result {
			log.Info("tag table not exists", slog.Any("tagName", req.Tag))

			render.JSON(w, r, resp.Error("failed get table by tag, pleas check table exists"))

			return
		}

		result, err = comixGetter.CheckComixExists(req.Tag, req.Name)
		if err != nil {

			log.Error("Cannot check exists comix in table", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot check exists comix in table"))

			return
		}

		if !result {
			log.Error("Comix not exists")

			render.JSON(w, r, resp.Error("Comix not exists"))

			return
		}

		comix, err := comixGetter.GetComixByName(req.Tag, req.Name)
		if err != nil {
			log.Error("Cannot get comix from bd", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot get comix from bd"))

			return
		}

		responseOK(w, r, req.Tag, req.Name, comix.Description, comix.UploadDate, comix.Views)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, tag string, name string, description string, uploadDate string, views int) {
	render.JSON(w, r, Response{
		Status:      200,
		Tag:         tag,
		Name:        name,
		Description: description,
		UploadDate:  uploadDate,
		Views:       views,
	})
}
