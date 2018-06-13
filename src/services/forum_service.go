package services

import (
	"fmt"

	"github.com/Alex-Kuz/tp-database/src/models"
	"github.com/jackc/pgx"
)

type ForumService struct {
	db        *PostgresDatabase
	tableName string
}

func MakeForumService(pgdb *PostgresDatabase) ForumService {
	return ForumService{db: pgdb, tableName: "forums"}
}

func (fs *ForumService) TableName() string {
	return fs.tableName
}

func (fs *ForumService) GetForumBySlug(slug string) *models.Forum {
	query := "SELECT slug::text, author::text, title::text, threads, posts FROM forums WHERE LOWER(slug) = LOWER($1)"

	rows := fs.db.Query(query, slug)
	defer rows.Close()

	for rows.Next() {
		forum := new(models.Forum)
		err := rows.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Threads, &forum.Posts)
		if err != nil {
			panic(err)
		}
		return forum
	}
	return nil
}

func (fs *ForumService) SlugBySlug(slug string) *string {
	query := "SELECT slug::text FROM forums WHERE LOWER(slug) = LOWER($1)"

	rows := fs.db.Query(query, slug)
	defer rows.Close()

	for rows.Next() {
		str := new(string)
		err := rows.Scan(str)
		if err != nil {
			panic(err)
		}
		return str
	}
	return nil
}

func (fs *ForumService) IncThreadsCountBySlug(slug string) bool {
	UPDATE_QUERY := "UPDATE forums SET threads = threads + 1 WHERE LOWER(slug) = LOWER($1);"

	resultRows := fs.db.QueryRow(UPDATE_QUERY, slug)

	if err := resultRows.Scan(); err != nil && err != pgx.ErrNoRows {
		panic(err)
	}

	return true
}

func (fs *ForumService) AddForum(forum *models.Forum) (bool, *models.Forum) {

	if conflictForum := fs.GetForumBySlug(forum.Slug); conflictForum != nil {
		return false, conflictForum
	}

	INSERT_QUERY := "insert into forums (slug, author, title, threads, posts) values ($1, $2, $3, $4, $5);"

	resultRows := fs.db.QueryRow(INSERT_QUERY, forum.Slug, forum.User, forum.Title, forum.Threads, forum.Posts)

	if err := resultRows.Scan(); err != nil && err != pgx.ErrNoRows {
		panic(err)
	}

	return true, forum
}

func (fs *ForumService) GetUsers(forum *models.Forum, since, limit string, desc bool) []models.User {

	sinceStr := ""
	if since != "" {
		sinceStr = " AND LOWER(u.nickname) "
		if desc {
			sinceStr += "< "
		} else {
			sinceStr += "> "
		}
		sinceStr += "LOWER('" + since + "')"
	}

	order := " ASC"
	if desc {
		order = " DESC"
	}

	limitStr := ""
	if limit != "" {
		limitStr = " LIMIT " + limit
	}

	query := fmt.Sprintf(
		"SELECT nickname::text, fullname::text, about::text, email::text FROM users u JOIN forum_users uf ON LOWER(u.nickname) = LOWER(uf.username)"+
			" WHERE LOWER(uf.forum) = LOWER('%s') %s ORDER BY LOWER(u.nickname) %s %s;",
		forum.Slug, sinceStr, order, limitStr)

	rows := fs.db.Query(query)
	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}
	return users
}

func (fs *ForumService) IncrementPostsCountBySlug(forumSlug string, postsCount int) {

	updateQuery := "UPDATE forums SET posts = posts + $2 WHERE LOWER(slug) = LOWER($1);"

	resultRows := fs.db.QueryRow(updateQuery, forumSlug, postsCount)

	if err := resultRows.Scan(); err != nil && err != pgx.ErrNoRows {
		panic(err)
	}
}
