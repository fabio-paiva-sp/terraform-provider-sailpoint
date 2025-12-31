package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &sailpointProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sailpointProvider{
			version: version,
		}
	}
}

// sailpointProvider is the provider implementation.
type sailpointProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// sailpointProviderModel maps provider schema data to a Go type.
type sailpointProviderModel struct {
	BaseUrl      types.String `tfsdk:"base_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Experimental types.Bool   `tfsdk:"experimental"`
}

// Metadata returns the provider type name.
func (p *sailpointProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sailpoint"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *sailpointProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "The API URL - The API URL used to access your Identity Security Cloud tenant (ex. https://tenant.api.identitynow.com), this is used for the api calls made by certain commands.",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The PAT Client ID https://developer.sailpoint.com/docs/api/authentication/#generate-a-personal-access-token",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The PAT Client Secret https://developer.sailpoint.com/docs/api/authentication/#generate-a-personal-access-token",
			},
			"experimental": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether it's allowed to use experimental resources",
			},
		},
	}
}

// Configure prepares a SailPoint API client for data sources and resources.
func (p *sailpointProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring SailPoint provider")

	// Retrieve provider data from configuration
	var config sailpointProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown SailPoint API Base URL",
			"The provider cannot create the SailPoint API client as there is an unknown configuration value for the SailPoint API base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAIL_BASE_URL environment variable.",
		)
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown SailPoint PAT client_id",
			"The provider cannot create the SailPoint API client as there is an unknown configuration value for the SailPoint API Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAIL_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown SailPoint PAT client_secret",
			"The provider cannot create the SailPoint API client as there is an unknown configuration value for the SailPoint API Client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SAIL_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	baseUrl := os.Getenv("SAIL_BASE_URL")
	clientID := os.Getenv("SAIL_CLIENT_ID")
	clientSecret := os.Getenv("SAIL_CLIENT_SECRET")
	experimental := os.Getenv("SAIL_EXPERIMENTAL") == "true"

	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}

	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.Experimental.IsNull() {
		experimental = config.Experimental.ValueBool()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if baseUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing SailPoint API base_url",
			"The provider cannot create the SailPoint API client as there is a missing or empty value for the SailPoint API Base URL. "+
				"Set the base_url value in the configuration or use the SAIL_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing SailPoint PAT client_id",
			"The provider cannot create the SailPoint API client as there is a missing or empty value for the SailPoint PAT Client ID. "+
				"Set the client_id value in the configuration or use the SAIL_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing SailPoint PAT client_secret",
			"The provider cannot create the SailPoint API client as there is a missing or empty value for the SailPoint PAT Client secret. "+
				"Set the client_secret value in the configuration or use the SAIL_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sailpoint_base_url", baseUrl)
	ctx = tflog.SetField(ctx, "sailpoint_client_id", clientID)
	ctx = tflog.SetField(ctx, "sailpoint_client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sailpoint_client_secret")

	tflog.Debug(ctx, "Creating SailPoint API client")

	// Create a new SailPoint client using the configuration values
	configuration := sailpoint.NewConfiguration(sailpoint.ClientConfiguration{
		BaseURL:      baseUrl,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("%s/oauth/token", baseUrl), // token URL seems to be required when passing the parameters to the client configuration
	})
	if experimental {
		configuration.Experimental = true
		tflog.Debug(ctx, "Allowing the client to use experimental resources")
	}
	client := sailpoint.NewAPIClient(configuration)

	// Make the SailPoint client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *sailpointProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewManagedClustersDataSource,
		NewManagedClusterDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *sailpointProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
