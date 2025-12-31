package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

var (
	managedClusterKeyPairAttrTypes = map[string]attr.Type{
		"public_key":             types.StringType,
		"public_key_thumbprint":  types.StringType,
		"public_key_certificate": types.StringType,
	}
	managedClusterAttributesAttrTypes = map[string]attr.Type{
		"queue": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":   types.StringType,
				"region": types.StringType,
			},
		},
		"key_store": types.StringType,
	}
	managedClusterRedisAttrTypes = map[string]attr.Type{
		"redis_host": types.StringType,
		"redis_port": types.Int32Type,
	}
	managedClusterEncryptionConfigAttrTypes = map[string]attr.Type{
		"format": types.StringType,
	}
	managedClusterSchemaAttributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Required: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"pod": schema.StringAttribute{
			Computed: true,
		},
		"org": schema.StringAttribute{
			Computed: true,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},
		"configuration": schema.MapAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"key_pair": schema.ObjectAttribute{
			Optional:       true,
			Computed:       false,
			AttributeTypes: managedClusterKeyPairAttrTypes,
		},
		"attributes": schema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterAttributesAttrTypes,
		},
		"redis": schema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterRedisAttrTypes,
		},
		"description": schema.StringAttribute{
			Computed: true,
		},
		"client_type": schema.StringAttribute{
			Computed: true,
		},
		"ccg_version": schema.StringAttribute{
			Computed: true,
		},
		"pinned_config": schema.BoolAttribute{
			Computed: true,
		},
		"operational": schema.BoolAttribute{
			Computed: true,
		},
		"status": schema.StringAttribute{
			Computed: true,
		},
		"public_key_certificate": schema.StringAttribute{
			Computed: true,
		},
		"public_key_thumbprint": schema.StringAttribute{
			Computed: true,
		},
		"public_key": schema.StringAttribute{
			Computed: true,
		},
		"alert_key": schema.StringAttribute{
			Computed: true,
		},
		"client_ids": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_count": schema.Int32Attribute{
			Computed: true,
		},
		"cc_id": schema.StringAttribute{
			Computed: true,
		},
		"created_at": schema.StringAttribute{
			Computed: true,
		},
		"updated_at": schema.StringAttribute{
			Computed: true,
			Optional: true,
		},
		"encryption_configuration": schema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterEncryptionConfigAttrTypes,
		},
	}
)

type managedClusterDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Pod                     types.String `tfsdk:"pod"`
	Org                     types.String `tfsdk:"org"`
	Type                    types.String `tfsdk:"type"`
	Configuration           types.Map    `tfsdk:"configuration"`
	KeyPair                 types.Object `tfsdk:"key_pair"`
	Attributes              types.Object `tfsdk:"attributes"`
	Redis                   types.Object `tfsdk:"redis"`
	Description             types.String `tfsdk:"description"`
	ClientType              types.String `tfsdk:"client_type"`
	CcgVersion              types.String `tfsdk:"ccg_version"`
	PinnedConfig            types.Bool   `tfsdk:"pinned_config"`
	Operational             types.Bool   `tfsdk:"operational"`
	Status                  types.String `tfsdk:"status"`
	PublicKeyCertificate    types.String `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint     types.String `tfsdk:"public_key_thumbprint"`
	PublicKey               types.String `tfsdk:"public_key"`
	AlertKey                types.String `tfsdk:"alert_key"`
	ClientIds               types.List   `tfsdk:"client_ids"`
	ServiceCount            types.Int32  `tfsdk:"service_count"`
	CcID                    types.String `tfsdk:"cc_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
	EncryptionConfiguration types.Object `tfsdk:"encryption_configuration"`
}

type managedClustersDataSourceModel struct {
	ManagedClusters []managedClusterDataSourceModel `tfsdk:"managed_clusters"`
	Filters         types.String                    `tfsdk:"filters"`
}

type managedClusterEncyprionConfigurationModel struct {
	Format types.String `tfsdk:"format"`
}

type managedClusterKeyPairModel struct {
	PublicKey            types.String `tfsdk:"public_key"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
}

type managedClusterAttributesQueueModel struct {
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

type managedClusterAttributesModel struct {
	Queue    managedClusterAttributesQueueModel `tfsdk:"queue"`
	KeyStore types.String                       `tfsdk:"key_store"`
}

type managedClusterRedisModel struct {
	RedisHost types.String `tfsdk:"redis_host"`
	RedisPort types.Int32  `tfsdk:"redis_port"`
}

func parseManagedClusterData(ctx context.Context, cluster api_v2025.ManagedCluster) (managedClusterDataSourceModel, diag.Diagnostics) {

	configuration, diags := types.MapValueFrom(ctx, types.StringType, cluster.Configuration)
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	clientIds, diags := types.ListValueFrom(ctx, types.StringType, cluster.ClientIds)
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	tflog.Trace(ctx, "Reading cluster configuration property", map[string]any{"configuration": cluster.Configuration})

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
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	attributesData := managedClusterAttributesModel{
		Queue: managedClusterAttributesQueueModel{
			Name:   types.StringValue(cluster.GetAttributes().Queue.GetName()),
			Region: types.StringValue(cluster.GetAttributes().Queue.GetRegion()),
		},
		KeyStore: types.StringPointerValue(extractNullableString(cluster.Attributes.GetKeystoreOk())),
	}

	attributesObject, diags := types.ObjectValueFrom(ctx, managedClusterAttributesAttrTypes, attributesData)
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	redisData := managedClusterRedisModel{
		RedisHost: types.StringPointerValue(cluster.GetRedis().RedisHost),
		RedisPort: types.Int32PointerValue(cluster.GetRedis().RedisPort),
	}

	redisObject, diags := types.ObjectValueFrom(ctx, managedClusterRedisAttrTypes, redisData)
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	encryptionConfigData := managedClusterEncyprionConfigurationModel{
		Format: types.StringPointerValue(cluster.GetEncryptionConfiguration().Format),
	}

	encryptionConfigObject, diags := types.ObjectValueFrom(ctx, managedClusterEncryptionConfigAttrTypes, encryptionConfigData)
	if diags != nil {
		return managedClusterDataSourceModel{}, diags
	}

	obj := managedClusterDataSourceModel{
		ID:                      types.StringValue(cluster.GetId()),
		Name:                    types.StringValue(cluster.GetName()),
		Pod:                     types.StringValue(cluster.GetPod()),
		Org:                     types.StringValue(cluster.GetOrg()),
		Type:                    types.StringValue(string(cluster.GetType())),
		Description:             types.StringValue(cluster.GetDescription()),
		ClientType:              types.StringValue(string(cluster.GetClientType())),
		CcgVersion:              types.StringValue(cluster.GetCcgVersion()),
		PinnedConfig:            types.BoolValue(cluster.GetPinnedConfig()),
		Operational:             types.BoolValue(cluster.GetOperational()),
		Status:                  types.StringValue(cluster.GetStatus()),
		PublicKeyCertificate:    types.StringValue(cluster.GetPublicKeyCertificate()),
		PublicKeyThumbprint:     types.StringValue(cluster.GetPublicKeyThumbprint()),
		PublicKey:               types.StringValue(cluster.GetPublicKey()),
		AlertKey:                types.StringValue(cluster.GetAlertKey()),
		ClientIds:               clientIds,
		ServiceCount:            types.Int32Value(cluster.GetServiceCount()),
		CcID:                    types.StringValue(cluster.GetCcId()),
		CreatedAt:               types.StringPointerValue(&createdAt),
		UpdatedAt:               types.StringPointerValue(&updatedAt),
		Configuration:           configuration,
		KeyPair:                 keyPairObject,
		Attributes:              attributesObject,
		Redis:                   redisObject,
		EncryptionConfiguration: encryptionConfigObject,
	}
	return obj, diags
}
