package comment

import (
	"catalog/internal/lib/handlers/response"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
