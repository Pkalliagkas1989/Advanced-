package handlers

import (
	"net/http"

	"forum/middleware"
	"forum/repository"
	"forum/utils"
)

type MyCommentsHandler struct {
	CommentRepo *repository.CommentRepository
}

func NewMyCommentsHandler(commentRepo *repository.CommentRepository) *MyCommentsHandler {
	return &MyCommentsHandler{CommentRepo: commentRepo}
}

func (h *MyCommentsHandler) GetMyComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	comments, err := h.CommentRepo.GetCommentsByUserWithPost(user.ID)
	if err != nil {
		utils.ErrorResponse(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, comments, http.StatusOK)
}
