package akamai

import (
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/networklists-v2"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	//        "github.com/hashicorp/terraform/helper/resource"
)

func resourceNetworkList() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkListCreate,
		Read:   resourceNetworkListRead,
		Delete: resourceNetworkListDelete,
		Update: resourceNetworkListUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IP",
					"GEO",
				}, false),
			},
			"description": {
				Type: schema.TypeString,
				Required: true,
			},
			"cidrblocks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			// TODO validation len > 0
			},
			"syncpoint": {
				Type: schema.TypeInt,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetworkListDelete(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceNetworkListDelete")

	networklistid := d.Id()
	if networklistid == "" {
		return nil
	}

	message, err := networklists.DeleteNetworkList(networklistid)
	if err != nil {
		return err
	}

	if message.Status == 200 {
		d.SetId("")
	}
	return nil
}

func resourceNetworkListUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceNetworkListUpdate")

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	listtype := d.Get("type").(string)
	syncpoint := d.Get("syncpoint").(int)
	aset := d.Get("cidrblocks").(*schema.Set)
	alist := aset.List()
	list := make([]string, len(alist))
	for i, v := range alist {
		list[i] = v.(string)
	}

	var networklist networklists.NetworkList
	networklist.Name = name
	networklist.Type = listtype
	networklist.Description = description
	networklist.UniqueID = d.Id()
	networklist.SyncPoint = syncpoint
	networklist.List = list

	// Update network list
	newnetworklist, err := networklists.UpdateNetworkList(networklist)
	if err != nil {
		return err
	}
	
	d.Set("syncpoint", newnetworklist.SyncPoint)

	return nil
}


func resourceNetworkListCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceNetworkListCreate")

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	listtype := d.Get("type").(string)
	aset := d.Get("cidrblocks").(*schema.Set)
	alist := aset.List()
	list := make([]string, len(alist))
	for i, v := range alist {
		list[i] = v.(string)
	}

	var networklist networklists.NetworkList
	networklist.Name = name
	networklist.Type = listtype
	networklist.Description = description
	networklist.List = list

	// Create new network list
	newnetworklist, err := networklists.CreateNetworkList(networklist)
	if err != nil {
		return err
	}

	d.SetId(newnetworklist.UniqueID)
	d.Set("syncpoint", newnetworklist.SyncPoint)

	return nil
}

func resourceNetworkListRead(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] enter resourceNetworkListRead")

	networklistid := d.Id()
	if networklistid == "" {
		return nil
	}

	networklist, err := networklists.GetNetworkList(networklistid)
	if err != nil {
		return err
	}

	d.Set("name", networklist.Name)
	d.Set("description", networklist.Description)
	d.Set("type", networklist.Type)
	d.Set("cidrblocks", networklist.List)
	d.Set("syncpoint", networklist.SyncPoint)
	return nil
}

