package db

import (
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type table string

func (t table) insert(cols ...string) string {
	markerSql := "$1"
	if len(cols) > 1 {
		for i := range cols[1:] {
			markerSql += ", $" + strconv.Itoa(i+2)
		}
	}
	return "INSERT INTO " + string(t) + " (" + strings.Join(cols, ", ") + ") VALUES (" + markerSql + ")"
}

func (t table) queryIfExists(db *sqlx.DB, where string, params ...any) (bool, error) {
	stmt, err := db.Preparex("SELECT (COUNT(*) > 0) FROM " + string(t) + " WHERE " + where)
	if err != nil {
		return false, err
	}
	var exists bool
	err = stmt.QueryRowx(params...).Scan(&exists)
	return exists, err
}
