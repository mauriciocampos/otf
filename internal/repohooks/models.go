// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repohooks

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

type Repohook struct {
	RepohookID    pgtype.UUID
	VCSID         pgtype.Text
	Secret        pgtype.Text
	RepoPath      pgtype.Text
	VCSProviderID resource.TfeID
}
