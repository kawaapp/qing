package sql

import (
	"github.com/kawaapp/kawaqing/store/datasource/sql/mysql"
	"github.com/kawaapp/kawaqing/store/datasource/sql/sqlite"
)

// Supported database drivers
const (
	DriverSqlite   = "sqlite3"
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
)

// Lookup returns the named sql statement compatible with
// the specified database driver.
func Lookup(driver string, name string) string {
	switch driver {
	case DriverMySQL:
		return mysql.Lookup(name)
	case DriverSqlite:
		fallthrough
	default:
		return sqlite.Lookup(name)
	}
}
