package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/me0888/twitter/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleCreateUser(writer http.ResponseWriter, request *http.Request) {
	var user *models.UserInput

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	user.Password = string(hash)

	item, err := s.usersSvc.Save(request.Context(), user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, item, http.StatusOK)

}

func (s *Server) handleLogin(writer http.ResponseWriter, request *http.Request) {
	var loginInput models.LoginInput
	var loginOutput models.LoginOutput

	if err := json.NewDecoder(request.Body).Decode(&loginInput); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.usersSvc.Token(request.Context(), loginInput.Email, loginInput.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	loginOutput.Token = token

	writeJSON(writer, loginOutput, http.StatusOK)

}

func (s *Server) handleGetUserByID(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	resp, err := s.usersSvc.User(request.Context(), id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleUpdateUser(writer http.ResponseWriter, request *http.Request) {
	var user *models.UserInput

	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	oldUser, err := s.usersSvc.User(request.Context(), id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if len(user.Password) > 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		user.Password = string(hash)
	}

	if len(user.Email) == 0 {
		user.Email = oldUser.Email
	}

	if len(user.Password) == 0 {
		user.Username = oldUser.UserName
	}

	item, err := s.usersSvc.Update(request.Context(), user, id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, item, http.StatusOK)

}

func (s *Server) handleFollow(writer http.ResponseWriter, request *http.Request) {

	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.usersSvc.Follow(request.Context(), id, username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}

func (s *Server) handleSearchUsers(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	username := request.URL.Query().Get("search")

	resp, err := s.usersSvc.Users(request.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)
}

func (s *Server) handleFollowers(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.usersSvc.Followers(request.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)
}

func (s *Server) handleFollowees(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := s.usersSvc.Followees(request.Context(), username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)
}
