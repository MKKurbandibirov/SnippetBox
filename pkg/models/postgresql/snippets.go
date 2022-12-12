package postgresql

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expired FROM snippets WHERE expired > now() AND id = $1`

	row := sm.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (sm *SnippetModel) Insert(title, content string, expired int) (int, error) {
	stmt := `INSERT INTO snippets(title, content, created, expired) VALUES ($1, $2, current_date, current_date + 7) RETURNING id`

	var id int
	if err := sm.DB.QueryRow(stmt, title, content).Scan(&id); err != nil {
		println("hhdhdhdhdhdhdhdhhdhdhdhhdh")
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expired FROM snippets WHERE snippets.expired > now() ORDER BY created DESC LIMIT 10`

	rows, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
