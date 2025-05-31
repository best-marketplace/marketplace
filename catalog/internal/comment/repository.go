package comment

import (
	"catalog/internal/models"
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

func (rep *CommentRepository) ViewCommentInProduct(
	ctx context.Context,
	productID string,
	offset int,
	limit int,
) ([]*models.CommentListView, []string, error) {
	const query = `
        SELECT id, user_id, comment
        FROM comments
        WHERE product_id = $1
        OFFSET $2
        LIMIT $3;
    `

	rows, err := rep.db.QueryContext(ctx, query, productID, offset, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("querying comments: %w", err)
	}
	defer rows.Close()

	var (
		comments []*models.CommentListView
		ids      []string
	)

	for rows.Next() {
		var comment models.CommentListView
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Comment); err != nil {
			return nil, nil, fmt.Errorf("scanning comment: %w", err)
		}
		comments = append(comments, &comment)
		ids = append(ids, comment.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("reading rows: %w", err)
	}

	return comments, ids, nil
}
