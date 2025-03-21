// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package module

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

const deleteModuleByID = `-- name: DeleteModuleByID :one
DELETE
FROM modules
WHERE module_id = $1
RETURNING module_id
`

func (q *Queries) DeleteModuleByID(ctx context.Context, db DBTX, moduleID resource.TfeID) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteModuleByID, moduleID)
	var module_id resource.TfeID
	err := row.Scan(&module_id)
	return module_id, err
}

const deleteModuleVersionByID = `-- name: DeleteModuleVersionByID :one
DELETE
FROM module_versions
WHERE module_version_id = $1
RETURNING module_version_id
`

func (q *Queries) DeleteModuleVersionByID(ctx context.Context, db DBTX, moduleVersionID resource.TfeID) (resource.TfeID, error) {
	row := db.QueryRow(ctx, deleteModuleVersionByID, moduleVersionID)
	var module_version_id resource.TfeID
	err := row.Scan(&module_version_id)
	return module_version_id, err
}

const findModuleByConnection = `-- name: FindModuleByConnection :one
SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    r.vcs_provider_id,
    r.repo_path,
    (
        SELECT array_agg(v.*)::module_versions[]
        FROM module_versions v
        WHERE v.module_id = m.module_id
        GROUP BY v.module_id
    ) AS module_versions
FROM modules m
JOIN repo_connections r USING (module_id)
WHERE r.vcs_provider_id = $1
AND   r.repo_path = $2
`

type FindModuleByConnectionParams struct {
	VCSProviderID resource.TfeID
	RepoPath      pgtype.Text
}

type FindModuleByConnectionRow struct {
	ModuleID         resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
	VCSProviderID    resource.TfeID
	RepoPath         pgtype.Text
	ModuleVersions   []ModuleVersionModel
}

func (q *Queries) FindModuleByConnection(ctx context.Context, db DBTX, arg FindModuleByConnectionParams) (FindModuleByConnectionRow, error) {
	row := db.QueryRow(ctx, findModuleByConnection, arg.VCSProviderID, arg.RepoPath)
	var i FindModuleByConnectionRow
	err := row.Scan(
		&i.ModuleID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.OrganizationName,
		&i.VCSProviderID,
		&i.RepoPath,
		&i.ModuleVersions,
	)
	return i, err
}

const findModuleByID = `-- name: FindModuleByID :one
SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    r.vcs_provider_id,
    r.repo_path,
    (
        SELECT array_agg(v.*)::module_versions[]
        FROM module_versions v
        WHERE v.module_id = m.module_id
        GROUP BY v.module_id
    ) AS module_versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.module_id = $1
`

type FindModuleByIDRow struct {
	ModuleID         resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
	VCSProviderID    resource.TfeID
	RepoPath         pgtype.Text
	ModuleVersions   []ModuleVersionModel
}

func (q *Queries) FindModuleByID(ctx context.Context, db DBTX, id resource.TfeID) (FindModuleByIDRow, error) {
	row := db.QueryRow(ctx, findModuleByID, id)
	var i FindModuleByIDRow
	err := row.Scan(
		&i.ModuleID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.OrganizationName,
		&i.VCSProviderID,
		&i.RepoPath,
		&i.ModuleVersions,
	)
	return i, err
}

const findModuleByModuleVersionID = `-- name: FindModuleByModuleVersionID :one
SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    r.vcs_provider_id,
    r.repo_path,
    (
        SELECT array_agg(v.*)::module_versions[]
        FROM module_versions v
        WHERE v.module_id = m.module_id
        GROUP BY v.module_id
    ) AS module_versions
FROM modules m
JOIN module_versions mv USING (module_id)
LEFT JOIN repo_connections r USING (module_id)
WHERE mv.module_version_id = $1
`

type FindModuleByModuleVersionIDRow struct {
	ModuleID         resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
	VCSProviderID    resource.TfeID
	RepoPath         pgtype.Text
	ModuleVersions   []ModuleVersionModel
}

func (q *Queries) FindModuleByModuleVersionID(ctx context.Context, db DBTX, moduleVersionID resource.TfeID) (FindModuleByModuleVersionIDRow, error) {
	row := db.QueryRow(ctx, findModuleByModuleVersionID, moduleVersionID)
	var i FindModuleByModuleVersionIDRow
	err := row.Scan(
		&i.ModuleID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.OrganizationName,
		&i.VCSProviderID,
		&i.RepoPath,
		&i.ModuleVersions,
	)
	return i, err
}

const findModuleByName = `-- name: FindModuleByName :one
SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    r.vcs_provider_id,
    r.repo_path,
    (
        SELECT array_agg(v.*)::module_versions[]
        FROM module_versions v
        WHERE v.module_id = m.module_id
        GROUP BY v.module_id
    ) AS module_versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.organization_name = $1
AND   m.name = $2
AND   m.provider = $3
`

type FindModuleByNameParams struct {
	OrganizationName resource.OrganizationName
	Name             pgtype.Text
	Provider         pgtype.Text
}

type FindModuleByNameRow struct {
	ModuleID         resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
	VCSProviderID    resource.TfeID
	RepoPath         pgtype.Text
	ModuleVersions   []ModuleVersionModel
}

func (q *Queries) FindModuleByName(ctx context.Context, db DBTX, arg FindModuleByNameParams) (FindModuleByNameRow, error) {
	row := db.QueryRow(ctx, findModuleByName, arg.OrganizationName, arg.Name, arg.Provider)
	var i FindModuleByNameRow
	err := row.Scan(
		&i.ModuleID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.OrganizationName,
		&i.VCSProviderID,
		&i.RepoPath,
		&i.ModuleVersions,
	)
	return i, err
}

const findModuleTarball = `-- name: FindModuleTarball :one
SELECT tarball
FROM module_tarballs
WHERE module_version_id = $1
`

func (q *Queries) FindModuleTarball(ctx context.Context, db DBTX, moduleVersionID resource.TfeID) ([]byte, error) {
	row := db.QueryRow(ctx, findModuleTarball, moduleVersionID)
	var tarball []byte
	err := row.Scan(&tarball)
	return tarball, err
}

const insertModule = `-- name: InsertModule :exec
INSERT INTO modules (
    module_id,
    created_at,
    updated_at,
    name,
    provider,
    status,
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

type InsertModuleParams struct {
	ID               resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
}

func (q *Queries) InsertModule(ctx context.Context, db DBTX, arg InsertModuleParams) error {
	_, err := db.Exec(ctx, insertModule,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Provider,
		arg.Status,
		arg.OrganizationName,
	)
	return err
}

const insertModuleTarball = `-- name: InsertModuleTarball :one
INSERT INTO module_tarballs (
    tarball,
    module_version_id
) VALUES (
    $1,
    $2
)
RETURNING module_version_id
`

type InsertModuleTarballParams struct {
	Tarball         []byte
	ModuleVersionID resource.TfeID
}

func (q *Queries) InsertModuleTarball(ctx context.Context, db DBTX, arg InsertModuleTarballParams) (resource.TfeID, error) {
	row := db.QueryRow(ctx, insertModuleTarball, arg.Tarball, arg.ModuleVersionID)
	var module_version_id resource.TfeID
	err := row.Scan(&module_version_id)
	return module_version_id, err
}

const insertModuleVersion = `-- name: InsertModuleVersion :one
INSERT INTO module_versions (
    module_version_id,
    version,
    created_at,
    updated_at,
    module_id,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING module_version_id, version, created_at, updated_at, status, status_error, module_id
`

type InsertModuleVersionParams struct {
	ModuleVersionID resource.TfeID
	Version         pgtype.Text
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	ModuleID        resource.TfeID
	Status          pgtype.Text
}

func (q *Queries) InsertModuleVersion(ctx context.Context, db DBTX, arg InsertModuleVersionParams) (ModuleVersionModel, error) {
	row := db.QueryRow(ctx, insertModuleVersion,
		arg.ModuleVersionID,
		arg.Version,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.ModuleID,
		arg.Status,
	)
	var i ModuleVersionModel
	err := row.Scan(
		&i.ModuleVersionID,
		&i.Version,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Status,
		&i.StatusError,
		&i.ModuleID,
	)
	return i, err
}

const listModulesByOrganization = `-- name: ListModulesByOrganization :many
SELECT
    m.module_id,
    m.created_at,
    m.updated_at,
    m.name,
    m.provider,
    m.status,
    m.organization_name,
    r.vcs_provider_id,
    r.repo_path,
    (
        SELECT array_agg(v.*)::module_versions[]
        FROM module_versions v
        WHERE v.module_id = m.module_id
        GROUP BY v.module_id
    ) AS module_versions
FROM modules m
LEFT JOIN repo_connections r USING (module_id)
WHERE m.organization_name = $1
`

type ListModulesByOrganizationRow struct {
	ModuleID         resource.TfeID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Provider         pgtype.Text
	Status           pgtype.Text
	OrganizationName resource.OrganizationName
	VCSProviderID    resource.TfeID
	RepoPath         pgtype.Text
	ModuleVersions   []ModuleVersionModel
}

func (q *Queries) ListModulesByOrganization(ctx context.Context, db DBTX, organizationName resource.OrganizationName) ([]ListModulesByOrganizationRow, error) {
	rows, err := db.Query(ctx, listModulesByOrganization, organizationName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListModulesByOrganizationRow
	for rows.Next() {
		var i ListModulesByOrganizationRow
		if err := rows.Scan(
			&i.ModuleID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Provider,
			&i.Status,
			&i.OrganizationName,
			&i.VCSProviderID,
			&i.RepoPath,
			&i.ModuleVersions,
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

const updateModuleStatusByID = `-- name: UpdateModuleStatusByID :one
UPDATE modules
SET status = $1
WHERE module_id = $2
RETURNING module_id
`

type UpdateModuleStatusByIDParams struct {
	Status   pgtype.Text
	ModuleID resource.TfeID
}

func (q *Queries) UpdateModuleStatusByID(ctx context.Context, db DBTX, arg UpdateModuleStatusByIDParams) (resource.TfeID, error) {
	row := db.QueryRow(ctx, updateModuleStatusByID, arg.Status, arg.ModuleID)
	var module_id resource.TfeID
	err := row.Scan(&module_id)
	return module_id, err
}

const updateModuleVersionStatusByID = `-- name: UpdateModuleVersionStatusByID :one
UPDATE module_versions
SET
    status = $1,
    status_error = $2
WHERE module_version_id = $3
RETURNING module_version_id, version, created_at, updated_at, status, status_error, module_id
`

type UpdateModuleVersionStatusByIDParams struct {
	Status          pgtype.Text
	StatusError     pgtype.Text
	ModuleVersionID resource.TfeID
}

func (q *Queries) UpdateModuleVersionStatusByID(ctx context.Context, db DBTX, arg UpdateModuleVersionStatusByIDParams) (ModuleVersionModel, error) {
	row := db.QueryRow(ctx, updateModuleVersionStatusByID, arg.Status, arg.StatusError, arg.ModuleVersionID)
	var i ModuleVersionModel
	err := row.Scan(
		&i.ModuleVersionID,
		&i.Version,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Status,
		&i.StatusError,
		&i.ModuleID,
	)
	return i, err
}
