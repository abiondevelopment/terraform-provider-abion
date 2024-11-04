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
	_ datasource.DataSource              = &dnsSRVRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsSRVRecordDataSource{}
)

// NewDnsSRVRecordDataSource is a helper function to simplify the provider implementation.
func NewDnsSRVRecordDataSource() datasource.DataSource {
	return &dnsSRVRecordDataSource{}
}

// dnsSRVRecordDataSource is the data source implementation.
type dnsSRVRecordDataSource struct {
	client     *abionclient.Client
	recordType utils.RecordType
}

// dnsSRVRecordModel maps the data source schema data.
type dnsSRVRecordModel struct {
	Zone    types.String    `tfsdk:"zone"`
	Name    types.String    `tfsdk:"name"`
	Records []SRVRecordData `tfsdk:"records"`
}

type SRVRecordData struct {
	Priority types.Int32  `tfsdk:"priority"`
	Weight   types.Int32  `tfsdk:"weight"`
	Port     types.Int32  `tfsdk:"port"`
	Target   types.String `tfsdk:"target"`
	TTL      types.Int32  `tfsdk:"ttl"`
	Comments types.String `tfsdk:"comments"`
}

func (d *dnsSRVRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.recordType = utils.RecordTypeSRV
}

// Metadata returns the data source type name.
func (d *dnsSRVRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_srv_record"
}

// Schema defines the schema for the data source.
func (d *dnsSRVRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get DNS SRV records of the zone.",
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
				Description: "The list of SRV records. Records are sorted to avoid constant changing plans",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target": schema.StringAttribute{
							Required:    true,
							Description: "The hostname of the machine providing the service.",
						},
						"port": schema.Int32Attribute{
							Required:    true,
							Description: "The TCP or UDP port on which the service is running.",
						},
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "The priority specifying the priority of the target host. Lower values indicate higher priority.",
						},
						"weight": schema.Int32Attribute{
							Required:    true,
							Description: "A relative weight for records with the same priority. Higher weights are more preferred.",
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
func (d *dnsSRVRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsSRVRecordModel

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

	state.Records = []SRVRecordData{}
	for _, record := range records {

		priority, weight, port, target := splitSrvRecordData(record)

		state.Records = append(state.Records, SRVRecordData{
			Priority: types.Int32Value(int32(priority)),
			Weight:   types.Int32Value(int32(weight)),
			Port:     types.Int32Value(int32(port)),
			Target:   types.StringValue(target),
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

func splitSrvRecordData(record abionclient.Record) (int, int, int, string) {
	// SRV record data follow the pattern "priority weight port target".
	parts := strings.Fields(record.Data)
	priority, _ := strconv.Atoi(parts[0])
	weight, _ := strconv.Atoi(parts[1])
	port, _ := strconv.Atoi(parts[2])
	target := parts[3]
	return priority, weight, port, target
}
