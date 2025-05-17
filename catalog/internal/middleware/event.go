package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Event struct {
	URL       string    `json:"url"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type EventProducer interface {
	Send(ctx context.Context, event any) error
}

func LogEventMiddleware(log *slog.Logger, producer EventProducer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			stringPath := urlPath[5 : len(urlPath)-1]
			event := Event{
				URL:       stringPath,
				Action:    "open",
				Timestamp: time.Now(),
			}
			log.Info(fmt.Sprintf("LogEventMiddleware: urlpath: %s", urlPath))

			go func() {
				ctx := context.Background()
				if err := producer.Send(ctx, event); err != nil {
					log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
