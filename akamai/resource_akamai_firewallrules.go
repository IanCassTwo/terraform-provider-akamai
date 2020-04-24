package akamai

import (
	"fmt"
	"log"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/firewallrules-v1"

        "github.com/hashicorp/terraform/helper/schema"
//        "github.com/hashicorp/terraform/helper/resource"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Delete: resourceFirewallRuleDelete,

		Schema: map[string]*schema.Schema{
			"serviceid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"servicename": {
				Type:     schema.TypeString,
				Computed: true,
			},
                        "cidrblocks": {
                                Type:         schema.TypeSet,
                                Computed:     true,
                                Elem: &schema.Schema {
                                        Type: schema.TypeString,
                                },
                        },
		},
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },
	}
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
        log.Print("[DEBUG] enter resourceFirewallRuleCreate")

	serviceid := d.Get("serviceid").(int)
	email := d.Get("email").(string)

	// Get existing subscriptions
	listsubscriptionsresponse, err := firewallrules.ListSubscriptions()
	if err != nil {
		return err
	}

	// Rebuild subscriptions without this one
	var subscriptions = make([]firewallrules.Subscriptions, 0)
	for _, s := range listsubscriptionsresponse.Subscriptions {
		if (s.ServiceID == serviceid && s.Email == email) {
			continue
		}
		subscriptions = append(subscriptions, s)
	}
	
	// Create update request
	var updatesubscriptionsrequest firewallrules.UpdateSubscriptionsRequest
	updatesubscriptionsrequest.Subscriptions = subscriptions

	// Update
	_, err = firewallrules.UpdateSubscriptions(updatesubscriptionsrequest)
	if err != nil {
		return err
	}

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
        log.Print("[DEBUG] enter resourceFirewallRuleCreate")

	serviceid := d.Get("serviceid").(int)
	email := d.Get("email").(string)

	// Get existing subscriptions
	listsubscriptionsresponse, err := firewallrules.ListSubscriptions()
	if err != nil {
		return err
	}

	// Create new subscription
	var newsubscription firewallrules.Subscriptions
	newsubscription.ServiceID = serviceid
	newsubscription.Email = email

	// Add to existing subscriptions
	subscriptions := listsubscriptionsresponse.Subscriptions
	subscriptions = append(subscriptions, newsubscription)

	// Create update request
	var updatesubscriptionsrequest firewallrules.UpdateSubscriptionsRequest
	updatesubscriptionsrequest.Subscriptions = subscriptions

	// Update
	_, err = firewallrules.UpdateSubscriptions(updatesubscriptionsrequest)
	if err != nil {
		return err
	}

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
        log.Print("[DEBUG] enter resourceFirewallRuleRead")

	serviceid := d.Get("serviceid").(int)
	email := d.Get("email").(string)
	
	subscriptions, err := firewallrules.ListSubscriptions()
	if err != nil {
		return err
	}

	for _, s := range subscriptions.Subscriptions {
		if (s.ServiceID == serviceid && s.Email == email) {
			// Found a subscription to this service
			d.Set("servicename", s.ServiceName)
			d.SetId(fmt.Sprintf("%d:%s", serviceid, email))
			getCidrs(d)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func getCidrs(d  *schema.ResourceData) error {
	log.Print("[DEBUG] enter getCidrs")
        serviceid := d.Get("serviceid").(int)

	cidrs := make([]string, 0)

	cidrblocks, err := firewallrules.ListCidrBlocks()
	if err != nil {
		return err
	}

	for _, s := range *cidrblocks {
		if s.ServiceID == serviceid {
			// Found CIDR block
			cidrs = append(cidrs, fmt.Sprintf("%s%s", s.Cidr, s.CidrMask))
		}
	}

	d.Set("cidrblocks", cidrs)

	return nil
	
}

