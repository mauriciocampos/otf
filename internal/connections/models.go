// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package connections

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

type RepoConnection struct {
	ModuleID      *resource.TfeID
	WorkspaceID   *resource.TfeID
	RepoPath      pgtype.Text
	VCSProviderID resource.TfeID
}
