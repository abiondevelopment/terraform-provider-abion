// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abionclient "terraform-provider-abion/internal/client"
	"terraform-provider-abion/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsPTRRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsPTRRecordDataSource{}
)

// NewDnsPTRRecordDataSource is a helper function to simplify the provider implementation.
func NewDnsPTRRecordDataSource() datasource.DataSource {
	return &dnsPTRRecordDataSource{}
}

// dnsPTRRecordDataSource is the data source implementation.
type dnsPTRRecordDataSource struct {
	client     *abionclient.Client
	recordType utils.RecordType
}

// dnsPTRRecordModel maps the data source schema data.
type dnsPTRRecordModel struct {
	Zone    types.String    `tfsdk:"zone"`
	Name    types.String    `tfsdk:"name"`
	Records []PTRRecordData `tfsdk:"records"`
}

type PTRRecordData struct {
	PTR      types.String `tfsdk:"ptr"`
	TTL      types.Int32  `tfsdk:"ttl"`
	Comments types.String `tfsdk:"comments"`
}

func (d *dnsPTRRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*abionclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *abionclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
	d.recordType = utils.RecordTypePTR
}

// Metadata returns the data source type name.
func (d *dnsPTRRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_ptr_record"
}

// Schema defines the schema for the data source.
func (d *dnsPTRRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get DNS PTR records of the zone.",
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "The zone the record belongs to.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name to fetch records for. For example `@`, `www`, `ftp`, `www.east`. The `@` character represents the root of the zone.",
			},
			"records": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of PTR records. Records are sorted to avoid constant changing plans",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ptr": schema.StringAttribute{
							Required:    true,
							Description: "The canonical name this record will point to.",
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

// Read refreshes the Terraform state with the latest data.
func (d *dnsPTRRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsPTRRecordModel

	// Load zone_name from the configuration into state
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "zone", state.Zone.ValueString())
	tflog.Debug(ctx, "Getting zone details")

	// Get the zone details from Abion API
	zone, err := d.client.GetZone(ctx, state.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Zone from Abion API",
			err.Error(),
		)
		return
	}

	recordTypes := zone.Data.Attributes.Records[state.Name.ValueString()]
	if len(recordTypes) == 0 {
		resp.Diagnostics.AddError(
			"No records exist on "+state.Name.ValueString()+" level.",
			"",
		)
		return
	}

	records := recordTypes[d.recordType.String()]
	if len(records) == 0 {
		resp.Diagnostics.AddError(
			"No "+d.recordType.String()+" records exist on "+state.Name.ValueString()+" level.",
			"",
		)
		return
	}

	state.Records = []PTRRecordData{}
	for _, record := range records {
		state.Records = append(state.Records, PTRRecordData{
			PTR:      types.StringValue(record.Data),
			Comments: utils.StringPointerToTerraformString(record.Comments),
			TTL:      utils.IntPointerToInt32(record.TTL),
		})
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
