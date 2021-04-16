package databasemigration

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDatabaseMigrationProject() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDatabaseMigrationProjectRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validateDatabaseMigrationProjectName,
			},

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validateDatabaseMigrationServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"source_platform": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"target_platform": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDatabaseMigrationProjectRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	serviceName := d.Get("service_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Database Migration Project (Project Name %q / Service Name %q / Group Name %q) was not found", name, serviceName, resourceGroup)
		}
		return fmt.Errorf("Error reading Database Migration Project (Project Name %q / Service Name %q / Group Name %q): %+v", name, serviceName, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Database Migration Project (Project Name %q / Service Name %q / Group Name %q) ID", name, serviceName, resourceGroup)
	}
	d.SetId(*resp.ID)

	d.Set("resource_group_name", resourceGroup)

	location := ""
	if resp.Location != nil {
		location = azure.NormalizeLocation(*resp.Location)
	}
	d.Set("location", location)

	if prop := resp.ProjectProperties; prop != nil {
		d.Set("source_platform", string(prop.SourcePlatform))
		d.Set("target_platform", string(prop.TargetPlatform))
	}

	return nil
}
