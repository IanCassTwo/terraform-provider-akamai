package akamai

import (
        "encoding/json"
	"time"
        "log"
        "strings"
	"strconv"
	"fmt"

        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
//        "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
        "github.com/hashicorp/terraform/helper/schema"
//        "github.com/hashicorp/terraform/helper/validation"
)

func setAdminContact(d *schema.ResourceData, enrollment *cps.Enrollment) {

        log.Print("[DEBUG] enter setAdminContact")

	acset := d.Get("admincontact").(*schema.Set)
	aclist := acset.List()
	acmap := aclist[0].(map[string]interface{})

	var admincontact cps.Contact

	afirstname := acmap["firstname"].(string)
	alastname := acmap["lastname"].(string)
	atitle := acmap["title"].(string)
	aorganization := acmap["organization"].(string)
	aemail := acmap["email"].(string)
	aphone := acmap["phone"].(string)
	aaddresslineone := acmap["addresslineone"].(string)
	aaddresslinetwo := acmap["addresslinetwo"].(string)
	acity := acmap["city"].(string)
	aregion := acmap["region"].(string)
	apostalcode := acmap["postalcode"].(string)
	acountry := acmap["countrycode"].(string)

	admincontact.FirstName = &afirstname
	admincontact.LastName = &alastname
	admincontact.Title = &atitle
	admincontact.Organization = &aorganization
	admincontact.Email = &aemail
	admincontact.Phone = &aphone
	admincontact.AddressLineOne = &aaddresslineone
	admincontact.AddressLineTwo = &aaddresslinetwo
	admincontact.City = &acity
	admincontact.Region = &aregion
	admincontact.PostalCode = &apostalcode
	admincontact.Country = &acountry

	enrollment.AdminContact = &admincontact
}

func getAdminContact(d *schema.ResourceData, enrollment *cps.Enrollment) {

        log.Print("[DEBUG] enter getAdminContact")

	acmap := make(map[string]interface{})

	acmap["firstname"] = *enrollment.AdminContact.FirstName
	acmap["lastname"] = *enrollment.AdminContact.LastName
	acmap["title"] = *enrollment.AdminContact.Title
	acmap["organization"] = *enrollment.AdminContact.Organization
	acmap["email"] = *enrollment.AdminContact.Email
	acmap["phone"] = *enrollment.AdminContact.Phone
	acmap["addresslineone"] = *enrollment.AdminContact.AddressLineOne
	acmap["addresslinetwo"] = *enrollment.AdminContact.AddressLineTwo
	acmap["city"] = *enrollment.AdminContact.City
	acmap["region"] = *enrollment.AdminContact.Region
	acmap["postalcode"] = *enrollment.AdminContact.PostalCode
	acmap["countrycode"] = *enrollment.AdminContact.Country

	aclist := make([]interface{}, 1)
	aclist[0] = acmap
	d.Set("admincontact", aclist)
}

func setTechContact(d *schema.ResourceData, enrollment *cps.Enrollment) {

        log.Print("[DEBUG] enter setTechContact")
	// There must be a better way to do this?

	tcset := d.Get("techcontact").(*schema.Set)
	tclist := tcset.List()
	tcmap := tclist[0].(map[string]interface{})

	var techcontact cps.Contact

	tfirstname := tcmap["firstname"].(string)
	tlastname := tcmap["lastname"].(string)
	ttitle := tcmap["title"].(string)
	torganization := tcmap["organization"].(string)
	temail := tcmap["email"].(string)
	tphone := tcmap["phone"].(string)
	taddresslineone := tcmap["addresslineone"].(string)
	taddresslinetwo := tcmap["addresslinetwo"].(string)
	tcity := tcmap["city"].(string)
	tregion := tcmap["region"].(string)
	tpostalcode := tcmap["postalcode"].(string)
	tcountry := tcmap["countrycode"].(string)

	techcontact.FirstName = &tfirstname
	techcontact.LastName = &tlastname
	techcontact.Title = &ttitle
	techcontact.Organization = &torganization
	techcontact.Email = &temail
	techcontact.Phone = &tphone
	techcontact.AddressLineOne = &taddresslineone
	techcontact.AddressLineTwo = &taddresslinetwo
	techcontact.City = &tcity
	techcontact.Region = &tregion
	techcontact.PostalCode = &tpostalcode
	techcontact.Country = &tcountry

	enrollment.TechContact = &techcontact
}

func getTechContact(d *schema.ResourceData, enrollment *cps.Enrollment) {

        log.Print("[DEBUG] enter getTechContact")

	tcmap := make(map[string]interface{})

	tcmap["firstname"] = *enrollment.TechContact.FirstName
	tcmap["lastname"] = *enrollment.TechContact.LastName
	tcmap["title"] = *enrollment.TechContact.Title
	tcmap["organization"] = *enrollment.TechContact.Organization
	tcmap["email"] = *enrollment.TechContact.Email
	tcmap["phone"] = *enrollment.TechContact.Phone
	tcmap["addresslineone"] = *enrollment.TechContact.AddressLineOne
	tcmap["addresslinetwo"] = *enrollment.TechContact.AddressLineTwo
	tcmap["city"] = *enrollment.TechContact.City
	tcmap["region"] = *enrollment.TechContact.Region
	tcmap["postalcode"] = *enrollment.TechContact.PostalCode
	tcmap["countrycode"] = *enrollment.TechContact.Country

	tclist := make([]interface{}, 1)
	tclist[0] = tcmap
	d.Set("techcontact", tclist)
}

func setOrganization(d *schema.ResourceData, enrollment *cps.Enrollment) {
        log.Print("[DEBUG] enter setOrganization")

	// There must be a better way to do this?

	orgset := d.Get("organization").(*schema.Set)
	orglist := orgset.List()
	orgmap := orglist[0].(map[string]interface{})

	var organization cps.Organization

	orgname := orgmap["name"].(string)
	orgphone := orgmap["phone"].(string)
	orgaddresslineone := orgmap["addresslineone"].(string)
	orgaddresslinetwo := orgmap["addresslinetwo"].(string)
	orgcity := orgmap["city"].(string)
	orgregion := orgmap["region"].(string)
	orgpostalcode := orgmap["postalcode"].(string)
	orgcountry := orgmap["countrycode"].(string)

	organization.Name = &orgname
	organization.Phone = &orgphone
	organization.AddressLineOne = &orgaddresslineone
	organization.AddressLineTwo = &orgaddresslinetwo
	organization.City = &orgcity
	organization.Region = &orgregion
	organization.PostalCode = &orgpostalcode
	organization.Country = &orgcountry

	enrollment.Organization = &organization
}

func getOrganization(d *schema.ResourceData, enrollment *cps.Enrollment) {
        log.Print("[DEBUG] enter getOrganization")

	orgmap := make(map[string]interface{})

	orgmap["name"] = *enrollment.Organization.Name
	orgmap["phone"] = *enrollment.Organization.Phone
	orgmap["addresslineone"] = *enrollment.Organization.AddressLineOne
	orgmap["addresslinetwo"] = *enrollment.Organization.AddressLineTwo
	orgmap["city"] = *enrollment.Organization.City
	orgmap["region"] = *enrollment.Organization.Region
	orgmap["postalcode"] = *enrollment.Organization.PostalCode
	orgmap["countrycode"] = *enrollment.Organization.Country

	orglist := make([]interface{}, 1)
	orglist[0] = orgmap
	d.Set("organization", orglist)
}

func setThirdPartyCertType(d *schema.ResourceData, enrollment *cps.Enrollment) {
        log.Print("[DEBUG] enter setCertType")
        enrollment.CertificateType = cps.ThirdPartyCertificate
}

func setThirdPartyValidationType(d *schema.ResourceData, enrollment *cps.Enrollment) {
        log.Print("[DEBUG] enter setValidationType")
        enrollment.ValidationType = cps.ThirdPartyValidation
        enrollment.RegistrationAuthority = cps.ThirdPartyRA
}

func setSanCertType(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setCertType")
	enrollment.CertificateType = cps.SanCertificate
}

func getCertType(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter getCertType")
	//TODO reject anything other than SAN
	if enrollment.CertificateType == cps.SanCertificate {
		d.Set("certificatetype", "san")
	} else if enrollment.CertificateType == cps.SymantecCertificate {
		d.Set("certificatetype", "single")
	} else if enrollment.CertificateType == cps.WildCardCertificate {
		d.Set("certificatetype", "wildcard")
	} else if enrollment.CertificateType == cps.WildCardSanCertificate {
		d.Set("certificatetype", "wildcard-san")
	} else if enrollment.CertificateType == cps.ThirdPartyCertificate {
		d.Set("certificatetype", "third-party")
	}
}

func setDVValidationType(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setValidationType")
	enrollment.ValidationType = cps.DomainValidation
	enrollment.RegistrationAuthority = cps.LetsEncryptRA
}

func getValidationType(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter getValidationType")

	log.Print("[DEBUG] validation type ", enrollment.ValidationType)
	// TODO reject an import of anything other than DV
	if enrollment.ValidationType == cps.DomainValidation {
		d.Set("validationtype", "dv")
	} else if enrollment.ValidationType == cps.OrganizationValidation {
		d.Set("validationtype", "ov")
	} else if enrollment.ValidationType == cps.ExtendedValidation {
		d.Set("validationtype", "ev")
	} else if enrollment.ValidationType == cps.ThirdPartyValidation {
		d.Set("validationtype", "third-party")
	}
}

func setNetworkConfig(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setNetworkConfig")
	var networkconfig cps.NetworkConfiguration
	networkconfig.Geography = "core"

	var securenetwork cps.TLSType
	if d.Get("securenetwork") == "standard-tls" {
		securenetwork = cps.StandardTLS
	} else if d.Get("securenetwork") == "enhanced-tls" {
		securenetwork = cps.EnhancedTLS
	}
	
	networkconfig.SNIOnly = d.Get("snionly").(bool)
	networkconfig.SecureNetwork = securenetwork
	networkconfig.MustHaveCiphers = "ak-akamai-default"
	networkconfig.PreferredCiphers = "ak-akamai-default"
	enrollment.NetworkConfiguration = &networkconfig

	if networkconfig.SNIOnly {
		var sni cps.DomainNameSettings
		sni.CloneDomainNames = true

		aset := d.Get("alternativenames").(*schema.Set)
		alist := aset.List()
		if len(alist) > 0 {
			alternativenames := make([]string, len(alist))
			for i, v := range alist {
				alternativenames[i] = v.(string)

			}
			sni.DomainNames = &alternativenames
		}
		
		enrollment.NetworkConfiguration.DomainNameSettings = &sni
	}
}

func getNetworkConfig(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter getNetworkConfig")
	
	if enrollment.NetworkConfiguration.SecureNetwork == cps.EnhancedTLS {
		d.Set("securenetwork", "enhanced-tls")
	} else {
		d.Set("securenetwork", "standard-tls")
	}

	d.Set("snionly", enrollment.NetworkConfiguration.SNIOnly)
}

func setChangeManagement(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setChangeManagement")
	enrollment.ChangeManagement = false
}

func setSignatureAuthority(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setSignatureAuthority")
	var signatureauthority = cps.SHA256
	enrollment.SignatureAuthority = &signatureauthority
}

func setCSR(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter setCSR")
	orgset := d.Get("organization").(*schema.Set)
	orglist := orgset.List()
	csrmap := orglist[0].(map[string]interface{})
	var csr cps.CSR
	csr.CommonName = d.Get("commonname").(string)

	aset := d.Get("alternativenames").(*schema.Set)
	alist := aset.List()
	if len(alist) > 0 {
		alternativenames := make([]string, len(alist))
		for i, v := range alist {
			alternativenames[i] = v.(string)
			
		}
		csr.AlternativeNames = &alternativenames
	}

	csrcity := csrmap["city"].(string)
	csrstate := csrmap["region"].(string)
	csrcountrycode := csrmap["countrycode"].(string)
	csrorganization := csrmap["name"].(string)
	csrorganizationalunit := ""
	csr.City = &csrcity
	csr.State = &csrstate
	csr.CountryCode = &csrcountrycode
	csr.Organization = &csrorganization
	csr.OrganizationalUnit = &csrorganizationalunit
	enrollment.CertificateSigningRequest = &csr
}

func getCSR(d *schema.ResourceData, enrollment *cps.Enrollment) {
	log.Print("[DEBUG] enter getCSR")

	d.Set("commonname", enrollment.CertificateSigningRequest.CommonName)
	d.Set("alternativenames", enrollment.CertificateSigningRequest.AlternativeNames)
}

func setThirdParty(d *schema.ResourceData, enrollment *cps.Enrollment) {
	var thirdparty cps.ThirdParty
	thirdparty.ExcludeSANS = false
	enrollment.ThirdParty = &thirdparty
	enrollment.EnableMultiStacked = false
}

func getEnrollmentIdFromLocation(s string) string {
	parts := strings.Split(s, "/")
	for _, c := range parts {
		if _, err := strconv.Atoi(c); err == nil {
			log.Print("Found id = %s\n", c)
			return c
		}
	}
	return ""
}

func setId(d *schema.ResourceData, enrollmentid string) {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("contract"), enrollmentid))
}

func getEnrollmentIdFromId(d *schema.ResourceData) (string) {
	parts := strings.Split(d.Id(), ":")
	return parts[1]
}

func getContractIdFromId(d *schema.ResourceData) (string) {
	parts := strings.Split(d.Id(), ":")
	return parts[0]
}

func awaitCertVerification(enrollment cps.Enrollment) error {

        currentstatus, err := enrollment.GetChangeStatus()
        if err != nil {
                return err
        }

        // wait until Akamai cert verication is complete
        for currentstatus.StatusInfo.Status == "verify-cert-request" {
                time.Sleep(10 * time.Second)

                var err error
                currentstatus, err = enrollment.GetChangeStatus()
                if err != nil {
                        return err
                }
                s,_ := json.MarshalIndent(currentstatus, "", "\t")
                log.Print("[DEBUG] current status ", string(s))

                // Did we have a validation error?
                if currentstatus.StatusInfo.Error.Description != "" {
                        return fmt.Errorf(currentstatus.StatusInfo.Error.Description)
                }
        }
        return nil
}
