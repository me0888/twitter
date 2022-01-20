package app

import (
	"net/http"
)

func (s *Server) handleReadTweets(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	resp, err := s.postsSvc.ReadTweets(request.Context(), id)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, resp, http.StatusOK)

}
