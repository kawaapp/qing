package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
	"database/sql"
)

func (db *datasource) GetUserList(params model.QueryParams, page, size int) ([]*model.User, error) {
	users := make([]*model.User, 0)
	query, args := sqlUserQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &users, query, args...)
	return users, err
}

func (db *datasource) GetUserCount(params model.QueryParams) (int, error) {
	query, args := sqlUserQuery("SELECT COUNT(id)", params, 0, 0)
	num, err := db.Count(query, args...)
	return num, err
}

func (db *datasource) GetUser(id int64) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.Load(db, "users", usr, id)
	return usr, err
}

func (db *datasource) GetUserByLogin(login string) (*model.User, error) {
	stmt := sqlUserFindLogin
	data := new(model.User)
	err := meddler.QueryRow(db, data, stmt, login)
	return data, err
}

func (db *datasource) GetUserByPhone(phone string) (*model.User, error) {
	stmt := sqlFindUserByPhone
	data := new(model.User)
	err := meddler.QueryRow(db, data, stmt, phone)
	return data, err
}

func (db *datasource) GetUserIdList(ids []int64) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	stmt := sqlUserFindByIds
	q := joinIntArray(ids)
	data := make([]*model.User, 0)
	err := meddler.QueryAll(db, &data, fmt.Sprintf(stmt, q))
	return data, err
}

func (db *datasource) CreateUser(usr *model.User) error {
	usr.CreatedAt = UnixNow()
	usr.UpdatedAt = UnixNow()
	return meddler.Insert(db, "users", usr)
}

func (db *datasource) UpdateUser(usr *model.User) error {
	usr.UpdatedAt = UnixNow()
	return meddler.Update(db, "users", usr)
}

func (db *datasource) UpdateUserSign(uid int64, count int) error {
	stmt := `UPDATE users SET sign_count=? WHERE id=?`
	_, err := db.Exec(stmt, count, uid)
	return err
}

func (db *datasource) UpdateUserExp(uid int64, exp int) error {
	stmt := `UPDATE users SET exp_count=? WHERE id=?`
	_, err := db.Exec(stmt, exp, uid)
	return err
}

func (db *datasource) DeleteUser(id int64) error {
	stmt := sqlUserDelete
	_, err := db.Exec(stmt, id)
	return err
}

// bind
func (db *datasource) GetUserByBind(bindId string) (*model.User, error) {
	stmt := sqlFindUserByBindId
	user := new(model.User)
	err := meddler.QueryRow(db, user, stmt, bindId)
	return user, err
}

func (db *datasource) GetUserByUnion(unionId string) (*model.User, error) {
	stmt := sqlFindUserByUnionId
	user := new(model.User)
	err := meddler.QueryRow(db, user, stmt, unionId)
	return user, err
}

func (db *datasource) CreateBind(bind *model.UserBind) error {
	bind.CreatedAt = UnixNow()
	return meddler.Insert(db, "user_binds", bind)
}

func (db *datasource) DeleteBind(bindId string) error {
	stmt := sqlBindDelete
	_, err := db.Exec(stmt, bindId)
	return err
}

func (db *datasource) DeleteBindByUser(uid int64) error  {
	stmt := sqlDeleteBindByUser
	_, err := db.Exec(stmt, uid)
	return err
}

// bind user , transaction
func (db *datasource) CreateBindUser(user *model.User, bind *model.UserBind) error {
	return db.Transact(func(tx *sql.Tx) error {
		user.CreatedAt = UnixNow()
		user.UpdatedAt = UnixNow()
		if err := meddler.Insert(tx, "users", user); err != nil {
			return err
		}
		bind.UserId = user.ID
		bind.CreatedAt = UnixNow()
		return meddler.Insert(tx, "user_binds", bind)
	})
}

func (db *datasource) DeleteBindUser(uid int64) error {
	return db.Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec(sqlUserDelete, uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(sqlDeleteBindByUser, uid)
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	})
}

func (db *datasource) getUserByIds(ids []int64) (map[int64]*model.User, error) {
	if len(ids) == 0 {
		return map[int64]*model.User{}, nil
	}

	stmt := sqlUserFindByIds
	q := joinIntArray(ids)
	data := make([]*model.User, 0)
	err := meddler.QueryAll(db, &data, fmt.Sprintf(stmt, q))

	users := make(map[int64]*model.User, 0)
	for _, user := range data {
		users[user.ID] = user
	}
	return users, err
}

const sqlUserDelete = `DELETE FROM users WHERE id=?;`
const sqlBindDelete = `DELETE FROM user_binds WHERE bind_id=?;`
const sqlDeleteBindByUser = `DELETE FROM user_binds WHERE user_id=?;`

const sqlUserSelect = `
SELECT
	users.id,
	users.created_at,
	sign_count,
	exp_count,
	last_login,
	login,
	nickname,
	email,
	phone,
	avatar,
	summary,
	hash,
	password_hash
`

// user
func sqlUserQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM users"

	where := ""
	if q, ok := params["login"]; ok {
		where += " AND login=?"
		args = append(args, q)
	}

	if 	q, ok := params["nickname"]; ok {
		where += " AND nickname LIKE ?"
		args = append(args, "%" + q + "%")
	}

	if _, ok := params["silence"]; ok {
		where += " AND silenced_at > ?"
	}

	if _, ok := params["block"]; ok {
		where += " AND blocked_at > 0"
	}

	if len(where) > 0 {
		query += " WHERE 1=1" + where
	}

	if size > 0 {
		query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", size, page * size)
	}
	return
}

const sqlUserFindByIds = sqlUserSelect + `
FROM users
WHERE id IN (%s);`

const sqlUserFindLogin = sqlUserSelect + `
FROM users
WHERE login = ?;`

const sqlFindUserByPhone = sqlUserSelect + `
FROM users
WHERE phone=? LIMIT 1
;`

const sqlFindUserByBindId = sqlUserSelect + `
FROM users
INNER JOIN user_binds
ON (users.id = user_binds.user_id AND user_binds.bind_id=?)
;`

const sqlFindUserByUnionId = sqlUserSelect + `
FROM users
INNER JOIN user_binds
ON (users.id = user_binds.user_id AND user_binds.union_id=?)
;`