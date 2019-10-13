package datasource

import (
	"github.com/kawaapp/kawaqing/store/datasource/sql"
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

func (db *datasource) PutKV(key, value string) error {
	stmt := sql.Lookup(db.driver,"meta.put")
	_, err := db.Exec(stmt, key, value, value)
	return err
}

func (db *datasource) GetKV(key string) (string, error) {
	var (
		stmt = sqlMetaGet
		value string
	)
	row := db.QueryRow(stmt, key)
	err := row.Scan(&value)
	return value, err
}

func (db *datasource) GetMetaValue(key string) (string, error) {
	return db.GetKV(key)
}

func (db *datasource) DelKV(key string) (error)  {
	_, err := db.Exec(sqlMetaDelete, key)
	return err
}

func (db *datasource) GetMetaData() (map[string]string, error)  {
	stmt := sqlMetaList
	data := make([]*model.Pair, 0)
	err := meddler.QueryAll(db, &data, stmt)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, v := range data {
		m[v.Key] = v.Value
	}
	return m, nil
}

func (db *datasource) PutMetaData(kvs map[string]interface{}) (error) {
	for k, v := range kvs {
		str := fmt.Sprintf("%v", v)
		db.PutKV(k, str)
	}
	return nil
}


const sqlMetaGet = `
SELECT
	kv_value
FROM metadata
WHERE kv_key=?
;`

const sqlMetaDelete = `
DELETE
FROM metadata
WHERE kv_key=?
;`

const sqlMetaList = `
SELECT
	kv_key,
	kv_value
FROM metadata
;`