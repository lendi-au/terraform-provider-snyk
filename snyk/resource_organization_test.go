package snyk

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lendi-au/terraform-provider-snyk/snyk/api"
)

var testAccProviders = map[string]*schema.Provider{
	"snyk": Provider("test")(),
}

func TestAccOrganization(t *testing.T) {
	var org = new(api.Organization)

	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrgDestroy(rName),
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccOrg(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckOrgExists("snyk_organization.org_test_org", org),
					// verify remote values
					testAccCheckOrgValues(org, rName),
					// verify local values
					resource.TestCheckResourceAttr("snyk_organization.org_test_org", "name", rName),
				),
			},
		},
	})
}

func testAccOrg(name string) string {
	return fmt.Sprintf(`
	resource "snyk_organization" "org_test_org" {
		name = "%s"
	}`, name)
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

func testAccCheckOrgDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the client options from the test setup
		so := testAccProviders["snyk"].Meta().(api.SnykOptions)

		exists, err := api.OrganizationExistsByName(so, name)

		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("organization %s still exists", name)
		}

		return nil
	}
}

func testAccCheckOrgValues(org *api.Organization, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if org.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, org.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckOrgExists(n string, org *api.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the client options from the test setup
		so := testAccProviders["snyk"].Meta().(api.SnykOptions)
		orgId := rs.Primary.ID

		res, err := api.GetOrganization(so, orgId)

		if err != nil {
			return err
		}

		// assign the response Widget attribute to the widget pointer
		*org = *res

		return nil
	}
}
