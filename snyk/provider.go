package snyk

import (
	"context"

	"github.com/lendi-au/terraform-provider-snyk/snyk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
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
				"snyk_integration":  resourceIntegration(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"snyk_organization": dataSourceOrganization(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		var diags diag.Diagnostics

		config := api.SnykOptions{
			GroupId:   d.Get("group_id").(string),
			ApiKey:    d.Get("api_key").(string),
			UserAgent: p.UserAgent("terraform-provider-snyk", version),
		}

		return config, diags
	}
}
