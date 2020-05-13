package akamai

import (
	"log"
	"fmt"
	"strconv"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/customdiff"
)

func resourceAppSecConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppSecConfigCreate,
		Read:   resourceAppSecConfigRead,
		Delete: resourceAppSecConfigDelete,
		Update: resourceAppSecConfigUpdate,

		Schema: map[string]*schema.Schema{
			"configid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"hostnames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(

			customdiff.ComputedIf("version", func (d *schema.ResourceDiff, meta interface{}) bool {
			    // Any change to "content" causes a new "version_id" to be allocated.
			    return d.HasChange("hostnames")
			}),
		),
	}
}

func resourceAppSecConfigCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entering resourceAppSecConfigCreate")
	return resourceAppSecConfigUpdate(d, meta)
}

func resourceAppSecConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entering resourceAppSecConfigDelete")
	// No delete operation exists.
	d.SetId("")
	return schema.Noop(d, meta)
}

func resourceAppSecConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entering resourceAppSecConfigUpdate")

	configid := d.Get("configid").(int)

	versionlist, err := appsec.ListConfigurationVersions(configid)
	if err != nil {
		return fmt.Errorf("Error listing configuration versions: %s", err)
	}

	currentversion := versionlist.LastCreatedVersion
	prodversion := versionlist.ProductionActiveVersion
	stagingversion := versionlist.StagingActiveVersion

	if currentversion == prodversion || currentversion == stagingversion {
		var configurationclone appsec.ConfigurationClone
		configurationclone.CreateFromVersion = currentversion
		configurationclone.RuleUpdate = false //FIXME: make optional

		version, err := appsec.CloneConfigurationVersion(configid, configurationclone)
		if err != nil {
			return fmt.Errorf("Error cloning a configuration : %s", err)
		}

		currentversion = version.Version
	}

	var selectedhostnames appsec.SelectedHostnames

	aset := d.Get("hostnames").(*schema.Set)
	alist := aset.List()
	list := make([]appsec.HostnameList, len(alist))
	for i, v := range alist {
		list[i].Hostname = v.(string)
	}

	selectedhostnames.HostnameList = list
	_, err = appsec.UpdateSelectedHostnames(configid, currentversion, selectedhostnames)
	if err != nil {
		return fmt.Errorf("Error updating selected hostnames : %s", err)
	}

	d.SetId(fmt.Sprintf("%d", configid))

	return resourceAppSecConfigRead(d, meta)
}

func resourceAppSecConfigRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entered resourceAppSecConfigRead")

	configid, err := strconv.Atoi(d.Id())
	versionlist, err := appsec.ListConfigurationVersions(configid)
	if err != nil {
		return fmt.Errorf("Error listing configuration versions : %s", err)
	}

	version := versionlist.LastCreatedVersion
	d.Set("name", versionlist.ConfigName)
	d.Set("version", version)
	d.Set("configid", d.Id())

	selectedhostnames, err := appsec.ListSelectedHostnames(configid, version)
	if err != nil {
		return fmt.Errorf("Error listing selected hostnames : %s", err)
	}

	hostnames := make([]string, len(selectedhostnames.HostnameList))
	for i, v := range selectedhostnames.HostnameList {
		hostnames[i] = v.Hostname
	}

	d.Set("hostnames", hostnames)

	return nil
}
