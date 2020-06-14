package server

import (
	"database/sql"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/post"
	"redditclone/pkg/user"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
	repo   *Repository
}

type Repository struct {
	UserRepository *user.Repository
	PostRepository *post.Repository
}

func New(mysql *sql.DB, mongodb *mgo.Database) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static"))))

	return &Server{
		router: router,
		logger: logrus.New(),
		repo: &Repository{
			UserRepository: user.NewRepository(mysql),
			PostRepository: post.NewRepository(mongodb.C("posts")),
		},
	}
}

func (s *Server) Start(addr string) error {
	if err := s.initLogger(); err != nil {
		return err
	}
	s.initRouter()

	apiHandler := middleware.API(s.router)
	apiHandler = middleware.Auth(s.repo.UserRepository, apiHandler)
	apiHandler = middleware.Panic(s.logger, apiHandler)

	s.logger.Info("server started on ", addr)
	return http.ListenAndServe(addr, apiHandler)
}

func (s *Server) initLogger() error {
	lvl, err := logrus.ParseLevel("debug")
	if err != nil {
		return err
	}
	s.logger.SetLevel(lvl)
	return nil
}

func (s *Server) initRouter() {
	frontendHandler := handlers.NewFrontendHandler()
	s.router.HandleFunc("/", frontendHandler.IndexPage).Methods("GET")

	userHandler := handlers.NewUserHandler(s.logger, s.repo.UserRepository, s.repo.PostRepository)
	s.router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	s.router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	s.router.HandleFunc("/api/user/{username}", userHandler.Posts).Methods("GET")

	postHandler := handlers.NewPostHandler(s.logger, s.repo.PostRepository)
	s.router.HandleFunc("/api/posts/", postHandler.Index).Methods("GET")
	s.router.HandleFunc("/api/posts", postHandler.Create).Methods("POST")
	s.router.HandleFunc("/api/post/{id}", postHandler.Delete).Methods("DELETE")
	s.router.HandleFunc("/api/post/{id}", postHandler.Detail).Methods("GET")
	s.router.HandleFunc("/api/post/{id}/unvote", postHandler.UpOrDownvote).Methods("GET")
	s.router.HandleFunc("/api/post/{id}/upvote", postHandler.UpOrDownvote).Methods("GET")
	s.router.HandleFunc("/api/post/{id}/downvote", postHandler.UpOrDownvote).Methods("GET")
	s.router.HandleFunc("/api/post/{id}", postHandler.AddComment).Methods("POST")
	s.router.HandleFunc("/api/post/{id}/{comment_id}", postHandler.DeleteComment).Methods("DELETE")
	s.router.HandleFunc("/api/posts/{category}", postHandler.CategoryList).Methods("GET")
}
