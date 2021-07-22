package snyk

import (
	"context"
	"errors"
	"terraform-provider-snyk/snyk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		UpdateContext: resourceOrganizationUpdate,
		DeleteContext: resourceOrganizationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notifications": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: getOrganizationNotificationSchema(),
				},
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

	notState, err := getStateNotificationOptions(d)

	if err != nil {
		return diag.FromErr(err)
	}

	_, err = api.SetOrgNotificationSettings(so, d.Id(), notState)

	if err != nil {
		return diag.FromErr(err)
	}

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

	nots, err := api.GetOrgNotificationSettings(so, id)

	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch old value of new issue types (if it existed), as the API sets it to "" when it is disabled
	// Sorry for jank, doesn't work well with nested structs

	if _, existed := d.GetOk("notifications"); existed {
		if n, ok := d.Get("notifications").([]interface{}); ok {
			if mn, ok := n[0].(map[string]interface{}); ok {
				if t, ok := mn["new_issues"].([]interface{}); ok {
					if mt, ok := t[0].(map[string]interface{}); ok {
						if !nots.NewIssuesRemediations.Enabled {
							nots.NewIssuesRemediations.IssueType = mt["type"].(string)
						}
					}
				}

			}
		}
	}

	d.Set("name", org.Name)
	setStateNotificationOptions(nots, d)

	return diags
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	so := m.(api.SnykOptions)
	id := d.Id()

	state, err := getStateNotificationOptions(d)

	if err != nil {
		return diag.FromErr(err)
	}

	_, err = api.SetOrgNotificationSettings(so, id, state)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOrganizationRead(ctx, d, m)
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

// Schema Definitions

func getOrganizationNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"new_issues": {
			Type:     schema.TypeList,
			MaxItems: 1,
			MinItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: getNewIssuesNotificationSchema(),
			},
		},
		"project_imports": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"test_limits": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"weekly_report": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}

func getNewIssuesNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"severity": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "high",
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "none",
		},
	}
}

// Helper functions
func getStateNotificationOptions(d *schema.ResourceData) (api.OrganizationNotifications, error) {
	if notState, ok := d.Get("notifications").([]interface{}); ok {
		if notStateMap, ok := notState[0].(map[string]interface{}); ok {
			if newIssues, ok := notStateMap["new_issues"].([]interface{}); ok {
				if newIssuesMap, ok := newIssues[0].(map[string]interface{}); ok {
					return api.OrganizationNotifications{
						NewIssuesRemediations: api.NewIssuesRemediationsOption{
							Enabled:       newIssuesMap["enabled"].(bool),
							IssueSeverity: newIssuesMap["severity"].(string),
							IssueType:     newIssuesMap["type"].(string),
						},
						ProjectImported: api.ProjectImportedOption{
							Enabled: notStateMap["project_imports"].(bool),
						},
						TestLimit: api.TestLimitOption{
							Enabled: notStateMap["test_limits"].(bool),
						},
						WeeklyReport: api.WeeklyReportOption{
							Enabled: notStateMap["weekly_report"].(bool),
						},
					}, nil
				}
			}
		}
	}

	return api.OrganizationNotifications{}, errors.New("failed getting notification options from state")

}

func setStateNotificationOptions(nots api.OrganizationNotifications, d *schema.ResourceData) {
	n := make([]map[string]interface{}, 1)

	newNotificationsState := make(map[string]interface{})
	newIssuesState := make([]map[string]interface{}, 1)
	newIssuesOptions := make(map[string]interface{})
	newIssuesOptions["enabled"] = nots.NewIssuesRemediations.Enabled
	newIssuesOptions["severity"] = nots.NewIssuesRemediations.IssueSeverity
	newIssuesOptions["type"] = nots.NewIssuesRemediations.IssueType
	newIssuesState[0] = newIssuesOptions

	newNotificationsState["new_issues"] = newIssuesState
	newNotificationsState["project_imports"] = nots.ProjectImported.Enabled
	newNotificationsState["test_limits"] = nots.TestLimit.Enabled
	newNotificationsState["weekly_report"] = nots.WeeklyReport.Enabled

	n[0] = newNotificationsState

	d.Set("notifications", n)
}
