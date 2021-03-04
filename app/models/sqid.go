package models

type SequentialIdentifier struct {
	ID int64 `db:"id" json:"id"`
}

func (si SequentialIdentifier) IDColumn() string {
	return "ID"
}

func (si SequentialIdentifier) IsNew() bool {
	return si.ID == 0
}
