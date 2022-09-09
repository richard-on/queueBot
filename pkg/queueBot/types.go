package queueBot

import "database/sql"

var BotState = Initial

const (
	Initial = iota
	SubjectSelect
	QueueSelect
	QueueAction
	AdminMode
)

type User struct {
	ID           int64
	TgUsername   string
	GroupID      int64
	SubgroupID   sql.NullInt64
	TgFirstName  sql.NullString
	TgLastName   sql.NullString
	FirstName    sql.NullString
	LastName     sql.NullString
	GroupName    string
	SubGroupName string
}

type Subjects struct {
	ID                int64
	SubjectName       string
	IsSubgroupSubject sql.NullBool
	GroupID           int64
	SubGroupID        sql.NullInt64
}

type QueueInfo struct {
	QueueId   int64
	SubjectId int64
	Name      string
}

type QueuePrint struct {
	Username  string
	FirstName string
	LastName  string
	Position  sql.NullInt64
}
