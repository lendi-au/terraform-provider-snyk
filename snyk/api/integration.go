package api

import (
	"encoding/json"
	"fmt"
)

type Integration struct {
	Id          string                 `json:"id,omitempty"`
	OrgId       string                 `json:"-"`
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

func CreateIntegration(so SnykOptions, orgId string, intType string, creds IntegrationCredentials) (*Integration, error) {
	path := fmt.Sprintf("/org/%s/integrations", orgId)

	i := Integration{
		Type:        intType,
		Credentials: creds,
	}

	body, _ := json.Marshal(i)

	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var newInt map[string]string
	err = json.NewDecoder(res.Body).Decode(&newInt)

	if err != nil {
		return nil, err
	}

	returnData := &Integration{
		Id:          newInt["id"],
		OrgId:       orgId,
		Type:        intType,
		Credentials: creds,
	}

	return returnData, nil
}

func GetIntegration(so SnykOptions, orgId string, intType string) (*Integration, error) {
	id, err := getIntegrationIdByType(so, orgId, intType)

	if err != nil {
		return nil, err
	}

	return &Integration{
		Id:    id,
		OrgId: orgId,
		Type:  intType,
	}, nil
}

func getIntegrationIdByType(so SnykOptions, orgId string, intType string) (string, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s", orgId, intType)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var data map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return "", err
	}

	return data["id"], nil
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

func UpdateIntegration(so SnykOptions, orgId string, intType string, creds IntegrationCredentials) (*Integration, error) {

	id, err := getIntegrationIdByType(so, orgId, intType)

	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/org/%s/integrations/%s", orgId, id)

	patchData := &Integration{
		Type:        intType,
		Credentials: creds,
	}

	body, _ := json.Marshal(patchData)

	_, err = clientDo(so, "PUT", path, body)

	if err != nil {
		return nil, err
	}

	returnData := &Integration{
		Id:          id,
		OrgId:       orgId,
		Type:        intType,
		Credentials: creds,
	}

	return returnData, nil
}

func DeleteIntegration(so SnykOptions, orgId string, intType string) error {

	id, err := getIntegrationIdByType(so, orgId, intType)

	if err != nil {
		return err
	}

	path := fmt.Sprintf("/org/%s/integrations/%s/authentication", orgId, id)

	_, err = clientDo(so, "DELETE", path, nil)

	return err
}
