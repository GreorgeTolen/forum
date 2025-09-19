package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, c *entity.Comment) (int64, error)
	GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error)
	DeleteComment(ctx context.Context, id int64, requesterID int64) error
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

type commentService struct{ repo repository.CommentRepository }

func (s *commentService) CreateComment(ctx context.Context, c *entity.Comment) (int64, error) {
	if c.PostID == 0 || c.AuthorID == 0 || c.Content == "" {
		return 0, errors.New("invalid input")
	}
	return s.repo.CreateComment(ctx, c)
}
func (s *commentService) GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error) {
	if postID == 0 {
		return nil, errors.New("post id required")
	}
	return s.repo.GetCommentsByPost(ctx, postID)
}
func (s *commentService) DeleteComment(ctx context.Context, id int64, requesterID int64) error {
	if id == 0 {
		return errors.New("id required")
	}
	return s.repo.DeleteComment(ctx, id)
}
