package Postgresql

import (
	"Se09.com/pkg/models"
	"context"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stml := "Insert into snippets (title, content, created, expires) VALUES  ($1, $2. &3, $4) returning id"

	var LastId int
	err := m.DB.QueryRow(context.Background(), stml, title, content, time.Now(), time.Now().AddDate(0, 0, 10)).Scan(&LastId)
	if err != nil {
		return 0, err
	}
	return int(LastId), err
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	stml := "Select id,title, content, created, expires From snippets Where expires > CLOCK_TIMESTAMP() and id=$1"
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
	stml := "Select id,title, content, created, expires From snippets Where expires > CLOCK_TIMESTAMP() order by created desc limit 10"
	rows, err := m.DB.Query(context.Background(), stml)
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
