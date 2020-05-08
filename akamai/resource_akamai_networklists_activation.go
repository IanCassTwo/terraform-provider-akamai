package akamai

import (
	"fmt"
	"log"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/networklists-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/helper/resource"
)

func resourceNetworkListActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkListActivationCreate,
		Read:   resourceNetworkListActivationRead,
		Delete: resourceNetworkListActivationDelete,
		Schema: map[string]*schema.Schema{
			"networklistid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STAGING",
					"PRODUCTION",
				}, false),
			},
			"contacts": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceNetworkListActivationCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	networklistid := d.Get("networklistid").(string)
	network := networklists.Staging
	if d.Get("network") == "PRODUCTION" {
		network = networklists.Production
	}

	aset := d.Get("contacts").(*schema.Set)
	alist := aset.List()
	contacts := make([]string, len(alist))
	for i, v := range alist {
		contacts[i] = v.(string)
	}

	var activationrequest networklists.ActivationRequest
	activationrequest.UniqueID = networklistid
	activationrequest.Network = network
	activationrequest.Comments = "activated by Terraform"
	activationrequest.NotificationRecipients = contacts

	// Check activation status before we try to activate
        activationstatus, err := networklists.GetActivationStatus(networklistid, network)
        if err != nil {
                return err
        }

	// Only try to activate again if it's not already active
	// We can get here if we've just imported a network list or if we interrupted a previous 
	// activation before it finished
	if activationstatus.ActivationStatus != "ACTIVE" {
		activationstatus, err = networklists.ActivateNetworkList(activationrequest)
		if err != nil {
			return err
		}
		
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

			activationstatus, err := networklists.GetActivationStatus(networklistid, network)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error : %s", err))
			}

			if activationstatus.ActivationStatus == "MODIFIED" {
				return resource.NonRetryableError(fmt.Errorf("MODIFIED: Indicates that a previous syncPoint version of the network list is currently active, and any subsequent modifications may need to be activated."))
			}

			if activationstatus.ActivationStatus == "PENDING_DEACTIVATION" {
				return resource.NonRetryableError(fmt.Errorf("PENDING_DEACTIVATION: An activation for another syncPoint version of the network list has launched, but it has not yet fully rendered this version INACTIVE."))
			}

			if activationstatus.ActivationStatus == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("FAILED: The network list has failed to activate."))
			}

			if activationstatus.ActivationStatus == "ACTIVE" {
				d.SetId(fmt.Sprintf("%d", activationstatus.ActivationID))
				return resource.NonRetryableError(nil)
			}
			return resource.RetryableError(fmt.Errorf("Awaiting activation"))

		})

		if isResourceTimeoutError(err) {
			return fmt.Errorf("Timed out waiting for activation")
		}

		if err != nil {
			return fmt.Errorf("Error : %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%d", activationstatus.ActivationID))
	d.Set("status", activationstatus.ActivationStatus)

	d.Partial(false)
	return nil
}

func resourceNetworkListActivationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DEACTIVATE PROPERTY")
	d.SetId("")
	return nil
}

func resourceNetworkListActivationRead(d *schema.ResourceData, meta interface{}) error {

	networklistid := d.Get("networklistid").(string)
	network := networklists.Staging
	if d.Get("network") == "PRODUCTION" {
		network = networklists.Production
	}

        activationstatus, err := networklists.GetActivationStatus(networklistid, network)
	if err != nil {
		return err
	}
	
	d.SetId(fmt.Sprintf("%d", activationstatus.ActivationID))
	d.Set("status", activationstatus.ActivationStatus)

	return nil
}
