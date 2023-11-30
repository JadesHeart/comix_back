package insert

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"time"
)

type Request struct {
	Password    string `json:"password" validate:"required"`
	TagName     string `json:"tagName" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ComixAdder interface {
	AddComixByTagName(tagName string, name string, description string, currentDate string) error
	AddComixToAllComixTable(tagName string, name string, description string, currentDate string) error
	CheckComixExists(tagName string, name string) (bool, error)
	TagExist(tagName string) (bool, error)
	CheckPass(inputPass string) (bool, error)
}

func New(log *slog.Logger, comixAdder ComixAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.comix.insert.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("cannot decode insert json body"))

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

		res, err := comixAdder.CheckPass(req.Password)
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

		result, err := comixAdder.TagExist(req.TagName)
		if err != nil {
			log.Error("failed get table by tag", sl.Err(err))

			render.JSON(w, r, resp.Error("failed get table by tag, pleas check table exists"))

			return
		}
		if !result {
			log.Info("tag table not exists", slog.Any("tagName", req.TagName))

			render.JSON(w, r, resp.Error("failed get table by tag, pleas check table exists"))

			return
		}

		result, err = comixAdder.CheckComixExists(req.TagName, req.Name)
		if err != nil {

			log.Error("Cannot check exists comix in table", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot check exists comix in table"))

			return
		}

		if result {
			log.Error("Comix already exists")

			render.JSON(w, r, resp.Error("Comix already exists"))

			return
		}
		currentDate := time.Now().Format("2006-01-02")
		err = comixAdder.AddComixByTagName(req.TagName, req.Name, req.Description, currentDate)
		if err != nil {
			log.Error("can not added comix", sl.Err(err))

			render.JSON(w, r, resp.Error("Comix already exists"))

			return
		}

		err = comixAdder.AddComixToAllComixTable(req.TagName, req.Name, req.Description, currentDate)
		if err != nil {
			log.Error("can not added comix", sl.Err(err))

			render.JSON(w, r, resp.Error("Comix already exists"))

			return
		}

		path := fmt.Sprintf("internal/storage/web/photos/%s/%s", req.TagName, req.Name)
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Error("failed create tag-name folder,check if you created the tag?", sl.Err(err))

			render.JSON(w, r, resp.Error("failed create tag-name folder,check if you created the tag?"))

			return
		}

		responseOK(w, r)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: 200,
	})
}
