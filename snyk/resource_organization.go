package snyk

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lendi-au/terraform-provider-snyk/snyk/api"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		DeleteContext: resourceOrganizationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	d.Set("name", org.Name)

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

