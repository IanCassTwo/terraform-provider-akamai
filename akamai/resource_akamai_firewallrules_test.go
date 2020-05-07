package akamai

import (
	"fmt"
	"log"
	"strings"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/firewallrules-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"strings"
	"testing"
)

var testAccAkamaiFirewallRulesConfig = fmt.Sprintf(`
provider "akamai" {
        firewallrules_section = "papi"
}

resource "akamai_firewallrule" "map" {
        serviceid = 9
        email = "test@akamai.com"
}
`)

func TestAccAkamaiFirewallRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiFirewallRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiFirewallRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiFirewallRulesExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiFirewallRulesDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_firewallrule" {
			continue
		}
		log.Printf("[DEBUG] [Akamai FirewallRules] Delete for edgehostname not supported  [%v]", rs.Primary.ID)
	}
	return nil
}

func parseFirewallID(id string) (int, string, error) {
	idComp := strings.Split(id, ":")
	if len(idComp) < 2 {
		return -1, "", errors.New("Invalid ID")
	}

	serviceid, err := strconv.Atoi(idComp[0])
	if err != nil {
		return -1, "", err
	}

	email := idComp[1]
	return serviceid, email, nil
}

func testAccCheckAkamaiFirewallRulesExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_edge_hostname" {
			continue
		}
		log.Printf("[DEBUG] [Akamai FirewallRules] Searching for firewall rule [%v]", rs.Primary.ID)

		serviceid, email, err := parseFirewallID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subscriptions, err := firewallrules.ListSubscriptions()
		if err != nil {
			return err
		}

		for _, s := range subscriptions.Subscriptions {
			if (s.ServiceID == serviceid && s.Email == email) {
				// Found a subscription to this service
				return nil
			}
		}
		return fmt.Errorf("error looking up firewall rule subscription")
	}
	return nil
}
