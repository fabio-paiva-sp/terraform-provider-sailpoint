package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

var (
	_ datasource.DataSource              = &managedClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &managedClusterDataSource{}
)

func NewManagedClusterDataSource() datasource.DataSource {
	return &managedClusterDataSource{}
}

type managedClusterDataSource struct {
	client *sailpoint.APIClient
}

func (d *managedClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

func (d *managedClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: managedClusterDataSourceSchemaAttributes,
	}
}

func (d *managedClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	tflog.Info(ctx, "Configuring SailPoint ManagedCluster data resource")

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

func (d *managedClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Managed Cluster")
	var state managedClusterSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	id := state.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Managed Cluster",
			"ID cannot be empty",
		)
		return
	}
	tflog.Debug(ctx, "Reading Managed Cluster filters", map[string]any{"id": id})

	cluster, res, err := d.client.V2025.ManagedClustersAPI.GetManagedCluster(ctx, id).Execute()

	if err != nil {
		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)
		tflog.Error(ctx, "Error reading managed cluster", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		resp.Diagnostics.AddError(
			"Unable to Read Managed Cluster",
			err.Error(),
		)
		return
	}

	state, diags := serializeManagedClusterData(ctx, *cluster)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
