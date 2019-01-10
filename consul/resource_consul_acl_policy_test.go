package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLPolicyDestroy(s *terraform.State) error {
	client, err := getMasterClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl" {
			continue
		}
		secret, _, err := client.ACL().Info(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if secret != nil {
			return fmt.Errorf("ACL %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLPolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_policy.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "rules", "node_prefix \"\" { policy = \"read\" }"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "datacenters.#", "1"),
				),
			},
		},
	})
}

const testResourceACLPolicyConfigBasic = testAccMasterProviderConfiguration + `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node_prefix \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}`
