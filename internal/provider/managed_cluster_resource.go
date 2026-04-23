package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &mangedClusterResource{}
	_ resource.ResourceWithConfigure   = &mangedClusterResource{}
	_ resource.ResourceWithImportState = &mangedClusterResource{}
)

// NewManagedClusterResource is a helper function to simplify the provider implementation.
func NewManagedClusterResource() resource.Resource {
	return &mangedClusterResource{}
}

// mangedClusterResource is the resource implementation.
type mangedClusterResource struct {
	client *sailpoint.APIClient
}

// Metadata returns the resource type name.
func (r *mangedClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

// Schema defines the schema for the resource.
func (r *mangedClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: managedClusterResourceSchemaAttributes,
	}
}

func (r *mangedClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	tflog.Info(ctx, "Configuring SailPoint ManagedCluster resource")

	client, ok := req.ProviderData.(*sailpoint.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *mangedClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating managed cluster resource")

	// Retrieve values from plan
	var plan managedClusterSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var configuration *map[string]string
	configuration = &map[string]string{}
	for k, v := range plan.Configuration.Elements() {
		// Dereference the pointer before indexing the map
		var result string
		switch val := v.(type) {
		case types.String:
			s := val.ValueString()
			result = s
		default:
			s := val.String()
			s = strings.Trim(s, "\"")
			tflog.Warn(ctx, "Creation - Striping quotes from string", map[string]any{"k": k, "v": v, "val": val, "s": s})
			result = s
		}
		(*configuration)[k] = result
	}
	var description api_v2025.NullableString
	if !plan.Description.IsNull() {
		description.Set(plan.Description.ValueStringPointer())
	}

	// Generate API request body from plan
	managedCluster := api_v2025.ManagedClusterRequest{
		Name:          plan.Name.ValueString(),
		Type:          (*api_v2025.ManagedClusterTypes)(plan.Type.ValueStringPointer()),
		Description:   description,
		Configuration: configuration,
	}
	tflog.Info(ctx, "Creating managed cluster with the values", map[string]any{"cluster": managedCluster})

	// Create new cluster
	cluster, res, err := r.client.V2025.ManagedClustersAPI.CreateManagedCluster(ctx).ManagedClusterRequest(managedCluster).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			decodedBytes, _ := base64.StdEncoding.DecodeString(string(bodyBytes))
			tflog.Error(ctx, "error creating cluster", map[string]any{"error": err.Error(), "response_body": decodedBytes})
			if decodedBytes != nil && string(decodedBytes) != "" {
				resp.Diagnostics.AddError(
					"unable to create Managed Cluster",
					string(decodedBytes),
				)
			} else {
				resp.Diagnostics.AddError(
					"unable to create Managed Cluster",
					string(bodyBytes),
				)
			}
		}
		resp.Diagnostics.AddError(
			"unable to create Managed Cluster",
			err.Error(),
		)
		return
	}

	// Get refreshed managed cluster value from Sailpoint API
	cluster, res, err = r.client.V2025.ManagedClustersAPI.GetManagedCluster(ctx, cluster.Id).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			tflog.Error(ctx, "error reading cluster resource", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		}
		resp.Diagnostics.AddError(
			"unable to read Managed Cluster resource",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state, diags := serializeManagedClusterData(ctx, *cluster)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	planConfig := make(map[string]string)
	plan.Configuration.ElementsAs(ctx, &planConfig, false)
	stateConfig := make(map[string]string)
	state.Configuration.ElementsAs(ctx, &stateConfig, false)

	filteredConfig := make(map[string]string)
	for k, v := range stateConfig {
		// If it's in the plan, keep it
		if _, exists := planConfig[k]; exists {
			tflog.Debug(ctx, "keeping config attr", map[string]any{"k": k})
			filteredConfig[k] = v
		}
	}
	state.Configuration, diags = types.MapValueFrom(ctx, types.StringType, filteredConfig)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "finish creating managed cluster resource")
}

// Read refreshes the Terraform state with the latest data.
func (r *mangedClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading managed cluster resource")
	// Get current state
	var state managedClusterSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateConfig := make(map[string]string)
	state.Configuration.ElementsAs(ctx, &stateConfig, false)

	// Get refreshed managed cluster value from Sailpoint API
	cluster, res, err := r.client.V2025.ManagedClustersAPI.GetManagedCluster(ctx, state.ID.ValueString()).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			tflog.Error(ctx, "error reading cluster resource", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		}
		resp.Diagnostics.AddError(
			"unable to read Managed Cluster resource",
			err.Error(),
		)
		return
	}

	state, diags = serializeManagedClusterData(ctx, *cluster)

	filteredConfig := make(map[string]string)
	for k, v := range *cluster.Configuration {
		// If it's in the state keep it
		if _, exists := stateConfig[k]; exists {
			filteredConfig[k] = v
		}
	}

	tflog.Debug(ctx, "Managed cluster filtered configuration: ", map[string]any{"filteredConfig": filteredConfig, "stateConfig": stateConfig})

	state.Configuration, _ = types.MapValueFrom(ctx, types.StringType, filteredConfig)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "finish reading managed cluster resource")
}

func buildConfigurationPatchOps(ctx context.Context, newConfig map[string]attr.Value, oldConfig map[string]attr.Value) ([]api_v2025.JsonPatchOperation, error) {
	tflog.Debug(ctx, "building configuration block", map[string]any{"newConfig": newConfig, "oldConfig": oldConfig})
	ops := make([]api_v2025.JsonPatchOperation, 0)

	for k, v := range newConfig {
		exists := oldConfig[k] != nil

		tflog.Debug(ctx, "building configuration block for new config", map[string]any{"k": k, "v": v, "exists": exists})
		op := api_v2025.NewJsonPatchOperation("replace", fmt.Sprintf("/configuration/%s", k))

		if !exists {
			op.SetOp("add")
		}
		switch val := v.(type) {
		case types.String:
			s := val.ValueString()
			op.SetValue(api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&s))
		default:
			s := val.String()
			s = strings.Trim(s, "\"")
			tflog.Warn(ctx, "Striping quotes from string", map[string]any{"k": k, "v": v, "val": val, "s": s})
			op.SetValue(api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&s))
		}

		ops = append(ops, *op)
	}

	for k := range oldConfig {
		exists := newConfig[k] != nil
		tflog.Debug(ctx, "building configuration block for old config", map[string]any{"k": k, "exists": exists})
		if !exists {
			op := api_v2025.NewJsonPatchOperation("remove", fmt.Sprintf("/configuration/%s", k))
			ops = append(ops, *op)
		}
	}

	return ops, nil
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *mangedClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating managed cluster resource")

	var (
		plan  managedClusterSourceModel
		state managedClusterSourceModel
	)

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "updating managed cluster resource with ID", map[string]any{"id": plan.ID.ValueString()})

	jpOps := api_v2025.NewJsonPatch()
	operations := make([]api_v2025.JsonPatchOperation, 0)

	if !plan.Name.Equal(state.Name) {
		op := api_v2025.NewJsonPatchOperation("replace", "/name")
		op.SetValue(api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(plan.Name.ValueStringPointer()))
		operations = append(operations, *op)
	}
	if !plan.Description.Equal(state.Description) {
		op := api_v2025.NewJsonPatchOperation("replace", "/description")
		op.SetValue(api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(plan.Description.ValueStringPointer()))
		operations = append(operations, *op)
	}
	if !plan.Configuration.IsUnknown() && !plan.Configuration.Equal(state.Configuration) {
		ops, err := buildConfigurationPatchOps(ctx, plan.Configuration.Elements(), state.Configuration.Elements())
		if err != nil {
			resp.Diagnostics.AddError(
				"unable to update Managed Cluster",
				err.Error(),
			)
			return
		}
		operations = append(operations, ops...)
	}
	jpOps.SetOperations(operations)

	tflog.Debug(ctx, "patching managed cluster with the values", map[string]any{"patchOps": operations})

	// Update cluster
	cluster, res, err := r.client.V2025.ManagedClustersAPI.UpdateManagedCluster(ctx, plan.ID.ValueString()).JsonPatchOperation(jpOps.GetOperations()).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			tflog.Error(ctx, "error updating cluster", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		}
		resp.Diagnostics.AddError(
			"unable to update Managed Cluster",
			err.Error(),
		)
		return
	}

	// Get refreshed managed cluster value from Sailpoint API
	cluster, res, err = r.client.V2025.ManagedClustersAPI.GetManagedCluster(ctx, plan.ID.ValueString()).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			tflog.Error(ctx, "error reading cluster resource", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		}
		resp.Diagnostics.AddError(
			"unable to read Managed Cluster resource",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	tflog.Debug(ctx, "serializing cluster updated data", map[string]any{"cluster": cluster})
	state, diags = serializeManagedClusterData(ctx, *cluster)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	planConfig := make(map[string]string)
	plan.Configuration.ElementsAs(ctx, &planConfig, false)
	stateConfig := make(map[string]string)
	state.Configuration.ElementsAs(ctx, &stateConfig, false)

	filteredConfig := make(map[string]string)
	for k, v := range stateConfig {
		// If it's in the plan, keep it
		if _, exists := planConfig[k]; exists {
			tflog.Debug(ctx, "keeping config attr", map[string]any{"k": k})
			filteredConfig[k] = v
		}
	}
	state.Configuration, diags = types.MapValueFrom(ctx, types.StringType, filteredConfig)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "persisting state with the values", map[string]any{"state": state, "configuration": state.Configuration})
	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "finish updating managed cluster resource")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *mangedClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting managed cluster resource")

	var state managedClusterSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "deleting managed cluster resource with ID", map[string]any{"id": state.ID.ValueString()})

	res, err := r.client.V2025.ManagedClustersAPI.DeleteManagedCluster(ctx, state.ID.ValueString()).Execute()

	if err != nil {
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			bodyBytes, _ := io.ReadAll(res.Body)
			tflog.Error(ctx, "error deleting cluster resource", map[string]any{"error": err.Error(), "response_body": bodyBytes})
		}
		resp.Diagnostics.AddError(
			"unable to delete Managed Cluster resource",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "finish deleting managed cluster resource")
}

func (r *mangedClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
