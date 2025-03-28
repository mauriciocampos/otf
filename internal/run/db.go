package run

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/configversion"
	"github.com/leg100/otf/internal/resource"
	"github.com/leg100/otf/internal/runstatus"
	"github.com/leg100/otf/internal/sql"
	"github.com/leg100/otf/internal/workspace"
)

var q = &Queries{}

type (
	// pgdb is a database of runs on postgres
	pgdb struct {
		*sql.DB // provides access to generated SQL queries
	}

	// pgresult is the result of a database query for a run.
	pgresult struct {
		RunID                  resource.TfeID
		CreatedAt              pgtype.Timestamptz
		CancelSignaledAt       pgtype.Timestamptz
		IsDestroy              pgtype.Bool
		PositionInQueue        pgtype.Int4
		Refresh                pgtype.Bool
		RefreshOnly            pgtype.Bool
		Source                 pgtype.Text
		Status                 pgtype.Text
		PlanStatus             pgtype.Text
		ApplyStatus            pgtype.Text
		ReplaceAddrs           []pgtype.Text
		TargetAddrs            []pgtype.Text
		AutoApply              pgtype.Bool
		PlanResourceReport     *Report
		PlanOutputReport       *Report
		ApplyResourceReport    *Report
		ConfigurationVersionID resource.TfeID
		WorkspaceID            resource.TfeID
		PlanOnly               pgtype.Bool
		CreatedBy              pgtype.Text
		TerraformVersion       pgtype.Text
		AllowEmptyApply        pgtype.Bool
		ExecutionMode          pgtype.Text
		Latest                 pgtype.Bool
		OrganizationName       resource.OrganizationName
		CostEstimationEnabled  pgtype.Bool
		RunStatusTimestamps    []RunStatusTimestampModel
		PlanStatusTimestamps   []PhaseStatusTimestampModel
		ApplyStatusTimestamps  []PhaseStatusTimestampModel
		RunVariables           []RunVariable
		IngressAttributes      *IngressAttributeModel
	}
)

func (result pgresult) toRun() *Run {
	run := Run{
		ID:                     result.RunID,
		CreatedAt:              result.CreatedAt.Time.UTC(),
		IsDestroy:              result.IsDestroy.Bool,
		PositionInQueue:        int(result.PositionInQueue.Int32),
		Refresh:                result.Refresh.Bool,
		RefreshOnly:            result.RefreshOnly.Bool,
		Source:                 Source(result.Source.String),
		Status:                 runstatus.Status(result.Status.String),
		ReplaceAddrs:           sql.FromStringArray(result.ReplaceAddrs),
		TargetAddrs:            sql.FromStringArray(result.TargetAddrs),
		AutoApply:              result.AutoApply.Bool,
		PlanOnly:               result.PlanOnly.Bool,
		AllowEmptyApply:        result.AllowEmptyApply.Bool,
		TerraformVersion:       result.TerraformVersion.String,
		ExecutionMode:          workspace.ExecutionMode(result.ExecutionMode.String),
		Latest:                 result.Latest.Bool,
		Organization:           result.OrganizationName,
		WorkspaceID:            result.WorkspaceID,
		ConfigurationVersionID: result.ConfigurationVersionID,
		CostEstimationEnabled:  result.CostEstimationEnabled.Bool,
		Plan: Phase{
			RunID:          result.RunID,
			PhaseType:      internal.PlanPhase,
			Status:         PhaseStatus(result.PlanStatus.String),
			ResourceReport: reportFromDB(result.PlanResourceReport),
			OutputReport:   reportFromDB(result.PlanOutputReport),
		},
		Apply: Phase{
			RunID:          result.RunID,
			PhaseType:      internal.ApplyPhase,
			Status:         PhaseStatus(result.ApplyStatus.String),
			ResourceReport: reportFromDB(result.ApplyResourceReport),
		},
	}
	// convert run timestamps from db result and sort them according to
	// timestamp (earliest first)
	run.StatusTimestamps = make([]StatusTimestamp, len(result.RunStatusTimestamps))
	for i, rst := range result.RunStatusTimestamps {
		run.StatusTimestamps[i] = StatusTimestamp{
			Status:    runstatus.Status(rst.Status.String),
			Timestamp: rst.Timestamp.Time.UTC(),
		}
	}
	sort.Slice(run.StatusTimestamps, func(i, j int) bool {
		return run.StatusTimestamps[i].Timestamp.Before(run.StatusTimestamps[j].Timestamp)
	})
	// convert plan timestamps from db result and sort them according to
	// timestamp (earliest first)
	run.Plan.StatusTimestamps = make([]PhaseStatusTimestamp, len(result.PlanStatusTimestamps))
	for i, pst := range result.PlanStatusTimestamps {
		run.Plan.StatusTimestamps[i] = PhaseStatusTimestamp{
			Status:    PhaseStatus(pst.Status.String),
			Timestamp: pst.Timestamp.Time.UTC(),
		}
	}
	sort.Slice(run.Plan.StatusTimestamps, func(i, j int) bool {
		return run.Plan.StatusTimestamps[i].Timestamp.Before(run.Plan.StatusTimestamps[j].Timestamp)
	})
	// convert apply timestamps from db result and sort them according to
	// timestamp (earliest first)
	run.Apply.StatusTimestamps = make([]PhaseStatusTimestamp, len(result.ApplyStatusTimestamps))
	for i, ast := range result.ApplyStatusTimestamps {
		run.Apply.StatusTimestamps[i] = PhaseStatusTimestamp{
			Status:    PhaseStatus(ast.Status.String),
			Timestamp: ast.Timestamp.Time.UTC(),
		}
	}
	sort.Slice(run.Apply.StatusTimestamps, func(i, j int) bool {
		return run.Apply.StatusTimestamps[i].Timestamp.Before(run.Apply.StatusTimestamps[j].Timestamp)
	})
	if len(result.RunVariables) > 0 {
		run.Variables = make([]Variable, len(result.RunVariables))
		for i, v := range result.RunVariables {
			run.Variables[i] = Variable{Key: v.Key.String, Value: v.Value.String}
		}
	}
	if result.CreatedBy.Valid {
		run.CreatedBy = &result.CreatedBy.String
	}
	if result.CancelSignaledAt.Valid {
		run.CancelSignaledAt = internal.Time(result.CancelSignaledAt.Time.UTC())
	}
	if result.IngressAttributes != nil {
		run.IngressAttributes = configversion.NewIngressFromRow((*configversion.IngressAttributeModel)(result.IngressAttributes))
	}
	return &run
}

// CreateRun persists a Run to the DB.
func (db *pgdb) CreateRun(ctx context.Context, run *Run) error {
	return db.Tx(ctx, func(ctx context.Context, conn sql.Connection) error {
		err := q.InsertRun(ctx, conn, InsertRunParams{
			ID:                     run.ID,
			CreatedAt:              sql.Timestamptz(run.CreatedAt),
			IsDestroy:              sql.Bool(run.IsDestroy),
			PositionInQueue:        sql.Int4(0),
			Refresh:                sql.Bool(run.Refresh),
			RefreshOnly:            sql.Bool(run.RefreshOnly),
			Source:                 sql.String(string(run.Source)),
			Status:                 sql.String(string(run.Status)),
			ReplaceAddrs:           sql.StringArray(run.ReplaceAddrs),
			TargetAddrs:            sql.StringArray(run.TargetAddrs),
			AutoApply:              sql.Bool(run.AutoApply),
			PlanOnly:               sql.Bool(run.PlanOnly),
			AllowEmptyApply:        sql.Bool(run.AllowEmptyApply),
			TerraformVersion:       sql.String(run.TerraformVersion),
			ConfigurationVersionID: run.ConfigurationVersionID,
			WorkspaceID:            run.WorkspaceID,
			CreatedBy:              sql.StringPtr(run.CreatedBy),
		})
		for _, v := range run.Variables {
			err = q.InsertRunVariable(ctx, conn, InsertRunVariableParams{
				RunID: run.ID,
				Key:   sql.String(v.Key),
				Value: sql.String(v.Value),
			})
			if err != nil {
				return fmt.Errorf("inserting run variable: %w", err)
			}
		}
		if err != nil {
			return fmt.Errorf("inserting run: %w", err)
		}
		err = q.InsertPlan(ctx, conn, InsertPlanParams{
			RunID:  run.ID,
			Status: sql.String(string(run.Plan.Status)),
		})
		if err != nil {
			return fmt.Errorf("inserting plan: %w", err)
		}
		err = q.InsertApply(ctx, conn, InsertApplyParams{
			RunID:  run.ID,
			Status: sql.String(string(run.Apply.Status)),
		})
		if err != nil {
			return fmt.Errorf("inserting apply: %w", err)
		}
		if err := db.insertRunStatusTimestamp(ctx, run); err != nil {
			return fmt.Errorf("inserting run status timestamp: %w", err)
		}
		if err := db.insertPhaseStatusTimestamp(ctx, run.Plan); err != nil {
			return fmt.Errorf("inserting plan status timestamp: %w", err)
		}
		if err := db.insertPhaseStatusTimestamp(ctx, run.Apply); err != nil {
			return fmt.Errorf("inserting apply status timestamp: %w", err)
		}
		return nil
	})
}

// UpdateStatus updates the run status as well as its plan and/or apply.
func (db *pgdb) UpdateStatus(ctx context.Context, runID resource.TfeID, fn func(context.Context, *Run) error) (*Run, error) {
	var runStatus runstatus.Status
	var planStatus PhaseStatus
	var applyStatus PhaseStatus
	var cancelSignaledAt *time.Time

	return sql.Updater(
		ctx,
		db.DB,
		func(ctx context.Context, conn sql.Connection) (*Run, error) {
			result, err := q.FindRunByIDForUpdate(ctx, conn, runID)
			if err != nil {
				return nil, err
			}
			run := pgresult(result).toRun()
			// Make copies of run attributes before update
			runStatus = run.Status
			planStatus = run.Plan.Status
			applyStatus = run.Apply.Status
			cancelSignaledAt = run.CancelSignaledAt
			return run, nil
		},
		fn,
		func(ctx context.Context, conn sql.Connection, run *Run) error {
			if run.Status != runStatus {
				_, err := q.UpdateRunStatus(ctx, conn, UpdateRunStatusParams{
					Status: sql.String(string(run.Status)),
					ID:     run.ID,
				})
				if err != nil {
					return err
				}

				if err := db.insertRunStatusTimestamp(ctx, run); err != nil {
					return err
				}
			}

			if run.Plan.Status != planStatus {
				_, err := q.UpdatePlanStatusByID(ctx, conn, UpdatePlanStatusByIDParams{
					Status: sql.String(string(run.Plan.Status)),
					RunID:  run.ID,
				})
				if err != nil {
					return err
				}

				if err := db.insertPhaseStatusTimestamp(ctx, run.Plan); err != nil {
					return err
				}
			}

			if run.Apply.Status != applyStatus {
				_, err := q.UpdateApplyStatusByID(ctx, db.Conn(ctx), UpdateApplyStatusByIDParams{
					Status: sql.String(string(run.Apply.Status)),
					RunID:  run.ID,
				})
				if err != nil {
					return err
				}

				if err := db.insertPhaseStatusTimestamp(ctx, run.Apply); err != nil {
					return err
				}
			}

			if run.CancelSignaledAt != cancelSignaledAt && run.CancelSignaledAt != nil {
				_, err := q.UpdateCancelSignaledAt(ctx, db.Conn(ctx), UpdateCancelSignaledAtParams{
					CancelSignaledAt: sql.Timestamptz(*run.CancelSignaledAt),
					ID:               run.ID,
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (db *pgdb) CreatePlanReport(ctx context.Context, runID resource.TfeID, resource, output Report) error {
	_, err := q.UpdatePlannedChangesByID(ctx, db.Conn(ctx), UpdatePlannedChangesByIDParams{
		RunID:                runID,
		ResourceAdditions:    sql.Int4(resource.Additions),
		ResourceChanges:      sql.Int4(resource.Changes),
		ResourceDestructions: sql.Int4(resource.Destructions),
		OutputAdditions:      sql.Int4(output.Additions),
		OutputChanges:        sql.Int4(output.Changes),
		OutputDestructions:   sql.Int4(output.Destructions),
	})
	if err != nil {
		return sql.Error(err)
	}
	return err
}

func (db *pgdb) CreateApplyReport(ctx context.Context, runID resource.TfeID, report Report) error {
	_, err := q.UpdateAppliedChangesByID(ctx, db.Conn(ctx), UpdateAppliedChangesByIDParams{
		RunID:        runID,
		Additions:    sql.Int4(report.Additions),
		Changes:      sql.Int4(report.Changes),
		Destructions: sql.Int4(report.Destructions),
	})
	if err != nil {
		return sql.Error(err)
	}
	return err
}

func (db *pgdb) ListRuns(ctx context.Context, opts ListOptions) (*resource.Page[*Run], error) {
	organization := "%"
	if opts.Organization != nil {
		organization = opts.Organization.String()
	}
	workspaceName := "%"
	if opts.WorkspaceName != nil {
		workspaceName = *opts.WorkspaceName
	}
	workspaceID := "%"
	if opts.WorkspaceID != nil {
		workspaceID = opts.WorkspaceID.String()
	}
	sources := []string{"%"}
	if len(opts.Sources) > 0 {
		sources = internal.ToStringSlice(opts.Sources)
	}
	statuses := []string{"%"}
	if len(opts.Statuses) > 0 {
		statuses = internal.ToStringSlice(opts.Statuses)
	}
	planOnly := "%"
	if opts.PlanOnly != nil {
		planOnly = strconv.FormatBool(*opts.PlanOnly)
	}
	rows, err := q.FindRuns(ctx, db.Conn(ctx), FindRunsParams{
		OrganizationNames: sql.StringArray([]string{organization}),
		WorkspaceNames:    sql.StringArray([]string{workspaceName}),
		WorkspaceIds:      sql.StringArray([]string{workspaceID}),
		CommitSHA:         sql.StringPtr(opts.CommitSHA),
		VCSUsername:       sql.StringPtr(opts.VCSUsername),
		Sources:           sql.StringArray(sources),
		Statuses:          sql.StringArray(statuses),
		PlanOnly:          sql.StringArray([]string{planOnly}),
		Limit:             sql.GetLimit(opts.PageOptions),
		Offset:            sql.GetOffset(opts.PageOptions),
	})
	if err != nil {
		return nil, fmt.Errorf("querying runs: %w", err)
	}
	count, err := q.CountRuns(ctx, db.Conn(ctx), CountRunsParams{
		OrganizationNames: sql.StringArray([]string{organization}),
		WorkspaceNames:    sql.StringArray([]string{workspaceName}),
		WorkspaceIds:      sql.StringArray([]string{workspaceID}),
		CommitSHA:         sql.StringPtr(opts.CommitSHA),
		VCSUsername:       sql.StringPtr(opts.VCSUsername),
		Sources:           sql.StringArray(sources),
		Statuses:          sql.StringArray(statuses),
		PlanOnly:          sql.StringArray([]string{planOnly}),
	})
	if err != nil {
		return nil, fmt.Errorf("counting runs: %w", err)
	}

	items := make([]*Run, len(rows))
	for i, r := range rows {
		items[i] = pgresult(r).toRun()
	}
	return resource.NewPage(items, opts.PageOptions, internal.Int64(count)), nil
}

// GetRun retrieves a run using the get options
func (db *pgdb) GetRun(ctx context.Context, runID resource.TfeID) (*Run, error) {
	result, err := q.FindRunByID(ctx, db.Conn(ctx), runID)
	if err != nil {
		return nil, sql.Error(err)
	}
	return pgresult(result).toRun(), nil
}

// SetPlanFile writes a plan file to the db
func (db *pgdb) SetPlanFile(ctx context.Context, runID resource.TfeID, file []byte, format PlanFormat) error {
	switch format {
	case PlanFormatBinary:
		_, err := q.UpdatePlanBinByID(ctx, db.Conn(ctx), UpdatePlanBinByIDParams{
			PlanBin: file,
			RunID:   runID,
		})
		return err
	case PlanFormatJSON:
		_, err := q.UpdatePlanJSONByID(ctx, db.Conn(ctx), UpdatePlanJSONByIDParams{
			PlanJSON: file,
			RunID:    runID,
		})
		return err
	default:
		return fmt.Errorf("unknown plan format: %s", string(format))
	}
}

// GetPlanFile retrieves a plan file for the run
func (db *pgdb) GetPlanFile(ctx context.Context, runID resource.TfeID, format PlanFormat) ([]byte, error) {
	switch format {
	case PlanFormatBinary:
		return q.GetPlanBinByID(ctx, db.Conn(ctx), runID)
	case PlanFormatJSON:
		return q.GetPlanJSONByID(ctx, db.Conn(ctx), runID)
	default:
		return nil, fmt.Errorf("unknown plan format: %s", string(format))
	}
}

// GetLockFile retrieves the lock file for the run
func (db *pgdb) GetLockFile(ctx context.Context, runID resource.TfeID) ([]byte, error) {
	return q.GetLockFileByID(ctx, db.Conn(ctx), runID)
}

// SetLockFile sets the lock file for the run
func (db *pgdb) SetLockFile(ctx context.Context, runID resource.TfeID, lockFile []byte) error {
	_, err := q.PutLockFile(ctx, db.Conn(ctx), PutLockFileParams{
		LockFile: lockFile,
		RunID:    runID,
	})
	return err
}

// DeleteRun deletes a run from the DB
func (db *pgdb) DeleteRun(ctx context.Context, id resource.TfeID) error {
	_, err := q.DeleteRunByID(ctx, db.Conn(ctx), id)
	return err
}

func (db *pgdb) insertRunStatusTimestamp(ctx context.Context, run *Run) error {
	ts, err := run.StatusTimestamp(run.Status)
	if err != nil {
		return err
	}
	err = q.InsertRunStatusTimestamp(ctx, db.Conn(ctx), InsertRunStatusTimestampParams{
		ID:        run.ID,
		Status:    sql.String(string(run.Status)),
		Timestamp: sql.Timestamptz(ts),
	})
	return err
}

func (db *pgdb) insertPhaseStatusTimestamp(ctx context.Context, phase Phase) error {
	ts, err := phase.StatusTimestamp(phase.Status)
	if err != nil {
		return err
	}
	err = q.InsertPhaseStatusTimestamp(ctx, db.Conn(ctx), InsertPhaseStatusTimestampParams{
		RunID:      phase.RunID,
		PhaseModel: sql.String(string(phase.PhaseType)),
		Status:     sql.String(string(phase.Status)),
		Timestamp:  sql.Timestamptz(ts),
	})
	return err
}
