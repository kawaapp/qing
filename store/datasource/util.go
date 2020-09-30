package datasource

import (
	"strconv"
	"time"
)

func joinIntArray(idx []int64) string {
	var (
		sz = len(idx)
		q  = ""
	)
	if sz > 1 {
		for i := 0; i < sz-1; i++ {
			q += strconv.Itoa(int(idx[i])) + ","
		}
	}
	if sz > 0 {
		q += strconv.Itoa(int(idx[sz-1]))
	}
	return q
}

func UnixNow() int64 {
	return time.Now().UTC().Unix()
}

// 时间
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// helper
func Count(db *datasource, stmt string, args ...interface{}) (int, error) {
	rows := db.QueryRow(stmt, args...)
	var count int
	err := rows.Scan(&count)
	return count, err
}

func Delete(db *datasource, stmt string, args ...interface{}) error {
	_, err := db.Exec(stmt, args...)
	return err
}