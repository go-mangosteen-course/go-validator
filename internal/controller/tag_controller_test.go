package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"mangosteen/api"
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

func TestUpdateTag(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := TagController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	u, _ := q.CreateUser(c, "1@qq.com")
	tag, err := q.CreateTag(context.Background(), queries.CreateTagParams{
		UserID: u.ID,
		Name:   "通勤",
		Sign:   "🎈",
		Kind:   "expenses",
	})
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"PATCH",
		fmt.Sprintf("/api/v1/tags/%d", tag.ID),
		strings.NewReader(`{
			"name": "吃饭"
		}`),
	)

	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.UpdateTagResponse
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, u.ID, j.Resource.UserID)
	assert.Equal(t, "吃饭", j.Resource.Name)
	assert.Equal(t, "🎈", j.Resource.Sign)
	assert.Equal(t, "expenses", j.Resource.Kind)
	assert.Nil(t, j.Resource.DeletedAt)
}

func TestDeleteTag(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := TagController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	u, _ := q.CreateUser(c, "1@qq.com")
	tag, err := q.CreateTag(context.Background(), queries.CreateTagParams{
		UserID: u.ID,
		Name:   "通勤",
		Sign:   "🎈",
		Kind:   "expenses",
	})
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/tags/%d", tag.ID),
		nil,
	)

	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	_, err = q.FindTag(c, queries.FindTagParams{
		ID:     tag.ID,
		UserID: u.ID,
	})
	assert.Error(t, err)
}

func TestGetPagedTags(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := TagController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/tags",
		nil,
	)
	u, _ := q.CreateUser(c, "1@qq.com")
	for i := 0; i < int(ctrl.PerPage*2); i++ {
		if _, err := q.CreateTag(c, queries.CreateTagParams{
			UserID: u.ID,
			Name:   fmt.Sprintf("通勤%d", i),
			Sign:   "🎈",
			Kind:   "expenses",
		}); err != nil {
			t.Error(err)
		}
	}
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.GetPagesTagsResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, ctrl.PerPage, int32(len(j.Resources)))
}

func TestGetTag(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := TagController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	u, _ := q.CreateUser(c, "1@qq.com")
	tag, _ := q.CreateTag(c, queries.CreateTagParams{
		UserID: u.ID,
		Name:   "通勤",
		Sign:   "🎈",
		Kind:   "expenses",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/tags/%d", tag.ID),
		nil,
	)
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var j api.GetTagResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, tag.Name, j.Resource.Name)
	assert.Equal(t, tag.Sign, j.Resource.Sign)
	assert.Equal(t, tag.Kind, j.Resource.Kind)
	assert.Equal(t, tag.UserID, j.Resource.UserID)
}
