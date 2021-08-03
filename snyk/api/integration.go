package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Integration struct {
	Id             string                 `json:"id,omitempty"`
	OrganizationId string                 `json:"-"`
	Type           string                 `json:"type"`
	Credentials    IntegrationCredentials `json:"credentials"`
}

type IntegrationCredentials struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RegistryBase string `json:"registryBase,omitempty"`
	Url          string `json:"url,omitempty"`
	Token        string `json:"token,omitempty"`
	Region       string `json:"region,omitempty"`
	RoleArn      string `json:"roleArn,omitempty"`
}

func CreateIntegration(so SnykOptions, i Integration) (Integration, error) {
	// Send request to url to create integration of type, fill out ID after creation
	path := fmt.Sprintf("/org/%s/integrations", i.OrganizationId)

	body, _ := json.Marshal(i)

	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return Integration{}, err
	}

	defer res.Body.Close()

	var newInt Integration
	err = json.NewDecoder(res.Body).Decode(&newInt)

	if err != nil {
		return Integration{}, err
	}

	i.Id = newInt.Id

	return i, nil
}

func GetIntegration(so SnykOptions, i Integration) (Integration, error) {
	path := fmt.Sprintf("/org/%s/integrations", i.OrganizationId)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return Integration{}, err
	}

	defer res.Body.Close()

	var listing map[string]string
	err = json.NewDecoder(res.Body).Decode(&listing)

	if err != nil {
		return Integration{}, err
	}

	exists := listing[i.Type] != ""

	if !exists {
		return i, errors.New("integration not found")
	}

	i.Id = listing[i.Type]

	return i, nil
}

func UpdateIntegration(so SnykOptions, i Integration) (Integration, error) {
	id := i.Id
	//Unsetting so JSON won't parse it
	i.Id = ""
	path := fmt.Sprintf("/org/%s/integrations/%s", i.OrganizationId, id)

	body, _ := json.Marshal(i)

	_, err := clientDo(so, "PUT", path, body)

	if err != nil {
		return Integration{}, err
	}

	i.Id = id

	return i, nil
}

func DeleteIntegration(so SnykOptions, i Integration) error {
	path := fmt.Sprintf("/org/%s/integrations/%s/authentication", i.OrganizationId, i.Id)

	_, err := clientDo(so, "DELETE", path, nil)

	return err
}
