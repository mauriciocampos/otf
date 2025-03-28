-- name: InsertOrganization :exec
INSERT INTO organizations (
    organization_id,
    created_at,
    updated_at,
    name,
    email,
    collaborator_auth_policy,
    cost_estimation_enabled,
    session_remember,
    session_timeout,
    allow_force_delete_workspaces
) VALUES (
    sqlc.arg('id'),
    sqlc.arg('created_at'),
    sqlc.arg('updated_at'),
    sqlc.arg('name'),
    sqlc.arg('email'),
    sqlc.arg('collaborator_auth_policy'),
    sqlc.arg('cost_estimation_enabled'),
    sqlc.arg('session_remember'),
    sqlc.arg('session_timeout'),
    sqlc.arg('allow_force_delete_workspaces')
);

-- name: FindOrganizationNameByWorkspaceID :one
SELECT organization_name
FROM workspaces
WHERE workspace_id = sqlc.arg('workspace_id')
;

-- name: FindOrganizationByName :one
SELECT * FROM organizations WHERE name = sqlc.arg('name');

-- name: FindOrganizationByID :one
SELECT * FROM organizations WHERE organization_id = sqlc.arg('organization_id');

-- name: FindOrganizationByNameForUpdate :one
SELECT *
FROM organizations
WHERE name = sqlc.arg('name')
FOR UPDATE
;

-- name: FindOrganizations :many
SELECT *
FROM organizations
WHERE name LIKE ANY(sqlc.arg('names')::text[])
ORDER BY updated_at DESC
LIMIT sqlc.arg('limit')::int OFFSET sqlc.arg('offset')::int
;

-- name: CountOrganizations :one
SELECT count(*)
FROM organizations
WHERE name LIKE ANY(sqlc.arg('names')::text[])
;

-- name: UpdateOrganizationByName :one
UPDATE organizations
SET
    name = sqlc.arg('new_name'),
    email = sqlc.arg('email'),
    collaborator_auth_policy = sqlc.arg('collaborator_auth_policy'),
    cost_estimation_enabled = sqlc.arg('cost_estimation_enabled'),
    session_remember = sqlc.arg('session_remember'),
    session_timeout = sqlc.arg('session_timeout'),
    allow_force_delete_workspaces = sqlc.arg('allow_force_delete_workspaces'),
    updated_at = sqlc.arg('updated_at')
WHERE name = sqlc.arg('name')
RETURNING organization_id;

-- name: DeleteOrganizationByName :one
DELETE
FROM organizations
WHERE name = sqlc.arg('name')
RETURNING organization_id;

-- name: UpsertOrganizationToken :exec
INSERT INTO organization_tokens (
    organization_token_id,
    created_at,
    organization_name,
    expiry
) VALUES (
    sqlc.arg('organization_token_id'),
    sqlc.arg('created_at'),
    sqlc.arg('organization_name'),
    sqlc.arg('expiry')
) ON CONFLICT (organization_name) DO UPDATE
  SET created_at            = sqlc.arg('created_at'),
      organization_token_id = sqlc.arg('organization_token_id'),
      expiry                = sqlc.arg('expiry');

-- name: FindOrganizationTokens :many
SELECT *
FROM organization_tokens
WHERE organization_name = sqlc.arg('organization_name');

-- name: FindOrganizationTokensByName :one
SELECT *
FROM organization_tokens
WHERE organization_name = sqlc.arg('organization_name');

-- name: FindOrganizationTokensByID :one
SELECT *
FROM organization_tokens
WHERE organization_token_id = sqlc.arg('organization_token_id');

-- name: DeleteOrganiationTokenByName :one
DELETE
FROM organization_tokens
WHERE organization_name = sqlc.arg('organization_name')
RETURNING organization_token_id;
