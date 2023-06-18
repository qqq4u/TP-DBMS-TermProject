package usecase

import (
	"context"
	"errors"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/pkg/forum"
	"strconv"
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
func (u *ForumUsecase) GetForum(ctx context.Context, slug string) (models.Forum, error) {
	return u.repo.GetForum(ctx, slug)
}

func (u *ForumUsecase) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	return u.repo.CreateThread(ctx, thread)
}

func (u *ForumUsecase) GetThreads(ctx context.Context, slug, limit, since, desc string) ([]models.Thread, error) {
	_, err := u.repo.GetForum(ctx, slug)
	if errors.Is(err, models.ErrorNotFound) {
		return nil, err
	}

	return u.repo.GetThreads(ctx, slug, limit, since, desc)
}

func (u *ForumUsecase) CheckThreadByIdOrSlug(ctx context.Context, slugOrId string) (models.Thread, error) {
	intValue, err := strconv.Atoi(slugOrId)
	if err != nil {
		return u.repo.GetThreadBySlug(ctx, slugOrId)
	} else {
		return u.repo.GetThreadById(ctx, intValue)
	}
}
func (u *ForumUsecase) CreatePosts(ctx context.Context, posts models.PostsList, thread models.Thread) (models.PostsList, error) {
	return u.repo.CreatePosts(ctx, posts, thread)
}
