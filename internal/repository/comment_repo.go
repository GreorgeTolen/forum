package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, c *entity.Comment) (int64, error)
	GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	ForceDeleteComment(ctx context.Context, id int64) error
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

type commentRepository struct{ db *sql.DB }

func (r *commentRepository) CreateComment(ctx context.Context, c *entity.Comment) (int64, error) {
	return 0, nil
}
func (r *commentRepository) GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error) {
	return nil, nil
}
func (r *commentRepository) DeleteComment(ctx context.Context, id int64) error      { return nil }
func (r *commentRepository) ForceDeleteComment(ctx context.Context, id int64) error { return nil }
