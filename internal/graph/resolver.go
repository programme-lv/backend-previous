package graph

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *sqlx.DB
}
