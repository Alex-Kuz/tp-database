package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Alex-Kuz/tp-database/src/models"
)

type ThreadService struct {
	db        *PostgresDatabase
	tableName string
}

func MakeThreadService(pgdb *PostgresDatabase) ThreadService {
	return ThreadService{db: pgdb, tableName: "threads"}
}




func (ts *ThreadService) AddThread(thread *models.Thread) (bool, *models.Thread) {

	INSERT_QUERY:=
		"insert into " + ts.tableName +
			" (%s author, forum, created, title, message, votes)" +
				" values (%s $1, $2, $3, $4, $5, $6) returning id;"

	if thread.Slug == "" {
		INSERT_QUERY = fmt.Sprintf(INSERT_QUERY, "", "")
	} else {
		INSERT_QUERY = fmt.Sprintf(INSERT_QUERY, "slug, ", "'" + thread.Slug + "', ")
	}

	layout := time.RFC3339Nano
	t, err := time.Parse(layout, thread.Created)
	if err != nil {
		fmt.Println(err)
	}

	thread.Created = t.UTC().Format(time.RFC3339Nano)
	fmt.Println("AddThread: since:", thread.Created)

	fmt.Println("AddThread: INSERT_QUERY:", INSERT_QUERY)

	err = ts.db.QueryRow(INSERT_QUERY, thread.Author, thread.Forum,
		thread.Created, thread.Title, thread.Message, thread.Votes).Scan(&thread.ID)

	if err != nil {
		fmt.Println("AddThread:  error after id:", err.Error())
		panic(err)
	}

	fmt.Println("AddForum: id:", thread.ID)

	return true, thread
}

func (ts *ThreadService) SelectThreads(slug, limit, since string, desc bool) (bool, []models.Thread) {

	fmt.Println("")

	limitStr := ""
	if limit != "" {
		fmt.Println("SelectThreads: limit:", limit)
		limitStr = "LIMIT " + limit
	}


	comp := " >= "
	order := "ORDER BY th.created "
	if desc {
		comp = " <= "
		order += "DESC "
	} else {
		order += "ASC "
	}

	offsetStr := ""
	if since != "" {
		t, err := time.Parse(time.RFC3339Nano, since)
		if err != nil {
			fmt.Println(err)
		}

		since = t.UTC().Format(time.RFC3339Nano)
		fmt.Println("SelectThreads: since:", since)
		offsetStr = "AND th.created " + comp + " '" + since + "'"
	}


	query := fmt.Sprintf(
		"SELECT id, slug, author, forum, created, title, message, votes FROM %s th WHERE LOWER(th.forum) = LOWER('%s') %s %s %s;",
		ts.tableName, slug, offsetStr, order,  limitStr)

	fmt.Println("SelectThreads: query:", query)


	rows := ts.db.Query(query)
	threads := make([]models.Thread, 0)

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Created,
			&thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			panic(err)
		}
		threads = append(threads, thread)
	}

	if len(threads) == 0 {
		return false, threads
	}

	return true, threads
}

func (ts *ThreadService) GetThreadBySlug(slug string) *models.Thread {

	fmt.Println("GetThreadBySlug: query start")

	query := fmt.Sprintf(
		"SELECT id, slug, author, forum, created, title, message, votes FROM %s WHERE LOWER(slug) = LOWER('%s');",
			ts.tableName, slug)

	fmt.Println("GetThreadBySlug: query:", query)
	fmt.Println("-----------------------------start------------------------------####################")

	rows := ts.db.Query(query)
	fmt.Println("------------------------------end-------------------------------####################")

	for rows.Next() {
		thread := new(models.Thread)
		err := rows.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Created,
			&thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		return thread
	}

	return nil
}

func (ts *ThreadService) GetThreadById(id uint64) *models.Thread {

	fmt.Println("GetThreadById: query start")

	query := fmt.Sprintf(
		"SELECT id, slug, author, forum, created, title, message, votes FROM %s WHERE id = %s;",
		ts.tableName, strconv.FormatUint(id, 10))

	fmt.Println("GetThreadBySlug: query:", query)
	fmt.Println("-----------------------------start------------------------------####################")

	rows := ts.db.Query(query)
	fmt.Println("------------------------------end-------------------------------####################")

	for rows.Next() {
		thread := new(models.Thread)
		err := rows.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Created,
			&thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		return thread
	}

	return nil
}

func (ts *ThreadService) Vote(thread *models.Thread, vote models.Vote) *models.Thread {

	fmt.Println("Vote: query start")

	addVoteStr := "+ 1"
	if vote.Voice == -1 {
		addVoteStr = "- 1"
	}

	query := fmt.Sprintf(
		"UPDATE %s SET votes = votes %s WHERE id = %s returning votes;",
		ts.tableName, addVoteStr, strconv.FormatUint(thread.ID, 10))

	voice, voteId := ts.getVote(vote.Nickname, thread.ID)
	if voice != nil {
		if *voice == vote.Voice {
			fmt.Println("Vote: this vote is existed:", query)
			return thread
		} else {
			voiceUpdate := "UPDATE votes SET voice = $1 WHERE id = $2;"

			fmt.Println("<< Vote  update: vote, thread.id:", vote, thread.ID)
			fmt.Println("-----------------------------start------------------------------####################")

			ts.db.Query(voiceUpdate, vote.Voice, voteId)

			addVoteStr = "+ 2"
			if vote.Voice == -1 {
				addVoteStr = "- 2"
			}
		}
	} else {

		voteInsert := "INSERT INTO votes (username, voice, thread) VALUES ($1, $2, $3) returning id;"

		fmt.Println("<< Vote  insert: vote, thread.id:", vote, thread.ID)
		fmt.Println("-----------------------------start------------------------------####################")

		var id uint64
		err := ts.db.QueryRow(voteInsert, vote.Nickname, vote.Voice, thread.ID).Scan(&id)
		if err != nil {
			fmt.Println("Vote:  error:", err.Error())
			panic(err)
		}
	}

	fmt.Println("-----------------------------start-2-----------------------------####################")
	err := ts.db.QueryRow(query).Scan(&thread.Votes)
	if err != nil {
		fmt.Println("Vote:  error:", err.Error())
		panic(err)
	}

	fmt.Println("------------------------------end-------------------------------####################")

	return thread
}

func (ts *ThreadService) getVote(username string, threadId uint64) (*int32, *uint64) {
	fmt.Println("getVote: query start")

	query := fmt.Sprintf(
		"SELECT id, voice FROM votes WHERE thread = %s AND username = '%s';",
		strconv.FormatUint(threadId, 10), username)

	fmt.Println("getVote: query:", query)
	fmt.Println("-----------------------------start------------------------------####################")

	rows := ts.db.Query(query)
	fmt.Println("------------------------------end-------------------------------####################")

	for rows.Next() {
		voice := new(int32)
		id := new(uint64)
		err := rows.Scan(id, voice)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		return voice, id
	}
	return nil, nil
}