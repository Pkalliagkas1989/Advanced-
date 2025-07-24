package handlers

import (
	"encoding/json"
	"net/http"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// PostHandler handles post related endpoints
type PostHandler struct {
	PostRepo *repository.PostRepository
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(repo *repository.PostRepository) *PostHandler {
	return &PostHandler{PostRepo: repo}
}

// CreatePost creates a new post for the authenticated user
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CategoryIDs []int  `json:"category_ids"` // Instead of CategoryID
		Title       string `json:"title"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if len(req.CategoryIDs) == 0 || req.Title == "" || req.Content == "" {
		utils.ErrorResponse(w, "At least one category, title and content are required", http.StatusBadRequest)
		return
	}

	post := models.Post{
		UserID:  user.ID,
		Title:   req.Title,
		Content: req.Content,
	}

	created, err := h.PostRepo.Create(post, req.CategoryIDs)
	if err != nil {
		utils.ErrorResponse(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, created, http.StatusCreated)
}

// UpdatePost allows the author to update the title or content of an existing post
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		PostID  string  `json:"post_id"`
		Title   *string `json:"title,omitempty"`
		Content *string `json:"content,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PostID == "" {
		utils.ErrorResponse(w, "Post ID required", http.StatusBadRequest)
		return
	}

	if req.Title == nil && req.Content == nil {
		utils.ErrorResponse(w, "No fields to update", http.StatusBadRequest)
		return
	}

	if err := h.PostRepo.Update(req.PostID, user.ID, req.Title, req.Content); err != nil {
		utils.ErrorResponse(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
}
