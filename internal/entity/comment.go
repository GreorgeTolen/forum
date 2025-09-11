package entity

import (
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	AuthorID  int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     int
	Dislikes  int
}
