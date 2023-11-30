package save

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
)

type Request struct {
	Password    string `json:"password" validate:"required"`
	TagName     string `json:"tagName" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}

type ComixSaver interface {
	CreateNewTag(tagName string) error
	TagExist(tagName string) (bool, error)
	CheckPass(inputPass string) (bool, error)
	AddComixTagToAllTags(tagName string) error
	AddTagDescription(tag string, description string) error
}

func New(log *slog.Logger, comixSaver ComixSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.comix.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

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

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("failed validate", sl.Err(err))

			render.JSON(w, r, resp.ValidateErrors(validateErr))
			return
		}

		res, err := comixSaver.CheckPass(req.Password)
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

		result, err := comixSaver.TagExist(req.TagName)
		if err != nil {

			log.Error("failed get table by tag", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add new tag"))

			return
		}

		if result {

			log.Info("tag table already exists", slog.Any("tagName", req.TagName))

			render.JSON(w, r, resp.Error("tag already exists"))

			return
		}

		err = comixSaver.CreateNewTag(req.TagName)
		if err != nil {

			log.Error("failed to add new tag", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add new tag"))

			return
		}

		err = comixSaver.AddComixTagToAllTags(req.TagName)
		if err != nil {

			log.Error("failed add new tag to all_tags table", sl.Err(err))

			render.JSON(w, r, resp.Error("failed add new tag to all_tags table"))

			return
		}

		path := fmt.Sprintf("internal/storage/web/photos/%s", req.TagName)
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Error("failed create tag folder", sl.Err(err))

			render.JSON(w, r, resp.Error("failed create tag folder"))

			return
		}

		err = comixSaver.AddTagDescription(req.TagName, req.Description)
		if err != nil {
			log.Error("failed add tag description", sl.Err(err))

			render.JSON(w, r, resp.Error("failed add tag description"))

			return
		}

		log.Info("tag added and folder successful created", slog.Any(path+"TagName", req.TagName))

		responseOK(w, r)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: 200,
	})
}
