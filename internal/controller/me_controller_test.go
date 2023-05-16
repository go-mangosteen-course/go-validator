package controller

import (
	"encoding/json"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMe(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/me",
		strings.NewReader(""),
	)
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestGetMeWithJWT(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)

	mc := MeController{}
	mc.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/me",
		strings.NewReader(""),
	)
	u, err := q.CreateUser(c, "1@qq.com")
	if err != nil {
		t.Fatal(err)
	}
	jwtString, err := jwt_helper.GenerateJWT(int(u.ID))
	if err != nil {
		t.Fatal(err)
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + jwtString},
	}
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j struct {
		Resource struct {
			ID    float64 `json:"id"`
			Email string  `json:"email"`
		} `json:"resource"`
	}
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, int32(j.Resource.ID))
	assert.Equal(t, u.Email, j.Resource.Email)
}
