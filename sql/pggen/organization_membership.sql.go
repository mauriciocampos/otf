// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const insertOrganizationMembershipSQL = `INSERT INTO organization_memberships (
    user_id,
    organization_id
) VALUES (
    $1,
    $2
)
RETURNING *;`

type InsertOrganizationMembershipRow struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
}

// InsertOrganizationMembership implements Querier.InsertOrganizationMembership.
func (q *DBQuerier) InsertOrganizationMembership(ctx context.Context, userID string, organizationID string) (InsertOrganizationMembershipRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertOrganizationMembership")
	row := q.conn.QueryRow(ctx, insertOrganizationMembershipSQL, userID, organizationID)
	var item InsertOrganizationMembershipRow
	if err := row.Scan(&item.UserID, &item.OrganizationID); err != nil {
		return item, fmt.Errorf("query InsertOrganizationMembership: %w", err)
	}
	return item, nil
}

// InsertOrganizationMembershipBatch implements Querier.InsertOrganizationMembershipBatch.
func (q *DBQuerier) InsertOrganizationMembershipBatch(batch genericBatch, userID string, organizationID string) {
	batch.Queue(insertOrganizationMembershipSQL, userID, organizationID)
}

// InsertOrganizationMembershipScan implements Querier.InsertOrganizationMembershipScan.
func (q *DBQuerier) InsertOrganizationMembershipScan(results pgx.BatchResults) (InsertOrganizationMembershipRow, error) {
	row := results.QueryRow()
	var item InsertOrganizationMembershipRow
	if err := row.Scan(&item.UserID, &item.OrganizationID); err != nil {
		return item, fmt.Errorf("scan InsertOrganizationMembershipBatch row: %w", err)
	}
	return item, nil
}

const deleteOrganizationMembershipSQL = `DELETE
FROM organization_memberships
WHERE
    user_id = $1 AND
    organization_id = $2
;`

// DeleteOrganizationMembership implements Querier.DeleteOrganizationMembership.
func (q *DBQuerier) DeleteOrganizationMembership(ctx context.Context, userID string, organizationID string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteOrganizationMembership")
	cmdTag, err := q.conn.Exec(ctx, deleteOrganizationMembershipSQL, userID, organizationID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query DeleteOrganizationMembership: %w", err)
	}
	return cmdTag, err
}

// DeleteOrganizationMembershipBatch implements Querier.DeleteOrganizationMembershipBatch.
func (q *DBQuerier) DeleteOrganizationMembershipBatch(batch genericBatch, userID string, organizationID string) {
	batch.Queue(deleteOrganizationMembershipSQL, userID, organizationID)
}

// DeleteOrganizationMembershipScan implements Querier.DeleteOrganizationMembershipScan.
func (q *DBQuerier) DeleteOrganizationMembershipScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec DeleteOrganizationMembershipBatch: %w", err)
	}
	return cmdTag, err
}
