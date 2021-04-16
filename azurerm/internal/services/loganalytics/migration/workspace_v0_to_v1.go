package migration

import (
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func WorkspaceV0ToV1() pluginsdk.StateUpgrader {
	return pluginsdk.StateUpgrader{
		Version: 0,
		Type:    workspaceV0V1Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: workspaceUpgradeV0ToV1,
	}
}

func workspaceV0V1Schema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"internet_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"internet_query_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"daily_quota_gb": {
				Type:     pluginsdk.TypeFloat,
				Optional: true,
				Default:  -1.0,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func workspaceUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[DEBUG] Migrating IDs to correct casing for Log Analytics Workspace")
	name := rawState["name"].(string)
	resourceGroup := rawState["resource_group_name"].(string)
	id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, name)

	rawState["id"] = id.ID()
	return rawState, nil
}
