package controller

import (
	"encoding/json"
	"log"
	"mangosteen/config/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	teardownTest := setupTestCase(t)
	defer teardownTest(t)

	sc := SessionController{}
	sc.RegisterRoutes(r.Group("/api"))

	email := "1@qq.com"
	code := "1234"
	if _, err := q.CreateValidationCode(c, queries.CreateValidationCodeParams{
		Email: email, Code: code,
	}); err != nil {
		log.Fatalln(err)
	}
	user, err := q.CreateUser(c, email)
	if err != nil {
		log.Fatalln(err)
	}
	w := httptest.NewRecorder()
	j := gin.H{
		"email": email,
		"code":  code,
	}
	bytes, _ := json.Marshal(j)
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/session",
		strings.NewReader(string(bytes)),
	)
	r.ServeHTTP(w, req)
	var responseBody struct {
		JWT    string `json:"jwt"`
		UserID int32  `json:"userId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Error("jwt is not a string")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, user.ID, responseBody.UserID)
}
