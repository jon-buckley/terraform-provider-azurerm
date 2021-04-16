package migration

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
)

func KeyVaultV1ToV2Upgrader() pluginsdk.StateUpgrader {
	return pluginsdk.StateUpgrader{
		Type:    keyVaultV1Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: keyVaultV1ToV2Upgrade,
		Version: 1,
	}
}

func keyVaultV1ToV2Upgrade(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// @tombuildsstuff: this is an int in the schema but was previously set into the
	// state as `*int32` - so using `.(int)` causes:
	// panic: interface conversion: interface {} is float64, not int
	// so I guess we try both?
	oldVal := 0
	if v, ok := rawState["soft_delete_retention_days"]; ok {
		if val, ok := v.(float64); ok {
			oldVal = int(val)
		}
		if val, ok := v.(int); ok {
			oldVal = val
		}
	}

	if oldVal == 0 {
		// The Azure API force-upgraded all Key Vaults so that Soft Delete was enabled on 2020-12-15
		// As a part of this, all Key Vaults got a default "retention days" of 90 days, however
		// for both newly created and upgraded key vaults, this value isn't returned unless it's
		// explicitly set by a user.
		//
		// As such we have to default this to a value of 90 days, which whilst assuming this default is
		// less than ideal, unfortunately there's little choice otherwise as this isn't returned.
		// Whilst the API Documentation doesn't show this default, it's listed here:
		//
		// > Once a secret, key, certificate, or key vault is deleted, it will remain recoverable
		// > for a configurable period of 7 to 90 calendar days. If no configuration is specified
		// > the default recovery period will be set to 90 days
		// https://docs.microsoft.com/en-us/azure/key-vault/general/soft-delete-overview
		//
		// Notably this value cannot be updated once it's initially been configured, meaning that we
		// must not send this during creation if it's the default value, to allow users to change
		// this value later. This also means we can't use Terraform's "ForceNew" here without breaking
		// the one-time change.
		//
		// Hopefully in time this behaviour is fixed, but for now I'm not sure what else we can do.
		//
		// - @tombuildsstuff
		rawState["soft_delete_retention_days"] = 90
	}
	return rawState, nil
}

func keyVaultV1Schema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VaultName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"access_policy": {
				Type:       pluginsdk.TypeList,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				MaxItems:   1024,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"object_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"application_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"certificate_permissions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"key_permissions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"secret_permissions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"storage_permissions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"enabled_for_deployment": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enabled_for_disk_encryption": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enabled_for_template_deployment": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enable_rbac_authorization": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"bypass": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},
						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      set.HashStringIgnoreCase,
						},
					},
				},
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"soft_delete_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"soft_delete_retention_days": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"contact": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"phone": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),

			// Computed
			"vault_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}
