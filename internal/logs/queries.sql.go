// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package logs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

const findLogChunkByID = `-- name: FindLogChunkByID :one
SELECT
    chunk_id,
    run_id,
    phase,
    chunk,
    _offset AS offset
FROM logs
WHERE chunk_id = $1
`

type FindLogChunkByIDRow struct {
	ChunkID resource.TfeID
	RunID   resource.TfeID
	Phase   pgtype.Text
	Chunk   []byte
	Offset  pgtype.Int4
}

func (q *Queries) FindLogChunkByID(ctx context.Context, db DBTX, chunkID resource.TfeID) (FindLogChunkByIDRow, error) {
	row := db.QueryRow(ctx, findLogChunkByID, chunkID)
	var i FindLogChunkByIDRow
	err := row.Scan(
		&i.ChunkID,
		&i.RunID,
		&i.Phase,
		&i.Chunk,
		&i.Offset,
	)
	return i, err
}

const findLogs = `-- name: FindLogs :one
SELECT
    string_agg(chunk, '')
FROM (
    SELECT run_id, phase, chunk
    FROM logs
    WHERE run_id = $1
    AND   phase  = $2
    ORDER BY _offset
) c
GROUP BY run_id, phase
`

type FindLogsParams struct {
	RunID resource.TfeID
	Phase pgtype.Text
}

// FindLogs retrieves all the logs for the given run and phase.
func (q *Queries) FindLogs(ctx context.Context, db DBTX, arg FindLogsParams) ([]byte, error) {
	row := db.QueryRow(ctx, findLogs, arg.RunID, arg.Phase)
	var string_agg []byte
	err := row.Scan(&string_agg)
	return string_agg, err
}

const insertLogChunk = `-- name: InsertLogChunk :exec
INSERT INTO logs (
    chunk_id,
    run_id,
    phase,
    chunk,
    _offset
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
`

type InsertLogChunkParams struct {
	ChunkID resource.TfeID
	RunID   resource.TfeID
	Phase   pgtype.Text
	Chunk   []byte
	Offset  pgtype.Int4
}

func (q *Queries) InsertLogChunk(ctx context.Context, db DBTX, arg InsertLogChunkParams) error {
	_, err := db.Exec(ctx, insertLogChunk,
		arg.ChunkID,
		arg.RunID,
		arg.Phase,
		arg.Chunk,
		arg.Offset,
	)
	return err
}
