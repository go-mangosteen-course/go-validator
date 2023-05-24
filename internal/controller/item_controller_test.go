package controller

import (
	"encoding/json"
	"mangosteen/config/queries"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/items",
		strings.NewReader(`{
			"amount": 100,
			"kind": "expenses",
			"happened_at": "2020-01-01T00:00:00Z",
			"tag_ids": [1, 2, 3]
		}`),
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource queries.Item
	}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, j.Resource.UserID)
	assert.Equal(t, int32(100), j.Resource.Amount)
}

func TestCreateItemWithError(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ic := ItemController{}
	ic.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/items",
		strings.NewReader(`{
			"kind": "expenses",
			"happened_at": "2020-01-01T00:00:00Z",
			"tag_ids": [1, 2, 3]
		}`),
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	jwtString, _ := jwt_helper.GenerateJWT(int(u.ID))
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 422, w.Code)
}
