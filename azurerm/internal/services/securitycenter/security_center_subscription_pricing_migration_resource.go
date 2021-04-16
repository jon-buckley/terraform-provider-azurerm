package securitycenter

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func ResourceSecurityCenterSubscriptionPricingV0() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.PricingTierFree),
					string(security.PricingTierStandard),
				}, false),
			},
		},
	}
}

func ResourceSecurityCenterSubscriptionPricingUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Println("[DEBUG] Migrating ResourceType from v0 to v1 format")
	oldId := rawState["id"].(string)
	newId := strings.Replace(oldId, "/default", "/VirtualMachines", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
