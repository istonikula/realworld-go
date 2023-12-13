package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
		var db = db()
		defer deleteUsers(db)
		var client = TestClient{router(db), nil}

		r := client.Post("/api/users", rest.UserRegistration{
			Email:    testUser.Email,
			Username: testUser.Username,
			Password: testUser.Password,
		})

		assert.Equal(t, http.StatusCreated, r.Code)

		var act rest.UserResponse
		assert.NoError(t, json.Unmarshal(r.Body.Bytes(), &act))
		exp := rest.User{Email: testUser.Email, Username: testUser.Username, Token: "ignore", Bio: nil, Image: nil}
		assertUserIgnoreToken(t, exp, act.User)
	})
}

func deleteUsers(db *sqlx.DB) {
	db.MustExec("DELETE FROM users")
}

func assertUserIgnoreToken(t *testing.T, exp, act rest.User) {
	exp.Token = "ignore"
	act.Token = "ignore"
	assert.Equal(t, exp, act)
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
