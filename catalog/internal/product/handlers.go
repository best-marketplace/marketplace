package product

import (
	"catalog/internal/lib/handlers/response"
	"catalog/internal/models"
	"context"
	"log/slog"
	"net/http"
	"strconv"
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
