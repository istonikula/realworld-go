package apitest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Router *gin.Engine
	Token  *string
}

func (c Client) Get(path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("GET", path, nil)
	c.maybeToken(r)

	return serve(c.Router, r)
}

func (c Client) Post(path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	r, _ := http.NewRequest("POST", path, bytes.NewBuffer(b))
	c.maybeToken(r)

	return serve(c.Router, r)
}

func serve(n *gin.Engine, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)
	return w
}

func (c Client) maybeToken(r *http.Request) {
	if c.Token != nil {
		r.Header.Add("Authorization", "Token "+*c.Token)
	}
}
