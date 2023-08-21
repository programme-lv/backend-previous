package testdb

import (
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

type DBMigrations struct {
	rootDir string
}

func cloneDBMigrations() (*DBMigrations, error) {
	tmpDir, err := os.MkdirTemp("", "proglv-db-migrations")
	if err != nil {
		return nil, err
	}
	repoUrl := "https://github.com/programme-lv/database"

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	res := &DBMigrations{
		rootDir: tmpDir,
	}
	return res, nil
}

func (dbm *DBMigrations) getFlywayMigrationsDir() string {
	return filepath.Join(dbm.rootDir, "flyway-migrations")
}

func (dbm *DBMigrations) erase() {
	err := os.RemoveAll(dbm.rootDir)
	if err != nil {
		log.Printf("Failed to remove tmp dir: %v", err)
	}
}
