package model

type User struct {
	UserID       int64
	ChatID       int64
	TgUsername   string
	GroupID      int64
	SubgroupID   int64
	TgFirstName  string
	TgLastName   string
	FirstName    string
	LastName     string
	IsRegistered bool
}
