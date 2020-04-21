package akamai

import (
        "log"
	"fmt"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
        "github.com/hashicorp/terraform/helper/schema"
//        "github.com/hashicorp/terraform/helper/validation"
)

func resourceCPSThirdPartyEnrollment() *schema.Resource {
        return &schema.Resource {
                Create: resourceCPSThirdPartyEnrollmentCreate,
                Read:   resourceCPSThirdPartyEnrollmentRead,
                Update: resourceCPSThirdPartyEnrollmentUpdate,
                Delete: resourceCPSThirdPartyEnrollmentDelete,
//                Exists: resourceCPSThirdPartyEnrollmentExists,
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },
                Schema: map[string]*schema.Schema{
			// TODO all validate functions
                        "admincontact": &schema.Schema {
                                Type: schema.TypeSet,
                                Required:     true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource {
					Schema: map[string]*schema.Schema {
						"firstname": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"lastname": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"title": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"organization": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"email": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"phone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslineone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslinetwo": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"city": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"region": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"postalcode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"countrycode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
					},
				},
                        },
                        "techcontact": &schema.Schema {
                                Type: schema.TypeSet,
                                Required:     true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource {
					Schema: map[string]*schema.Schema {
						"firstname": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"lastname": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"title": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"organization": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"email": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"phone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslineone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslinetwo": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"city": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"region": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"postalcode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"countrycode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
					},
				},
                        },
                        "organization": &schema.Schema {
                                Type: schema.TypeSet,
                                Required:     true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource {
					Schema: map[string]*schema.Schema {
						"name": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"phone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslineone": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"addresslinetwo": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"city": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"region": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"postalcode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
						"countrycode": &schema.Schema {
							Type: schema.TypeString,
							Required:     true,
						},
					},
				},
                        },
			"contract": {
                                Type:         schema.TypeString,
                                Required:     true,
                                ForceNew:     true,
                        },

			"securenetwork": {
                                Type:         schema.TypeString,
                                Required:     true,
                                ForceNew:     true,
                        },
			"snionly": {
                                Type:         schema.TypeBool,
                                Required:     true,
                                ForceNew:     true,
                        },
			"commonname": {
                                Type:         schema.TypeString,
                                Required:     true,
                                ForceNew:     true,
                        },
			"alternativenames": {
                                Type:         schema.TypeSet,
				Elem: &schema.Schema {
					Type: schema.TypeString,
				},
                                Optional:     true,
                        },
			"certificatetype": {
                                Type:         schema.TypeString,
                                Computed:     true,
                        },

			"validationtype": {
                                Type:         schema.TypeString,
                                Computed:     true,
                        },
                        "csr": {
                                Type:         schema.TypeString,
                                Computed:     true,
                        },
                },
        }
}

func resourceCPSThirdPartyEnrollmentCreate(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] enter resourceCPSThirdPartyEnrollmentCreate")

	d.Partial(true)

	var enrollment cps.Enrollment

	setAdminContact(d, &enrollment)
	setTechContact(d, &enrollment)
	setSanCertType(d, &enrollment)
	setThirdPartyValidationType(d, &enrollment)
	setNetworkConfig(d, &enrollment)
	setSignatureAuthority(d, &enrollment)
	setChangeManagement(d, &enrollment)
	setCSR(d, &enrollment)
	setOrganization(d, &enrollment)
	setThirdParty(d, &enrollment)

	var queryparams cps.CreateEnrollmentQueryParams
	queryparams.ContractID = d.Get("contract").(string)

	enrollmentresponse, err := enrollment.Create(queryparams)
	if err != nil {
		d.SetId("")
		return err
	}

	// no error so far, so set an Id & update partials
	enrollmentid := getEnrollmentIdFromLocation(enrollmentresponse.Location)
	setId(d, enrollmentid)
	d.SetPartial("commonname")
	d.SetPartial("organization")
	d.SetPartial("admincontact")
	d.SetPartial("techcontact")
	d.SetPartial("snionly")
	d.SetPartial("certificatetype")
	d.SetPartial("validationtype")
	d.SetPartial("securenetwork")
	d.SetPartial("alternativenames")
	d.SetPartial("contract")

        // Wait for validation to complete
	err = awaitCertVerification(enrollment)
        if err != nil {
                return err
        }

	d.Partial(false)
	
	return resourceCPSThirdPartyEnrollmentRead(d, meta)
}

func resourceCPSThirdPartyEnrollmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceCPSThirdPartyEnrollmentUpdate")
	d.Partial(true)

	var enrollment cps.Enrollment
	enrollment.Location = cps.GetLocation(getEnrollmentIdFromId(d))

	setAdminContact(d, &enrollment)
	setTechContact(d, &enrollment)
	setSanCertType(d, &enrollment)
	setThirdPartyValidationType(d, &enrollment)
	setNetworkConfig(d, &enrollment)
	setSignatureAuthority(d, &enrollment)
	setChangeManagement(d, &enrollment)
	setCSR(d, &enrollment)
	setOrganization(d, &enrollment)
	setThirdParty(d, &enrollment)

	enrollmentresponse, err := enrollment.Update()
	if err != nil {
		d.SetId("")
		return err
	}

	enrollmentid := getEnrollmentIdFromLocation(enrollmentresponse.Location)
	setId(d, enrollmentid)
	d.SetPartial("commonname")
	d.SetPartial("organization")
	d.SetPartial("admincontact")
	d.SetPartial("techcontact")
	d.SetPartial("snionly")
	d.SetPartial("certificatetype")
	d.SetPartial("validationtype")
	d.SetPartial("securenetwork")
	d.SetPartial("alternativenames")
	d.SetPartial("contract")

        // Wait for validation to complete
	err = awaitCertVerification(enrollment)
        if err != nil {
                return err
        }

	d.Partial(false)
	
	return resourceCPSThirdPartyEnrollmentRead(d, meta)
}

func resourceCPSThirdPartyEnrollmentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceCPSThirdPartyEnrollmentDelete")
	var enrollment cps.Enrollment
	enrollment.Location = cps.GetLocation(getEnrollmentIdFromId(d))
	_, err := enrollment.Delete()
	if err != nil {
		d.SetId("")
		return err
	}

	//TODO - get status & poll for actual deletion
	return nil
}

func resourceCPSThirdPartyEnrollmentRead(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] enter resourceCPSThirdPartyEnrollmentRead")
	var enrollment cps.Enrollment
	enrollment.Location = cps.GetLocation(getEnrollmentIdFromId(d))
	err := enrollment.GetEnrollment()
	if err != nil {
		d.SetId("")
		apierror := err.(client.APIError)
		if apierror.Status == 404 {
			// enrollment not found
			return nil
		}
		return err
	}

	getAdminContact(d, &enrollment)
	getTechContact(d, &enrollment)
	getNetworkConfig(d, &enrollment)
	getCSR(d, &enrollment)
	getOrganization(d, &enrollment)
	getCertType(d, &enrollment)
	getValidationType(d, &enrollment)
	d.Set("contract", getContractIdFromId(d))

	err = validateThirdPartyImport(d)
	if err != nil {
		d.SetId("")
		return err
	}

        // See if there are any pending changes
	currentstatus, err := enrollment.GetChangeStatus()
	if err != nil {
		return err
	}

	// Override with the real thing if they are present
	if currentstatus != nil {
		// Retrieve the CSR
		if len(currentstatus.AllowedInput) != 0 {
			if currentstatus.AllowedInput[0].Type == "third-party-csr" {
				thirdpartycsr, _ := enrollment.GetThirdPartyCSR()
				if thirdpartycsr != nil {
					getThirdPartyCSR(d, *thirdpartycsr)
				}
			}
		}
	}
	
	return nil
}

func validateThirdPartyImport(d *schema.ResourceData) error {
	v := d.Get("validationtype")
	if v != "third-party" {
		d.SetId("")
		return fmt.Errorf("Error reading certificate, validation type must be dv for this provider, not %s", d.Get("validationtype"))
	}

	c := d.Get("certificatetype")
	if c != "third-party" {
		d.SetId("")
		return fmt.Errorf("Error reading certificate, certifivate type must be san for this provider not %s", d.Get("certificatetype"))
	}
	return nil
}

func getThirdPartyCSR(d *schema.ResourceData, csr cps.ThirdPartyCSR) {
        log.Print("[DEBUG] enter getThirdPartyCSR")
        d.Set("csr", csr.Csr)
}

