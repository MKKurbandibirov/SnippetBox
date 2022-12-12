package postgresql

import (
	"database/sql"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
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

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
