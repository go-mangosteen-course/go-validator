package controller

import (
	"context"
	"mangosteen/config"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
	q *queries.Queries
	c context.Context
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	r = gin.Default()
	config.LoadAppConfig()
	database.Connect()

	q = database.NewQuery()
	c = context.Background()

	if err := q.DeleteAllUsers(c); err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		database.Close()
	}

}
