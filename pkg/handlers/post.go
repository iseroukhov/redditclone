package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/pkg/post"
	"redditclone/pkg/response"
	"redditclone/pkg/user"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	logger *logrus.Logger
	repo   *post.Repository
}

func NewPostHandler(logger *logrus.Logger, repo *post.Repository) *PostHandler {
	return &PostHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repo.GetAll()
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, posts, http.StatusOK)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	pst := post.NewPost()
	if err := json.NewDecoder(r.Body).Decode(pst); err != nil {
		h.logger.Errorf("can't create post, request body: %s", r.Body)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	usr, err := user.FromContext(r.Context())
	if err != nil || usr == nil {
		h.logger.Warn("can't create post", err)
		response.Error(w, user.Unauthorized, http.StatusUnauthorized)
		return
	}
	pst.Author = usr
	pst.Votes = append(pst.Votes, post.NewVote(usr.ID, 1))
	pst, err = h.repo.Create(pst)
	if err != nil {
		h.logger.Errorf("can't create new post: %s", err.Error())
		response.Error(w, err, http.StatusUnprocessableEntity)
		return
	}
	h.logger.Infof("create new post: %s", pst.ID)
	response.JSON(w, pst, http.StatusCreated)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" || !bson.IsObjectIdHex(id) {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	pst, err := h.repo.GetByID(id)
	if err != nil || pst == nil {
		response.Error(w, post.InvalidID, http.StatusNotFound)
		return
	}
	usr, err := user.FromContext(r.Context())
	if err != nil {
		response.Error(w, err, http.StatusUnauthorized)
		return
	}
	if usr == nil || pst.Author.ID != usr.ID {
		response.Error(w, user.Unauthorized, http.StatusUnauthorized)
		return
	}
	err = h.repo.Delete(pst.ID)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, response.Message{Content: "success"}, http.StatusOK)
}

func (h *PostHandler) Detail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" || !bson.IsObjectIdHex(id) {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	pst, err := h.repo.GetByID(id)
	if pst == nil || err != nil {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	if err := h.repo.IncViews(pst); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, pst, http.StatusOK)
}

func (h *PostHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug, ok := params["category"]
	if !ok || slug == "" {
		response.Error(w, post.InvalidCategoryID, http.StatusBadRequest)
		return
	}
	posts, err := h.repo.GetByCategory(slug)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, posts, http.StatusOK)
}

func (h *PostHandler) UpOrDownvote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" || !bson.IsObjectIdHex(id) {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	usr, err := user.FromContext(r.Context())
	if err != nil || usr == nil {
		h.logger.Errorf("can't upvote post, unauthorized")
		response.Error(w, user.Unauthorized, http.StatusUnauthorized)
		return
	}
	pst, err := h.repo.GetByID(id)
	if err != nil || pst == nil {
		response.Error(w, post.NotFound, http.StatusNotFound)
		return
	}
	var find bool
	for i, v := range pst.Votes {
		if v.UserID == usr.ID {
			find = true
			pst.Votes = append(pst.Votes[:i], pst.Votes[i+1:]...)
			break
		}
	}
	if !find {
		value := -1
		if strings.HasSuffix(r.URL.Path, "upvote") {
			value = 1
		}
		pst.Votes = append(pst.Votes, post.NewVote(usr.ID, value))
	}
	pst, err = h.repo.Update(pst)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
	}
	response.JSON(w, pst, http.StatusOK)
}

type comment struct {
	Comment string `json:"comment"`
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" || !bson.IsObjectIdHex(id) {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	usr, err := user.FromContext(r.Context())
	if err != nil || usr == nil {
		h.logger.Errorf("can't upvote post; unauthorized")
		response.Error(w, user.Unauthorized, http.StatusUnauthorized)
		return
	}
	pst, err := h.repo.GetByID(id)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	if pst == nil {
		response.Error(w, post.NotFound, http.StatusNotFound)
		return
	}
	comm := &comment{}
	if err := json.NewDecoder(r.Body).Decode(comm); err != nil {
		h.logger.Errorf("can't create comment, request body: %s", r.Body)
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	c := post.NewComment(comm.Comment, usr)
	pst.Comments = append(pst.Comments, c)
	_, err = h.repo.Update(pst)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
	}
	response.JSON(w, pst, http.StatusCreated)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" || !bson.IsObjectIdHex(id) {
		response.Error(w, post.InvalidID, http.StatusBadRequest)
		return
	}
	commentId, ok := params["comment_id"]
	if !ok || id == "" || !bson.IsObjectIdHex(commentId) {
		response.Error(w, post.InvalidCommentID, http.StatusBadRequest)
		return
	}
	usr, err := user.FromContext(r.Context())
	if err != nil || usr == nil {
		h.logger.Errorf("can't delete comment; unauthorized")
		response.Error(w, err, http.StatusUnauthorized)
		return
	}
	pst, err := h.repo.GetByID(id)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	if pst == nil {
		response.Error(w, post.InvalidID, http.StatusNotFound)
		return
	}
	for i, c := range pst.Comments {
		if c.ID == bson.ObjectIdHex(commentId) {
			if c.Author.ID != usr.ID {
				h.logger.Errorf("can't delete comment, unauthorized")
				response.Error(w, user.Unauthorized, http.StatusUnauthorized)
				return
			}
			pst.Comments = append(pst.Comments[:i], pst.Comments[i+1:]...)
		}
	}
	_, err = h.repo.Update(pst)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
	}
	response.JSON(w, pst, http.StatusOK)
}
