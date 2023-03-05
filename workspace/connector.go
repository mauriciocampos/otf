// Package workspace is responsible for workspaces
package workspace

import (
	"context"

	"github.com/google/uuid"
	"github.com/leg100/otf"
	"github.com/leg100/otf/sql"
	"github.com/pkg/errors"
)

type Connector struct {
	otf.HookService        // for registering and unregistering connections to webhooks
	otf.WorkspaceService   // for retrieving workspace
	otf.VCSProviderService // for retrieving cloud client
}

func (c *Connector) Connect(ctx context.Context, workspaceID string, opts otf.ConnectWorkspaceOptions) error {
	provider, err := c.GetVCSProvider(ctx, opts.ProviderID)
	if err != nil {
		return errors.Wrap(err, "retrieving provider")
	}

	client, err := provider.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "creating vcs client")
	}

	repo, err := client.GetRepository(ctx, opts.Identifier)
	if err != nil {
		return errors.Wrap(err, "retrieving repository info")
	}

	branch := opts.Branch
	if branch == "" {
		branch = repo.Branch
	}

	hookCallback := func(ctx context.Context, tx otf.Database, hookID uuid.UUID) error {
		err := sql.CreateWorkspaceRepo(ctx, tx, workspaceID, otf.WorkspaceRepo{
			Branch:     branch,
			ProviderID: opts.ProviderID,
			WebhookID:  hookID,
		})
		if err != nil {
			return errors.Wrap(err, "creating workspace")
		}
		return nil
	}
	err = c.Hook(ctx, otf.HookOptions{
		Identifier:   opts.Identifier,
		Cloud:        provider.CloudConfig().Name,
		HookCallback: hookCallback,
		Client:       client,
	})
	if err != nil {
		return errors.Wrap(err, "registering webhook connection")
	}
	return nil
}

// Disconnect a repo from a workspace. The repo's webhook is deleted if no other
// workspace is connected to the repo.
func (c *Connector) Disconnect(ctx context.Context, workspaceID string) (*otf.Workspace, error) {
	ws, err := c.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	repo := ws.Repo()
	client, err := c.GetVCSClient(ctx, repo.ProviderID)
	if err != nil {
		return nil, err
	}

	// Perform multiple operations within a transaction:
	// 1. delete workspace repo from db
	// 2. delete webhook from db
	unhookCallback := func(ctx context.Context, tx otf.Database) error {
		return sql.DeleteWorkspaceRepo(ctx, tx, workspaceID)
	}
	err = c.Unhook(ctx, otf.UnhookOptions{
		HookID:         repo.WebhookID,
		Client:         client,
		UnhookCallback: unhookCallback,
	})
	return ws, err
}
