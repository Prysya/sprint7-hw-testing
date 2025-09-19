package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[city])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		if v.count == 0 {
			assert.Equal(t, "", body)
		} else {
			list := strings.Split(body, ",")
			assert.Equal(t, v.want, len(list))
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"ФаСоль", 0},
		{"кофе", 2},
		{"КофЕ", 2},
		{"вилка", 1},
		{"ВИЛКА", 1},
		{"ложка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		if v.want == 0 {
			assert.Equal(t, "", body)
		} else {
			list := strings.Split(body, ",")
			assert.Equal(t, v.want, len(list))

			for _, name := range list {
				assert.Contains(t, strings.ToLower(name), strings.ToLower(v.search))
			}
		}
	}
}
