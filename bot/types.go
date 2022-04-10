package bot

import (
	"time"
)

var botState = initial

const (
	initial = iota
	subjectSelect
	queueSelect
	queueAction
	adminMode
)

type Subjects struct {
	Id       int64  `json:"id"`
	Alias    string `json:"alias"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

type QueueInfo struct {
	SubjectId int64
	QueueId   int64
	Name      string
}

type QueueList struct {
	queueId   string
	subjectId int64
	name      string
	date      time.Time
}

type QueuePrint struct {
	Username  string
	FirstName string
	LastName  string
	Position  int
}
