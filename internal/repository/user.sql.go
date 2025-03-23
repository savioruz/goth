// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, password, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error) {
	row := db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL RETURNING id, email, password, created_at, updated_at, deleted_at
`

func (q *Queries) DeleteUser(ctx context.Context, db DBTX, id pgtype.UUID) (User, error) {
	row := db.QueryRow(ctx, deleteUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, password, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, db DBTX, email string) (User, error) {
	row := db.QueryRow(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, password, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC
`

func (q *Queries) ListUsers(ctx context.Context, db DBTX) ([]User, error) {
	rows, err := db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET email = $2, password = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING id, email, password, created_at, updated_at, deleted_at
`

type UpdateUserParams struct {
	ID       pgtype.UUID `json:"id"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

func (q *Queries) UpdateUser(ctx context.Context, db DBTX, arg UpdateUserParams) (User, error) {
	row := db.QueryRow(ctx, updateUser, arg.ID, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
