package comment

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

const (
	messageUnauthorized  = "Unauthorized"
	errCreateComment     = "failed to create comment"
	successCreateComment = "create comment successfully"
)

type RequestCreateComment struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Comment   string `json:"comment"`
}

type CommentCreater interface {
	CreateComment(context.Context, string, string, string) error
}

func CreateComment(log *slog.Logger, CommentCreater CommentCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "comment.handlers.CreateComment"

		log := log.With(
			slog.String("op", op),
		)

		var req RequestCreateComment
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(fmt.Errorf("failed to decode request body: Error: %w", err).Error())

			response.RespondWithError(w, log, http.StatusBadRequest, "invalid request")
			return
		}

		log = log.With(
			slog.String("user_id", req.UserID),
		)

		if err := CommentCreater.CreateComment(r.Context(), req.UserID, req.ProductID, req.Comment); err != nil {
			log.Error(fmt.Errorf("%s Error: %w", errCreateComment, err).Error())

			response.RespondWithError(w, log, http.StatusInternalServerError, errCreateComment)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info(successCreateComment)

	}
}

type CommentListViewer interface {
	ViewCommentInProduct(context.Context, string, int, int) ([]*models.CommentListView, error)
}

func ViewCommentInProduct(log *slog.Logger, CommentListViewer CommentListViewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "comment.handlers.ViewCommentInProduct"

		idStr := r.URL.Query().Get("product_id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			log.Warn(op+": invalid UUID", slog.String("id", idStr), slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusBadRequest, "invalid 'id' parameter")
			return
		}
		fmt.Println(productID)
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

		comments, err := CommentListViewer.ViewCommentInProduct(r.Context(), productID.String(), offset, limit)
		if err != nil {
			log.Error(op+": failed to get comments", slog.Any("err", err))
			response.RespondWithError(w, log, http.StatusInternalServerError, "failed to get comments")
			return
		}

		response.RespondWithJSON(w, log, http.StatusOK, comments)
	}
}
