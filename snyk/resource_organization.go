package snyk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lendi-au/terraform-provider-snyk/snyk/api"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		DeleteContext: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	so := m.(api.SnykOptions)
	name := d.Get("name").(string)

	org, err := api.CreateOrganization(so, name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(org.Id)
	d.Set("created", org.Created)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)
	d.Set("url", org.Url)

	return resourceOrganizationRead(ctx, d, m)
}

func resourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	so := m.(api.SnykOptions)
	id := d.Id()

	org, err := api.GetOrganization(so, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("created", org.Created)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)
	d.Set("url", org.Url)

	return diags
}

func resourceOrganizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	so := m.(api.SnykOptions)
	id := d.Id()

	err := api.DeleteOrganization(so, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
