package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		Attributes: managedClusterSchemaAttributes,
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
	var state managedClusterDataSourceModel

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

	configuration, diags := types.MapValueFrom(ctx, types.StringType, cluster.Configuration)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	clientIds, diags := types.ListValueFrom(ctx, types.StringType, cluster.ClientIds)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "Reading cluster configuration property", map[string]any{"configuration": cluster.Configuration})

	var createdAt string
	createdAtDate, _ := cluster.GetCreatedAtOk()

	if createdAtDate != nil {
		createdAt = createdAtDate.String()
		tflog.Trace(ctx, "Reading cluster created at property string", map[string]any{"date": createdAtDate.String()})
	}

	var updatedAt string
	updatedAtDate, _ := cluster.GetUpdatedAtOk()

	if updatedAtDate != nil {
		updatedAt = updatedAtDate.String()
		tflog.Trace(ctx, "Reading cluster updated at property string", map[string]any{"date": updatedAtDate.String()})
	}

	keyPairData := managedClusterKeyPairModel{
		PublicKey:            types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyOk())),
		PublicKeyThumbprint:  types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyThumbprintOk())),
		PublicKeyCertificate: types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyCertificateOk())),
	}

	keyPairObject, diags := types.ObjectValueFrom(ctx, managedClusterKeyPairAttrTypes, keyPairData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.KeyPair = keyPairObject

	attributesData := managedClusterAttributesModel{
		Queue: managedClusterAttributesQueueModel{
			Name:   types.StringValue(cluster.GetAttributes().Queue.GetName()),
			Region: types.StringValue(cluster.GetAttributes().Queue.GetRegion()),
		},
		KeyStore: types.StringPointerValue(extractNullableString(cluster.Attributes.GetKeystoreOk())),
	}

	attributesObject, diags := types.ObjectValueFrom(ctx, managedClusterAttributesAttrTypes, attributesData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Attributes = attributesObject

	redisData := managedClusterRedisModel{
		RedisHost: types.StringPointerValue(cluster.GetRedis().RedisHost),
		RedisPort: types.Int32PointerValue(cluster.GetRedis().RedisPort),
	}

	redisObject, diags := types.ObjectValueFrom(ctx, managedClusterRedisAttrTypes, redisData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Redis = redisObject

	encryptionConfigData := managedClusterEncyprionConfigurationModel{
		Format: types.StringPointerValue(cluster.GetEncryptionConfiguration().Format),
	}

	encryptionConfigObject, diags := types.ObjectValueFrom(ctx, managedClusterEncryptionConfigAttrTypes, encryptionConfigData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.EncryptionConfiguration = encryptionConfigObject

	state.Name = types.StringValue(cluster.GetName())
	state.Pod = types.StringValue(cluster.GetPod())
	state.Org = types.StringValue(cluster.GetOrg())
	state.Type = types.StringValue(string(cluster.GetType()))
	state.Configuration = configuration
	state.Description = types.StringValue(cluster.GetDescription())
	state.ClientType = types.StringValue(string(cluster.GetClientType()))
	state.CcgVersion = types.StringValue(cluster.GetCcgVersion())
	state.PinnedConfig = types.BoolValue(cluster.GetPinnedConfig())
	state.Operational = types.BoolValue(cluster.GetOperational())
	state.Status = types.StringValue(cluster.GetStatus())
	state.PublicKeyCertificate = types.StringValue(cluster.GetPublicKeyCertificate())
	state.PublicKeyThumbprint = types.StringValue(cluster.GetPublicKeyThumbprint())
	state.PublicKey = types.StringValue(cluster.GetPublicKey())
	state.AlertKey = types.StringValue(cluster.GetAlertKey())
	state.ClientIds = clientIds
	state.ServiceCount = types.Int32Value(cluster.GetServiceCount())
	state.CcID = types.StringValue(cluster.GetCcId())
	state.CreatedAt = types.StringPointerValue(&createdAt)
	state.UpdatedAt = types.StringPointerValue(&updatedAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
