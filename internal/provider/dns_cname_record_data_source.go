// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abionclient "terraform-provider-abion/internal/client"
	"terraform-provider-abion/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsCNameRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsCNameRecordDataSource{}
)

// NewDnsCNameRecordDataSource is a helper function to simplify the provider implementation.
func NewDnsCNameRecordDataSource() datasource.DataSource {
	return &dnsCNameRecordDataSource{}
}

// dnsCNameRecordDataSource is the data source implementation.
type dnsCNameRecordDataSource struct {
	client     *abionclient.Client
	recordType utils.RecordType
}

// dnsCNameRecordModel maps the data source schema data.
type dnsCNameRecordModel struct {
	Zone   types.String     `tfsdk:"zone"`
	Name   types.String     `tfsdk:"name"`
	Record *CNameRecordData `tfsdk:"record"`
}

type CNameRecordData struct {
	CName    types.String `tfsdk:"cname"`
	TTL      types.Int32  `tfsdk:"ttl"`
	Comments types.String `tfsdk:"comments"`
}

func (d *dnsCNameRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.recordType = utils.RecordTypeCName
}

// Metadata returns the data source type name.
func (d *dnsCNameRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_cname_record"
}

// Schema defines the schema for the data source.
func (d *dnsCNameRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get DNS CNAME records of the zone.",
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "The zone the record belongs to.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name to fetch records for. For example `@`, `www`, `ftp`, `www.east`. The `@` character represents the root of the zone.",
			},
			"record": schema.SingleNestedAttribute{
				Description: "CNAME record details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"cname": schema.StringAttribute{
						Description: "The canonical name (CNAME) for the record.",
						Required:    true,
					},
					"ttl": schema.Int32Attribute{
						Description: "Time-to-live (TTL) for the record, in seconds.",
						Computed:    true,
						Optional:    true,
					},
					"comments": schema.StringAttribute{
						Description: "Comments for the record.",
						Computed:    true,
						Optional:    true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dnsCNameRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsCNameRecordModel

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

	// only one cname per level is supported, pick first element.
	state.Record = &CNameRecordData{
		CName:    types.StringValue(records[0].Data),
		Comments: utils.StringPointerToTerraformString(records[0].Comments),
		TTL:      utils.IntPointerToInt32(records[0].TTL),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
