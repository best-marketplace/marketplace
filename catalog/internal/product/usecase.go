package product

import (
	"catalog/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Event struct {
	URL       string    `json:"url,omitempty"`
	Ids       []string  `json:"ids,omitempty"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type EventProducer interface {
	Send(ctx context.Context, event any, topic string) error
}

type RepoProductView interface {
	ViewListProducts(context.Context, int, int) ([]*models.ProductListView, []string, error)
	ViewProduct(context.Context, string) (*models.ProductView, error)
}

type Useacase struct {
	repoProductView RepoProductView
	eventProducer   EventProducer
	log             *slog.Logger
}

func NewUseacase(log *slog.Logger, repoProductView RepoProductView, eventProducer EventProducer) *Useacase {
	return &Useacase{
		repoProductView: repoProductView,
		eventProducer:   eventProducer,
		log:             log,
	}
}

func (u *Useacase) ViewListProducts(ctx context.Context, offset int, limit int) ([]*models.ProductListView, error) {
	const op = "product.usecase.ViewListProducts"

	products, ids, err := u.repoProductView.ViewListProducts(ctx, offset, limit)
	if err != nil {
		u.log.Error(op+": failed to get products", slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	event := Event{
		URL:       "products",
		Ids:       ids,
		Action:    "visibility",
		Timestamp: time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, event, "user-events"); err != nil {
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

func (u *Useacase) ViewProduct(ctx context.Context, id string) (*models.ProductView, error) {
	const op = "product.usecase.ViewProduct"

	product, err := u.repoProductView.ViewProduct(ctx, id)
	if err != nil {
		u.log.Error(op+": failed to get product", slog.String("id", id), slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if product == nil {
		u.log.Warn(op+": product not found", slog.String("id", id))
		return nil, fmt.Errorf("%s: product not found", op)
	}

	u.log.Info(op+": successfully retrieved product", slog.String("id", id))

	// event := Event{
	// 	Ids:       []string{id},
	// 	Action:    "view",
	// 	Timestamp: time.Now(),
	// }

	// go func() {
	// 	bgCtx := context.Background()
	// 	if err := u.eventProducer.Send(bgCtx, event); err != nil {
	// 		u.log.Error(op+": failed to send event", slog.Any("err", err))
	// 	}
	// }()

	return product, nil
}

type AddUseacase struct {
	repoProductAdder RepoProductAdder
	eventProducer    EventProducer
	log              *slog.Logger
}

func NewAddUseacase(log *slog.Logger, repoProductAdder RepoProductAdder, eventProducer EventProducer) *AddUseacase {
	return &AddUseacase{
		repoProductAdder: repoProductAdder,
		eventProducer:    eventProducer,
		log:              log,
	}
}

type RepoProductAdder interface {
	AddProduct(context.Context, string, string, string, string, int) error
}

func (u *AddUseacase) AddProduct(ctx context.Context, sellerID, name, categoryName, Description string, price int) error {
	const op = "product.usecase.AddProduct"

	if err := u.repoProductAdder.AddProduct(ctx, sellerID, name, categoryName, Description, price); err != nil {
		return fmt.Errorf("%s:: %w", op, err)

	}

	event := Event{
		URL:       "product",
		Action:    "create",
		Timestamp: time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, event, "user-events"); err != nil {
			u.log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
		}
	}()

	return nil
}
