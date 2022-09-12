package model

import "database/sql"

type User struct {
	UserID      int64
	ChatID      int64
	TgUsername  string
	GroupID     int64
	SubgroupID  sql.NullInt64
	TgFirstName sql.NullString
	TgLastName  sql.NullString
	FirstName   sql.NullString
	LastName    sql.NullString
}
