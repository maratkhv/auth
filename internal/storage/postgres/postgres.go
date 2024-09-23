package postgres

import (
	"auth/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

var (
	ErrUserNotExist = errors.New("user does not exist")
)

func New(connStr string) Database {
	db := sqlx.MustConnect("pgx", connStr)
	return Database{db: db}
}

func (d Database) InsertUser(login string, password []byte) (*models.User, error) {
	tx, err := d.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into users(login, password) values($1, $2)", login, password)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	var u models.User
	err = tx.QueryRow("select * from users where login=$1", login).Scan(&u.Id, &u.Login, &u.PassHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}

	err = tx.Commit()
	return &u, err
}

func (d Database) GetUser(login string) (*models.User, error) {
	var u models.User
	query, args, err := sq.Select("*").From("users").Where(sq.Eq{"login": login}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	err = d.db.QueryRow(query, args...).Scan(&u.Id, &u.Login, &u.PassHash)
	if errors.Is(err, sql.ErrNoRows) {
		return &u, ErrUserNotExist
	}
	return &u, err
}

func (d Database) Close() error {
	return d.db.Close()
}
