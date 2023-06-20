package forum

import (
	"context"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
)

type ForumUsecase interface {
	GetUser(ctx context.Context, nickname string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	GetUsers(ctx context.Context, slug, limit, since, desc string) ([]models.User, error)

	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, slug string) (models.Forum, error)

	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	UpdateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThreads(ctx context.Context, slug, limit, since, desc string) ([]models.Thread, error)

	CheckThreadByIdOrSlug(ctx context.Context, slugOrId string) (models.Thread, error)
	CreatePosts(ctx context.Context, posts models.PostsList, thread models.Thread) (models.PostsList, error)

	Vote(ctx context.Context, vote models.Vote) error

	GetPost(ctx context.Context, id string, related []string) (models.PostFull, error)
	GetThreadPosts(ctx context.Context, limit, since, desc, sort string, threadId int) ([]models.Post, error)
	UpdatePost(ctx context.Context, post models.PostUpdate) (models.Post, error)

	GetStatus() models.Status
	Clear()
}

type ForumRepository interface {
	GetUser(ctx context.Context, nickname string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)

	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetForum(ctx context.Context, slug string) (models.Forum, error)
	GetUsers(ctx context.Context, slug, limit, since, desc string) ([]models.User, error)

	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	UpdateThread(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetThreads(ctx context.Context, slug, limit, since, desc string) ([]models.Thread, error)

	GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error)
	GetThreadById(ctx context.Context, id int) (models.Thread, error)
	CreatePosts(ctx context.Context, posts models.PostsList, thread models.Thread) (models.PostsList, error)

	Vote(ctx context.Context, vote models.Vote) error

	GetPost(ctx context.Context, id int, related []string) (models.PostFull, error)
	GetThreadPosts(ctx context.Context, limit, since, desc, sort string, threadId int) ([]models.Post, error)
	UpdatePost(ctx context.Context, post models.PostUpdate) (models.Post, error)

	GetStatus() models.Status
	Clear()
}
