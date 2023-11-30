package get_comix_photo

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io/ioutil"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ViewsAdder interface {
	AddViews(tag string, name string) error
}

func New(log *slog.Logger, viewsAdder ViewsAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.comix.get_comix_photo.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tag := chi.URLParam(r, "tag")
		comixName := chi.URLParam(r, "name")

		filePath := "internal/storage/web/photos/" + tag + "/" + comixName

		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Error("failed read files", sl.Err(err))

			render.JSON(w, r, "failed read files")

			return
		}

		sortFiles(files)

		var images [][]byte

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") {

				filePath := fmt.Sprintf("%s/%s", filePath, file.Name())

				data, err := ioutil.ReadFile(filePath)

				if err != nil {
					log.Error("Unable to send file", sl.Err(err))

					render.JSON(w, r, "Unable to send file")

					return
				}

				images = append(images, data)
			}
		}

		response := map[string]interface{}{
			"images": images,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResponse)

		err = viewsAdder.AddViews(tag, comixName)
		if err != nil {
			log.Error("failed added view", sl.Err(err))

			render.JSON(w, r, "failed added view")

			return
		}
	}
}

func sortFiles(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		name1 := strings.TrimSuffix(files[i].Name(), filepath.Ext(files[i].Name()))
		name2 := strings.TrimSuffix(files[j].Name(), filepath.Ext(files[j].Name()))

		num1, err1 := strconv.Atoi(name1)
		num2, err2 := strconv.Atoi(name2)

		if err1 == nil && err2 == nil {
			return num1 < num2
		}

		return name1 < name2
	})
}
