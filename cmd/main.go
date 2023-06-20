package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	handler "github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum/delivery"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum/repo"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum/usecase"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	connectionString := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable&pool_max_conns=1000"
	pgxConn, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatal("Fail to connect to DB", err)
	}

	forumRepo := repo.NewForumRepo(pgxConn)
	forumUsecase := usecase.NewForumUsecase(forumRepo)
	forumHandler := handler.NewForumHandler(forumUsecase)

	apiSubrouter := router.PathPrefix("/api").Subrouter()
	{
		userSubrouter := apiSubrouter.PathPrefix("/user").Subrouter()
		{
			userSubrouter.HandleFunc("/{nickname}/profile", forumHandler.GetUser).Methods(http.MethodGet)
			userSubrouter.HandleFunc("/{nickname}/create", forumHandler.CreateUser).Methods(http.MethodPost)
			userSubrouter.HandleFunc("/{nickname}/profile", forumHandler.UpdateUser).Methods(http.MethodPost)
		}
		forumSubrouter := apiSubrouter.PathPrefix("/forum").Subrouter()
		{
			forumSubrouter.HandleFunc("/create", forumHandler.CreateForum).Methods(http.MethodPost)
			forumSubrouter.HandleFunc("/{slug}/details", forumHandler.GetForumDetails).Methods(http.MethodGet)
			forumSubrouter.HandleFunc("/{slug}/create", forumHandler.CreateThread).Methods(http.MethodPost)
			forumSubrouter.HandleFunc("/{slug}/threads", forumHandler.GetThreads).Methods(http.MethodGet)
			forumSubrouter.HandleFunc("/{slug}/users", forumHandler.GetUsers).Methods(http.MethodGet)
		}
		threadSubrouter := apiSubrouter.PathPrefix("/thread").Subrouter()
		{
			threadSubrouter.HandleFunc("/{slug_or_id}/create", forumHandler.CreatePosts).Methods(http.MethodPost)
			threadSubrouter.HandleFunc("/{slug_or_id}/vote", forumHandler.Vote).Methods(http.MethodPost)
			threadSubrouter.HandleFunc("/{slug_or_id}/details", forumHandler.GetThread).Methods(http.MethodGet)
			threadSubrouter.HandleFunc("/{slug_or_id}/details", forumHandler.UpdateThread).Methods(http.MethodPost)
			threadSubrouter.HandleFunc("/{slug_or_id}/posts", forumHandler.GetThreadPosts).Methods(http.MethodGet)
		}
		postSubrouter := apiSubrouter.PathPrefix("/post").Subrouter()
		{
			postSubrouter.HandleFunc("/{id}/details", forumHandler.GetPost).Methods(http.MethodGet)
			postSubrouter.HandleFunc("/{id}/details", forumHandler.UpdatePost).Methods(http.MethodPost)
		}
		serviceSubrouter := apiSubrouter.PathPrefix("/service").Subrouter()
		{
			serviceSubrouter.HandleFunc("/status", forumHandler.GetStatus).Methods(http.MethodGet)
			serviceSubrouter.HandleFunc("/clear", forumHandler.Clear).Methods(http.MethodPost)
		}
	}

	http.Handle("/", router)
	log.Print(http.ListenAndServe(":5000", router))
}
