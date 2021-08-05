package snyk

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lendi-au/terraform-provider-snyk/snyk/api"
)

func TestAccIntegration(t *testing.T) {
	var integration = new(api.Integration)

	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/helper/acctest
	rOrgName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	intType := "bitbucket-cloud"
	username := "test_user"
	password := "test_pass"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccIntegration(rOrgName, intType, username, password),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckIntegrationExists("snyk_integration.integ_test_integ", integration),
					// verify remote values
					testAccCheckIntegrationValues(integration, intType),
					// verify local values
					resource.TestCheckResourceAttr("snyk_integration.integ_test_integ", "type", intType),
				),
			},
		},
	})
}

// TODO: environment variables so you can use legit credentials
func testAccIntegration(name string, intType string, username string, password string) string {
	return fmt.Sprintf(`
	resource "snyk_organization" "integ_test_org" {
		name = "%s"
	}
	
	resource "snyk_integration" "integ_test_integ" {
		organization = snyk_organization.integ_test_org.id
		type = "%s"
		credentials {
			username = "%s"
			password = "%s"
		}
	}
	`, name, intType, username, password)
}

func testAccCheckIntegrationValues(i *api.Integration, intType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if i.Type != intType {
			return fmt.Errorf("bad integration type, expected \"%s\", got: %#v", intType, i.Type)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckIntegrationExists(n string, i *api.Integration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the client options from the test setup
		so := testAccProviders["snyk"].Meta().(api.SnykOptions)
		intType := rs.Primary.Attributes["type"]
		orgId := rs.Primary.Attributes["organization"]

		res, err := api.GetIntegration(so, orgId, intType)

		if err != nil {
			return err
		}

		// assign the response Widget attribute to the widget pointer
		*i = *res

		return nil
	}
}
