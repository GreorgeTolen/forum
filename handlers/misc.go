package handlers

import (
	"forum1/utils"
	"net/http"
)

func SearchPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "search_page.html", nil)
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
