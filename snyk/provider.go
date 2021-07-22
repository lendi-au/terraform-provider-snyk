package snyk

import (
	"context"
	"terraform-provider-snyk/snyk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNYK_API_GROUP", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SNYK_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snyk_organization": resourceOrganization(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"snyk_organization": dataSourceOrganization(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := api.SnykOptions{
		GroupId: d.Get("group_id").(string),
		ApiKey:  d.Get("api_key").(string),
	}

	return config, diags
}
