package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	managedClusterDataSourceSchemaAttributes = map[string]dataSchema.Attribute{
		"id": dataSchema.StringAttribute{
			Required: true,
		},
		"name": dataSchema.StringAttribute{
			Computed: true,
		},
		"pod": dataSchema.StringAttribute{
			Computed: true,
		},
		"org": dataSchema.StringAttribute{
			Computed: true,
		},
		"type": dataSchema.StringAttribute{
			Computed: true,
		},
		"configuration": dataSchema.MapAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"key_pair": dataSchema.ObjectAttribute{
			Optional:       true,
			Computed:       false,
			AttributeTypes: managedClusterKeyPairAttrTypes,
		},
		"attributes": dataSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterAttributesAttrTypes,
		},
		"redis": dataSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterRedisAttrTypes,
		},
		"description": dataSchema.StringAttribute{
			Computed: true,
		},
		"client_type": dataSchema.StringAttribute{
			Computed: true,
		},
		"ccg_version": dataSchema.StringAttribute{
			Computed: true,
		},
		"pinned_config": dataSchema.BoolAttribute{
			Computed: true,
		},
		"operational": dataSchema.BoolAttribute{
			Computed: true,
		},
		"status": dataSchema.StringAttribute{
			Computed: true,
		},
		"public_key_certificate": dataSchema.StringAttribute{
			Computed: true,
		},
		"public_key_thumbprint": dataSchema.StringAttribute{
			Computed: true,
		},
		"public_key": dataSchema.StringAttribute{
			Computed: true,
		},
		"alert_key": dataSchema.StringAttribute{
			Computed: true,
		},
		"client_ids": dataSchema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_count": dataSchema.Int32Attribute{
			Computed: true,
		},
		"cc_id": dataSchema.StringAttribute{
			Computed: true,
		},
		"created_at": dataSchema.StringAttribute{
			Computed: true,
		},
		"encryption_configuration": dataSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterEncryptionConfigAttrTypes,
		},
	}
	managedClusterResourceSchemaAttributes = map[string]resourceSchema.Attribute{
		"id": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": resourceSchema.StringAttribute{
			Required: true,
		},
		"description": resourceSchema.StringAttribute{
			Optional: true,
			Computed: true, // API converts null to empty string
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": resourceSchema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		// configuration is supposed to be dynamic but the SailPoint Go SDK maps it to string:string
		"configuration": resourceSchema.MapAttribute{
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Map{
				mapplanmodifier.UseStateForUnknown(),
			},
		},
		"pod": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"org": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"key_pair": resourceSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterKeyPairAttrTypes,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
		"attributes": resourceSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterAttributesAttrTypes,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
		"redis": resourceSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterRedisAttrTypes,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
		"client_type": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ccg_version": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"pinned_config": resourceSchema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"operational": resourceSchema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"status": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"public_key_certificate": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"public_key_thumbprint": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"public_key": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"alert_key": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"client_ids": resourceSchema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"service_count": resourceSchema.Int32Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
		},
		"cc_id": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_at": resourceSchema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"encryption_configuration": resourceSchema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: managedClusterEncryptionConfigAttrTypes,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
	}
)

type managedClusterSourceModel struct {
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
	EncryptionConfiguration types.Object `tfsdk:"encryption_configuration"`
}

type managedClustersDataSourceModel struct {
	ManagedClusters []managedClusterSourceModel `tfsdk:"managed_clusters"`
	Filters         types.String                `tfsdk:"filters"`
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

func serializeManagedClusterData(ctx context.Context, cluster api_v2025.ManagedCluster) (managedClusterSourceModel, diag.Diagnostics) {

	configuration, diags := types.MapValueFrom(ctx, types.StringType, cluster.Configuration)
	if diags != nil {
		return managedClusterSourceModel{}, diags
	}

	clientIds, diags := types.ListValueFrom(ctx, types.StringType, cluster.ClientIds)
	if diags != nil {
		return managedClusterSourceModel{}, diags
	}

	tflog.Trace(ctx, "Reading cluster configuration property", map[string]any{"configuration": cluster.Configuration})

	var createdAt string
	createdAtDate, _ := cluster.GetCreatedAtOk()

	if createdAtDate != nil {
		createdAt = createdAtDate.String()
		tflog.Trace(ctx, "Reading cluster created at property string", map[string]any{"date": createdAtDate.String()})
	}

	keyPairData := managedClusterKeyPairModel{
		PublicKey:            types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyOk())),
		PublicKeyThumbprint:  types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyThumbprintOk())),
		PublicKeyCertificate: types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyCertificateOk())),
	}

	keyPairObject, diags := types.ObjectValueFrom(ctx, managedClusterKeyPairAttrTypes, keyPairData)
	if diags != nil {
		return managedClusterSourceModel{}, diags
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
		return managedClusterSourceModel{}, diags
	}

	redisData := managedClusterRedisModel{
		RedisHost: types.StringPointerValue(cluster.GetRedis().RedisHost),
		RedisPort: types.Int32PointerValue(cluster.GetRedis().RedisPort),
	}

	redisObject, diags := types.ObjectValueFrom(ctx, managedClusterRedisAttrTypes, redisData)
	if diags != nil {
		return managedClusterSourceModel{}, diags
	}

	encryptionConfigData := managedClusterEncyprionConfigurationModel{
		Format: types.StringPointerValue(cluster.GetEncryptionConfiguration().Format),
	}

	encryptionConfigObject, diags := types.ObjectValueFrom(ctx, managedClusterEncryptionConfigAttrTypes, encryptionConfigData)
	if diags != nil {
		return managedClusterSourceModel{}, diags
	}

	obj := managedClusterSourceModel{
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
		Configuration:           configuration,
		KeyPair:                 keyPairObject,
		Attributes:              attributesObject,
		Redis:                   redisObject,
		EncryptionConfiguration: encryptionConfigObject,
	}
	return obj, diags
}
