package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
)

func (db *datasource) CreateFollow(f *model.Follow) error {
	f.CreatedAt = UnixNow()
	return meddler.Insert(db, "follows", f)
}

func (db *datasource) GetFollow(uid, fid int64) (*model.Follow, error) {
	f := new(model.Follow)
	err := meddler.QueryRow(db, f, sqlGetFollow, uid, fid)
	return f, err
}

func (db *datasource) DeleteFollow(uid, fid int64) error {
	return Delete(db, sqlDeleteFollow, uid, fid)
}

func (db *datasource) GetFollowerList(uid int64, page, size int) ([]*model.User, error) {
	if page == 0 {
		page = 1
	}
	users := make([]*model.User, 0)
	err := meddler.QueryAll(db, &users, sqlGetFollowerList, uid, size, (page-1)*size)
	return users, err
}

func (db *datasource) GetFollowingList(uid int64, page, size int) ([]*model.User, error) {
	if page == 0 {
		page = 1
	}
	users := make([]*model.User, 0)
	err := meddler.QueryAll(db, &users, sqlGetFollowingList, uid, size, (page-1)*size)
	return users, err
}

const sqlDeleteFollow = `
DELETE FROM follows WHERE user_id=? AND follower_id=?
`

const sqlGetFollow = `
SELECT
	id,
	created_at,
	user_id,
	follower_id
FROM follows WHERE user_id=? AND follower_id=?
`

const sqlGetFollowerList = sqlUserSelect + `
FROM users
JOIN follows ON(users.id = follows.follower_id)
WHERE follows.user_id=?
ORDER BY id DESC LIMIT ? OFFSET ?
`

const sqlGetFollowingList = sqlUserSelect + `
FROM users
JOIN follows ON(users.id = follows.user_id)
WHERE follows.follower_id=?
ORDER BY id DESC LIMIT ? OFFSET ?
`
