package api

import (
	"encoding/json"
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

type organizationCreateRequest struct {
	Name    string `json:"name"`
	GroupId string `json:"groupId"`
}

func GetOrganization(so SnykOptions, id string) (*Organization, error) {
	path := fmt.Sprintf("/group/%s/orgs", so.GroupId)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	group := map[string]json.RawMessage{}

	err = json.NewDecoder(res.Body).Decode(&group)

	if err != nil {
		return nil, err
	}

	var orgs []Organization
	json.Unmarshal(group["orgs"], &orgs)

	for _, element := range orgs {
		if element.Id == id {
			return &element, nil
		}
	}

	return nil, ErrNotFound
}

func OrganizationExistsByName(so SnykOptions, name string) (bool, error) {
	path := fmt.Sprintf("/group/%s/orgs", so.GroupId)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	group := map[string]json.RawMessage{}

	err = json.NewDecoder(res.Body).Decode(&group)

	if err != nil {
		return false, err
	}

	var orgs []Organization
	json.Unmarshal(group["orgs"], &orgs)

	for _, element := range orgs {
		if element.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func CreateOrganization(so SnykOptions, name string) (*Organization, error) {
	path := "/org"

	newOrg := organizationCreateRequest{
		Name:    name,
		GroupId: so.GroupId,
	}

	body, _ := json.Marshal(newOrg)

	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var org = new(Organization)
	err = json.NewDecoder(res.Body).Decode(org)

	if err != nil {
		return nil, err
	}

	return org, nil
}

func DeleteOrganization(so SnykOptions, id string) error {
	path := fmt.Sprintf("/org/%s", id)

	_, err := clientDo(so, "DELETE", path, nil)

	return err
}
