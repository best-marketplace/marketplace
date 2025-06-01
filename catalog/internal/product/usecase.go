package product

import (
	"catalog/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	URL        string    `json:"url,omitempty"`
	Ids        []string  `json:"ids,omitempty"`
	Action     string    `json:"action"`
	ProductID  string    `json:"product_id,omitempty"`
	Title      string    `json:"title,omitempty"`
	SellerName string    `json:"seller_name,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
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
	AddProduct(context.Context, string, string, string, string, int) (uuid.UUID, error)
}

func (u *AddUseacase) AddProduct(ctx context.Context, sellerName, name, categoryName, Description string, price int) error {
	const op = "product.usecase.AddProduct"
	productID, err := u.repoProductAdder.AddProduct(ctx, sellerName, name, categoryName, Description, price)
	if err != nil {
		return fmt.Errorf("%s:: %w", op, err)

	}

	eventUser := Event{
		URL:       "product",
		Action:    "create",
		Timestamp: time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, eventUser, "user-events"); err != nil {
			u.log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
		}
	}()

	eventProduct := Event{
		Action:     "product_created",
		ProductID:  productID.String(),
		Title:      name,
		SellerName: sellerName,
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, eventProduct, "product-events"); err != nil {
			u.log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
		}
	}()

	return nil
}

type RepoProductSearch interface {
	SearchProduct(context.Context, string) ([]*models.ProductView, []string, error)
}

type SearchUseacase struct {
	RepoProductSearch RepoProductSearch
	eventProducer     EventProducer
	log               *slog.Logger
}

func NewSearchUseacase(log *slog.Logger, RepoProductSearch RepoProductSearch, eventProducer EventProducer) *SearchUseacase {
	return &SearchUseacase{
		RepoProductSearch: RepoProductSearch,
		eventProducer:     eventProducer,
		log:               log,
	}
}

func (u *SearchUseacase) SearchProduct(ctx context.Context, req string) ([]*models.ProductView, error) {
	const op = "product.usecase.ViewListProducts"

	products, ids, err := u.RepoProductSearch.SearchProduct(ctx, req)
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
		slog.String("search_request", req),
	)

	return products, nil
}
