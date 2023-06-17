package forum

import (
	"context"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
)

type ForumUsecase interface {
	GetUser(ctx context.Context, nickname string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)

	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
}

type ForumRepository interface {
	GetUser(ctx context.Context, nickname string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)

	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
}
