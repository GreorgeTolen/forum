package app

import (
	"fmt"
	"forum1/db"
	_ "forum1/docs"
	"forum1/internal/handlers"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Run запускает сервер форума
func Run() {

	// Инициализация базы данных
	err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка подключения к базе:", err)
		return
	}
	defer db.CloseDB()

	// Основные страницы
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

	// Работа с постами и комментариями
	http.HandleFunc("/vote_post/", handlers.VotePost)
	http.HandleFunc("/add_comment/", handlers.AddComment)
	http.HandleFunc("/delete_comment/", handlers.DeleteComment)
	http.HandleFunc("/vote_comment/", handlers.VoteComment)

	// Swagger UI доступен по адресу: http://localhost:8080/swagger/index.html
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
