package get_number_of_comics

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "jadesheart/comix_back/internal/lib/api/response"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type Response struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
	Number int    `json:"NumberOfComix,omitempty"`
}

type NumberComixGetter interface {
	GetComixQuantity() (int, error)
}

func New(log *slog.Logger, numberComixGetter NumberComixGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get_number_of_comix.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		count, err := numberComixGetter.GetComixQuantity()
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
