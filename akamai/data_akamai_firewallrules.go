package akamai

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceAkamaiFirewallRule() *schema.Resource {
	return &schema.Resource{
		Read: dataAkamaiFirewallRuleRead,
		Schema: map[string]*schema.Schema{
			"serviceid": {
				Type:     schema.TypeInt,
				Required: true,
			},
                        "email": {
                                Type:     schema.TypeString,
                                Required: true,
                                ForceNew: true,
                        },
                        "ipv4only": {
                                Type:     schema.TypeBool,
                                Default: false,
                                ForceNew: true,
                                Optional: true,
                        },
                        "servicename": {
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                        "cidrblocks": {
                                Type:     schema.TypeSet,
                                Computed: true,
                                Elem: &schema.Schema{
                                        Type: schema.TypeString,
                                },
                        },
		},
	}
}

func dataAkamaiFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entered dataAkamaiFirewallRuleRead")
	return resourceFirewallRuleRead(d, meta)
}
