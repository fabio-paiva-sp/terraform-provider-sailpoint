package provider

// import (
// 	"context"

// 	"github.com/hashicorp/terraform-plugin-framework/attr"
// 	dataSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
// )

// var (
// 	// managedClientKeyPairAttrTypes = map[string]attr.Type{
// 	// 	"public_key":             types.StringType,
// 	// 	"public_key_thumbprint":  types.StringType,
// 	// 	"public_key_certificate": types.StringType,
// 	// }
// 	// managedClientAttributesAttrTypes = map[string]attr.Type{
// 	// 	"queue": types.ObjectType{
// 	// 		AttrTypes: map[string]attr.Type{
// 	// 			"name":   types.StringType,
// 	// 			"region": types.StringType,
// 	// 		},
// 	// 	},
// 	// 	"key_store": types.StringType,
// 	// }
// 	// managedClientRedisAttrTypes = map[string]attr.Type{
// 	// 	"redis_host": types.StringType,
// 	// 	"redis_port": types.Int32Type,
// 	// }
// 	// managedClientEncryptionConfigAttrTypes = map[string]attr.Type{
// 	// 	"format": types.StringType,
// 	// }
// 	managedClientDataSourceSchemaAttributes = map[string]dataSchema.Attribute{
// 		"id": dataSchema.StringAttribute{
// 			Required: true,
// 		},
// 		"name": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"description": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"alert_key": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"api_gateway_base_url": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"cookbook": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"cc_id": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"cluster_id": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"ip_address": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"pinned_config": dataSchema.BoolAttribute{
// 			Computed: true,
// 		},
// 		"operational": dataSchema.BoolAttribute{
// 			Computed: true,
// 		},
// 		"status": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"public_key_certificate": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"public_key_thumbprint": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"public_key": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"alert_key": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"client_ids": dataSchema.ListAttribute{
// 			Computed:    true,
// 			ElementType: types.StringType,
// 		},
// 		"service_count": dataSchema.Int32Attribute{
// 			Computed: true,
// 		},
// 		"cc_id": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"created_at": dataSchema.StringAttribute{
// 			Computed: true,
// 		},
// 		"encryption_configuration": dataSchema.ObjectAttribute{
// 			Computed:       true,
// 			AttributeTypes: managedClientEncryptionConfigAttrTypes,
// 		},
// 	}
// 	managedClientResourceSchemaAttributes = map[string]resourceSchema.Attribute{
// 		"id": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"name": resourceSchema.StringAttribute{
// 			Required: true,
// 		},
// 		"description": resourceSchema.StringAttribute{
// 			Optional: true,
// 			Computed: true, // API converts null to empty string
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"type": resourceSchema.StringAttribute{
// 			Optional: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.RequiresReplace(),
// 			},
// 		},
// 		// configuration is supposed to be dynamic but the SailPoint Go SDK maps it to string:string
// 		"configuration": resourceSchema.MapAttribute{
// 			Optional:    true,
// 			Computed:    true,
// 			ElementType: types.StringType,
// 			PlanModifiers: []planmodifier.Map{
// 				mapplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"pod": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"org": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"key_pair": resourceSchema.ObjectAttribute{
// 			Computed:       true,
// 			AttributeTypes: managedClientKeyPairAttrTypes,
// 			PlanModifiers: []planmodifier.Object{
// 				objectplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"attributes": resourceSchema.ObjectAttribute{
// 			Computed:       true,
// 			AttributeTypes: managedClientAttributesAttrTypes,
// 			PlanModifiers: []planmodifier.Object{
// 				objectplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"redis": resourceSchema.ObjectAttribute{
// 			Computed:       true,
// 			AttributeTypes: managedClientRedisAttrTypes,
// 			PlanModifiers: []planmodifier.Object{
// 				objectplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"client_type": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"ccg_version": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"pinned_config": resourceSchema.BoolAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.Bool{
// 				boolplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"operational": resourceSchema.BoolAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.Bool{
// 				boolplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"status": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"public_key_certificate": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"public_key_thumbprint": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"public_key": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"alert_key": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"client_ids": resourceSchema.ListAttribute{
// 			Computed:    true,
// 			ElementType: types.StringType,
// 			PlanModifiers: []planmodifier.List{
// 				listplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"service_count": resourceSchema.Int32Attribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.Int32{
// 				int32planmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"cc_id": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"created_at": resourceSchema.StringAttribute{
// 			Computed: true,
// 			PlanModifiers: []planmodifier.String{
// 				stringplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 		"encryption_configuration": resourceSchema.ObjectAttribute{
// 			Computed:       true,
// 			AttributeTypes: managedClientEncryptionConfigAttrTypes,
// 			PlanModifiers: []planmodifier.Object{
// 				objectplanmodifier.UseStateForUnknown(),
// 			},
// 		},
// 	}
// )

// type managedClientSourceModel struct {
// 	ID                      types.String `tfsdk:"id"`
// 	Name                    types.String `tfsdk:"name"`
// 	Pod                     types.String `tfsdk:"pod"`
// 	Org                     types.String `tfsdk:"org"`
// 	Type                    types.String `tfsdk:"type"`
// 	Configuration           types.Map    `tfsdk:"configuration"`
// 	KeyPair                 types.Object `tfsdk:"key_pair"`
// 	Attributes              types.Object `tfsdk:"attributes"`
// 	Redis                   types.Object `tfsdk:"redis"`
// 	Description             types.String `tfsdk:"description"`
// 	ClientType              types.String `tfsdk:"client_type"`
// 	CcgVersion              types.String `tfsdk:"ccg_version"`
// 	PinnedConfig            types.Bool   `tfsdk:"pinned_config"`
// 	Operational             types.Bool   `tfsdk:"operational"`
// 	Status                  types.String `tfsdk:"status"`
// 	PublicKeyCertificate    types.String `tfsdk:"public_key_certificate"`
// 	PublicKeyThumbprint     types.String `tfsdk:"public_key_thumbprint"`
// 	PublicKey               types.String `tfsdk:"public_key"`
// 	AlertKey                types.String `tfsdk:"alert_key"`
// 	ClientIds               types.List   `tfsdk:"client_ids"`
// 	ServiceCount            types.Int32  `tfsdk:"service_count"`
// 	CcID                    types.String `tfsdk:"cc_id"`
// 	CreatedAt               types.String `tfsdk:"created_at"`
// 	EncryptionConfiguration types.Object `tfsdk:"encryption_configuration"`
// }

// type managedClientsDataSourceModel struct {
// 	ManagedClients []managedClientSourceModel `tfsdk:"managed_clusters"`
// 	Filters         types.String                `tfsdk:"filters"`
// }

// type managedClientEncyprionConfigurationModel struct {
// 	Format types.String `tfsdk:"format"`
// }

// type managedClientKeyPairModel struct {
// 	PublicKey            types.String `tfsdk:"public_key"`
// 	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
// 	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
// }

// type managedClientAttributesQueueModel struct {
// 	Name   types.String `tfsdk:"name"`
// 	Region types.String `tfsdk:"region"`
// }

// type managedClientAttributesModel struct {
// 	Queue    managedClientAttributesQueueModel `tfsdk:"queue"`
// 	KeyStore types.String                       `tfsdk:"key_store"`
// }

// type managedClientRedisModel struct {
// 	RedisHost types.String `tfsdk:"redis_host"`
// 	RedisPort types.Int32  `tfsdk:"redis_port"`
// }

// func serializeManagedClientData(ctx context.Context, cluster api_v2025.ManagedClient) (managedClientSourceModel, diag.Diagnostics) {

// 	configuration, diags := types.MapValueFrom(ctx, types.StringType, cluster.Configuration)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	clientIds, diags := types.ListValueFrom(ctx, types.StringType, cluster.ClientIds)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	tflog.Trace(ctx, "Reading cluster configuration property", map[string]any{"configuration": cluster.Configuration})

// 	var createdAt string
// 	createdAtDate, _ := cluster.GetCreatedAtOk()

// 	if createdAtDate != nil {
// 		createdAt = createdAtDate.String()
// 		tflog.Trace(ctx, "Reading cluster created at property string", map[string]any{"date": createdAtDate.String()})
// 	}

// 	keyPairData := managedClientKeyPairModel{
// 		PublicKey:            types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyOk())),
// 		PublicKeyThumbprint:  types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyThumbprintOk())),
// 		PublicKeyCertificate: types.StringPointerValue(extractNullableString(cluster.KeyPair.GetPublicKeyCertificateOk())),
// 	}

// 	keyPairObject, diags := types.ObjectValueFrom(ctx, managedClientKeyPairAttrTypes, keyPairData)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	attributesData := managedClientAttributesModel{
// 		Queue: managedClientAttributesQueueModel{
// 			Name:   types.StringValue(cluster.GetAttributes().Queue.GetName()),
// 			Region: types.StringValue(cluster.GetAttributes().Queue.GetRegion()),
// 		},
// 		KeyStore: types.StringPointerValue(extractNullableString(cluster.Attributes.GetKeystoreOk())),
// 	}

// 	attributesObject, diags := types.ObjectValueFrom(ctx, managedClientAttributesAttrTypes, attributesData)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	redisData := managedClientRedisModel{
// 		RedisHost: types.StringPointerValue(cluster.GetRedis().RedisHost),
// 		RedisPort: types.Int32PointerValue(cluster.GetRedis().RedisPort),
// 	}

// 	redisObject, diags := types.ObjectValueFrom(ctx, managedClientRedisAttrTypes, redisData)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	encryptionConfigData := managedClientEncyprionConfigurationModel{
// 		Format: types.StringPointerValue(cluster.GetEncryptionConfiguration().Format),
// 	}

// 	encryptionConfigObject, diags := types.ObjectValueFrom(ctx, managedClientEncryptionConfigAttrTypes, encryptionConfigData)
// 	if diags != nil {
// 		return managedClientSourceModel{}, diags
// 	}

// 	obj := managedClientSourceModel{
// 		ID:                      types.StringValue(cluster.GetId()),
// 		Name:                    types.StringValue(cluster.GetName()),
// 		Pod:                     types.StringValue(cluster.GetPod()),
// 		Org:                     types.StringValue(cluster.GetOrg()),
// 		Type:                    types.StringValue(string(cluster.GetType())),
// 		Description:             types.StringValue(cluster.GetDescription()),
// 		ClientType:              types.StringValue(string(cluster.GetClientType())),
// 		CcgVersion:              types.StringValue(cluster.GetCcgVersion()),
// 		PinnedConfig:            types.BoolValue(cluster.GetPinnedConfig()),
// 		Operational:             types.BoolValue(cluster.GetOperational()),
// 		Status:                  types.StringValue(cluster.GetStatus()),
// 		PublicKeyCertificate:    types.StringValue(cluster.GetPublicKeyCertificate()),
// 		PublicKeyThumbprint:     types.StringValue(cluster.GetPublicKeyThumbprint()),
// 		PublicKey:               types.StringValue(cluster.GetPublicKey()),
// 		AlertKey:                types.StringValue(cluster.GetAlertKey()),
// 		ClientIds:               clientIds,
// 		ServiceCount:            types.Int32Value(cluster.GetServiceCount()),
// 		CcID:                    types.StringValue(cluster.GetCcId()),
// 		CreatedAt:               types.StringPointerValue(&createdAt),
// 		Configuration:           configuration,
// 		KeyPair:                 keyPairObject,
// 		Attributes:              attributesObject,
// 		Redis:                   redisObject,
// 		EncryptionConfiguration: encryptionConfigObject,
// 	}
// 	return obj, diags
// }
