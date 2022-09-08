package initEnv

import (
	"fmt"
	"os"
)

var SentryDsn string
var Host string
var Port string
var User string
var Password string
var DbName string
var DbInfo string
var InitDbInfo string

func Init() {
	SentryDsn = os.Getenv("SENTRY_DSN")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")
	User = os.Getenv("USER")
	Password = os.Getenv("PASSWORD")
	DbName = os.Getenv("DBNAME")

	DbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("HOST"),
		os.Getenv("PORT"), os.Getenv("DBNAME"))

	InitDbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("HOST"), os.Getenv("PORT"))
}
