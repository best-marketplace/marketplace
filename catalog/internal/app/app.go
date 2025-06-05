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

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// type test struct {
// }

// func (t *test) Send(ctx context.Context, event any, topic string) error {
// 	return nil
// }

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
	// producer := &test{}

	esHost := cfg.Elasticsearch.Host
	if esHost == "" {
		esHost = "elasticsearch:9200"
	}
	log.Info("Configuring Elasticsearch client", slog.String("host", esHost))

	esConfig := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s", esHost)},
	}

	esClient, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		log.Error("Failed to create Elasticsearch client", "error", err)
	}

	productRepo := product.NewRepository(storage.DB, esClient)

	productViewUsecase := product.NewUseacase(log, productRepo, producer)
	productAddUsecase := product.NewAddUseacase(log, productRepo, producer)
	commentRepo := comment.NewRepository(storage.DB)
	commentCreateUsecase := comment.NewCreateUseacase(log, commentRepo, producer)
	commentViewUsecase := comment.NewViewUseacase(log, commentRepo, producer)

	productSearch := product.NewSearchUseacase(log, productRepo, producer)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(m.MetricsMiddleware)

	router.Use(m.LogEventMiddleware(log, producer))
	router.Handle("/metrics", promhttp.Handler())

	router.Get("/products/", product.ViewListProducts(log, productViewUsecase))
	router.Get("/product/", product.ViewProduct(log, productViewUsecase))
	router.Post("/product/", product.AddProduct(log, productAddUsecase))
	router.Post("/comment/", comment.CreateComment(log, commentCreateUsecase))
	router.Get("/comments/", comment.ViewCommentInProduct(log, commentViewUsecase))
	router.Get("/search/", product.SearchProduct(log, productSearch))

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.ServerPort),
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
