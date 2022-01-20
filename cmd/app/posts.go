package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/me0888/twitter/pkg/models"
)

func (s *Server) handleCreateTweet(writer http.ResponseWriter, request *http.Request) {
	var createPostInput models.CreatePostInput

	if err := json.NewDecoder(request.Body).Decode(&createPostInput); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	resp, err := s.postsSvc.CreateTweet(request.Context(), id, createPostInput.Content)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleGetTweetByID(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.GetTweet(request.Context(), tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleUpdateTweet(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	var updatePostInput models.Tweet
	if err := json.NewDecoder(request.Body).Decode(&updatePostInput); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tweet, err := s.postsSvc.GetTweet(request.Context(), strconv.FormatInt(updatePostInput.ID, 10))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tweet.Content = updatePostInput.Content
	tweet.UpdatedAt = time.Now()

	_, err = s.postsSvc.UpdateTweet(request.Context(), id, tweet)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, tweet, http.StatusOK)

}

func (s *Server) handleDeleteTweet(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tweet, err := s.postsSvc.DeleteTweet(request.Context(), id, tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, tweet, http.StatusOK)

}

func (s *Server) handleLikeTweet(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.TweetLike(request.Context(), id, tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleRetweetTweet(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.TweetRetweet(request.Context(), id, tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleTweetLikedUsers(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.TweetLikes(request.Context(), tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleTweetRetweetedUsers(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	tweetId, ok := mux.Vars(request)["tweet_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.TweetRetweetedUsers(request.Context(), tweetId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleGetTweets(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.postsSvc.GetTweets(request.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}
