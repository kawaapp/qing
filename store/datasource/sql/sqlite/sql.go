package sqlite

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"counter.new-user-daily":    newUserDaily,
	"counter.new-post-daily":    newPostDaily,
	"counter.active-user-daily": activeUserDaily,
}

const newUserDaily = `
	SELECT
		COUNT(id) AS count, strftime('%Y-%m-%d', created_at, 'unixepoch') AS date
	FROM users
	WHERE created_at >= ? AND created_at <= ?
	GROUP BY strftime('%Y-%m-%d', created_at, 'unixepoch')
	ORDER BY date ASC
;`

const newPostDaily = `
	SELECT
		COUNT(id) AS count, strftime('%Y-%m-%d', created_at, 'unixepoch') AS date
	FROM discussions
	WHERE created_at >= ? AND created_at <= ?
	GROUP BY strftime('%Y-%m-%d', created_at, 'unixepoch')
	ORDER BY date ASC
;`

const activeUserDaily = `
	SELECT
		COUNT(*) AS count, strftime('%Y-%m-%d', last_login, 'unixepoch') AS date
	FROM users
	WHERE last_login >= ? AND last_login <= ?
	GROUP BY strftime('%Y-%m-%d', last_login, 'unixepoch')
	ORDER BY date ASC
;`