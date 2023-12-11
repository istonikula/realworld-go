package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	rwt "github.com/istonikula/realworld-go/realworld-testing"
)

var testUser = struct {
	Email    string
	Username string
	Password string
}{
	Email:    "foo@bar.com",
	Username: "foo",
	Password: "plain",
}

func TestUsers(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		var client = TestClient{router(db()), nil}

		act := client.Post("/api/users", UserRegistration{
			Email:    testUser.Email,
			Username: testUser.Username,
			Password: testUser.Password,
		})

		rwt.Equals(t, act.Code, http.StatusCreated)
	})
}

type TestClient struct {
	router *gin.Engine
	token  *string
}

func (c TestClient) Get(path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("GET", path, nil)
	c.maybeToken(r)

	return serve(c.router, r)
}

func (c TestClient) Post(path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	r, _ := http.NewRequest("POST", path, bytes.NewBuffer(b))
	c.maybeToken(r)

	return serve(c.router, r)
}

func serve(n *gin.Engine, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)
	return w
}

func (c TestClient) maybeToken(r *http.Request) {
	if c.token != nil {
		r.Header.Add("Authorization", "Token"+*c.token)
	}
}
