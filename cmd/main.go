package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"jadesheart/comix_back/internal/config"
	"jadesheart/comix_back/internal/http-server/handlers/comix/delete_comix"
	"jadesheart/comix_back/internal/http-server/handlers/comix/edit_comix"
	"jadesheart/comix_back/internal/http-server/handlers/comix/find_comix"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_all_tag_comix"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_all_tags"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_comix"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_comix_for_main_page"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_comix_photo"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_number_of_comics"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_number_of_comix_form_name"
	get_number_of_comics_from_tag "jadesheart/comix_back/internal/http-server/handlers/comix/get_number_of_comix_form_tag"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_photo"
	"jadesheart/comix_back/internal/http-server/handlers/comix/get_tag_description"
	"jadesheart/comix_back/internal/http-server/handlers/comix/insert"
	"jadesheart/comix_back/internal/http-server/handlers/comix/insert_photo"
	"jadesheart/comix_back/internal/http-server/handlers/comix/save"
	mnLogger "jadesheart/comix_back/internal/http-server/middleware/logger"
	"jadesheart/comix_back/internal/lib/logger/handlers/slogpretty"
	"jadesheart/comix_back/internal/lib/logger/sl"
	"jadesheart/comix_back/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("Successful init database")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mnLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	// Добавляем обработчик CORS в цепочку middleware
	router.Use(corsHandler.Handler)
	router.Post("/newtag", save.New(logger, storage))
	router.Post("/newcomix", insert.New(logger, storage))
	router.Post("/insertphoto", insert_photo.New(logger, storage))
	router.Post("/getcomix", get_comix.New(logger, storage))
	router.Post("/gettagdescription", get_tag_description.New(logger, storage))
	router.Post("/getmainpagecomix", get_comix_for_main_page.New(logger, storage))
	router.Post("/getalltagcomix", get_all_tag_comix.New(logger, storage))
	router.Post("/alltags", get_all_tags.New(logger, storage))
	router.Post("/deletecomix", delete_comix.New(logger, storage))
	router.Post("/findcomix", find_comix.New(logger, storage))
	router.Post("/editcomix", edit_comix.New(logger, storage))
	router.Post("/getquantitycomix", get_number_of_comics.New(logger, storage))
	router.Post("/getquantitytag", get_number_of_comics_from_tag.New(logger, storage))
	router.Post("/getquantityname", get_number_of_comix_form_name.New(logger, storage))
	router.Get("/{folder1}/{folder2}/{fileName}", get_photo.New(logger))
	router.Get("/comix/{tag}/{name}/", get_comix_photo.New(logger, storage))

	logger.Info("starting server", slog.String("addres", cfg.Address))
	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case "local":
		logger = setupPrettySlog()
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
