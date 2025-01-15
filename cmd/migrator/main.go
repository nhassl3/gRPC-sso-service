package main

import (
	"errors"
	"flag"
	"fmt"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate"
	// Драйвер для выполнения миграции в SQLite3
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "", "Path to a table containing the migration files")
	flag.Parse()

	if storagePath == "" || migrationsPath == "" {
		panic("storage-path and migrationsPath is required")
	}

	databaseUrl := fmt.Sprintf("sqlite3://%s", storagePath)
	if migrationsTable != "" {
		databaseUrl = fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable)
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseUrl,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Nothing to migrate")
			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied successfully")
}
