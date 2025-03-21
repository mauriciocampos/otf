// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package vcsprovider

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

const deleteVCSProviderByID = `-- name: DeleteVCSProviderByID :one
DELETE
FROM vcs_providers
WHERE vcs_provider_id = $1
RETURNING vcs_provider_id
`

func (q *Queries) DeleteVCSProviderByID(ctx context.Context, db DBTX, vcsProviderID resource.TfeID) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteVCSProviderByID, vcsProviderID)
	var vcs_provider_id resource.TfeID
	err := row.Scan(&vcs_provider_id)
	return vcs_provider_id, err
}

const findVCSProvider = `-- name: FindVCSProvider :one
SELECT
    v.vcs_provider_id, v.token, v.created_at, v.name, v.vcs_kind, v.organization_name, v.github_app_id,
    (ga.*)::"github_apps" AS github_app,
    (gi.*)::"github_app_installs" AS github_app_install
FROM vcs_providers v
LEFT JOIN (github_app_installs gi JOIN github_apps ga USING (github_app_id)) USING (vcs_provider_id)
WHERE v.vcs_provider_id = $1
`

type FindVCSProviderRow struct {
	VCSProviderID    resource.TfeID
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	OrganizationName resource.OrganizationName
	GithubAppID      pgtype.Int8
	GithubApp        *GithubApp
	GithubAppInstall *GithubAppInstall
}

func (q *Queries) FindVCSProvider(ctx context.Context, db DBTX, vcsProviderID resource.TfeID) (FindVCSProviderRow, error) {
	row := db.QueryRow(ctx, findVCSProvider, vcsProviderID)
	var i FindVCSProviderRow
	err := row.Scan(
		&i.VCSProviderID,
		&i.Token,
		&i.CreatedAt,
		&i.Name,
		&i.VCSKind,
		&i.OrganizationName,
		&i.GithubAppID,
		&i.GithubApp,
		&i.GithubAppInstall,
	)
	return i, err
}

const findVCSProviderForUpdate = `-- name: FindVCSProviderForUpdate :one
SELECT
    v.vcs_provider_id, v.token, v.created_at, v.name, v.vcs_kind, v.organization_name, v.github_app_id,
    (ga.*)::"github_apps" AS github_app,
    (gi.*)::"github_app_installs" AS github_app_install
FROM vcs_providers v
LEFT JOIN (github_app_installs gi JOIN github_apps ga USING (github_app_id)) USING (vcs_provider_id)
WHERE v.vcs_provider_id = $1
FOR UPDATE OF v
`

type FindVCSProviderForUpdateRow struct {
	VCSProviderID    resource.TfeID
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	OrganizationName resource.OrganizationName
	GithubAppID      pgtype.Int8
	GithubApp        *GithubApp
	GithubAppInstall *GithubAppInstall
}

func (q *Queries) FindVCSProviderForUpdate(ctx context.Context, db DBTX, vcsProviderID resource.TfeID) (FindVCSProviderForUpdateRow, error) {
	row := db.QueryRow(ctx, findVCSProviderForUpdate, vcsProviderID)
	var i FindVCSProviderForUpdateRow
	err := row.Scan(
		&i.VCSProviderID,
		&i.Token,
		&i.CreatedAt,
		&i.Name,
		&i.VCSKind,
		&i.OrganizationName,
		&i.GithubAppID,
		&i.GithubApp,
		&i.GithubAppInstall,
	)
	return i, err
}

const findVCSProviders = `-- name: FindVCSProviders :many
SELECT
    v.vcs_provider_id, v.token, v.created_at, v.name, v.vcs_kind, v.organization_name, v.github_app_id,
    (ga.*)::"github_apps" AS github_app,
    (gi.*)::"github_app_installs" AS github_app_install
FROM vcs_providers v
LEFT JOIN (github_app_installs gi JOIN github_apps ga USING (github_app_id)) USING (vcs_provider_id)
`

type FindVCSProvidersRow struct {
	VCSProviderID    resource.TfeID
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	OrganizationName resource.OrganizationName
	GithubAppID      pgtype.Int8
	GithubApp        *GithubApp
	GithubAppInstall *GithubAppInstall
}

func (q *Queries) FindVCSProviders(ctx context.Context, db DBTX) ([]FindVCSProvidersRow, error) {
	rows, err := db.Query(ctx, findVCSProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindVCSProvidersRow
	for rows.Next() {
		var i FindVCSProvidersRow
		if err := rows.Scan(
			&i.VCSProviderID,
			&i.Token,
			&i.CreatedAt,
			&i.Name,
			&i.VCSKind,
			&i.OrganizationName,
			&i.GithubAppID,
			&i.GithubApp,
			&i.GithubAppInstall,
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

const findVCSProvidersByGithubAppInstallID = `-- name: FindVCSProvidersByGithubAppInstallID :many
SELECT
    v.vcs_provider_id, v.token, v.created_at, v.name, v.vcs_kind, v.organization_name, v.github_app_id,
    (ga.*)::"github_apps" AS github_app,
    (gi.*)::"github_app_installs" AS github_app_install
FROM vcs_providers v
JOIN (github_app_installs gi JOIN github_apps ga USING (github_app_id)) USING (vcs_provider_id)
WHERE gi.install_id = $1
`

type FindVCSProvidersByGithubAppInstallIDRow struct {
	VCSProviderID    resource.TfeID
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	OrganizationName resource.OrganizationName
	GithubAppID      pgtype.Int8
	GithubApp        *GithubApp
	GithubAppInstall *GithubAppInstall
}

func (q *Queries) FindVCSProvidersByGithubAppInstallID(ctx context.Context, db DBTX, installID pgtype.Int8) ([]FindVCSProvidersByGithubAppInstallIDRow, error) {
	rows, err := db.Query(ctx, findVCSProvidersByGithubAppInstallID, installID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindVCSProvidersByGithubAppInstallIDRow
	for rows.Next() {
		var i FindVCSProvidersByGithubAppInstallIDRow
		if err := rows.Scan(
			&i.VCSProviderID,
			&i.Token,
			&i.CreatedAt,
			&i.Name,
			&i.VCSKind,
			&i.OrganizationName,
			&i.GithubAppID,
			&i.GithubApp,
			&i.GithubAppInstall,
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

const findVCSProvidersByOrganization = `-- name: FindVCSProvidersByOrganization :many
SELECT
    v.vcs_provider_id, v.token, v.created_at, v.name, v.vcs_kind, v.organization_name, v.github_app_id,
    (ga.*)::"github_apps" AS github_app,
    (gi.*)::"github_app_installs" AS github_app_install
FROM vcs_providers v
LEFT JOIN (github_app_installs gi JOIN github_apps ga USING (github_app_id)) USING (vcs_provider_id)
WHERE v.organization_name = $1
`

type FindVCSProvidersByOrganizationRow struct {
	VCSProviderID    resource.TfeID
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	OrganizationName resource.OrganizationName
	GithubAppID      pgtype.Int8
	GithubApp        *GithubApp
	GithubAppInstall *GithubAppInstall
}

func (q *Queries) FindVCSProvidersByOrganization(ctx context.Context, db DBTX, organizationName resource.OrganizationName) ([]FindVCSProvidersByOrganizationRow, error) {
	rows, err := db.Query(ctx, findVCSProvidersByOrganization, organizationName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindVCSProvidersByOrganizationRow
	for rows.Next() {
		var i FindVCSProvidersByOrganizationRow
		if err := rows.Scan(
			&i.VCSProviderID,
			&i.Token,
			&i.CreatedAt,
			&i.Name,
			&i.VCSKind,
			&i.OrganizationName,
			&i.GithubAppID,
			&i.GithubApp,
			&i.GithubAppInstall,
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

const insertVCSProvider = `-- name: InsertVCSProvider :exec
INSERT INTO vcs_providers (
    vcs_provider_id,
    created_at,
    name,
    vcs_kind,
    token,
    github_app_id,
    organization_name
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
`

type InsertVCSProviderParams struct {
	VCSProviderID    resource.TfeID
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	VCSKind          pgtype.Text
	Token            pgtype.Text
	GithubAppID      pgtype.Int8
	OrganizationName resource.OrganizationName
}

func (q *Queries) InsertVCSProvider(ctx context.Context, db DBTX, arg InsertVCSProviderParams) error {
	_, err := db.Exec(ctx, insertVCSProvider,
		arg.VCSProviderID,
		arg.CreatedAt,
		arg.Name,
		arg.VCSKind,
		arg.Token,
		arg.GithubAppID,
		arg.OrganizationName,
	)
	return err
}

const updateVCSProvider = `-- name: UpdateVCSProvider :one
UPDATE vcs_providers
SET name = $1, token = $2
WHERE vcs_provider_id = $3
RETURNING vcs_provider_id, token, created_at, name, vcs_kind, organization_name, github_app_id
`

type UpdateVCSProviderParams struct {
	Name          pgtype.Text
	Token         pgtype.Text
	VCSProviderID resource.TfeID
}

func (q *Queries) UpdateVCSProvider(ctx context.Context, db DBTX, arg UpdateVCSProviderParams) (Model, error) {
	row := db.QueryRow(ctx, updateVCSProvider, arg.Name, arg.Token, arg.VCSProviderID)
	var i Model
	err := row.Scan(
		&i.VCSProviderID,
		&i.Token,
		&i.CreatedAt,
		&i.Name,
		&i.VCSKind,
		&i.OrganizationName,
		&i.GithubAppID,
	)
	return i, err
}
