package authz

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/rbac"
	"github.com/leg100/otf/internal/resource"
)

// SiteAuthorizer authorizes access to site-wide actions
type SiteAuthorizer struct {
	logr.Logger
}

func (a *SiteAuthorizer) CanAccess(ctx context.Context, action rbac.Action, _ resource.ID) (Subject, error) {
	subj, err := SubjectFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if subj.CanAccessSite(action) {
		return subj, nil
	}
	a.Error(nil, "unauthorized action", "action", action, "subject", subj)
	return nil, internal.ErrAccessNotPermitted
}
