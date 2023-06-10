package controller

import (
	"context"
	"mangosteen/config/queries"
	"mangosteen/internal"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
	q *queries.Queries
	c context.Context
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r = gin.New()
	r.Use(gin.Recovery())
	internal.InitRouter(r)

	q = database.NewQuery()
	c = context.Background()

	if err := q.DeleteAllUsers(c); err != nil {
		t.Fatal(err)
	}
	if err := q.DeleteAllItems(c); err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		database.Close()
	}
}

func signIn(t *testing.T, userID int32, req *http.Request) {
	jwtString, _ := jwt_helper.GenerateJWT(int(userID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
}
