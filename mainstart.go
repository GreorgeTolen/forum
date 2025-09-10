package main

import (
	"fmt"
	"forum1/db"
	"forum1/handlers"
	"net/http"
)

func main() {
	// Подключаем базу данных
	err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка подключения к базе:", err)
		return
	}
	defer db.CloseDB()

	// Роуты
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/login_page/", handlers.LoginPage)
	http.HandleFunc("/register_page/", handlers.RegisterPage)
	http.HandleFunc("/profile_page/", handlers.ProfilePage)
	http.HandleFunc("/create_post_page/", handlers.CreatePostPage)
	http.HandleFunc("/board_page/", handlers.BoardPage)
	http.HandleFunc("/boards_list_page/", handlers.BoardsListPage)
	http.HandleFunc("/post_page/", handlers.PostPage)
	http.HandleFunc("/post_image/", handlers.PostImage)
	http.HandleFunc("/edit_post_page/", handlers.EditPostPage)
	// Комментарии и голоса
	http.HandleFunc("/vote_post/", handlers.VotePost)
	http.HandleFunc("/add_comment/", handlers.AddComment)
	http.HandleFunc("/delete_comment/", handlers.DeleteComment)
	http.HandleFunc("/vote_comment/", handlers.VoteComment)

	// Запуск сервера
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
