package insert_photo

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io/ioutil"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	Password  string                  `form:"password" validate:"required"`
	TagName   string                  `form:"tag" validate:"required"`
	ComixName string                  `form:"name" validate:"required"`
	Photo     []*multipart.FileHeader `form:"photo" validate:"required"`
}

type Response struct {
	resp.Response
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type PasswordVerifier interface {
	CheckPass(inputPass string) (bool, error)
}

func New(log *slog.Logger, passwordVerifier PasswordVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.comix.insert_photo.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		req.TagName = strings.TrimSpace(r.FormValue("tag"))
		req.Password = r.FormValue("password")
		req.ComixName = r.FormValue("name")
		req.Photo = r.MultipartForm.File["photo"]

		log.Info("request body decoded", slog.Any("tag", req.TagName))

		validate := ValidateComixImg(req)
		if validate.Status == resp.StatusError {
			log.Error("failed validate", sl.Err(errors.New(validate.Error)))

			render.JSON(w, r, resp.Error(validate.Error))

			return
		}

		log.Info("All data valid", slog.Any("name", req.ComixName))

		res, err := passwordVerifier.CheckPass(req.Password)
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

		files := req.Photo
		for i, file := range files {
			f, err := file.Open()
			if err != nil {
				log.Error("failed open file", sl.Err(err))

				render.JSON(w, r, resp.Error("failed open file"))

				return
			}
			defer f.Close()

			data, err := ioutil.ReadAll(f)
			if err != nil {
				log.Error("failed read file", sl.Err(err))

				render.JSON(w, r, resp.Error("failed read file"))

				return
			}

			path := fmt.Sprintf("internal/storage/web/photos/%s/%s/", req.TagName, req.ComixName)

			_, err = os.Stat(path)
			if os.IsNotExist(err) {

				log.Error("Directory does not exist", sl.Err(err))

				render.JSON(w, r, resp.Error("Directory does not exist, check if you created the tag?"))

				return

			} else if err != nil {

				log.Error("failed get directory: ", sl.Err(err))

				render.JSON(w, r, resp.Error("failed get directory"))

				return

			}

			err = ioutil.WriteFile(path+strconv.Itoa(i+1)+".jpg", data, 0666)
			if err != nil {
				log.Error("failed write file", sl.Err(err))

				render.JSON(w, r, resp.Error("failed write file"))

				return
			}
		}

		responseOK(w, r)

	}
}

func ValidateComixImg(req Request) resp.Response {
	var errMsg []string

	tag := req.TagName
	name := req.ComixName
	password := req.Password
	re := regexp.MustCompile("^[a-zA-Z]+$")

	if reflect.TypeOf(password).Kind() != reflect.String {
		errMsg = append(errMsg, fmt.Sprintf("password is not a valued"))
	}

	if reflect.TypeOf(name).Kind() != reflect.String {
		errMsg = append(errMsg, fmt.Sprintf("name is not a string"))
	}

	if !re.MatchString(tag) {
		errMsg = append(errMsg, fmt.Sprintf("tag name is not a string"))
	}

	files := req.Photo
	for _, file := range files {

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			continue
		} else {
			errMsg = append(errMsg, fmt.Sprintf("file is not a image or not supported: %s", file.Filename))
		}
	}

	if len(errMsg) == 0 {
		return resp.Response{
			Status: resp.StatusOK,
		}
	} else {
		return resp.Response{
			Status: resp.StatusError,
			Error:  strings.Join(errMsg, ", "),
		}
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Status:   "successful insert comix photo",
	})
}
