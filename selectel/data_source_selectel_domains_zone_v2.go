package selectel

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrZoneNotFound       = errors.New("zone not found")
	ErrFoundMultipleZones = errors.New("found multiple zones")
)

func dataSourceDomainsZoneV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainsZoneV2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delegation_checked_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_check_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_delegated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO: add it to api
			// "disabled": {
			// 	Type:     schema.TypeBool,
			// 	Computed: true,
			// },
		},
	}
}

func dataSourceDomainsZoneV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneName := d.Get("name").(string)

	log.Print(msgGet(objectZone, zoneName))

	optsForSearchZone := &map[string]string{
		"filter": zoneName,
	}
	zones, err := client.ListZones(ctx, optsForSearchZone)
	if err != nil {
		return diag.FromErr(errGettingObject(objectDomain, zoneName, err))
	}
	if zones.GetCount() == 0 {
		// return err_not_found
		return diag.FromErr(errGettingObject(objectDomain, zoneName, ErrZoneNotFound))
	}
	if zones.GetCount() > 1 {
		// return err_many_zones
		return diag.FromErr(errGettingObject(objectDomain, zoneName, ErrFoundMultipleZones))
	}
	zone := zones.GetItems()[0]
	d.SetId(zone.UUID)
	d.Set("name", zone.Name)
	d.Set("comment", zone.Comment)
	d.Set("created_at", zone.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", zone.UpdatedAt.Format(time.RFC3339))
	d.Set("delegation_checked_at", zone.DelegationCheckedAt.Format(time.RFC3339))
	d.Set("last_check_status", zone.LastCheckStatus)
	d.Set("last_delegated_at", zone.LastDelegatedAt.Format(time.RFC3339))
	d.Set("project_id", zone.ProjectID)

	return nil
}
