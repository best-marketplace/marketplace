package product

import (
	"catalog/internal/lib/handlers/response"
	"catalog/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type ProductListViewer interface {
	ViewListProducts(context.Context, int, int) ([]*models.ProductListView, error)
}

func ViewListProducts(log *slog.Logger, productListViewer ProductListViewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.ViewListProducts"

		q := r.URL.Query()

		offsetStr := q.Get("offset")
		limitStr := q.Get("limit")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			log.Warn(op+": invalid offset", slog.String("offset", offsetStr), slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusBadRequest, "invalid 'offset' parameter")
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			log.Warn(op+": invalid limit", slog.String("limit", limitStr), slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusBadRequest, "invalid 'limit' parameter")
			return
		}

		products, err := productListViewer.ViewListProducts(r.Context(), offset, limit)
		if err != nil {
			log.Error(op+": failed to get products", slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusInternalServerError, "failed to get products")
			return
		}

		response.RespondWithJSON(w, log, http.StatusOK, products)
	}
}

type ProductViewer interface {
	ViewProduct(context.Context, string) (*models.ProductView, error)
}

func ViewProduct(log *slog.Logger, productViewer ProductViewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.ViewProduct"

		idStr := r.URL.Query().Get("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Warn(op+": invalid UUID", slog.String("id", idStr), slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusBadRequest, "invalid 'id' parameter")
			return
		}

		product, err := productViewer.ViewProduct(r.Context(), id.String())
		if err != nil {
			log.Error(op+": failed to get product", slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusInternalServerError, "failed to get product")
			return
		}

		if product == nil {
			response.RespondWithError(w, log, http.StatusNotFound, "product not found")
			return
		}

		response.RespondWithJSON(w, log, http.StatusOK, product)
	}
}

const (
	warnGetUserID        = "failed to get user ID from context"
	messageUnauthorized  = "Unauthorized"
	errAddProduct        = "failed to add product"
	ErrInsufficientFunds = "insufficient balance to send coin"
	successAddproduct    = "add product successfully"
)

type RequestAddProduct struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	CategoryName string `json:"categoryName"`
}

type ProductAdder interface {
	AddProduct(context.Context, string, string, string, string, int) error
}

func AddProduct(log *slog.Logger, productAdder ProductAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "product.handlers.AddProduct"

		// userID, ok := middleware.GetUserID(r)
		// if !ok {
		// 	log.Warn(warnGetUserID)

		// 	response.RespondWithError(w, log, http.StatusUnauthorized, messageUnauthorized)
		// }
		userID := uuid.New().String()

		log := log.With(
			slog.String("op", op),
			slog.String("user_id", userID),
		)

		var req RequestAddProduct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(fmt.Errorf("failed to decode request body: Error: %w", err).Error())

			response.RespondWithError(w, log, http.StatusBadRequest, "invalid request")
			return
		}

		if err := productAdder.AddProduct(r.Context(), userID, req.Name, req.Description, req.CategoryName, req.Price); err != nil {
			log.Error(fmt.Errorf("%s Error: %w", errAddProduct, err).Error())

			// if errors.Is(err, wallet.ErrInsufficientFunds) {
			// 	response.RespondWithError(w, log, http.StatusBadRequest, ErrInsufficientFunds)
			// 	return
			// }

			response.RespondWithError(w, log, http.StatusInternalServerError, errAddProduct)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info(successAddproduct)

	}
}
