package config

import (
	"fmt"
	"os"
)

func Init() {
	TgToken = os.Getenv("TOKENDEV")
	SentryDsn = os.Getenv("SENTRY_DSN")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")
	User = os.Getenv("USER")
	Password = os.Getenv("PASSWORD")
	DbName = os.Getenv("DBNAME")

	DbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		User, Password, Host, Port, DbName)

	InitDbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		User, Password, Host, Port)
}
