package repository

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetAllComments() ([]models.Comment, error) {
	rows, err := r.db.Query(`
		SELECT comment_id, post_id, user_id, content, created_at, updated_at 
		FROM comments ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Create inserts a new comment into the database
func (r *CommentRepository) Create(comment models.Comment) (*models.Comment, error) {
	comment.ID = utils.GenerateUUID()
	comment.CreatedAt = time.Now()
	_, err := r.db.Exec(`INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
		comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsByUserWithPost returns all comments created by a specific user
// along with the title of the post each comment belongs to.
func (r *CommentRepository) GetCommentsByUserWithPost(userID string) ([]models.CommentWithPost, error) {
	rows, err := r.db.Query(`
        SELECT c.comment_id, c.post_id, p.title, c.content, c.created_at
        FROM comments c
        JOIN posts p ON c.post_id = p.post_id
        WHERE c.user_id = ?
        ORDER BY c.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentWithPost
	for rows.Next() {
		var c models.CommentWithPost
		if err := rows.Scan(&c.ID, &c.PostID, &c.PostTitle, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
