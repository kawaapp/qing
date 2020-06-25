package ddl

import (
	"database/sql"
	"github.com/kawaapp/kawaqing/store/datasource/ddl/mysql"
	"github.com/kawaapp/kawaqing/store/datasource/ddl/postgres"
	"github.com/kawaapp/kawaqing/store/datasource/ddl/sqlite"
)

// Supported database drivers
const (
	DriverSqlite   = "sqlite3"
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
)

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(driver string, db *sql.DB) error {
	if err := checkPriorMigration(db); err != nil {
		return err
	}
	switch driver {
	case DriverMySQL:
		return mysql.Migrate(db)
	case DriverSqlite:
		fallthrough
	default:
		return sqlite.Migrate(db)
	}
	return nil
}

func checkPriorMigration(db *sql.DB) error {
	return nil
}
