package handler

import (
	"forum1/internal/service"
	"forum1/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PageHandler struct {
	posts  service.PostService
	boards service.BoardService
}

func NewPageHandler(p service.PostService, b service.BoardService) *PageHandler {
	return &PageHandler{posts: p, boards: b}
}

func (h *PageHandler) HomePageHTML(w http.ResponseWriter, r *http.Request) {
	// Load boards for sidebar/home
	boards, _ := h.boards.List(r.Context())
	data := map[string]interface{}{
		"Boards": boards,
	}
	utils.RenderTemplate(w, "home_page.html", data)
}

func (h *PageHandler) BoardsListPage(w http.ResponseWriter, r *http.Request) {
	boards, _ := h.boards.List(r.Context())
	data := map[string]interface{}{"Boards": boards}
	utils.RenderTemplate(w, "boards_list_page.html", data)
}

func (h *PageHandler) BoardPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	b, err := h.boards.GetBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	posts, _ := h.posts.GetPostsByBoard(r.Context(), int64(b.ID))
	data := map[string]interface{}{"Board": b, "Posts": posts}
	utils.RenderTemplate(w, "board_page.html", data)
}

func (h *PageHandler) PostPageHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	post, _ := h.posts.GetPostByID(r.Context(), id)
	data := map[string]interface{}{"Post": post}
	utils.RenderTemplate(w, "post_page.html", data)
}

func (h *PageHandler) ProfilePageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "profile_page.html", map[string]interface{}{})
}

func (h *PageHandler) LoginPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login_page.html", map[string]interface{}{})
}

func (h *PageHandler) RegisterPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "register_page.html", map[string]interface{}{})
}

func (h *PageHandler) CreatePostPageHTML(w http.ResponseWriter, r *http.Request) {
	boards, _ := h.boards.List(r.Context())
	// Template expects to range over root (.)
	utils.RenderTemplate(w, "create_post_page.html", boards)
}

func (h *PageHandler) BoardsSearchPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "boards_search_page.html", map[string]interface{}{})
}

func (h *PageHandler) SearchPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "search_page.html", map[string]interface{}{})
}

func (h *PageHandler) SettingsPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "settings_page.html", map[string]interface{}{})
}

func (h *PageHandler) MessagesPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "messages_page.html", map[string]interface{}{})
}

func (h *PageHandler) NotificationsPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "notifications_page.html", map[string]interface{}{})
}
