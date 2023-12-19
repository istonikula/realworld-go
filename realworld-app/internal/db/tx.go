package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TxMgr struct {
	DB *sqlx.DB
}

func (m *TxMgr) Read(fn func(*sqlx.Tx) error) error {
	return m.run(true, fn)
}

func (m *TxMgr) Write(fn func(*sqlx.Tx) error) error {
	return m.run(false, fn)
}

func (m *TxMgr) run(readOnly bool, fn func(*sqlx.Tx) error) error {
	tx, err := m.DB.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: readOnly})
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
