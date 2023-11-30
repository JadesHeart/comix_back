package get_all_tags

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type Response struct {
	Status  int      `json:"status,omitempty"`
	Error   string   `json:"error,omitempty"`
	TagList []string `json:"tagList"`
}

type ComixGetter interface {
	GetAllTags() ([]string, error)
}

func New(log *slog.Logger, comixGetter ComixGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.comix.get_all_tags.New"

		log.With(
			slog.String("op", op),
			slog.String("request", middleware.GetReqID(r.Context())),
		)

		tagsList, err := comixGetter.GetAllTags()
		if err != nil {

			log.Error("Cannot get all tags", sl.Err(err))

			render.JSON(w, r, resp.Error("Cannot get all tags"))

			return
		}

		responseOK(w, r, tagsList)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, tagsList []string) {
	render.JSON(w, r, Response{
		Status:  200,
		TagList: tagsList,
	})
}
