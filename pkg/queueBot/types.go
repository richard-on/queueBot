package queueBot

var BotState = Initial

const (
	Initial = iota
	SubjectSelect
	QueueSelect
	QueueAction
	AdminMode
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

type QueuePrint struct {
	Username  string
	FirstName string
	LastName  string
	Position  int
}
