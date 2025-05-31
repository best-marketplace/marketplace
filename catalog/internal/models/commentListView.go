package models

type CommentListView struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Comment string `json:"comment"`
}
