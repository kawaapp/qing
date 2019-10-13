package postgres

import (
	"database/sql"
)

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	return nil
}

// 这里的设计还是比较屌的，先创建一个数据库迁移表，然后再用这个表
// 记录数据库迁移的记录
func createTable(db *sql.DB) error {
	for _, sql := range index {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

// Not implemented, yet.
var index = []string {
}
