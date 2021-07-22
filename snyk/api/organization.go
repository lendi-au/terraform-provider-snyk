package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Organization struct {
	Id      string    `json:"id,omitempty"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	Url     string    `json:"url"`
	Created time.Time `json:"created,omitempty"`
}

type OrganizationNotifications struct {
	NewIssuesRemediations NewIssuesRemediationsOption `json:"new-issues-remediations"`
	ProjectImported       ProjectImportedOption       `json:"project-imported"`
	TestLimit             TestLimitOption             `json:"test-limit"`
	WeeklyReport          WeeklyReportOption          `json:"weekly-report"`
}

type NewIssuesRemediationsOption struct {
	Enabled       bool   `json:"enabled"`
	IssueSeverity string `json:"issueSeverity"`
	IssueType     string `json:"issueType"`
}

type ProjectImportedOption struct {
	Enabled bool `json:"enabled"`
}

type TestLimitOption struct {
	Enabled bool `json:"enabled"`
}

type WeeklyReportOption struct {
	Enabled bool `json:"enabled"`
}

type organizationCreateRequest struct {
	Name    string `json:"name"`
	GroupId string `json:"groupId"`
}

func GetOrganization(so SnykOptions, id string) (Organization, error) {
	path := fmt.Sprintf("/group/%s/orgs", so.GroupId)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return Organization{}, err
	}

	defer res.Body.Close()

	group := map[string]json.RawMessage{}

	err = json.NewDecoder(res.Body).Decode(&group)

	if err != nil {
		return Organization{}, err
	}

	var orgs []Organization
	json.Unmarshal(group["orgs"], &orgs)

	for _, element := range orgs {
		if element.Id == id {
			return element, nil
		}
	}

	return Organization{}, errors.New("Organization not found")

}

func GetOrgNotificationSettings(so SnykOptions, id string) (OrganizationNotifications, error) {
	path := fmt.Sprintf("/org/%s/notification-settings", id)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return OrganizationNotifications{}, err
	}

	defer res.Body.Close()

	var notifications OrganizationNotifications

	err = json.NewDecoder(res.Body).Decode(&notifications)

	if err != nil {
		return OrganizationNotifications{}, err
	}

	return notifications, nil
}

func SetOrgNotificationSettings(so SnykOptions, id string, nots OrganizationNotifications) (OrganizationNotifications, error) {
	path := fmt.Sprintf("/org/%s/notification-settings", id)

	body, _ := json.Marshal(nots)

	res, err := clientDo(so, "PUT", path, body)

	if err != nil {
		return OrganizationNotifications{}, err
	}

	defer res.Body.Close()

	var notifications OrganizationNotifications

	err = json.NewDecoder(res.Body).Decode(&notifications)

	if err != nil {
		return OrganizationNotifications{}, err
	}

	return notifications, nil
}

func CreateOrganization(so SnykOptions, name string) (Organization, error) {
	path := "/org"

	newOrg := organizationCreateRequest{
		Name:    name,
		GroupId: so.GroupId,
	}

	body, _ := json.Marshal(newOrg)

	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return Organization{}, err
	}

	defer res.Body.Close()

	var org Organization
	err = json.NewDecoder(res.Body).Decode(&org)

	if err != nil {
		return Organization{}, err
	}

	return org, nil
}

func DeleteOrganization(so SnykOptions, id string) error {
	path := fmt.Sprintf("/org/%s", id)

	_, err := clientDo(so, "DELETE", path, nil)

	return err
}
