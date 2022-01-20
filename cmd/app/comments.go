package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/me0888/twitter/pkg/models"
)

func (s *Server) handleGetCommentByID(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}
	commentId, ok := mux.Vars(request)["comment_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.commentsSvc.GetComment(request.Context(), commentId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleUpdateComment(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	var updateCommentInput models.Comment
	if err := json.NewDecoder(request.Body).Decode(&updateCommentInput); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := s.commentsSvc.GetComment(request.Context(), strconv.FormatInt(updateCommentInput.ID, 10))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	comment.Content = updateCommentInput.Content
	comment.UpdatedAt = time.Now()

	_, err = s.commentsSvc.UpdateComment(request.Context(), id, comment)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, comment, http.StatusOK)

}

func (s *Server) handleDeleteComment(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	commentId, ok := mux.Vars(request)["comment_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	comment, err := s.commentsSvc.DeleteComment(request.Context(), id, commentId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, comment, http.StatusOK)

}

func (s *Server) handleLikeComment(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}
	commentId, ok := mux.Vars(request)["comment_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.commentsSvc.CommentLike(request.Context(), id, commentId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleGetCommentsLikedUsers(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	commentId, ok := mux.Vars(request)["comment_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.commentsSvc.GetCommetsLikedUsers(request.Context(), commentId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleCreateComment(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	var in models.CreateCommentInput
	if err := json.NewDecoder(request.Body).Decode(&in); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.commentsSvc.CreateComment(request.Context(), id, tweetId, in.Content)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleGetTweetComments(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}
	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.commentsSvc.GetComments(request.Context(), tweetId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}
