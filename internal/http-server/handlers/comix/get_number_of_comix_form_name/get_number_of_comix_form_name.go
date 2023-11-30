package get_number_of_comix_form_name

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
	Name string `json:"name" validator:"required"`
}
type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
	Number int    `json:"NumberOfComix,omitempty"`
}

type NumberComixGetter interface {
	GetComixQuantityFromName(name string) (int, error)
}

func New(log *slog.Logger, numberComixGetter NumberComixGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get_number_of_comix_form_name.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
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

		count, err := numberComixGetter.GetComixQuantityFromName(req.Name)
		if err != nil {
			log.Error("Failed get number of comix", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed get number of comix"))

			return
		}

		responseOK(w, r, (count/16)+1)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, count int) {
	render.JSON(w, r, Response{
		Status: 200,
		Number: count,
	})
}
