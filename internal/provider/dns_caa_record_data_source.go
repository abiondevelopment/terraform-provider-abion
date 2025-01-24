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
	"strconv"
	"strings"
	abionclient "terraform-provider-abion/internal/client"
	"terraform-provider-abion/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsCAARecordDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsCAARecordDataSource{}
)

// NewDnsCAARecordDataSource is a helper function to simplify the provider implementation.
func NewDnsCAARecordDataSource() datasource.DataSource {
	return &dnsCAARecordDataSource{}
}

// dnsCAARecordDataSource is the data source implementation.
type dnsCAARecordDataSource struct {
	client     *abionclient.Client
	recordType utils.RecordType
}

// dnsCAARecordModel maps the data source schema data.
type dnsCAARecordModel struct {
	Zone    types.String    `tfsdk:"zone"`
	Name    types.String    `tfsdk:"name"`
	Records []CAARecordData `tfsdk:"records"`
}

type CAARecordData struct {
	Flag     types.Int32  `tfsdk:"flag"`
	Tag      types.String `tfsdk:"tag"`
	Value    types.String `tfsdk:"value"`
	TTL      types.Int32  `tfsdk:"ttl"`
	Comments types.String `tfsdk:"comments"`
}

func (d *dnsCAARecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.recordType = utils.RecordTypeCAA
}

// Metadata returns the data source type name.
func (d *dnsCAARecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_caa_record"
}

// Schema defines the schema for the data source.
func (d *dnsCAARecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get DNS CAA records of the zone.",
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
				Description: "The list of CAA records. Records are sorted to avoid constant changing plans",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flag": schema.Int32Attribute{
							Required:    true,
							Description: "A numeric value (usually 0) indicating special processing instructions.",
						},
						"tag": schema.StringAttribute{
							Required:    true,
							Description: "Specifies the property or directive. Commonly used tags are: issue, issuewild, iodef",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "The value depends on the tag",
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
func (d *dnsCAARecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsCAARecordModel

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

	state.Records = []CAARecordData{}
	for _, record := range records {

		flag, tag, value := splitCAARecordData(record)

		state.Records = append(state.Records, CAARecordData{
			Flag:     types.Int32Value(int32(flag)),
			Tag:      types.StringValue(tag),
			Value:    types.StringValue(value),
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

func splitCAARecordData(record abionclient.Record) (int, string, string) {
	// CAA record data follow the pattern "flag tag value".
	parts := strings.Fields(record.Data)
	flag, _ := strconv.Atoi(parts[0])
	tag := parts[1]
	value := strings.ReplaceAll(parts[2], "\"", "")
	return flag, tag, value
}
