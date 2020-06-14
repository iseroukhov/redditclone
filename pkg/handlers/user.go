package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"redditclone/pkg/post"
	"redditclone/pkg/response"

	"github.com/gorilla/mux"

	"redditclone/pkg/jwt"
	//"redditclone/internal/post"
	"redditclone/pkg/user"

	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	logger   *logrus.Logger
	repo     *user.Repository
	postRepo *post.Repository
}

func NewUserHandler(logger *logrus.Logger, userRepo *user.Repository, postRepo *post.Repository) *UserHandler {
	return &UserHandler{
		logger:   logger,
		repo:     userRepo,
		postRepo: postRepo,
	}
}

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

var (
	UsernameRequired = errors.New("field username is required")
	PasswordRequired = errors.New("field password is required")
	UsernameInvalid  = errors.New("invalid username")
)

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.logger.Errorf("can't create new user, request body: %s", r.Body)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	if req.Username == "" {
		response.Error(w, UsernameRequired, http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		response.Error(w, PasswordRequired, http.StatusBadRequest)
		return
	}
	usr, err := h.repo.Create(user.NewUser(req.Username, req.Password))
	if err != nil {
		h.logger.Errorf("can't create new user: %s", err)
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	}
	token, err := jwt.GetToken(usr.ID, usr.Username)
	if err != nil {
		h.logger.Errorf("can't generate user token, user: %s", usr)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("create new user: %s", usr.Username)
	response.JSON(w, &tokenResponse{Token: token}, http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.logger.Errorf("can't login user, request body: %s", r.Body)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	if req.Username == "" {
		response.Error(w, UsernameRequired, http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		response.Error(w, PasswordRequired, http.StatusBadRequest)
		return
	}
	usr, err := h.repo.Authorize(user.NewUser(req.Username, req.Password))
	if err != nil {
		response.Error(w, err, http.StatusUnauthorized)
		return
	}
	token, err := jwt.GetToken(usr.ID, usr.Username)
	if err != nil {
		h.logger.Errorf("can't generate user token, user: %s", usr)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("authorize user: %#v", usr)
	response.JSON(w, &tokenResponse{Token: token}, http.StatusOK)
}

func (h *UserHandler) Posts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username, ok := params["username"]
	if !ok || username == "" {
		response.Error(w, UsernameInvalid, http.StatusBadGateway)
		return
	}
	posts, err := h.postRepo.GetByAuthor(username)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, posts, http.StatusOK)
}
