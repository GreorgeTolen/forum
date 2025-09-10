package handlers

import (
	"forum1/models"
	"forum1/utils"
	"net/http"
	"time"
)

var Posts = []models.Post{
	{ID: 1, BoardID: 1, Title: "Добро пожаловать!", Content: "Это первый пост на форуме.", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, BoardID: 4, Title: "Новости форума", Content: "Скоро будут новые функции!", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Ошибка при получении постов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts []models.Post
	}{
		Posts: posts,
	}
	utils.RenderTemplate(w, "home_page.html", data)
}
