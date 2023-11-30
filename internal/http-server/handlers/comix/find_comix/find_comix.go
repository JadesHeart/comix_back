package find_comix

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
	Name          string `json:"name" validator:"required"`
	PageToDisplay int    `json:"pageNumber" validator:"required"`
}
type Response struct {
	resp.Response
	Status            int                          `json:"status,omitempty"`
	Error             string                       `json:"error,omitempty"`
	ComixFromAllComix []postgres.ComixFromAllComix `json:"comixFromForMainPage"`
}

type ComixGetter interface {
	FindComixFromAllComix(name string, pageToDisplay int) ([]postgres.ComixFromAllComix, error)
}

func New(log *slog.Logger, comixGetter ComixGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.find_comix.New"

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

		comix, err := comixGetter.FindComixFromAllComix(req.Name, req.PageToDisplay)
		if err != nil {
			log.Error("Cannot get comix from bd", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot get comix from bd"))

			return
		}

		responseOK(w, r, comix)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, comix []postgres.ComixFromAllComix) {
	render.JSON(w, r, Response{
		Response:          resp.OK(),
		Status:            200,
		ComixFromAllComix: comix,
	})
}
