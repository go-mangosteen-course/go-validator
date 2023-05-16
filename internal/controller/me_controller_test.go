package controller

import (
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
