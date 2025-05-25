package models

type Quote struct {
	ID     int64  `json:"id" db:"id"`
	Author Author `json:"author" db:"author"`
	Text   string `json:"quote" db:"text"`
}

type Author struct {
	Name string `json:"name" db:"name"`
}
