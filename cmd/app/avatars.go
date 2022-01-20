package app

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (s *Server) handleUploadAvatar(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	file, handler, err := request.FormFile("avatar")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	extension := strings.Split(handler.Filename, ".")[1]
	file_route := "avatars/" + strconv.FormatInt(id, 10) + "." + extension

	f, err := os.OpenFile(file_route, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	avatar, err := s.usersSvc.UpdateAvatar(request.Context(), file_route, id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, avatar, http.StatusOK)

}

func (s *Server) handleGetAvatar(writer http.ResponseWriter, request *http.Request) {
	id := s.Auth(writer, request)
	if id == 0 {
		return
	}

	avatar, err := s.usersSvc.GetAvatar(request.Context(), id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	OpenFile, err := os.Open(avatar)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = io.Copy(writer, OpenFile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

}
