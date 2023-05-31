package controller

import (
	"encoding/json"
	"mangosteen/config/queries"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := TagController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/api/v1/tags",
		strings.NewReader(`{
			"name": "通勤",
			"kind": "expenses",
			"sign": "👟"
		}`),
	)

	u, _ := q.CreateUser(c, "1@qq.com")
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource queries.Tag
	}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, j.Resource.UserID)
	assert.Equal(t, "通勤", j.Resource.Name)
	assert.Nil(t, j.Resource.DeletedAt)
}
