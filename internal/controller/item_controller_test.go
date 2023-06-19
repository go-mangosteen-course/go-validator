package controller

import (
	"encoding/json"
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/nav-inc/datetime"
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

func TestGetPagedItems(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := ItemController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items",
		nil,
	)
	u, _ := q.CreateUser(c, "1@qq.com")
	for i := 0; i < int(ctrl.PerPage*2); i++ {
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: time.Now(),
		}); err != nil {
			t.Error(err)
		}
	}
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.GetPagesItemsResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, ctrl.PerPage, int32(len(j.Resources)))
}

func TestGetBalance(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := ItemController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/balance",
		nil,
	)
	u, _ := q.CreateUser(c, "1@qq.com")
	for i := 0; i < int(ctrl.PerPage*2); i++ {
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: time.Now(),
		}); err != nil {
			t.Error(err)
		}
	}
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 10000*int(ctrl.PerPage*2), j.Expenses)
	assert.Equal(t, 0, j.Income)
	assert.Equal(t, -10000*int(ctrl.PerPage*2), j.Balance)
}

func TestGetBalanceWithTime(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := ItemController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/balance?happened_after="+url.QueryEscape("2020-01-01T00:00:00+0800")+
			"&happened_before="+url.QueryEscape("2020-01-02T00:00:00+0800"),
		nil,
	)
	u, _ := q.CreateUser(c, "1@qq.com")
	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2019-12-31T23:59:00+08:00", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2020-01-01T12:00:00+0800", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 3; i++ {
		d, _ := datetime.Parse("2020-01-10T12:00:00+0800", time.Local)
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.GetBalanceResponse
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		t.Error("json.Unmarshal fail", err)
	}
	assert.Equal(t, 10000*3, j.Expenses)
	assert.Equal(t, 0, j.Income)
	assert.Equal(t, -10000*3, j.Balance)
}

func TestGetSummary(t *testing.T) {
	done := setupTestCase(t)
	defer done(t)

	ctrl := ItemController{}
	ctrl.RegisterRoutes(r.Group("/api"))

	qs := url.Values{
		"happened_after":  []string{"2020-01-15T00:00:00+08:00"},
		"happened_before": []string{"2020-01-18T00:00:00+08:00"},
		"kind":            []string{"expenses"},
		"group_by":        []string{"happened_at"},
	}.Encode()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/summary?"+qs,
		nil,
	)
	u, _ := q.CreateUser(c, "1@qq.com")
	for i := 0; i < 3; i++ {
		d, _ := time.Parse(time.RFC3339, "2020-01-15T00:00:00+08:00")
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 3; i++ {
		d, _ := time.Parse(time.RFC3339, "2020-01-17T00:00:00+08:00")
		if _, err := q.CreateItem(c, queries.CreateItemParams{
			UserID:     u.ID,
			Amount:     10000,
			Kind:       "expenses",
			TagIds:     []int32{1},
			HappenedAt: d,
		}); err != nil {
			t.Error(err)
		}
	}
	signIn(t, u.ID, req)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var j api.GetSummaryResponse
	json.Unmarshal([]byte(body), &j)
	assert.Equal(t, 60000, j.Total)
	assert.Equal(t, 2, len(j.Groups))
	assert.Equal(t, 30000, j.Groups[0].Amount)
	assert.Equal(t, 30000, j.Groups[1].Amount)
	assert.Equal(t, "2020-01-15", j.Groups[0].HappenedAt)
	assert.Equal(t, "2020-01-17", j.Groups[1].HappenedAt)
}
