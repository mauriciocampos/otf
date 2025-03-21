// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

const deleteTeamMembership = `-- name: DeleteTeamMembership :many
WITH
    users AS (
        SELECT username
        FROM unnest($2::text[]) t(username)
    )
DELETE
FROM team_memberships tm
USING users
WHERE
    tm.username = users.username AND
    tm.team_id  = $1
RETURNING tm.username
`

type DeleteTeamMembershipParams struct {
	TeamID    resource.TfeID
	Usernames []pgtype.Text
}

func (q *Queries) DeleteTeamMembership(ctx context.Context, db DBTX, arg DeleteTeamMembershipParams) ([]pgtype.Text, error) {
	rows, err := db.Query(ctx, deleteTeamMembership, arg.TeamID, arg.Usernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var username pgtype.Text
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteTokenByID = `-- name: DeleteTokenByID :one
DELETE
FROM tokens
WHERE token_id = $1
RETURNING token_id
`

func (q *Queries) DeleteTokenByID(ctx context.Context, db DBTX, tokenID resource.TfeID) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteTokenByID, tokenID)
	var token_id resource.TfeID
	err := row.Scan(&token_id)
	return token_id, err
}

const deleteUserByID = `-- name: DeleteUserByID :one
DELETE
FROM users
WHERE user_id = $1
RETURNING user_id
`

func (q *Queries) DeleteUserByID(ctx context.Context, db DBTX, userID resource.TfeID) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteUserByID, userID)
	var user_id resource.TfeID
	err := row.Scan(&user_id)
	return user_id, err
}

const deleteUserByUsername = `-- name: DeleteUserByUsername :one
DELETE
FROM users
WHERE username = $1
RETURNING user_id
`

func (q *Queries) DeleteUserByUsername(ctx context.Context, db DBTX, username pgtype.Text) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteUserByUsername, username)
	var user_id resource.TfeID
	err := row.Scan(&user_id)
	return user_id, err
}

const findTokenByID = `-- name: FindTokenByID :one
SELECT token_id, created_at, description, username
FROM tokens
WHERE token_id = $1
`

func (q *Queries) FindTokenByID(ctx context.Context, db DBTX, tokenID resource.TfeID) (Token, error) {
	row := db.QueryRow(ctx, findTokenByID, tokenID)
	var i Token
	err := row.Scan(
		&i.TokenID,
		&i.CreatedAt,
		&i.Description,
		&i.Username,
	)
	return i, err
}

const findTokensByUsername = `-- name: FindTokensByUsername :many
SELECT token_id, created_at, description, username
FROM tokens
WHERE username = $1
`

func (q *Queries) FindTokensByUsername(ctx context.Context, db DBTX, username pgtype.Text) ([]Token, error) {
	rows, err := db.Query(ctx, findTokensByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err := rows.Scan(
			&i.TokenID,
			&i.CreatedAt,
			&i.Description,
			&i.Username,
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

const findUserByAuthenticationTokenID = `-- name: FindUserByAuthenticationTokenID :one
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
JOIN tokens t ON u.username = t.username
WHERE t.token_id = $1
`

type FindUserByAuthenticationTokenIDRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUserByAuthenticationTokenID(ctx context.Context, db DBTX, tokenID resource.TfeID) (FindUserByAuthenticationTokenIDRow, error) {
	row := db.QueryRow(ctx, findUserByAuthenticationTokenID, tokenID)
	var i FindUserByAuthenticationTokenIDRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SiteAdmin,
		&i.Teams,
	)
	return i, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
WHERE u.user_id = $1
`

type FindUserByIDRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUserByID(ctx context.Context, db DBTX, userID resource.TfeID) (FindUserByIDRow, error) {
	row := db.QueryRow(ctx, findUserByID, userID)
	var i FindUserByIDRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SiteAdmin,
		&i.Teams,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
WHERE u.username = $1
`

type FindUserByUsernameRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUserByUsername(ctx context.Context, db DBTX, username pgtype.Text) (FindUserByUsernameRow, error) {
	row := db.QueryRow(ctx, findUserByUsername, username)
	var i FindUserByUsernameRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SiteAdmin,
		&i.Teams,
	)
	return i, err
}

const findUsers = `-- name: FindUsers :many
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
`

type FindUsersRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUsers(ctx context.Context, db DBTX) ([]FindUsersRow, error) {
	rows, err := db.Query(ctx, findUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUsersRow
	for rows.Next() {
		var i FindUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SiteAdmin,
			&i.Teams,
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

const findUsersByOrganization = `-- name: FindUsersByOrganization :many
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
JOIN team_memberships tm USING (username)
JOIN teams t USING (team_id)
WHERE t.organization_name = $1
GROUP BY u.user_id
`

type FindUsersByOrganizationRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUsersByOrganization(ctx context.Context, db DBTX, organizationName resource.OrganizationName) ([]FindUsersByOrganizationRow, error) {
	rows, err := db.Query(ctx, findUsersByOrganization, organizationName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUsersByOrganizationRow
	for rows.Next() {
		var i FindUsersByOrganizationRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SiteAdmin,
			&i.Teams,
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

const findUsersByTeamID = `-- name: FindUsersByTeamID :many
SELECT
    u.user_id, u.username, u.created_at, u.updated_at, u.site_admin,
    (
        SELECT array_agg(t.*)::teams[]
        FROM teams t
        JOIN team_memberships tm USING (team_id)
        WHERE tm.username = u.username
        GROUP BY tm.username
    ) AS teams
FROM users u
JOIN team_memberships tm USING (username)
JOIN teams t USING (team_id)
WHERE t.team_id = $1
GROUP BY u.user_id
`

type FindUsersByTeamIDRow struct {
	UserID    resource.TfeID
	Username  pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	SiteAdmin pgtype.Bool
	Teams     []TeamModel
}

func (q *Queries) FindUsersByTeamID(ctx context.Context, db DBTX, teamID resource.TfeID) ([]FindUsersByTeamIDRow, error) {
	rows, err := db.Query(ctx, findUsersByTeamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUsersByTeamIDRow
	for rows.Next() {
		var i FindUsersByTeamIDRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SiteAdmin,
			&i.Teams,
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

const insertTeamMembership = `-- name: InsertTeamMembership :many

WITH
    users AS (
        SELECT username
        FROM unnest($2::text[]) t(username)
    )
INSERT INTO team_memberships (username, team_id)
SELECT username, $1
FROM users
RETURNING username
`

type InsertTeamMembershipParams struct {
	TeamID    resource.TfeID
	Usernames []pgtype.Text
}

// team membership
func (q *Queries) InsertTeamMembership(ctx context.Context, db DBTX, arg InsertTeamMembershipParams) ([]pgtype.Text, error) {
	rows, err := db.Query(ctx, insertTeamMembership, arg.TeamID, arg.Usernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var username pgtype.Text
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertToken = `-- name: InsertToken :exec

INSERT INTO tokens (
    token_id,
    created_at,
    description,
    username
) VALUES (
    $1,
    $2,
    $3,
    $4
)
`

type InsertTokenParams struct {
	TokenID     resource.TfeID
	CreatedAt   pgtype.Timestamptz
	Description pgtype.Text
	Username    pgtype.Text
}

// user tokens
func (q *Queries) InsertToken(ctx context.Context, db DBTX, arg InsertTokenParams) error {
	_, err := db.Exec(ctx, insertToken,
		arg.TokenID,
		arg.CreatedAt,
		arg.Description,
		arg.Username,
	)
	return err
}

const insertUser = `-- name: InsertUser :exec
INSERT INTO users (
    user_id,
    created_at,
    updated_at,
    username
) VALUES (
    $1,
    $2,
    $3,
    $4
)
`

type InsertUserParams struct {
	ID        resource.TfeID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Username  pgtype.Text
}

func (q *Queries) InsertUser(ctx context.Context, db DBTX, arg InsertUserParams) error {
	_, err := db.Exec(ctx, insertUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Username,
	)
	return err
}

const resetUserSiteAdmins = `-- name: ResetUserSiteAdmins :many
UPDATE users
SET site_admin = false
WHERE site_admin = true
RETURNING username
`

func (q *Queries) ResetUserSiteAdmins(ctx context.Context, db DBTX) ([]pgtype.Text, error) {
	rows, err := db.Query(ctx, resetUserSiteAdmins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var username pgtype.Text
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserSiteAdmins = `-- name: UpdateUserSiteAdmins :many
UPDATE users
SET site_admin = true
WHERE username = ANY($1::text[])
RETURNING username
`

func (q *Queries) UpdateUserSiteAdmins(ctx context.Context, db DBTX, usernames []pgtype.Text) ([]pgtype.Text, error) {
	rows, err := db.Query(ctx, updateUserSiteAdmins, usernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var username pgtype.Text
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
