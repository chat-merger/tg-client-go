package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	DataSourceName string
}

func InitSqlite(cfg Config) (*sql.DB, error) {
	// https://github.com/glebarez/go-sqlite#connection-string-examples
	db, err := sql.Open("sqlite", cfg.DataSourceName+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %s", err)
	}
	return db, nil
}
