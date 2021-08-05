package snyk

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders = map[string]*schema.Provider{
	"snyk": Provider("test")(),
}

func testAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{
		"SNYK_API_GROUP",
		"SNYK_API_KEY",
	}

	for _, env := range requiredEnvVars {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("env variable %s required for acceptance tests", env)
		}
	}
}
