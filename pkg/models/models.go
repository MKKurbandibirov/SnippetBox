package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no such record")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
}
