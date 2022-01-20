package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/me0888/twitter/pkg/comments"
	"github.com/me0888/twitter/pkg/posts"
	"github.com/me0888/twitter/pkg/users"
)

type Server struct {
	mux         *mux.Router
	usersSvc    *users.Service
	postsSvc    *posts.Service
	commentsSvc *comments.Service
}

func NewServer(mux *mux.Router, usersSvc *users.Service, postsSvc *posts.Service, commentsSvc *comments.Service) *Server {
	return &Server{mux: mux, usersSvc: usersSvc, postsSvc: postsSvc, commentsSvc: commentsSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func (s *Server) Init() {

	s.mux.HandleFunc("/users", s.handleCreateUser).Methods(POST)
	s.mux.HandleFunc("/login", s.handleLogin).Methods(POST)
	s.mux.HandleFunc("/user", s.handleGetUserByID).Methods(GET)
	s.mux.HandleFunc("/user", s.handleUpdateUser).Methods(PUT)
	s.mux.HandleFunc("/users", s.handleSearchUsers).Methods(GET)
	s.mux.HandleFunc("/users/{username}/follow", s.handleFollow).Methods(POST)
	s.mux.HandleFunc("/users/{username}/followers", s.handleFollowers).Methods(GET)
	s.mux.HandleFunc("/users/{username}/followees", s.handleFollowees).Methods(GET)
	s.mux.HandleFunc("/users/{username}/tweets", s.handleGetTweets).Methods(GET)

	s.mux.HandleFunc("/tweets", s.handleCreateTweet).Methods(POST)
	s.mux.HandleFunc("/tweets", s.handleUpdateTweet).Methods(PUT)
	s.mux.HandleFunc("/tweets/{tweet_id}", s.handleGetTweetByID).Methods(GET)
	s.mux.HandleFunc("/tweets/{tweet_id}", s.handleDeleteTweet).Methods(DELETE)
	s.mux.HandleFunc("/tweets/{tweet_id}/like", s.handleLikeTweet).Methods(POST)
	s.mux.HandleFunc("/tweets/{tweet_id}/liked_users", s.handleTweetLikedUsers).Methods(GET)
	s.mux.HandleFunc("/tweets/{tweet_id}/retweet", s.handleRetweetTweet).Methods(POST)
	s.mux.HandleFunc("/tweets/{tweet_id}/retweeted_users", s.handleTweetRetweetedUsers).Methods(GET)
	s.mux.HandleFunc("/tweets/{tweet_id}/comments", s.handleCreateComment).Methods(POST)
	s.mux.HandleFunc("/tweets/{tweet_id}/comments", s.handleGetTweetComments).Methods(GET)

	s.mux.HandleFunc("/comments", s.handleUpdateComment).Methods(PUT)
	s.mux.HandleFunc("/comments/{comment_id}", s.handleGetCommentByID).Methods(GET)
	s.mux.HandleFunc("/comments/{comment_id}", s.handleDeleteComment).Methods(DELETE)
	s.mux.HandleFunc("/comments/{comment_id}/like", s.handleLikeComment).Methods(POST)
	s.mux.HandleFunc("/comments/{comment_id}/liked_users", s.handleGetCommentsLikedUsers).Methods(GET)

	s.mux.HandleFunc("/avatar", s.handleUploadAvatar).Methods(POST)
	s.mux.HandleFunc("/avatar", s.handleGetAvatar).Methods(GET)

	s.mux.HandleFunc("/feed", s.handleReadTweets).Methods(GET)

}
