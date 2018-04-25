package services

import (
	"fmt"

	"github.com/Alex-Kuz/tp-database/src/models"
)


type ForumService struct {
	db        *PostgresDatabase
	tableName string
}

func MakeForumService(pgdb *PostgresDatabase) ForumService {
	return ForumService{db: pgdb, tableName: "forums"}
}



func (fs *ForumService) GetForumBySlug(slug string) *models.Forum {
	query := fmt.Sprintf(
		"SELECT slug, author, title, threads, posts FROM %s WHERE LOWER(slug) = LOWER('%s')",
			fs.tableName, slug)

	rows := fs.db.Query(query)

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
	query := fmt.Sprintf(
		"SELECT slug FROM %s WHERE LOWER(slug) = LOWER('%s')",
		fs.tableName, slug)

	rows := fs.db.Query(query)

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
	UPDATE_QUERY :=
		"UPDATE " + fs.tableName + " SET threads = threads + 1 WHERE LOWER($1) = LOWER(slug);"

	insertQuery, err := fs.db.Prepare(UPDATE_QUERY)
	//defer insertQuery.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	fmt.Println("IncThreadsCountBySlug:  slug =", slug)

	_, err = insertQuery.Exec(slug)
	if err != nil {
		fmt.Println("AddForum:  error:", err.Error())
		panic(err)
	}

	return true
}

func (fs *ForumService) AddForum(forum *models.Forum) (bool, *models.Forum) {

	if conflictForum := fs.GetForumBySlug(forum.Slug); conflictForum != nil {
		return false, conflictForum
	}

	INSERT_QUERY :=
		"insert into " + fs.tableName + " (slug, author, title, threads, posts) values ($1, $2, $3, $4, $5);"

	insertQuery, err := fs.db.Prepare(INSERT_QUERY)
	//defer insertQuery.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("AddForum:  forum.User =", forum.User)

	_, err = insertQuery.Exec(forum.Slug, forum.User, forum.Title, forum.Threads, forum.Posts)
	if err != nil {
		fmt.Println("AddForum:  error:", err.Error())
		panic(err)
	}

	return true, forum
}
