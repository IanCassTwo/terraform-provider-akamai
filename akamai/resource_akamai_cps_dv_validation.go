package akamai

import (
	"fmt"
	"log"
	"time"
	"encoding/json"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"

        "github.com/hashicorp/terraform/helper/schema"
        "github.com/hashicorp/terraform/helper/resource"
)

func resourceCPSDVValidation() *schema.Resource {
	return &schema.Resource{
		Create: resourceCPSDVValidationCreate,
		Read:   resourceCPSDVValidationRead,
		Delete: resourceCPSDVValidationDelete,

		Schema: map[string]*schema.Schema{
			"certificateid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
                        "alternativenames": {
                                Type:         schema.TypeSet,
                                Required:     true,
				ForceNew: true,
                                Elem: &schema.Schema {
                                        Type: schema.TypeString,
                                },
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

func resourceCPSDVValidationCreate(d *schema.ResourceData, meta interface{}) error {
        log.Print("DEBUG: enter resourceCPSDVValidationCreate")

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

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

        	currentstatus, err := enrollment.GetChangeStatus()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error : %s", err))
		}

                s,_ := json.MarshalIndent(currentstatus, "", "\t")
		log.Printf("[DEBUG] Status = ", string(s))

		if currentstatus.StatusInfo.State != "running" {
			_, err := enrollment.AcknowledgeDVChallenges()
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error : %s", err))
			}
		}

		// TODO list failure states & return an error
		if currentstatus.StatusInfo.Status != "coodinate-domain-validation" {
			return resource.NonRetryableError(resourceCPSDVValidationRead(d, meta))
		}
		return resource.RetryableError(fmt.Errorf("Awaiting validation"))

	})

	if isResourceTimeoutError(err) {
        	currentstatus, err := enrollment.GetChangeStatus()
		if err != nil {
			return err
		}

		return fmt.Errorf("Expected certificate to be issued but was in state %s", currentstatus.StatusInfo.Status)
	}

	if err != nil {
		return fmt.Errorf("Error : %s", err)
	}
	return nil
}

func resourceCPSDVValidationRead(d *schema.ResourceData, meta interface{}) error {
        log.Print("DEBUG: enter resourceCPSDVValidationRead")

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
	
	// FIXME check updateurl, use as Id
	if currentstatus != nil {
		d.Set("currentstatus", currentstatus.StatusInfo.Status)
		d.SetId(getEnrollmentIdFromCertificateId(d))
	} else {
		d.Set("currentstatus", "No outstanding changes")
		d.SetId("")
	}

	return nil
}

func resourceCPSDVValidationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
