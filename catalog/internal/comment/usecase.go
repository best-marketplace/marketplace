package comment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Action      string    `json:"action"`
	CommentID   string    `json:"comment_id,omitempty"`
	Comment     string    `json:"comment,omitempty"`
	URL         string    `json:"url,omitempty"`
	Ids         []string  `json:"ids,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}

type EventProducer interface {
	Send(ctx context.Context, event any, topic string) error
}

type CreateUseacase struct {
	repoCreateComment RepoCreateComment
	eventProducer     EventProducer
	log               *slog.Logger
}

func NewCreateUseacase(log *slog.Logger, repoCreateComment RepoCreateComment, eventProducer EventProducer) *CreateUseacase {
	return &CreateUseacase{
		repoCreateComment: repoCreateComment,
		eventProducer:     eventProducer,
		log:               log,
	}
}

type RepoCreateComment interface {
	CreateComment(context.Context, string, string, string) (uuid.UUID, error)
}

func (u *CreateUseacase) CreateComment(ctx context.Context, UserID, ProductID, Comment string) error {
	const op = "comment.usecase.CreateComment"
	commentID, err := u.repoCreateComment.CreateComment(ctx, UserID, ProductID, Comment)
	if err != nil {
		return fmt.Errorf("%s:: %w", op, err)

	}

	eventUser := Event{
		URL:       "comment",
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
		Action:    "comment_created",
		CommentID: commentID.String(),
		Comment:   Comment,
	}

	go func() {
		ctx := context.Background()
		if err := u.eventProducer.Send(ctx, eventProduct, "comment-events"); err != nil {
			u.log.Error("failed to send event to Kafka", slog.String("err", err.Error()))
		}
	}()

	return nil
}
