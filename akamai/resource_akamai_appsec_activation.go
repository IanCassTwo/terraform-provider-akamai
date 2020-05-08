package akamai

import (
	"fmt"
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAppsecActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppsecActivationCreate,
		Read:   resourceAppsecActivationRead,
		Update: resourceAppsecActivationUpdate,
		Delete: resourceAppsecActivationDelete,
		Schema: map[string]*schema.Schema{
			"configid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STAGING",
				ValidateFunc: validation.StringInSlice([]string{
					"STAGING",
					"PRODUCTION",
				}, false),
			},
			"activate": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"contact": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppsecActivationCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	configid := d.Get("configid").(int)

	versionlist, err := appsec.ListConfigurationVersions(configid)
	if err != nil {
		return err
	}

	version := d.Get("version").(int)
	if version == 0 {
        	version = versionlist.LastCreatedVersion
	}

	network := d.Get("network").(string)
	activate := d.Get("activate").(bool)

	// Check if we even have to do anything
	if activate {
		if network == "STAGING" && version == versionlist.StagingActiveVersion {
			d.Set("status", "ACTIVE")
			return nil
		}

		if network == "PRODUCTION" && version == versionlist.ProductionActiveVersion {
			d.Set("status", "ACTIVE")
			return nil
		}
	} else {
		if network == "STAGING" && version != versionlist.StagingActiveVersion {
			d.Set("status", "INACTIVE")
			return nil
		}

		if network == "PRODUCTION" && version != versionlist.ProductionActiveVersion {
			d.Set("status", "INACTIVE")
			return nil
		}
	}

	// Prepare our activation object
	aset := d.Get("contact").(*schema.Set)
	alist := aset.List()
	contactlist := make([]string, len(alist))
	for i, v := range alist {
		contactlist[i] = v.(string)
	}
	
	var activation appsec.Activation

	if activate {
		activation.Action = "ACTIVATE"
	} else {
		activation.Action = "DEACTIVATE"
	}

	activation.Network = network
	activation.Note = "Activated by Terraform"
	activation.NotificationEmails = contactlist

	activationconfigs := make([]appsec.ActivationConfig, 1)
	var activationconfig appsec.ActivationConfig
	activationconfig.ConfigID = configid
	activationconfig.ConfigVersion = version
	activationconfigs[0] = activationconfig

	activation.ActivationConfigs = activationconfigs

	// Do the deed!
	activationresponse, err := appsec.ActivateConfigurationVersion(activation)
	if err != nil {
		return err
	}

	// Pesky API returns different objects depending on status code!
	if activationresponse.ResponseCode != 200 {
		// TODO. I don't think this will happen if we're only submitting a single config version
		return nil
	}

	// Check activation status 
	activationstatus := activationresponse.ActivationStatus
	statusid := activationstatus.ActivationID

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		activationstatus, err := appsec.GetActivationStatus(statusid)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error : %s", err))
		}

		if activationstatus.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("FAILED: The network list has failed to activate."))
		}

		// Final states
		if activationstatus.Status == "ACTIVATED" || activationstatus.Status == "LIVE" || activationstatus.Status == "DEPLOYED" || activationstatus.Status == "REMOVED" || activationstatus.Status == "UNDEPLOYED" {
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

	d.Partial(false)
	return resourceAppsecActivationRead(d, meta)
}

func resourceAppsecActivationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DEACTIVATE PROPERTY")
	//FIXME deactivate!!
	d.SetId("")
	return nil
}

func resourceAppsecActivationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] enter resourceAppsecActivationRead")

	configid := d.Get("configid").(int)

	versionlist, err := appsec.ListConfigurationVersions(configid)
	if err != nil {
		return err
	}

	network := d.Get("network").(string)
	if network == "STAGING" {
		d.Set("version", versionlist.StagingActiveVersion)
		d.SetId(fmt.Sprintf("%d:%d", configid, versionlist.StagingActiveVersion))
	} else {
		d.Set("version", versionlist.ProductionActiveVersion)
		d.SetId(fmt.Sprintf("%d:%d", configid, versionlist.ProductionActiveVersion))
	}
	
	return nil
}

func resourceAppsecActivationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entered resourceAppsecActivationUpdate")
	return resourceAppsecActivationCreate(d, meta)
}

