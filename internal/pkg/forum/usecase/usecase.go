package usecase

import (
	"context"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum"
)

type ForumUsecase struct {
	repo forum.ForumRepository
}

func NewForumUsecase(repo forum.ForumRepository) *ForumUsecase {
	return &ForumUsecase{
		repo: repo,
	}
}

func (u *ForumUsecase) GetUser(ctx context.Context, nickname string) (models.User, error) {
	return u.repo.GetUser(ctx, nickname)
}

func (u *ForumUsecase) CreateUser(ctx context.Context, user models.User) ([]models.User, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *ForumUsecase) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	return u.repo.UpdateUser(ctx, user)
}

func (u *ForumUsecase) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	return u.repo.CreateForum(ctx, forum)
}
