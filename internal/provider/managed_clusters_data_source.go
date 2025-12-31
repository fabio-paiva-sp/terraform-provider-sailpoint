package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

var (
	_ datasource.DataSource              = &managedClustersDataSource{}
	_ datasource.DataSourceWithConfigure = &managedClustersDataSource{}
)

func NewManagedClustersDataSource() datasource.DataSource {
	return &managedClustersDataSource{}
}

type managedClustersDataSource struct {
	client *sailpoint.APIClient
}

func (d *managedClustersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_clusters"
}

func (d *managedClustersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filters": schema.StringAttribute{
				Optional:    true,
				Description: "Filter results using the standard syntax described in V3 API Standard Collection Parameters",
			},
			"managed_clusters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: managedClusterSchemaAttributes,
				},
			},
		},
	}
}

func (d *managedClustersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	tflog.Info(ctx, "Configuring SailPoint ManagedClusters data resource")

	client, ok := req.ProviderData.(*sailpoint.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func extractNullableString(str *string, ok bool) *string {
	if ok {
		return str
	}
	return nil
}

func (d *managedClustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Managed Clusters")
	var state managedClustersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	filters := state.Filters.ValueString()
	tflog.Debug(ctx, "Reading Managed Clusters filters", map[string]any{"filters": filters})

	results, res, err := sailpoint.PaginateWithDefaults[v2025.ManagedCluster](d.client.V2025.ManagedClustersAPI.GetManagedClusters(ctx).Filters(filters))

	if err != nil {
		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)
		tflog.Error(ctx, "Error reading managed clusters", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		resp.Diagnostics.AddError(
			"Unable to Read Managed Clusters",
			err.Error(),
		)
		return
	}

	state.ManagedClusters = make([]managedClusterDataSourceModel, 0)
	for _, cluster := range results {
		tflog.Debug(ctx, "Iterating through the clusters")

		clusterState, diags := parseManagedClusterData(ctx, cluster)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.ManagedClusters = append(state.ManagedClusters, clusterState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
