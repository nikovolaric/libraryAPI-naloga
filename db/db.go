package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "postgresql://neondb_owner:npg_khNAa4FVnWe9@ep-old-king-agl4bigp-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(("DB connected"))
}
