package api

import (
	"encoding/json"
	"fmt"
)

type Integration struct {
	Type        string                 `json:"type"`
	Credentials IntegrationCredentials `json:"credentials"`
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

func CreateIntegration(so SnykOptions, orgId string, intType string, creds IntegrationCredentials) (string, error) {
	// Send request to url to create integration of type, fill out ID after creation
	path := fmt.Sprintf("/org/%s/integrations", orgId)

	i := Integration{
		Type:        intType,
		Credentials: creds,
	}

	body, _ := json.Marshal(i)

	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var newInt map[string]string
	err = json.NewDecoder(res.Body).Decode(&newInt)

	if err != nil {
		return "", err
	}

	return newInt["id"], nil
}

func GetIntegrationDetails(so SnykOptions, orgId string, integrationId string) (string, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s", orgId, integrationId)

	_, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return "", err
	}

	return "", nil
}

func GetIntegrationByType(so SnykOptions, orgId string, intType string) (string, error) {
	path := fmt.Sprintf("/org/%s/integrations", orgId)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var listing map[string]string
	err = json.NewDecoder(res.Body).Decode(&listing)

	if err != nil {
		return "", err
	}

	return listing[intType], nil
}

func IntegrationExists(so SnykOptions, org string, intType string) (bool, error) {
	path := fmt.Sprintf("/org/%s/integrations", org)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	var listing map[string]string
	err = json.NewDecoder(res.Body).Decode(&listing)

	if err != nil {
		return false, err
	}

	exists := listing[intType] != ""

	return exists, nil
}

func UpdateIntegration(so SnykOptions, orgId string, id string, intType string, creds IntegrationCredentials) (string, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s", orgId, id)

	i := Integration{
		Type:        intType,
		Credentials: creds,
	}

	body, _ := json.Marshal(i)

	_, err := clientDo(so, "PUT", path, body)

	if err != nil {
		return "", err
	}

	return id, nil
}

func DeleteIntegration(so SnykOptions, orgId string, id string) error {
	path := fmt.Sprintf("/org/%s/integrations/%s/authentication", orgId, id)

	_, err := clientDo(so, "DELETE", path, nil)

	return err
}
