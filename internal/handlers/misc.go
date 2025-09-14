package handlers

import (
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/utils"
	"net/http"
	"strings"
)

func SearchPage(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if query == "" {
		utils.RenderTemplate(w, "search_page.html", nil)
		return
	}

	posts, _ := models.SearchPosts(query)
	boards, _ := models.SearchBoards(query)

	data := struct {
		Query  string
		Posts  []entity.Post
		Boards []entity.Board
	}{
		Query:  query,
		Posts:  posts,
		Boards: boards,
	}

	utils.RenderTemplate(w, "search_page.html", data)
}

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "settings_page.html", nil)
}

func NotificationsPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "notifications_page.html", nil)
}

func MessagesPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "messages_page.html", nil)
}
