package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/qqq4u/TP-DBMS-TermProject/internal/models"
	"strings"
	"time"
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
	GetUserByNickname       = `SELECT email, fullname, nickname, about FROM "user" WHERE nickname=$1 LIMIT 1;`
	GetUsersOnConflict      = `SELECT email, fullname, nickname, about FROM "user" WHERE email = $1 or nickname = $2`
	CreateUser              = `INSERT INTO "user" (email, fullname, nickname, about) VALUES ($1, $2, $3, $4) RETURNING nickname;`
	UpdateUser              = `UPDATE "user" SET fullname=$1, email=$2, about=$3 WHERE nickname = $4 RETURNING nickname, fullname, about, email;`
	CheckIfUserExists       = `SELECT nickname FROM "user" WHERE nickname =  $1`
	CheckIfForumExists      = `SELECT slug FROM "forum" WHERE slug = $1;`
	CreateForum             = `INSERT INTO "forum" (title, "user", slug) VALUES ($1, $2, $3) RETURNING slug;`
	GetForumBySlug          = `SELECT title, "user", slug, posts, threads FROM "forum" WHERE slug = $1 limit 1;`
	GetThreadBySlug         = `SELECT id, "author", message, title, created, forum, slug, votes FROM "thread" WHERE slug = $1 limit 1;`
	CreateThread            = `INSERT INTO "thread" (author, message, title, created, forum, slug, votes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	InserIntoUserForum      = `INSERT INTO forum_user(Nickname, Forum) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	GetThreadsWithSinceDesc = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 AND created <= $2 ORDER BY created DESC LIMIT $3;`
	GetThreadsWithSinceAsc  = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 and created >= $2 ORDER BY created ASC LIMIT $3;`
	GetThreadsDesc          = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 ORDER BY created DESC LIMIT $2;`
	GetThreadsAsc           = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 ORDER BY created ASC LIMIT $2;`
	SelectThreadById        = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE id=$1 LIMIT 1;`
	SelectThreadBySlug      = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE slug=$1 LIMIT 1;`
	InsertPostsStartQuery   = `INSERT INTO "post"(author, created, forum, message, parent, thread) VALUES`
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

func (r *ForumRepository) checkIfForumExists(ctx context.Context, slug string) (models.Forum, error) {
	result := models.Forum{}
	row := r.conn.QueryRow(ctx, CheckIfForumExists, slug)

	err := row.Scan(&result.Slug)

	if err != nil {
		return result, models.ErrorNotFound
	}
	return result, nil
}

func (r *ForumRepository) GetThread(ctx context.Context, slug string) (models.Thread, error) {
	var resultThread models.Thread

	row := r.conn.QueryRow(ctx, GetThreadBySlug, slug)

	err := row.Scan(&resultThread.ID, &resultThread.Author, &resultThread.Message, &resultThread.Title, &resultThread.Created, &resultThread.Forum, &resultThread.Slug, &resultThread.Votes)
	if err != nil {
		return models.Thread{}, models.ErrorNotFound
	}

	return resultThread, nil
}

func (r *ForumRepository) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	user, err := r.checkIfUserExists(ctx, thread.Author)
	if err != nil {
		return models.Thread{}, models.ErrorNotFound
	}

	forum, err := r.checkIfForumExists(ctx, thread.Forum)
	if err != nil {
		return models.Thread{}, models.ErrorNotFound
	}

	thread.Forum = forum.Slug
	thread.Author = user.Nickname

	if thread.Slug != "" {
		threadConflict, err := r.GetThreadBySlug(ctx, thread.Slug)
		if !errors.Is(err, models.ErrorNotFound) {
			return threadConflict, models.ErrorConflict
		}
	}

	row := r.conn.QueryRow(ctx, CreateThread, thread.Author, thread.Message, thread.Title, thread.Created, thread.Forum, thread.Slug, 0)
	err = row.Scan(&thread.ID)
	if err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok {
			switch pqError.Code {
			case DuplicatesKeyError:
				result, _ := r.GetThread(ctx, thread.Slug)
				return result, models.ErrorConflict
			case ForeingKeyError:
				result, _ := r.GetThread(ctx, thread.Slug)
				return result, models.ErrorConflict
			}

		}
	}

	r.Status.ThreadsCount++

	_, _ = r.conn.Exec(ctx, InserIntoUserForum, thread.Author, thread.Forum)

	return thread, nil
}

func (r *ForumRepository) GetThreads(ctx context.Context, slug, limit, since, desc string) ([]models.Thread, error) {
	threads := make([]models.Thread, 0)
	if since != "" {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetThreadsWithSinceDesc, slug, since, limit)
			defer rows.Close()
			for rows.Next() {
				tmpThread := models.Thread{}
				err := rows.Scan(&tmpThread.ID, &tmpThread.Title, &tmpThread.Author, &tmpThread.Forum, &tmpThread.Message,
					&tmpThread.Votes, &tmpThread.Slug, &tmpThread.Created)
				if err != nil {
					continue
				}
				threads = append(threads, tmpThread)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetThreadsWithSinceAsc, slug, since, limit)
			defer rows.Close()
			for rows.Next() {
				tmpThread := models.Thread{}
				err := rows.Scan(&tmpThread.ID, &tmpThread.Title, &tmpThread.Author, &tmpThread.Forum, &tmpThread.Message,
					&tmpThread.Votes, &tmpThread.Slug, &tmpThread.Created)
				if err != nil {
					continue
				}
				threads = append(threads, tmpThread)
			}
		}
	} else {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetThreadsDesc, slug, limit)
			defer rows.Close()
			for rows.Next() {
				tmpThread := models.Thread{}
				err := rows.Scan(&tmpThread.ID, &tmpThread.Title, &tmpThread.Author, &tmpThread.Forum, &tmpThread.Message,
					&tmpThread.Votes, &tmpThread.Slug, &tmpThread.Created)
				if err != nil {
					continue
				}
				threads = append(threads, tmpThread)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetThreadsAsc, slug, limit)
			defer rows.Close()
			for rows.Next() {
				tmpThread := models.Thread{}
				err := rows.Scan(&tmpThread.ID, &tmpThread.Title, &tmpThread.Author, &tmpThread.Forum, &tmpThread.Message,
					&tmpThread.Votes, &tmpThread.Slug, &tmpThread.Created)
				if err != nil {
					continue
				}
				threads = append(threads, tmpThread)
			}
		}
	}
	return threads, nil
}

func (r *ForumRepository) GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error) {
	thread := models.Thread{}
	row := r.conn.QueryRow(ctx, SelectThreadBySlug, slug)
	err := row.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum,
		&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrorNotFound
	}
	return thread, models.ErrorConflict
}
func (r *ForumRepository) GetThreadById(ctx context.Context, id int) (models.Thread, error) {
	thread := models.Thread{}
	row := r.conn.QueryRow(ctx, SelectThreadById, id)
	err := row.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum,
		&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrorNotFound
	}
	return thread, models.ErrorConflict
}
func (r *ForumRepository) CreatePosts(ctx context.Context, posts models.PostsList, thread models.Thread) (models.PostsList, error) {
	InsertPosts := InsertPostsStartQuery

	var values []interface{}
	created := time.Now()
	for i, post := range posts {
		value := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		InsertPosts += value
		values = append(values, post.Author)
		values = append(values, created)
		values = append(values, thread.Forum)
		values = append(values, post.Message)
		values = append(values, post.Parent)
		values = append(values, thread.ID)
	}

	InsertPosts = strings.TrimSuffix(InsertPosts, ",")
	InsertPosts += ` RETURNING id, created, forum, isEdited, thread;`

	rows, err := r.conn.Query(ctx, InsertPosts, values...)
	if err != nil {
		return nil, models.ErrorConflict
	}
	defer rows.Close()

	for i := range posts {
		if rows.Next() {
			err := rows.Scan(&posts[i].ID, &posts[i].Created, &posts[i].Forum, &posts[i].IsEdited, &posts[i].Thread)
			if err != nil {
				return nil, err
			}
		}
	}
	return posts, nil
}
