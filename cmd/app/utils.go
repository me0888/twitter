package app

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/me0888/twitter/pkg/comments"
	"github.com/me0888/twitter/pkg/posts"
	"github.com/me0888/twitter/pkg/users"
)

func writeJSON(writer http.ResponseWriter, item interface{}, code int) {

	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) Auth(writer http.ResponseWriter, request *http.Request) (id int64) {
	token := request.Header.Get("Authorization")
	id, err := s.usersSvc.IDByToken(request.Context(), token)
	if err != nil {
		status := map[string]string{"status": "Не авторизован"}
		writeJSON(writer, status, http.StatusUnauthorized)
		return
	}
	return id
}

func Execute(host string, port string, dns string) (err error) {
	connCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	pool, err := pgxpool.Connect(connCtx, dns)
	if err != nil {
		log.Println(err)
		return
	}

	defer pool.Close()

	mux := mux.NewRouter()
	usersSvc := users.NewService(pool)
	postsSvc := posts.NewService(pool)
	commentsSvc := comments.NewService(pool)
	server := NewServer(mux, usersSvc, postsSvc, commentsSvc)
	server.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	log.Println("Server starts at port: " + port)
	return srv.ListenAndServe()
}
