package connection

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var dsn = "landmark:landmark@csmsu@tcp(202.28.34.197:3306)/landmark"

func GetConnectionX() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	println("connection successful")

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return db, nil
}
