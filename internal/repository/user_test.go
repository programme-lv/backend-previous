package repository_test

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/programme-lv/backend/internal"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/repository"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(5 * time.Minute),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			log.Fatalf("Could not terminate container: %s", err)
		}
	}(postgresContainer, ctx)

	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432/tcp")
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// Run migrations
	migrations, err := migrate.New(
		"github://programme-lv/database/go-migrate", dsn)

	if err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}
	err = migrations.Up()
	if err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestUserRepoPostgreSQLImpl_CreateUserGetByID(t *testing.T) {
	var repo internal.UserRepo = repository.NewUserRepoPostgreSQLImpl(db)
	repoTx, err := repo.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("Could not begin transaction: %s", err)
	}
	defer repoTx.Rollback()

	username := "testuser"
	password := "password1"
	email := "test@gmail.com"
	firstName := "Test"
	lastName := "User"

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Could not hash password: %s", err)
	}

	userID, err := repo.CreateUser(username, encryptedPassword,
		email, firstName, lastName)
	if err != nil {
		t.Fatalf("Could not create user: %s", err)
	}

	user, err := repo.GetUserByID(userID)
	if err != nil {
		t.Fatalf("Could not get user by ID: %s", err)
	}

	if user.Username != username {
		t.Fatalf("Expected username to be %s, got %s", username, user.Username)
	}

	if user.Email != email {
		t.Fatalf("Expected email to be %s, got %s", email, user.Email)
	}

	if user.FirstName != firstName {
		t.Fatalf("Expected first name to be %s, got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Fatalf("Expected last name to be %s, got %s", lastName, user.LastName)
	}

	if user.ID != userID {
		t.Fatalf("Expected ID to be %d, got %d", userID, user.ID)
	}

	if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) != nil {
		t.Fatalf("Password does not match")
	}
}

func TestUserRepoPostgreSQLImpl_CreateUserGetByUsername(t *testing.T) {
	repo := repository.NewUserRepoPostgreSQLImpl(db)

	username := "testuser"
	password := "password1"
	email := "test@gmail.com"
	firstName := "Test"
	lastName := "User"

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Could not hash password: %s", err)
	}

	userID, err := repo.CreateUser(username, encryptedPassword,
		email, firstName, lastName)
	if err != nil {
		t.Fatalf("Could not create user: %s", err)
	}

	user, err := repo.GetUserByUsername(username)
	if err != nil {
		t.Fatalf("Could not get user by ID: %s", err)
	}

	if user.Username != username {
		t.Fatalf("Expected username to be %s, got %s", username, user.Username)
	}

	if user.Email != email {
		t.Fatalf("Expected email to be %s, got %s", email, user.Email)
	}

	if user.FirstName != firstName {
		t.Fatalf("Expected first name to be %s, got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Fatalf("Expected last name to be %s, got %s", lastName, user.LastName)
	}

	if user.ID != userID {
		t.Fatalf("Expected ID to be %d, got %d", userID, user.ID)
	}

	if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) != nil {
		t.Fatalf("Password does not match")
	}

}
