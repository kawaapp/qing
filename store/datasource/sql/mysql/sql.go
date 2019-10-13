package mysql

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string {
	"counter.new-user-daily":    newUserDaily,
	"counter.new-post-daily":    newPostDaily,
	"counter.active-user-daily": activeUserDaily,
}

const newUserDaily = `
	SELECT
		COUNT(id) AS count, DATE(FROM_UNIXTIME(created_at)) AS date
	FROM users
	WHERE created_at >= ? AND created_at <= ?
	GROUP BY DATE(FROM_UNIXTIME(created_at))
	ORDER BY date ASC
;`

const newPostDaily = `
	SELECT
		COUNT(*) AS count, DATE(FROM_UNIXTIME(created_at)) AS date
	FROM discussions
	WHERE created_at >= ? AND created_at <= ?
	GROUP BY DATE(FROM_UNIXTIME(created_at))
	ORDER BY date ASC
;`

const activeUserDaily = `
	SELECT
		COUNT(*) AS count, DATE(FROM_UNIXTIME(last_login)) AS date
	FROM users
	WHERE last_login >= ? AND last_login <= ?
	GROUP BY DATE(FROM_UNIXTIME(last_login))
	ORDER BY date ASC
;`