package get_photo

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.comix.get_photo.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		folder1 := chi.URLParam(r, "folder1")
		folder2 := chi.URLParam(r, "folder2")
		fileName := chi.URLParam(r, "fileName")

		filePath := "internal/storage/web/photos/" + folder1 + "/" + folder2[:len(folder2)-1] + "/" + fileName + ".jpg"

		file, err := os.Open(filePath)
		if err != nil {
			log.Error("fail not found", sl.Err(err))

			render.JSON(w, r, "fail not found")

			return
		}
		defer file.Close()

		contentType := http.DetectContentType(nil)
		w.Header().Set("Content-Type", contentType)

		_, err = io.Copy(w, file)
		if err != nil {
			log.Error("Unable to send file", sl.Err(err))

			render.JSON(w, r, "Unable to send file")

			return
		}
	}

}
