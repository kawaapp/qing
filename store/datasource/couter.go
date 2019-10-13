package datasource

import (
	"time"
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/store/datasource/sql"
)

// All users
func (db *datasource) GetTotalUser() (int, error) {
	return Count(db, `SELECT COUNT(id) FROM users;`)
}

// All discussions
func (db *datasource) GetTotalDiscussion() (int, error) {
	return Count(db, "SELECT COUNT(id) FROM discussions;")
}

// New users/day
func (db *datasource) GetNewUser(day time.Time) (int, error) {
	var (
		from = Bod(day).Unix()
		to = Bod(day).Add(24 * time.Hour).Unix()
	)
	stmt := "SELECT COUNT(id) FROM users WHERE created_at > ? AND created_at < ?;"
	return Count(db, stmt, from, to)
}

// New discussions/day
func (db *datasource) GetNewDiscussion(day time.Time) (int, error) {
	var (
		from = Bod(day).Unix()
		to = Bod(day).Add(24 * time.Hour).Unix()
	)
	stmt := "SELECT COUNT(id) FROM discussions WHERE created_at > ? AND created_at < ?;"
	return Count(db, stmt, from, to)
}

// Active users/day
func (db *datasource) GetUserActive(day time.Time) (int, error) {
	var (
		from = Bod(day).Unix()
		to = Bod(day).Add(24 * time.Hour).Unix()
	)
	stmt := "SELECT COUNT(id) FROM users WHERE last_login > ? AND last_login < ?;"
	return Count(db, stmt, from, to)
}

// New users/day, range [from, to]
func (db *datasource) GetNewUserDaily(from, to time.Time) ([]*model.DailyCount, error) {
	var (
		fromUnix = from.Unix()
		toUnix = to.Unix()
	)
	stmt := sql.Lookup(db.driver, "counter.new-user-daily")
	data := make([]*model.DailyCount, 0)
	err := meddler.QueryAll(db, &data, stmt, fromUnix, toUnix)
	return data, err
}

// New posts/day, range [from, to]
func (db *datasource) GetNewDiscussionDaily(from, to time.Time) ([]*model.DailyCount, error) {
	var (
		fromUnix = from.Unix()
		toUnix = to.Unix()
	)
	stmt := sql.Lookup(db.driver, "counter.new-post-daily")
	data := make([]*model.DailyCount, 0)
	err := meddler.QueryAll(db, &data, stmt, fromUnix, toUnix)
	return data, err
}

// Active users/day, range [from, to]
func (db *datasource) GetUserActiveDaily(from, to time.Time) ([]*model.DailyCount, error) {
	var (
		fromUnix = from.Unix()
		toUnix = to.Unix()
	)
	stmt := sql.Lookup(db.driver, "counter.active-user-daily")
	data := make([]*model.DailyCount, 0)
	err := meddler.QueryAll(db, &data, stmt, fromUnix, toUnix)
	return data, err
}

