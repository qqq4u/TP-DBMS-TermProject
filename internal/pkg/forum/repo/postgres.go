package repo

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
)

type ForumRepository struct {
	conn   *pgxpool.Pool
	Status models.Status
}

func NewForumRepo(connection *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		conn:   connection,
		Status: models.Status{},
	}
}

const (
	GetUserByNickname  = `SELECT email, fullname, nickname, about FROM "user" WHERE nickname=$1 LIMIT 1;`
	GetUsersOnConflict = `SELECT email, fullname, nickname, about FROM "user" WHERE email = $1 or nickname = $2`
	CreateUser         = `INSERT INTO "user" (email, fullname, nickname, about) VALUES ($1, $2, $3, $4) RETURNING nickname;`
	UpdateUser         = `UPDATE "user" SET fullname=$1, email=$2, about=$3 WHERE nickname = $4 RETURNING nickname, fullname, about, email;`
	CheckIfUserExists  = `SELECT nickname FROM "user" WHERE nickname =  $1`
	CreateForum        = `INSERT INTO "forum" (title, "user", slug) VALUES ($1, $2, $3) RETURNING slug;`
	GetForumBySlug     = `SELECT title, "user", slug, posts, threads FROM "forum" WHERE slug = $1;`
)

const (
	DuplicatesKeyError = "23505"
	ForeingKeyError    = "23503"
)

func (r *ForumRepository) GetUser(ctx context.Context, nickname string) (models.User, error) {
	var resultUser models.User

	row := r.conn.QueryRow(ctx, GetUserByNickname, nickname)

	err := row.Scan(&resultUser.Email, &resultUser.Fullname, &resultUser.Nickname, &resultUser.About)
	if err != nil {
		return models.User{}, models.ErrorNotFound
	}

	return resultUser, nil
}

func (r *ForumRepository) getUsersOnConflict(user models.User) []models.User {
	results := []models.User{}

	rows, _ := r.conn.Query(context.Background(), GetUsersOnConflict, user.Email, user.Nickname)
	defer rows.Close()

	for rows.Next() {
		result := models.User{}
		rows.Scan(&result.Email, &result.Fullname, &result.Nickname, &result.About)
		results = append(results, result)
	}
	return results
}

func (r *ForumRepository) CreateUser(ctx context.Context, user models.User) ([]models.User, error) {
	_, err := r.conn.Exec(ctx, CreateUser, user.Email, user.Fullname, user.Nickname, user.About)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case DuplicatesKeyError:
				us := r.getUsersOnConflict(user)
				return us, models.ErrorConflict
			}
		}
	}

	r.Status.UsersCount++
	return []models.User{user}, nil
}

func (r *ForumRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	updatedUser, err := r.GetUser(ctx, user.Nickname)
	if err != nil {
		return updatedUser, err
	}
	if user.Fullname != "" {
		updatedUser.Fullname = user.Fullname
	}
	if user.Email != "" {
		updatedUser.Email = user.Email
	}
	if user.About != "" {
		updatedUser.About = user.About
	}
	rows := r.conn.QueryRow(ctx, UpdateUser, updatedUser.Fullname, updatedUser.Email, updatedUser.About, updatedUser.Nickname)
	err = rows.Scan(&updatedUser.Nickname, &updatedUser.Fullname, &updatedUser.About, &updatedUser.Email)
	if pqError, ok := err.(*pgconn.PgError); ok {
		switch pqError.Code {
		case DuplicatesKeyError:
			return user, models.ErrorConflict
		case ForeingKeyError:
			return user, models.ErrorNotFound
		}
	}

	return updatedUser, nil
}

func (r *ForumRepository) checkIfUserExists(ctx context.Context, nickname string) (models.User, error) {
	result := models.User{}
	rows := r.conn.QueryRow(ctx, CheckIfUserExists, nickname)

	err := rows.Scan(&result.Nickname)
	if err != nil {
		return result, models.ErrorNotFound
	}
	return result, nil
}

func (r *ForumRepository) GetForumBySlug(ctx context.Context, slug string) (models.Forum, error) {
	result := models.Forum{}

	row := r.conn.QueryRow(ctx, GetForumBySlug, slug)
	err := row.Scan(&result.Title, &result.User, &result.Slug, &result.Posts, &result.Threads)

	if err != nil {
		return models.Forum{}, models.ErrorNotFound
	}

	return result, nil
}

func (r *ForumRepository) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	user, err := r.checkIfUserExists(ctx, forum.User)
	if err != nil {
		return models.Forum{}, models.ErrorNotFound
	}

	row := r.conn.QueryRow(ctx, CreateForum, forum.Title, user.Nickname, forum.Slug)
	err = row.Scan(&forum.Slug)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case DuplicatesKeyError:
				result, _ := r.GetForumBySlug(ctx, forum.Slug)
				return result, models.ErrorConflict
			}
		}
	}

	forum.User = user.Nickname
	r.Status.ForumsCount++
	return forum, nil
}

func (r *ForumRepository) GetForum(ctx context.Context, slug string) (models.Forum, error) {
	var resultForum models.Forum

	row := r.conn.QueryRow(ctx, GetForumBySlug, slug)

	err := row.Scan(&resultForum.Title, &resultForum.User, &resultForum.Slug, &resultForum.Posts, &resultForum.Threads)
	if err != nil {
		return models.Forum{}, models.ErrorNotFound
	}

	return resultForum, nil
}
