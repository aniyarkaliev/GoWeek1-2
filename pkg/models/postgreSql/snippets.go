package postgreSql

import (
	"Se09.com/pkg/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"strings"
	"time"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	//Query
	s := &models.Snippet{}
	stml := "SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets WHERE expires > now() AND id=$1"

	err := m.DB.QueryRow(context.Background(), stml, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	s.Content = strings.Replace(s.Content, "\\n", "\n", -2)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stml := "SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets " +
		"WHERE expires > now() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(context.Background(), stml)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, err
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stml := "INSERT INTO snippets (title, content, created, expires)" +
		"VALUES($1, $2, $3, $4) RETURNING id"
	days, error := strconv.Atoi(expires)
	if error != nil {
		return 0, error
	}
	fmt.Println(stml)

	var lastId int
	created := time.Now()
	expiresDate := time.Now().AddDate(0, 0, days)

	err := m.DB.QueryRow(context.Background(), stml, title, content, created, expiresDate).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return int(lastId), err
}
