package handlers

/*
import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"redditclone/internal/post"
	"redditclone/internal/user"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	userRepo = user.NewRepository()
	postRepo = post.NewRepository()
)

type UserHandlerTestCase struct {
	Body       string
	Response   string
	IsError    bool
	StatusCode int
}

func TestRegister(t *testing.T) {
	cases := []*UserHandlerTestCase{
		{
			Body:       `{"username":"testuser","password":"password"}`,
			Response:   `{"token":""}`,
			StatusCode: http.StatusCreated,
		},
		{
			Body:       `{"username":"testuser","password":"password"}`,
			Response:   `{"message":"already exists"}`,
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Body:       `{"username":"","password":"password"}`,
			Response:   `{"message":"field username is required"}`,
			StatusCode: http.StatusBadRequest,
		},
		{
			Body:       `{"username":"testuser","password":""}`,
			Response:   `{"message":"field password is required"}`,
			StatusCode: http.StatusBadRequest,
		},
		{
			Body:       `{"broken_json"}`,
			Response:   `{"message":"internal server error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}

	userHandler := NewUserHandler(logrus.New(), userRepo, postRepo)

	for caseNum, item := range cases {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(item.Body)))
		if err != nil {
			t.Fatal(err)
		}

		userHandler.Register(rec, req)

		if item.StatusCode != http.StatusCreated && strings.TrimRight(rec.Body.String(), "\n") != item.Response {
			t.Error(item.Response)
			t.Error(rec.Body.String())
			t.Errorf("[%d] unexpected body; expected: %s, got: %s", caseNum, item.Response, rec.Body.String())
		}

		if item.StatusCode != rec.Code {
			t.Errorf("[%d] unexpected status code; expected: %d, get: %d", caseNum, item.StatusCode, rec.Code)
		}
	}
}

func TestLogin(t *testing.T) {
	cases := []*UserHandlerTestCase{
		{
			Body:       `{"username":"testuser","password":"password"}`,
			Response:   `{"token":""}`,
			StatusCode: http.StatusOK,
		},
		{
			Body:       `{"username":"notfound","password":"password"}`,
			Response:   `{"message":"user not found"}`,
			StatusCode: http.StatusUnauthorized,
		},
		{
			Body:       `{"username":"","password":"password"}`,
			Response:   `{"message":"field username is required"}`,
			StatusCode: http.StatusBadRequest,
		},
		{
			Body:       `{"username":"testuser","password":""}`,
			Response:   `{"message":"field password is required"}`,
			StatusCode: http.StatusBadRequest,
		},
		{
			Body:       `{"broken_json"}`,
			Response:   `{"message":"internal server error"}`,
			StatusCode: http.StatusInternalServerError,
		},
		{
			Body:       `{"username":"createduser","password":"broken_password"}`,
			Response:   `{"message":"invalid password"}`,
			StatusCode: http.StatusUnauthorized,
		},
	}

	_, err := userRepo.Create("createduser", "password")
	if err != nil {
		t.Errorf("unexpected error, %s", err)
	}

	userHandler := NewUserHandler(logrus.New(), userRepo, postRepo)

	for caseNum, item := range cases {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer([]byte(item.Body)))

		userHandler.Login(rec, req)

		if item.StatusCode != http.StatusOK && strings.TrimRight(rec.Body.String(), "\n") != item.Response {
			t.Errorf("[%d] unexpected body; expected: %s, got: %s", caseNum, item.Response, rec.Body.String())
		}

		if item.StatusCode != rec.Code {
			t.Errorf("[%d] unexpected status code; expected: %d, got: %d", caseNum, item.StatusCode, rec.Code)
		}
	}
}

func TestPosts(t *testing.T) {
	cases := []*UserHandlerTestCase{
		{
			Body:       "",
			Response:   `{"message":"invalid username"}`,
			StatusCode: http.StatusBadGateway,
			IsError:    true,
		},
		{
			Body:       "author",
			Response:   "",
			StatusCode: http.StatusOK,
			IsError:    false,
		},
	}

	userHandler := NewUserHandler(logrus.New(), userRepo, postRepo)

	u, err := userHandler.userRepo.Create("author", "password")
	if err != nil {
		t.Errorf("unexpected error, %s", err)
	}

	_, err = userHandler.postRepo.Create(&post.Post{
		Type:     "text",
		Category: "programming",
		Title:    "test",
		Text:     "test",
		Author:   u,
	})
	if err != nil {
		t.Errorf("unexpected error, %s", err)
	}

	for caseNum, item := range cases {
		path := "/api/user/" + item.Body
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{"username": item.Body})
		userHandler.Posts(rec, req)

		if item.IsError && strings.TrimRight(rec.Body.String(), "\n") != item.Response {
			t.Errorf("[%d] unexpected body; expected: %s, got: %s", caseNum, item.Response, rec.Body.String())
		}

		if item.StatusCode != rec.Code {
			t.Errorf("[%d] unexpected status code; expected: %d, got: %d", caseNum, item.StatusCode, rec.Code)
		}
	}
}
*/
