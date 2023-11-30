package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	domainsV2 "github.com/selectel/domains-go/pkg/v2"
	v2 "github.com/selectel/domains-go/pkg/v2"
)

var (
	ErrRrsetNotFound       = errors.New("rrset not found")
	ErrFoundMultipleRRsets = errors.New("found multiple rrsets")
)

func dataSourceDomainsRrsetV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainsRrsetV2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainsRrsetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	rrsetName := d.Get("name").(string)
	zoneID := d.Get("zone_id").(string)
	rrsetType := d.Get("type").(string)

	log.Print(msgGet(objectRrset, rrsetName))
	// TODO: id Option and create and test if statement
	// TODO: type Option and create and test if statement
	optsForSearchRrset := &map[string]string{
		"name": rrsetName,
		"type": rrsetType,
	}

	rrsets, err := client.ListRRSets(ctx, zoneID, optsForSearchRrset)
	if err != nil {
		return diag.FromErr(errGettingObject(objectRrset, rrsetName, err))
	}

	r, err := regexp.Compile(fmt.Sprintf("^%s.?", rrsetName))
	if err != nil {
		return diag.FromErr(errGettingObject(objectRrset, rrsetName, err))
	}

	var rrset *domainsV2.RRSet
	for _, rrsetInResp := range rrsets.GetItems() {
		if match := r.MatchString(rrsetInResp.Name); match {
			rrset = rrsetInResp
			break
		}
	}
	if rrset == nil {
		return diag.FromErr(errGettingObject(objectRrset, rrsetName, ErrRrsetNotFound))
	}

	d.SetId(rrset.UUID)
	d.Set("name", rrset.Name)
	d.Set("comment", rrset.Comment)
	d.Set("managed_by", rrset.ManagedBy)
	d.Set("ttl", rrset.TTL)
	d.Set("type", rrset.Type)
	d.Set("zone_id", rrset.ZoneUUID)
	d.Set("records", generateListFromRecords(rrset.Records))

	return nil
}

// generateListFromRecords - generate terraform TypeList from records in rrset
func generateListFromRecords(records []v2.RecordItem) []interface{} {
	var recordsAsList []interface{}
	for _, record := range records {
		recordsAsList = append(recordsAsList, map[string]interface{}{
			"content":  record.Content,
			"disabled": record.Disabled,
		})
	}

	return recordsAsList
}
