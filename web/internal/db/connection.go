package db

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	Postgres *sql.DB
	Redis    *redis.Client
)

func Init() {
	dsn := os.Getenv("POSTGRES_DSN")

	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_DATABASE")
		user := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")

		dsn = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=disable"
	}

	db, _ := sql.Open("postgres", dsn)

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err == nil {
		m, _ := migrate.NewWithDatabaseInstance("file://./internal/db/migrations", "postgres", driver)
		m.Up()
	}

	Postgres = db

	addr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	Redis = rdb
}
