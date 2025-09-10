package handlers

import (
	"fmt"
	"forum1/models"
	"forum1/utils"
	"net/http"
	"strconv"
)

var Boards = []models.Board{
	{ID: 1, Title: "Расписание", Description: "Обсуждаем расписание этого года"},
	{ID: 2, Title: "Игры", Description: "Все о видеоиграх, консолях и ПК"},
	{ID: 3, Title: "Оффтопик", Description: "Свободное общение на любые темы"},
	{ID: 4, Title: "Новости", Description: "Обсуждение последних новостей"},
	{ID: 5, Title: "Рецензии", Description: "Ваши обзоры на фильмы, игры и книги"},
}

func BoardPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Не указана доска", http.StatusBadRequest)
		return
	}

	var board *models.Board
	for _, b := range Boards {
		if fmt.Sprint(b.ID) == id {
			board = &b
			break
		}
	}
	if board == nil {
		http.Error(w, "Доска не найдена", http.StatusNotFound)
		return
	}

	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Ошибка при получении постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var boardPosts []models.Post
	boardID, _ := strconv.Atoi(id)
	for _, p := range posts {
		if p.BoardID == boardID {
			boardPosts = append(boardPosts, p)
		}
	}

	data := struct {
		Board models.Board
		Posts []models.Post
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
