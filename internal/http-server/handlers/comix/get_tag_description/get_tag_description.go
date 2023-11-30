package get_tag_description

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
	TagName string `json:"tagName" validator:"required"`
}
type Response struct {
	Status      int    `json:"status,omitempty"`
	Error       string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

type TagDescriptionGetter interface {
	GetTagDescription(tagName string) (string, error)
}

func New(log *slog.Logger, tagDescriptionGetter TagDescriptionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.get_tag_description.New"

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

		description, err := tagDescriptionGetter.GetTagDescription(req.TagName)
		if err != nil {
			log.Error("INCORRECT PASSWORD", slog.Any("Failed get tag description ", sl.Err(err)))

			render.JSON(w, r, resp.Error("Failed get tag description"))

			return
		}

		responseOK(w, r, description)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, description string) {
	render.JSON(w, r, Response{
		Status:      200,
		Description: description,
	})
}
