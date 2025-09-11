package handlers

import (
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/utils"
	"net/http"
)

var Boards = []entity.Board{
	{ID: 1, Slug: "schedule", Title: "Расписание", Description: "Обсуждаем расписание этого года"},
	{ID: 2, Slug: "games", Title: "Игры", Description: "Все о видеоиграх, консолях и ПК"},
	{ID: 3, Slug: "offtopic", Title: "Оффтопик", Description: "Свободное общение на любые темы"},
	{ID: 4, Slug: "news", Title: "Новости", Description: "Обсуждение последних новостей"},
	{ID: 5, Slug: "reviews", Title: "Рецензии", Description: "Ваши обзоры на фильмы, игры и книги"},
}

func BoardPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "Не указана доска", http.StatusBadRequest)
		return
	}

	var board *entity.Board
	for _, b := range Boards {
		if b.Slug == slug {
			board = &b
			break
		}
	}
	if board == nil {
		http.Error(w, "Доска не найдена", http.StatusNotFound)
		return
	}

	posts, err := models.GetPostsByBoard(board.ID)
	if err != nil {
		http.Error(w, "Ошибка при получении постов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var boardPosts []entity.Post
	for _, p := range posts {
		if p.BoardID == board.ID {
			boardPosts = append(boardPosts, p)
		}
	}

	data := struct {
		Board entity.Board
		Posts []entity.Post
	}{
		Board: *board,
		Posts: boardPosts,
	}

	utils.RenderTemplate(w, "board_page.html", data)
}

func BoardsListPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "boards_list_page.html", Boards)
}

func BoardsSearchPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "boards_search_page.html", nil)
}
