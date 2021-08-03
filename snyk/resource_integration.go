package snyk

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lendi-au/terraform-provider-snyk/snyk/api"
)

func resourceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		ReadContext:   resourceIntegrationRead,
		UpdateContext: resourceIntegrationUpdate,
		DeleteContext: resourceIntegrationDelete,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getCredentialSchema(),
				},
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	so := m.(api.SnykOptions)

	credentials, err := getCredentialState(d)

	if err != nil {
		return diag.FromErr(err)
	}

	integrationData := api.Integration{
		OrganizationId: d.Get("organization").(string),
		Type:           d.Get("type").(string),
		Credentials:    credentials,
	}

	i, err := api.GetIntegration(so, integrationData)

	var newIntegration api.Integration
	if err != nil { // if integration not found, create it
		newIntegration, err = api.CreateIntegration(so, i)
	} else { // otherwise, reactivate credentials
		newIntegration, err = api.UpdateIntegration(so, i)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newIntegration.Id)
	d.Set("organization", newIntegration.OrganizationId)
	d.Set("type", newIntegration.Type)
	setCredentialState(credentials, d)

	return diags
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	so := m.(api.SnykOptions)

	credentials, err := getCredentialState(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// For now, we're just going to traverse the list of integrations in the org and verify it exists
	integrationData := api.Integration{
		Id:             d.Id(),
		OrganizationId: d.Get("organization").(string),
		Type:           d.Get("type").(string),
		Credentials:    credentials,
	}

	i, err := api.GetIntegration(so, integrationData)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(i.Id)
	d.Set("organization", i.OrganizationId)
	d.Set("type", i.Type)
	setCredentialState(credentials, d)

	return diags
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	so := m.(api.SnykOptions)

	credentials, err := getCredentialState(d)

	if err != nil {
		return diag.FromErr(err)
	}

	integrationData := api.Integration{
		Id:             d.Id(),
		OrganizationId: d.Get("organization").(string),
		Type:           d.Get("type").(string),
		Credentials:    credentials,
	}

	i, err := api.UpdateIntegration(so, integrationData)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i.Id)
	d.Set("organization", i.OrganizationId)
	d.Set("type", i.Type)
	setCredentialState(credentials, d)

	return diags
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	so := m.(api.SnykOptions)

	integrationData := api.Integration{
		Id:             d.Id(),
		OrganizationId: d.Get("organization").(string),
		Type:           d.Get("type").(string),
	}

	err := api.DeleteIntegration(so, integrationData)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"registry_base": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"url": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"token": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"role_arn": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func getCredentialState(d *schema.ResourceData) (api.IntegrationCredentials, error) {
	if credList, ok := d.Get("credentials").([]interface{}); ok {
		creds := credList[0].(map[string]interface{})
		return api.IntegrationCredentials{
			Username:     creds["username"].(string),
			Password:     creds["password"].(string),
			RegistryBase: creds["registry_base"].(string),
			Url:          creds["url"].(string),
			Token:        creds["token"].(string),
			Region:       creds["region"].(string),
			RoleArn:      creds["role_arn"].(string),
		}, nil
	}

	return api.IntegrationCredentials{}, errors.New("unable to fetch credentials from state")
}

func setCredentialState(creds api.IntegrationCredentials, d *schema.ResourceData) {

	stateList := make([]interface{}, 1)

	stateMap := make(map[string]interface{})

	if creds.Username != "" {
		stateMap["username"] = creds.Username
	}
	if creds.Password != "" {
		stateMap["password"] = creds.Password
	}
	if creds.RegistryBase != "" {
		stateMap["registry_base"] = creds.RegistryBase
	}
	if creds.Url != "" {
		stateMap["url"] = creds.Url
	}
	if creds.Token != "" {
		stateMap["token"] = creds.Token
	}
	if creds.Region != "" {
		stateMap["region"] = creds.Region
	}
	if creds.RoleArn != "" {
		stateMap["role_arn"] = creds.RoleArn
	}

	stateList[0] = stateMap

	d.Set("credentials", stateList)
}
