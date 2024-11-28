// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"strings"
	abionclient "terraform-provider-abion/internal/client"
	"terraform-provider-abion/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnsMXRecordResource{}
	_ resource.ResourceWithConfigure   = &dnsMXRecordResource{}
	_ resource.ResourceWithImportState = &dnsMXRecordResource{}
)

// NewDnsMXRecordResource is a helper function to simplify the provider implementation.
func NewDnsMXRecordResource() resource.Resource {
	return &dnsMXRecordResource{}
}

// dnsMXRecordResource is the resource implementation.
type dnsMXRecordResource struct {
	client     *abionclient.Client
	recordType utils.RecordType
}

func (r *dnsMXRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*abionclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *abionclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
	r.recordType = utils.RecordTypeMX
}

// Metadata returns the resource type name.
func (r *dnsMXRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_mx_record"
}

// Schema defines the schema for the resource.
func (r *dnsMXRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create, update and delete DNS MX records of a zone.",
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "The zone the record belongs to.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name to create records for. For example `@`, `www`, `ftp`, `www.east`. The `@` character represents the root of the zone.",
			},
			"records": schema.ListNestedAttribute{
				Required:    true,
				Description: "The list of MX records. Records are sorted to avoid constant changing plans",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Required:    true,
							Description: "The hostname of the mail server.",
						},
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "The priority in which order mail servers are tried.",
						},
						"ttl": schema.Int32Attribute{
							Optional:    true,
							Description: "Time-to-live (TTL) for the record, in seconds.",
						},
						"comments": schema.StringAttribute{
							Optional:    true,
							Description: "Comments for the record.",
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dnsMXRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnsMXRecordModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patchRequest := r.createMXRecordCreateUpdateRequest(plan)

	ctx = tflog.SetField(ctx, "zone", plan.Zone.ValueString())
	ctx = tflog.SetField(ctx, "name", plan.Name.ValueString())
	ctx = tflog.SetField(ctx, "record_type", r.recordType.String())
	tflog.Debug(ctx, "Creating zone "+r.recordType.String()+" record")

	// Update zone by adding the record
	_, err := r.client.PatchZone(ctx, plan.Zone.ValueString(), patchRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching zone",
			"Could not create record, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dnsMXRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state dnsMXRecordModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the zone details from Abion API
	zone, err := r.client.GetZone(ctx, state.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Abion Zone",
			err.Error(),
		)
		return
	}

	recordTypes := zone.Data.Attributes.Records[state.Name.ValueString()]

	state.Records = []MXRecordData{}
	if len(recordTypes) > 0 {
		records := recordTypes[r.recordType.String()]

		for _, record := range records {
			// MX record data follow the pattern "priority host".
			parts := strings.Fields(record.Data)
			priority, _ := strconv.Atoi(parts[0])

			state.Records = append(state.Records, MXRecordData{
				Priority: types.Int32Value(int32(priority)),
				Host:     types.StringValue(parts[1]),
				Comments: utils.StringPointerToTerraformString(record.Comments),
				TTL:      utils.IntPointerToInt32(record.TTL),
			})
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsMXRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dnsMXRecordModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from current state
	var state dnsMXRecordModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update or add the new records according to the plan
	patchRequest := r.createMXRecordCreateUpdateRequest(plan)

	if plan.Name.ValueString() != state.Name.ValueString() {
		// records has been moved from one level to another, remove the records from the old state level
		if patchRequest.Data.Attributes.Records[state.Name.ValueString()] == nil {
			patchRequest.Data.Attributes.Records[state.Name.ValueString()] = make(map[string][]abionclient.Record)
		}
		patchRequest.Data.Attributes.Records[state.Name.ValueString()][r.recordType.String()] = nil
	}

	ctx = tflog.SetField(ctx, "zone", plan.Zone.ValueString())
	ctx = tflog.SetField(ctx, "name", plan.Name.ValueString())
	ctx = tflog.SetField(ctx, "record_type", r.recordType.String())
	tflog.Debug(ctx, "Updating zone "+r.recordType.String()+" record")

	// Update zone by adding the record
	_, err := r.client.PatchZone(ctx, plan.Zone.ValueString(), patchRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching zone",
			"Could not update record, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsMXRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dnsMXRecordModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patchRequest := abionclient.CreateRecordPatchRequest(state.Zone.ValueString(), state.Name.ValueString(), r.recordType, nil)

	ctx = tflog.SetField(ctx, "zone", state.Zone.ValueString())
	ctx = tflog.SetField(ctx, "name", state.Name.ValueString())
	ctx = tflog.SetField(ctx, "record_type", r.recordType.String())
	tflog.Debug(ctx, "Deleting zone "+r.recordType.String()+" record")

	// Update zone by removing the records
	_, err := r.client.PatchZone(ctx, state.Zone.ValueString(), patchRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching zone",
			"Could not delete record, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *dnsMXRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp)
}

func (r *dnsMXRecordResource) createMXRecordCreateUpdateRequest(plan dnsMXRecordModel) abionclient.ZoneRequest {
	var data []abionclient.Record
	for _, record := range plan.Records {
		record := abionclient.Record{
			Data:     record.Priority.String() + " " + record.Host.ValueString(),
			TTL:      utils.Int32ToIntPointer(record.TTL),
			Comments: record.Comments.ValueStringPointer(),
		}
		data = append(data, record)
	}

	return abionclient.CreateRecordPatchRequest(plan.Zone.ValueString(), plan.Name.ValueString(), r.recordType, data)
}
