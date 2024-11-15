package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storer interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	SelectOne(ctx context.Context, query string, args ...any) (*sql.Row, error)
	Close() error
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbfile string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) SelectOne(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(args...)
	return row, row.Err()
}
