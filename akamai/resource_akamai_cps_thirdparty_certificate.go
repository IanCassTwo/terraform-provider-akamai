package akamai

import (
	"fmt"
	"log"
	"time"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"

        "github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSThirdPartyCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCPSThirdPartyCertificateCreate,
		Read:   resourceCPSThirdPartyCertificateRead,
		Delete: resourceCPSThirdPartyCertificateDelete,

		Schema: map[string]*schema.Schema{
			"certificateid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trustchain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceCPSThirdPartyCertificateCreate(d *schema.ResourceData, meta interface{}) error {
        log.Print("DEBUG: enter resourceCPSThirdPartyCertificateCreate")

        var enrollment cps.Enrollment
        enrollment.Location = cps.GetLocation(getEnrollmentIdFromCertificateId(d))
        err := enrollment.GetEnrollment()
        if err != nil {
                return err
        }

        currentstatus, err := enrollment.GetChangeStatus()
	if err != nil {
		return err
	}

	if currentstatus == nil {
		d.SetId("none")
		return nil
	}

	if currentstatus.StatusInfo.Status != "wait-upload-third-party" {
		d.SetId("none")
		return nil
	}

	if currentstatus.StatusInfo.State != "awaiting-input" {
		d.SetId("none")
		return nil
	}

	certificate := d.Get("certificate").(string)
	trustchain := d.Get("trustchain").(string)

	var thirdpartycert cps.ThirdPartyCert
	thirdpartycert.Certificate = certificate
	thirdpartycert.TrustChain = trustchain

	_, err = enrollment.SubmitThirdPartyCert(thirdpartycert)
	if err != nil {
		return fmt.Errorf("Error : %s", err)
	}
	return nil
}

func resourceCPSThirdPartyCertificateRead(d *schema.ResourceData, meta interface{}) error {
        log.Print("DEBUG: enter resourceCPSThirdPartyCertificateRead")

        var enrollment cps.Enrollment
        enrollment.Location = cps.GetLocation(getEnrollmentIdFromCertificateId(d))
        err := enrollment.GetEnrollment()
        if err != nil {
                return err
        }

	currentstatus, err := enrollment.GetChangeStatus()
	if err != nil {
		return err
	}
	
	if currentstatus != nil {
		d.Set("currentstatus", currentstatus.StatusInfo.Status)
		d.SetId(getEnrollmentIdFromCertificateId(d))
	} else {
		d.Set("currentstatus", "No outstanding changes")
		d.SetId("")
	}

	return nil
}

func resourceCPSThirdPartyCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
