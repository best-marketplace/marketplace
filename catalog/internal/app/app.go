package app

import (
	"catalog/internal/comment"
	"catalog/internal/config"
	"catalog/internal/database/postgresql"
	"catalog/internal/kafka"
	m "catalog/internal/middleware"
	"catalog/internal/product"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	httpServer *http.Server
	log        *slog.Logger
	storage    *postgresql.Storage
	broker     *kafka.KafkaProducer
}

func NewApp(log *slog.Logger, cfg *config.Config) *App {
	storage, err := postgresql.ConnectAndNew(log, &cfg.Database)
	if err != nil {
		log.Error("Failed to create DB:")
		os.Exit(1)
	}
	brokers := []string{"kafka:9092"}

	producer := kafka.NewKafkaProducer(brokers)

	productRepo := product.NewRepository(storage.DB)
	productViewUsecase := product.NewUseacase(log, productRepo, producer)
	productAddUsecase := product.NewAddUseacase(log, productRepo, producer)
	commentRepo := comment.NewRepository(storage.DB)
	commentCreateUsecase := comment.NewCreateUseacase(log, commentRepo, producer)
	commentViewUsecase := comment.NewViewUseacase(log, commentRepo, producer)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Use(m.LogEventMiddleware(log, producer))
		r.Get("/products/", product.ViewListProducts(log, productViewUsecase))
		r.Get("/product/", product.ViewProduct(log, productViewUsecase))
		r.Post("/product/", product.AddProduct(log, productAddUsecase))
		r.Post("/comment/", comment.CreateComment(log, commentCreateUsecase))
		r.Get("/comments/", comment.ViewCommentInProduct(log, commentViewUsecase))
	})
	// router.Post("/api/auth", registeruser.New(log, authUsecase))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &App{
		httpServer: srv,
		log:        log,
		storage:    storage,
		broker:     producer,
	}
}

func (a *App) Run() error {
	a.log.Info("Starting server ", slog.String("port", a.httpServer.Addr))
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("Shutting down server...")
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	a.storage.Stop()
	a.log.Info("Database connection closed.")

	a.broker.Close()
	a.log.Info("Broker connection closed.")

	return nil
}
