package migrator

import (
	"database/sql"
	"fmt"
	"gameapp/repository/mysql"
	"github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dbConfig      mysql.Config
	migrateSource *migrate.FileMigrationSource
}

// TODO - set migration table name
// TODO - get status of migrations

func New(config mysql.Config) *Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	return &Migrator{
		migrateSource: migrations,
		dbConfig:      config,
	}
}

func (m Migrator) Up() {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?parseTime=true",
			m.dbConfig.Username,
			m.dbConfig.Password,
			m.dbConfig.Host,
			m.dbConfig.Port,
			m.dbConfig.DBName,
		),
	)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db for migrations: %v", err))
	}

	n, err := migrate.Exec(db, "mysql", m.migrateSource, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v", err))
	}

	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?parseTime=true",
			m.dbConfig.Username,
			m.dbConfig.Password,
			m.dbConfig.Host,
			m.dbConfig.Port,
			m.dbConfig.DBName,
		),
	)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db for migrations: %v", err))
	}

	n, err := migrate.Exec(db, "mysql", m.migrateSource, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v", err))
	}

	fmt.Printf("rollbacked %d migrations!\n", n)
}

//func (m Migrator) Status() {
//	db, err := sql.Open("mysql",
//		fmt.Sprintf(
//			"%s:%s@(%s:%d)/%s?parseTime=true",
//			m.dbConfig.Username,
//			m.dbConfig.Password,
//			m.dbConfig.Host,
//			m.dbConfig.Port,
//			m.dbConfig.DBName,
//		),
//	)
//	if err != nil {
//		panic(fmt.Errorf("can't open mysql db for migrations: %v", err))
//	}
//
//	n, err := migrate.Exec(db, "mysql", m.migrateSource, migrate.Up)
//	if err != nil {
//		panic(fmt.Errorf("can't get status of migrations: %v", err))
//	}
//	fmt.Printf("Applied %d migrations!\n", n)
//}
