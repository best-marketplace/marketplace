package product

import (
	"catalog/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Event struct {
	Ids       []string  `json:"ids"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type EventProducer interface {
	Send(ctx context.Context, event any) error
}

type RepoProductListViewer interface {
	ViewListProducts(context.Context, int, int) ([]*models.ProductListView, []string, error)
}

type Useacase struct {
	repoProductListViewer RepoProductListViewer
	eventProducer         EventProducer
	log                   *slog.Logger
}

func NewUseacase(log *slog.Logger, repoProductListViewer RepoProductListViewer, eventProducer EventProducer) *Useacase {
	return &Useacase{
		repoProductListViewer: repoProductListViewer,
		eventProducer:         eventProducer,
		log:                   log,
	}
}

func (u *Useacase) ViewListProducts(ctx context.Context, offset int, limit int) ([]*models.ProductListView, error) {
	const op = "product.usecase.ViewListProducts"

	products, ids, err := u.repoProductListViewer.ViewListProducts(ctx, offset, limit)
	if err != nil {
		u.log.Error(op+": failed to get products", slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	event := Event{
		Ids:       ids,
		Action:    "visibility",
		Timestamp: time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, event); err != nil {
			u.log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
		}
	}()

	u.log.Info(op+": successfully retrieved products",
		slog.Int("count", len(products)),
		slog.Int("offset", offset),
		slog.Int("limit", limit),
	)

	return products, nil
}
