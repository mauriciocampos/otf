// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"time"

	"github.com/leg100/otf/internal/resource"
)

// RunEventList represents a list of run events.
type RunEventList struct {
	// Pagination is not supported by the API
	*Pagination
	Items []*RunEvent
}

// RunEvent represents a Terraform Enterprise run event.
type RunEvent struct {
	ID          resource.TfeID `jsonapi:"primary,run-events"`
	Action      string         `jsonapi:"attr,action"`
	CreatedAt   time.Time      `jsonapi:"attr,created-at,iso8601"`
	Description string         `jsonapi:"attr,description"`

	// Relations - Note that `target` is not supported yet
	Actor *User `jsonapi:"relation,actor"`
	// Comment *Comment `jsonapi:"relation,comment"`
}
