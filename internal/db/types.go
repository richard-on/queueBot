package db

import "database/sql"

type Group struct {
	ID         int64
	Name       string
	IsSubgroup bool
}

type Subject struct {
	ID                int64
	GroupID           int64
	IsSubgroupSubject sql.NullBool
	SubGroupID        sql.NullInt64
	SubjectName       string
}

type Queue struct {
	ID        int64
	SubjectId int64
	Name      string
}

type QueuePrint struct {
	Username  string
	FirstName string
	LastName  string
	Position  sql.NullInt64
}
