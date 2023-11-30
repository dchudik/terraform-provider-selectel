package selectel

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	domainsV2 "github.com/selectel/domains-go/pkg/v2"
)

func resourceDomainsZoneV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainsZoneV2Create,
		ReadContext:   resourceDomainsZoneV2Read,
		DeleteContext: resourceDomainsZoneV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			// TODO: add field disabled in API lib
			// "disabled": {
			// 	Type:     schema.TypeBool,
			// 	Computed: true,
			// },
		},
	}
}

func resourceDomainsZoneV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneName := d.Get("name").(string)
	createOpts := &domainsV2.Zone{
		Name: zoneName,
	}
	log.Print(msgCreate(objectZone, zoneName))
	zone, err := client.CreateZone(ctx, createOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectDomain, err))
	}
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

func resourceDomainsZoneV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneName := d.Get("name").(string)
	optsForSearchZone := &map[string]string{
		"filter": zoneName,
	}
	zones, err := client.ListZones(ctx, optsForSearchZone)
	if err != nil {
		return diag.FromErr(errGettingObject(objectDomain, zoneName, err))
	}
	if zones.GetCount() == 0 {
		return diag.FromErr(errGettingObject(objectDomain, zoneName, ErrZoneNotFound))
	}
	if zones.GetCount() > 1 {
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

func resourceDomainsZoneV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.DeleteZone(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errDeletingObject(objectDomain, d.Id(), err))
	}

	return nil
}
