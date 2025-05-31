package comment

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CommentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (rep *CommentRepository) CreateComment(ctx context.Context, UserID, ProductID, Comment string) (uuid.UUID, error) {
	const op = "comment.repository.CreateComment"

	commentID := uuid.New()

	const queryInsert = `
		INSERT INTO comments (id,user_id,product_id,comment)
		VALUES ($1, $2, $3, $4)
	`

	_, err := rep.db.ExecContext(ctx, queryInsert,
		commentID, UserID, ProductID, Comment,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: failed to insert comment: %w", op, err)
	}

	return commentID, nil
}
