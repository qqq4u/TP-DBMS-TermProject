package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	GetUserByNickname                     = `SELECT email, fullname, nickname, about FROM "user" WHERE nickname=$1 LIMIT 1;`
	GetUsersOnConflict                    = `SELECT email, fullname, nickname, about FROM "user" WHERE email = $1 or nickname = $2`
	CreateUser                            = `INSERT INTO "user" (email, fullname, nickname, about) VALUES ($1, $2, $3, $4) RETURNING nickname;`
	UpdateUser                            = `UPDATE "user" SET fullname=$1, email=$2, about=$3 WHERE nickname = $4 RETURNING nickname, fullname, about, email;`
	CheckIfUserExists                     = `SELECT nickname FROM "user" WHERE nickname =  $1`
	CheckIfForumExists                    = `SELECT slug FROM "forum" WHERE slug = $1;`
	CreateForum                           = `INSERT INTO "forum" (title, "user", slug) VALUES ($1, $2, $3) RETURNING slug;`
	GetForumBySlug                        = `SELECT title, "user", slug, posts, threads FROM "forum" WHERE slug = $1 LIMIT 1;`
	GetThreadBySlug                       = `SELECT id, author, message, title, created, forum, slug, votes FROM "thread" WHERE slug = $1 limit 1;`
	CreateThread                          = `INSERT INTO "thread" (author, message, title, created, forum, slug, votes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	GetThreadsWithSinceDesc               = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 AND created <= $2 ORDER BY created DESC LIMIT $3;`
	GetThreadsWithSinceAsc                = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 AND created >= $2 ORDER BY created ASC LIMIT $3;`
	GetThreadsDesc                        = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 ORDER BY created DESC LIMIT $2;`
	GetThreadsAsc                         = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE forum=$1 ORDER BY created ASC LIMIT $2;`
	SelectThreadById                      = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE id=$1 LIMIT 1;`
	SelectThreadBySlug                    = `SELECT id, title, author, forum, message, votes, slug, created FROM "thread" WHERE slug=$1 LIMIT 1;`
	InsertPostsStartQuery                 = `INSERT INTO "post"(author, created, forum, message, parent, thread) VALUES`
	UpdateVote                            = `UPDATE "vote" SET voice=$1 WHERE author=$2 AND thread=$3;`
	InsertVote                            = `INSERT INTO "vote"(author, voice, thread) VALUES ($1, $2, $3);`
	SelectPostById                        = `SELECT author, message, created, forum, isedited, parent, thread FROM "post" WHERE id = $1 LIMIT 1;`
	GetPostsWithSinceDesc                 = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread=$1 ORDER BY id DESC limit $2;`
	GetPostsWithSinceAsc                  = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread=$1 ORDER BY id ASC limit $2;`
	GetPostsDesc                          = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread=$1 AND id < $2 ORDER BY id DESC LIMIT $3;`
	GetPostsAsc                           = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread=$1 AND id > $2 ORDER BY id ASC LIMIT $3;`
	GetPostsTreeDesc                      = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread = $1 ORDER BY path, id DESC`
	GetPostsTreeAsc                       = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread = $1 ORDER BY path, id ASC`
	GetPostsTreeWithLimitDesc             = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread = $1 ORDER BY path DESC, id DESC LIMIT $2`
	GetPostsTreeWithLimitAsc              = `SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE thread = $1 ORDER BY path, id ASC LIMIT $2`
	GetPostsTreeWithLimitWithSinceDesc    = `SELECT "post".id, "post".author, "post".created, "post".forum, "post".isedited, "post".message, "post".parent, "post".thread FROM "post" JOIN "post" parent ON parent.id = $2 WHERE "post".path < parent.path AND "post".thread = $1 ORDER BY "post".path DESC, "post".id DESC LIMIT $3`
	GetPostsTreeWithLimitWithSinceAsc     = `SELECT "post".id, "post".author, "post".created, "post".forum, "post".isedited, "post".message, "post".parent, "post".thread FROM "post" JOIN "post" parent ON parent.id = $2 WHERE "post".path > parent.path AND "post".thread = $1 ORDER BY "post".path ASC, "post".id ASC LIMIT $3`
	SelectTreeSinceNilDesc                = `SELECT "post".id, "post".author, "post".created, "post".forum, "post".isedited, "post".message, "post".parent, "post".thread FROM "post" JOIN "post" parent ON parent.id = $2 WHERE "post".path < parent.path AND "post".thread = $1 ORDER BY "post".path DESC, "post".id DESC`
	SelectTreeSinceNilDescNil             = `SELECT "post".id, "post".author, "post".created, "post".forum, "post".isedited, "post".message, "post".parent, "post".thread FROM "post" JOIN "post" parent ON parent.id = $2 WHERE "post".path > parent.path AND "post".thread = $1 ORDER BY "post".path ASC, "post".id ASC`
	UpdateThreadWithoutIdentifier         = "UPDATE thread SET title=coalesce(nullif($1, ''), title), message=coalesce(nullif($2, ''), message) WHERE %s RETURNING *"
	GetUsersWithSinceDesc                 = `SELECT nickname, fullname, about, email FROM "user_forum" WHERE slug=$1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3;`
	GetUsersWithSinceAsc                  = `SELECT nickname, fullname, about, email FROM "user_forum" WHERE slug=$1 AND nickname > $2 ORDER BY nickname ASC LIMIT $3;`
	GetUsersDesc                          = `SELECT nickname, fullname, about, email FROM "user_forum" WHERE slug=$1 ORDER BY nickname DESC LIMIT $2;`
	GetUsersAsc                           = `SELECT nickname, fullname, about, email FROM "user_forum" WHERE slug=$1 ORDER BY nickname ASC LIMIT $2;`
	UpdatePostMessage                     = `UPDATE "post" SET message=coalesce(nullif($1, ''), message), isedited = CASE WHEN $1 = '' OR message = $1 THEN isedited ELSE TRUE END WHERE id=$2 RETURNING *`
	GetThreadFromPost                     = `SELECT thread FROM "post" WHERE id = $1;`
	DESTROY_DATABASE_DONT_TOCUH_DANGEROUS = `TRUNCATE TABLE "user", "forum", "thread", "post", "vote", "user_forum" CASCADE;`
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
	results := make([]models.User, 0)

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

	thread, _ = r.GetThreadById(ctx, thread.ID)
	r.Status.ThreadsCount++

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
		if _, err := r.checkIfUserExists(ctx, post.Author); err != nil {
			return []models.Post{}, models.ErrorNotFound
		}
		values = append(values, created)
		values = append(values, thread.Forum)
		values = append(values, post.Message)
		values = append(values, post.Parent)
		values = append(values, thread.ID)
		if post.Parent != 0 {
			prevParent := 0
			row := r.conn.QueryRow(ctx, GetThreadFromPost, post.Parent)
			err := row.Scan(&prevParent)
			if err != nil || prevParent != thread.ID {
				return []models.Post{}, models.ErrorConflict
			}
		}
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
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	r.Status.PostsCount++
	return posts, nil
}

func (r *ForumRepository) Vote(ctx context.Context, vote models.Vote) error {
	_, err := r.checkIfUserExists(ctx, vote.Nickname)
	if err != nil {
		return models.ErrorNotFound
	}
	_, err = r.conn.Exec(ctx, InsertVote, vote.Nickname, vote.Voice, vote.Thread)
	if pqError, ok := err.(*pgconn.PgError); ok {
		switch pqError.Code {
		case DuplicatesKeyError:
			_, _ = r.conn.Exec(ctx, UpdateVote, vote.Voice, vote.Nickname, vote.Thread)
		}
	}
	return nil
}

func (r *ForumRepository) GetPost(ctx context.Context, id int, related []string) (models.PostFull, error) {
	postTmp := models.Post{}
	postResult := models.PostFull{Author: nil, Forum: nil, Post: models.Post{}, Thread: nil}

	postTmp.ID = id

	row := r.conn.QueryRow(ctx, SelectPostById, id)
	err := row.Scan(&postTmp.Author, &postTmp.Message, &postTmp.Created, &postTmp.Forum, &postTmp.IsEdited, &postTmp.Parent, &postTmp.Thread)
	if err != nil {
		return postResult, models.ErrorNotFound
	}
	postResult.Post = postTmp
	for i := 0; i < len(related); i++ {
		switch related[i] {
		case "user":
			user, _ := r.GetUser(ctx, postTmp.Author)
			postResult.Author = &user
		case "forum":
			forum, _ := r.GetForum(ctx, postTmp.Forum)
			postResult.Forum = &forum
		case "thread":
			thread, _ := r.GetThreadById(ctx, postTmp.Thread)
			postResult.Thread = &thread
		}
	}

	return postResult, nil
}

func (r *ForumRepository) getThreadPostsFlat(ctx context.Context, limit, since, desc string, id int) []models.Post {
	result := make([]models.Post, 0)
	if since == "" {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetPostsWithSinceDesc, id, limit)
			defer rows.Close()
			for rows.Next() {
				post := models.Post{}
				rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
				result = append(result, post)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetPostsWithSinceAsc, id, limit)
			defer rows.Close()
			for rows.Next() {
				post := models.Post{}
				rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
				result = append(result, post)
			}
		}
	} else {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetPostsDesc, id, since, limit)
			defer rows.Close()
			for rows.Next() {
				post := models.Post{}
				rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
				result = append(result, post)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetPostsAsc, id, since, limit)
			defer rows.Close()
			for rows.Next() {
				post := models.Post{}
				rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
				result = append(result, post)
			}
		}
	}
	return result
}

func (r *ForumRepository) getThreadPostsTree(ctx context.Context, limit, since, desc string, id int) []models.Post {
	result := make([]models.Post, 0)
	var (
		rows       pgx.Rows
		queryFinal string
	)

	if limit == "" && since == "" {
		if desc == "true" {
			queryFinal = GetPostsTreeDesc
		} else {
			queryFinal = GetPostsTreeAsc
		}
		rows, _ = r.conn.Query(ctx, queryFinal, id)
	} else if limit != "" && since == "" {
		if desc == "true" {
			queryFinal = GetPostsTreeWithLimitDesc
		} else {
			queryFinal = GetPostsTreeWithLimitAsc
		}
		rows, _ = r.conn.Query(ctx, queryFinal, id, limit)
	} else if limit == "" && since != "" {
		if desc == "true" {
			queryFinal = SelectTreeSinceNilDesc
		} else {
			queryFinal = SelectTreeSinceNilDescNil
		}
		rows, _ = r.conn.Query(ctx, queryFinal, id, since)
	} else if limit != "" && since != "" {
		if desc == "true" {
			queryFinal = GetPostsTreeWithLimitWithSinceDesc
		} else {
			queryFinal = GetPostsTreeWithLimitWithSinceAsc
		}
		rows, _ = r.conn.Query(ctx, queryFinal, id, since, limit)
	}

	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		result = append(result, post)
	}
	return result
}

func (r *ForumRepository) getThreadPostsParentTree(ctx context.Context, limit, since, desc string, id int) []models.Post {
	result := make([]models.Post, 0)
	var rows pgx.Rows
	halfQuery := fmt.Sprintf(`SELECT id FROM "post" WHERE thread = %d AND parent = 0 `, id)
	if since != "" {
		if desc == "true" {
			halfQuery += ` AND path[1] < ` + fmt.Sprintf(`(SELECT path[1] FROM "post" WHERE id = %s) `, since)
		} else {
			halfQuery += ` AND path[1] > ` + fmt.Sprintf(`(SELECT path[1] FROM "post" WHERE id = %s) `, since)
		}
	}
	if desc == "true" {
		halfQuery += ` ORDER BY id DESC `
	} else {
		halfQuery += ` ORDER BY id ASC `
	}
	if limit != "" {
		halfQuery += " LIMIT " + limit
	}
	fullQuery := fmt.Sprintf(`SELECT id, author, created, forum, isedited, message, parent, thread FROM "post" WHERE path[1] = ANY (%s) `, halfQuery)
	if desc == "true" {
		fullQuery += ` ORDER BY path[1] DESC, path, id `
	} else {
		fullQuery += ` ORDER BY path[1] ASC, path, id `
	}

	rows, _ = r.conn.Query(ctx, fullQuery)
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message,
			&post.Parent, &post.Thread)
		result = append(result, post)
	}
	return result
}
func (r *ForumRepository) GetThreadPosts(ctx context.Context, limit, since, desc, sort string, threadId int) ([]models.Post, error) {
	result := make([]models.Post, 0)
	switch sort {
	case "flat":
		result = r.getThreadPostsFlat(ctx, limit, since, desc, threadId)
	case "tree":
		result = r.getThreadPostsTree(ctx, limit, since, desc, threadId)
	case "parent_tree":
		result = r.getThreadPostsParentTree(ctx, limit, since, desc, threadId)
	default:
		result = r.getThreadPostsFlat(ctx, limit, since, desc, threadId)
	}

	return result, nil
}
func (r *ForumRepository) UpdateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {
	result := models.Thread{}
	if thread.Slug == "" {
		resultQuery := fmt.Sprintf(UpdateThreadWithoutIdentifier, `id=$3`)
		row := r.conn.QueryRow(ctx, resultQuery, thread.Title, thread.Message, thread.ID)
		err := row.Scan(&result.ID, &result.Title, &result.Author,
			&result.Forum, &result.Message, &result.Votes, &result.Slug, &result.Created)
		if err != nil {
			return models.Thread{}, models.ErrorNotFound
		}
	} else {
		resultQuery := fmt.Sprintf(UpdateThreadWithoutIdentifier, `slug=$3`)
		row := r.conn.QueryRow(ctx, resultQuery, thread.Title, thread.Message, thread.Slug)
		err := row.Scan(&result.ID, &result.Title, &result.Author,
			&result.Forum, &result.Message, &result.Votes, &result.Slug, &result.Created)
		if err != nil {
			return models.Thread{}, models.ErrorNotFound
		}
	}
	return result, nil
}

func (r *ForumRepository) GetUsers(ctx context.Context, slug, limit, since, desc string) ([]models.User, error) {
	if limit == "" {
		limit = "100"
	}
	users := make([]models.User, 0)
	if since != "" {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetUsersWithSinceDesc, slug, since, limit)
			defer rows.Close()
			for rows.Next() {
				tmpUser := models.User{}
				err := rows.Scan(&tmpUser.Nickname, &tmpUser.Fullname, &tmpUser.About, &tmpUser.Email)
				if err != nil {
					continue
				}
				users = append(users, tmpUser)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetUsersWithSinceAsc, slug, since, limit)
			defer rows.Close()
			for rows.Next() {
				tmpUser := models.User{}
				err := rows.Scan(&tmpUser.Nickname, &tmpUser.Fullname, &tmpUser.About, &tmpUser.Email)
				if err != nil {
					continue
				}
				users = append(users, tmpUser)
			}
		}
	} else {
		if desc == "true" {
			rows, _ := r.conn.Query(ctx, GetUsersDesc, slug, limit)
			defer rows.Close()
			for rows.Next() {
				tmpUser := models.User{}
				err := rows.Scan(&tmpUser.Nickname, &tmpUser.Fullname, &tmpUser.About, &tmpUser.Email)
				if err != nil {
					continue
				}
				users = append(users, tmpUser)
			}
		} else {
			rows, _ := r.conn.Query(ctx, GetUsersAsc, slug, limit)
			defer rows.Close()
			for rows.Next() {
				tmpUser := models.User{}
				err := rows.Scan(&tmpUser.Nickname, &tmpUser.Fullname, &tmpUser.About, &tmpUser.Email)
				if err != nil {
					continue
				}
				users = append(users, tmpUser)
			}
		}
	}
	return users, nil
}
func (r *ForumRepository) UpdatePost(ctx context.Context, post models.PostUpdate) (models.Post, error) {
	row := r.conn.QueryRow(ctx, UpdatePostMessage, post.Message, post.ID)
	result := models.Post{}
	err := row.Scan(&result.ID, &result.Author, &result.Created, &result.Forum,
		&result.IsEdited, &result.Message, &result.Parent, &result.Thread, &result.Path)
	if err != nil {
		return result, models.ErrorNotFound
	}
	return result, nil
}

func (r *ForumRepository) GetStatus() models.Status {
	return models.Status{
		UsersCount:   r.Status.UsersCount,
		ForumsCount:  r.Status.ForumsCount,
		ThreadsCount: r.Status.ThreadsCount,
		PostsCount:   r.Status.PostsCount,
	}
}

func (r *ForumRepository) Clear() {
	r.Status = models.Status{}
	r.conn.Exec(context.Background(), DESTROY_DATABASE_DONT_TOCUH_DANGEROUS)
}
