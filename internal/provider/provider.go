// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strconv"
	abionclient "terraform-provider-abion/internal/client"
)

// Ensure AbionDnsProvider satisfies various provider interfaces.
var _ provider.Provider = &AbionDnsProvider{}
var _ provider.ProviderWithFunctions = &AbionDnsProvider{}

// AbionDnsProvider defines the provider implementation.
type AbionDnsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AbionProviderModel describes the provider data model.
type AbionProviderModel struct {
	Host    types.String `tfsdk:"host"`
	Apikey  types.String `tfsdk:"apikey"`
	Timeout types.Int32  `tfsdk:"timeout"`
}

// Metadata returns the provider type name.
func (p *AbionDnsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "abion"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *AbionDnsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The Abion API host URL. If not set, defaults to `https://api.abion.com`. " +
					"This value can also be set using the `ABION_API_HOST` environment variable. " +
					"The order of precedence: Terraform configuration value (highest priority) > " +
					"environment variable > default value.",
				Optional: true,
			},
			"apikey": schema.StringAttribute{
				MarkdownDescription: "The Abion API key. Contact [Abion](https://abion.com) for help on " +
					"how to create an account and an API key and whitelist IP addresses to be able to access the Abion API. " +
					"This value can also be set using the `ABION_API_KEY` environment variable. " +
					"The order of precedence: Terraform configuration value (highest priority) > " +
					"environment variable (lowest priority). ",
				Optional:  true,
				Sensitive: true,
			},
			"timeout": schema.Int32Attribute{
				MarkdownDescription: "The Abion API timeout in seconds. If not set, defaults to `60`. " +
					"This value can also be set using the `ABION_API_TIMEOUT` environment variable. " +
					"The order of precedence: Terraform configuration value (highest priority) > " +
					"environment variable > default value.",
				Optional: true,
			},
		},
	}
}

// Configure prepares a Abion API client for data sources and resources.
func (p *AbionDnsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Abion client.")

	// Retrieve provider data from configuration
	var config AbionProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Abion API Host",
			"The provider cannot create the Abion API client as there is an unknown configuration value for the Abion API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ABION_API_HOST environment variable.",
		)
	}

	if config.Timeout.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			"Unknown Abion API timeout",
			"The provider cannot create the Abion API client as there is an unknown configuration value for the Abion API timeout. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ABION_API_TIMEOUT environment variable.",
		)
	}

	if config.Apikey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown Abion API Key",
			"The provider cannot create the Abion API client as there is an unknown configuration value for the Abion API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ABION_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set. If neither set, use default abion api host
	host := os.Getenv("ABION_API_HOST")
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if host == "" {
		host = "https://api.abion.com"
	}

	// Environment variable ABION_API_TIMEOUT
	timeoutEnv := os.Getenv("ABION_API_TIMEOUT")

	var timeout int
	if !config.Timeout.IsNull() {
		// If Terraform configuration is set, use it
		timeout = int(config.Timeout.ValueInt32())
	} else if timeoutEnv != "" {
		// If ABION_API_TIMEOUT environment variable is set, convert it
		timeoutInt, err := strconv.Atoi(timeoutEnv)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"),
				"Invalid ABION_API_TIMEOUT value in environment",
				"Must be an integer.",
			)
		}
		timeout = timeoutInt
	} else {
		// Use the default value
		timeout = 60
	}

	apikey := os.Getenv("ABION_API_KEY")

	if !config.Apikey.IsNull() {
		apikey = config.Apikey.ValueString()
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing Abion API Key",
			"The provider cannot create the Abion API client as there is a missing or empty value for the Abion API Key. "+
				"Set the apikey value in the configuration or use the ABION_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "abion_host", host)
	ctx = tflog.SetField(ctx, "abion_apikey", apikey)
	ctx = tflog.SetField(ctx, "timeout", timeout)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "abion_apikey")

	tflog.Debug(ctx, "Creating Abion client")

	// Create a new Abion client using the configuration values
	client, err := abionclient.NewAbionClient(host, apikey, timeout)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Abion API Client",
			"An unexpected error occurred when creating the Abion API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Abion Client Error: "+err.Error(),
		)
		return
	}

	// Make the Abion client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Abion client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *AbionDnsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDnsARecordResource,
		NewDnsAAAARecordResource,
		NewDnsNSRecordResource,
		NewDnsTXTRecordResource,
		NewDnsCNameRecordResource,
		NewDnsMXRecordResource,
		NewDnsSRVRecordResource,
		NewDnsPTRRecordResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *AbionDnsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDnsARecordDataSource,
		NewDnsAAAARecordDataSource,
		NewDnsNSRecordDataSource,
		NewDnsTXTRecordDataSource,
		NewDnsCNameRecordDataSource,
		NewDnsMXRecordDataSource,
		NewDnsSRVRecordDataSource,
		NewDnsPTRRecordDataSource,
	}
}

// DataSources defines the functions implemented in the provider.
func (p *AbionDnsProvider) Functions(ctx context.Context) []func() function.Function {
	//return []func() function.Function{
	//	NewExampleFunction,
	//}
	return nil
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AbionDnsProvider{
			version: version,
		}
	}
}
