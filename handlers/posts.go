package handlers

import (
	"database/sql"
	"fmt"
	"forum1/db"
	"forum1/models"
	"forum1/utils"
	"io"
	"net/http"
	"strconv"
)

// Страница создания поста
func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	// Получаем ID пользователя из БД
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при проверке пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		boardIDStr := r.FormValue("board_id")
		title := r.FormValue("title")
		content := r.FormValue("content")

		// читаем файл изображения если он есть
		var imageBytes []byte
		r.ParseMultipartForm(10 << 20) // 10MB
		file, _, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			buf, _ := io.ReadAll(file)
			imageBytes = buf
		}

		if boardIDStr == "" || title == "" || content == "" {
			http.Error(w, "Заполните все поля", http.StatusBadRequest)
			return
		}

		boardID, err := strconv.Atoi(boardIDStr)
		if err != nil {
			http.Error(w, "Неверный ID доски", http.StatusBadRequest)
			return
		}

		post := &models.Post{
			BoardID:   boardID,
			Title:     title,
			Content:   content,
			AuthorID:  userID,
			ImageData: imageBytes,
		}

		err = models.CreatePost(post)
		if err != nil {
			http.Error(w, "Ошибка при сохранении поста: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Редирект на страницу доски
		http.Redirect(w, r, fmt.Sprintf("/board_page?id=%d", boardID), http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, "create_post_page.html", Boards)
}

// Просмотр поста
func PostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// уникальный просмотр (по username)
	if cookie, cerr := r.Cookie("user"); cerr == nil {
		var uid int
		if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&uid); err == nil {
			_, _ = db.DB.Exec(`INSERT INTO post_views (post_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, id, uid)
		}
	}
	var views int
	_ = db.DB.QueryRow(`SELECT COUNT(*) FROM post_views WHERE post_id=$1`, id).Scan(&views)

	// комментарии
	comments, _ := models.GetCommentsByPost(id)
	post.Comments = comments

	// лайки/дизлайки поста
	_ = db.DB.QueryRow(`SELECT COALESCE(SUM(CASE WHEN value=1 THEN 1 ELSE 0 END),0) AS likes,
		COALESCE(SUM(CASE WHEN value=-1 THEN 1 ELSE 0 END),0) AS dislikes
		FROM post_votes WHERE post_id=$1`, id).Scan(&post.Likes, &post.Dislikes)

	// прокинем просмотры через фейковое поле LinkURL (не трогаю структуры шаблонов) — нет, лучше в заголовок
	utils.RenderTemplate(w, "post_page.html", post)
}

// Отдача изображения поста из БД
func PostImage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}
	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(post.ImageData) == 0 {
		http.NotFound(w, r)
		return
	}
	contentType := http.DetectContentType(post.ImageData)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(post.ImageData)
}

// Редактирование поста
func EditPostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title != "" {
			post.Title = title
		}
		if content != "" {
			post.Content = content
		}

		err = models.UpdatePost(post)
		if err != nil {
			http.Error(w, "Ошибка при обновлении поста: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post_page?id=%d", post.ID), http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, "edit_post_page.html", post)
}

// Голос за пост
func VotePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("post_id")
	valStr := r.FormValue("value")
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	id, _ := strconv.Atoi(idStr)
	value, _ := strconv.Atoi(valStr)
	_, err = db.DB.Exec(`INSERT INTO post_votes (post_id, user_id, value) VALUES ($1,$2,$3) ON CONFLICT (post_id,user_id) DO UPDATE SET value=EXCLUDED.value`, id, userID, value)
	if err != nil {
		http.Error(w, "Ошибка голосования: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+idStr, http.StatusSeeOther)
}

// Добавить комментарий
func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")
	postID, _ := strconv.Atoi(postIDStr)
	c := &models.Comment{PostID: postID, AuthorID: userID, Content: content}
	if err := models.CreateComment(c); err != nil {
		http.Error(w, "Ошибка добавления комментария: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}

// Удалить комментарий (может автор поста или автор коммента)
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	postID, _ := strconv.Atoi(postIDStr)
	commentID, _ := strconv.Atoi(commentIDStr)
	post, err := models.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}
	if post.AuthorID == userID {
		_ = models.ForceDeleteComment(commentID)
	} else {
		_ = models.DeleteComment(commentID, userID)
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}

// Голос за комментарий
func VoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	valueStr := r.FormValue("value")
	commentID, _ := strconv.Atoi(commentIDStr)
	value, _ := strconv.Atoi(valueStr)
	if err := models.SetCommentVote(commentID, userID, value); err != nil {
		http.Error(w, "Ошибка голосования: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}
