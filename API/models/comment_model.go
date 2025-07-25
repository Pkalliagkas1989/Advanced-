package models

import "time"

type Comment struct {
	ID        string     `json:"id"`
	PostID    string     `json:"post_id"`
	UserID    string     `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CommentWithUser is a comment along with the username of its author
type CommentWithUser struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CommentWithPost represents a user's comment along with the title of the post
// it belongs to. This is useful for activity pages displaying a list of
// comments in context.
type CommentWithPost struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	PostTitle string    `json:"post_title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
