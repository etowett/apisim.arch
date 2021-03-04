package models

import (
	"time"

	null "gopkg.in/guregu/null.v4"
)

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt null.Time `db:"updated_at" json:"updated_at"`
}

func (t *Timestamps) Touch() {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
}
