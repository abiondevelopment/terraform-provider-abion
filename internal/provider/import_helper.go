// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"strings"
)

func importState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect the import ID to be in the format: "zone/name"
	// e.g., "example.com/@", or "example.com/www"
	parts := strings.Split(req.ID, "/")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format `zone/name`, e.g., `example.com/@`).",
		)
		return
	}

	// Import the state for zone and name
	zone := parts[0]
	if zone == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"`zone` is required.",
		)
		return
	}

	// Import the state for zone and name
	name := parts[1]
	if name == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"`name` is required.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone"), zone)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
}
