package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"libraryAPI/devops"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://library:library@localhost:5432/librarydb?sslmode=disable"
	}

	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB connected")

	migrateUp()
}

// migrateUp applies all pending migrations from the embedded SQL files.
func migrateUp() {
	src, err := iofs.New(devops.Migrations, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	driver, err := migratepg.WithInstance(DB, &migratepg.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	fmt.Println("migrations up to date")
}
