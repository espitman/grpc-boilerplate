package db

// import (
// 	"context"
// 	"database/sql"
// 	"entgo.io/ent/dialect"
// 	entsql "entgo.io/ent/dialect/sql"
// 	_ "github.com/jackc/pgx/v5/stdlib"
// )

// type DB struct {
// 	Client *ent.Client
// }

// func NewDB(dsn string) *DB {
// 	db, err := sql.Open("pgx", dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create an ent.Driver from `db`.
// 	drv := entsql.OpenDB(dialect.Postgres, db)
// 	client := ent.NewClient(ent.Driver(drv))

// 	if err := client.Schema.Create(context.Background()); err != nil {
// 		panic(err)
// 	}

// 	return &DB{
// 		Client: client,
// 	}

// }