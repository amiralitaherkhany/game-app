package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func New(cfg Config) *MySQLDB {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		),
	)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{
		db: db,
	}
}
