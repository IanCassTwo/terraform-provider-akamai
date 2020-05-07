package akamai

import (
        "log"
	"fmt"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
        "github.com/hashicorp/terraform/helper/schema"
//        "github.com/hashicorp/terraform/helper/validation"
)

func resourceCPSDVEnrollment() *schema.Resource {
        return &schema.Resource {
                Create: resourceCPSDVEnrollmentCreate,
                Read:   resourceCPSDVEnrollmentRead,
                Update: resourceCPSDVEnrollmentUpdate,
                Delete: resourceCPSDVEnrollmentDelete,
//                Exists: resourceCPSDVEnrollmentCExists,
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
                        "redirectchallenges": {
                                Type:         schema.TypeMap,
                                Computed:     true,
                        },
                        "httpchallenges": {
                                Type:         schema.TypeMap,
                                Computed:     true,
                        },
                        "dnschallenges": {
                                Type:         schema.TypeMap,
                                Computed:     true,
                        },

                },
        }
}

func resourceCPSDVEnrollmentCreate(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] enter resourceCPSDVEnrollmentCreate")

	d.Partial(true)

	var enrollment cps.Enrollment

	setAdminContact(d, &enrollment)
	setTechContact(d, &enrollment)
	setSanCertType(d, &enrollment)
	setDVValidationType(d, &enrollment)
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
	
	return resourceCPSDVEnrollmentRead(d, meta)
}

func resourceCPSDVEnrollmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceCPSDVEnrollmentUpdate")
	d.Partial(true)

	var enrollment cps.Enrollment
	enrollment.Location = cps.GetLocation(getEnrollmentIdFromId(d))

	setAdminContact(d, &enrollment)
	setTechContact(d, &enrollment)
	setSanCertType(d, &enrollment)
	setDVValidationType(d, &enrollment)
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
	
	return resourceCPSDVEnrollmentRead(d, meta)
}

func resourceCPSDVEnrollmentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceCPSDVEnrollmentDelete")
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

func resourceCPSDVEnrollmentRead(d *schema.ResourceData, meta interface{}) error {

	log.Print("[DEBUG] enter resourceCPSDVEnrollmentRead")
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

	err = validateDVImport(d)
	if err != nil {
		d.SetId("")
		return err
	}

        // See if there are any pending changes
	currentstatus, err := enrollment.GetChangeStatus()
	if err != nil {
		return err
	}

	// Initially set the challenges to blank
	emptyDVChallenges(d)

	// Override with the real thing if they are present
	if currentstatus != nil {
		// Retrieve the challenges
		if len(currentstatus.AllowedInput) != 0 {
			if currentstatus.AllowedInput[0].Type == "lets-encrypt-challenges" {
				domainvalidations, _ := enrollment.GetDVChallenges()
				if domainvalidations != nil {
					getDVChallenges(d, *domainvalidations)
				}
			}
		}
	}
	
	return nil
}

func validateDVImport(d *schema.ResourceData) error {
	v := d.Get("validationtype")
	if v != "dv" {
		d.SetId("")
		return fmt.Errorf("Error reading certificate, validation type must be dv for this provider, not %s", d.Get("validationtype"))
	}

	c := d.Get("certificatetype")
	if c != "san" {
		d.SetId("")
		return fmt.Errorf("Error reading certificate, certifivate type must be san for this provider not %s", d.Get("certificatetype"))
	}
	return nil
}

func emptyDVChallenges(d *schema.ResourceData) {
	// Set the challenges to blank
	httpchallenges := make(map[string]interface{})
	redirectchallenges := make(map[string]interface{})
	dnschallenges := make(map[string]interface{})
	d.Set("httpchallenges", httpchallenges)
	d.Set("redirectchallenges", redirectchallenges)
	d.Set("dnschallenges", dnschallenges)
}

func getDVChallenges(d *schema.ResourceData, domainvalidations cps.DomainValidations) {
        log.Print("[DEBUG] enter getDVChallenges")

        httpchallenges := make(map[string]interface{})
        redirectchallenges := make(map[string]interface{})
        dnschallenges := make(map[string]interface{})

        for _, element := range domainvalidations.Dv {
                if (element.ValidationStatus != "VALIDATED") {
                        for _, challenge := range element.Challenges {
                                if challenge.Status != "pending" {
                                        continue
                                }

                                if challenge.Type == "http-01" {
                                        httpchallenges[challenge.FullPath] = challenge.ResponseBody
                                        redirectchallenges[challenge.FullPath] = challenge.RedirectFullPath
                                }

                                if challenge.Type == "dns-01" {
                                        dnschallenges[challenge.FullPath] = challenge.ResponseBody
                                }
                        }
                }
        }

        d.Set("httpchallenges", httpchallenges)
        d.Set("redirectchallenges", redirectchallenges)
        d.Set("dnschallenges", dnschallenges)
}

