package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	domainsV2 "github.com/selectel/domains-go/pkg/v2"
)

var ErrProjectIDNotSetupForDNSV2 = errors.New("env variable SEL_PROJECT_ID or variable project_id must be set for the dns v2")

func getDomainsV2Client(meta interface{}) (domainsV2.DNSClient[domainsV2.Zone, domainsV2.RRSet], error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, ErrProjectIDNotSetupForDNSV2
	}

	selvpcClient, err := config.GetSelVPCClientWithProjectScope(config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("can't get selvpc client for domains v2: %w", err)
	}

	httpClient := &http.Client{}
	userAgent := "terraform-provider-selectel"
	defaultApiURL := "https://api.selectel.ru/domains/v2"
	hdrs := http.Header{}
	hdrs.Add("X-Auth-Token", selvpcClient.GetXAuthToken())
	hdrs.Add("User-Agent", userAgent)
	domainsClient := domainsV2.NewClient(defaultApiURL, httpClient, hdrs)

	return domainsClient, nil
}

func getZoneByName(ctx context.Context, client domainsV2.DNSClient[domainsV2.Zone, domainsV2.RRSet], zoneName string) (*domainsV2.Zone, error) {
	optsForSearchZone := &map[string]string{
		"name": zoneName,
	}
	zones, err := client.ListZones(ctx, optsForSearchZone)
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile(fmt.Sprintf("^%s.?", zoneName))
	if err != nil {
		return nil, err
	}

	for _, zone := range zones.GetItems() {
		if r.MatchString(zone.Name) {
			return zone, nil
		}
	}

	return nil, ErrZoneNotFound
}

func getRrsetByNameAndType(ctx context.Context, client domainsV2.DNSClient[domainsV2.Zone, domainsV2.RRSet], zoneID, rrsetName, rrsetType string) (*domainsV2.RRSet, error) {
	optsForSearchRrset := &map[string]string{
		"name": rrsetName,
		"type": rrsetType,
	}

	rrsets, err := client.ListRRSets(ctx, zoneID, optsForSearchRrset)
	if err != nil {
		return nil, errGettingObject(objectRrset, rrsetName, err)
	}

	r, err := regexp.Compile(fmt.Sprintf("^%s.?", rrsetName))
	if err != nil {
		return nil, errGettingObject(objectRrset, rrsetName, err)
	}

	var rrset *domainsV2.RRSet
	for _, rrsetInResp := range rrsets.GetItems() {
		match := r.MatchString(rrsetInResp.Name)
		if match && string(rrsetInResp.Type) == rrsetType {

			rrset = rrsetInResp
			break
		}
	}
	log.Println("Selected rrset:", rrset)
	if rrset == nil {
		return nil, errGettingObject(objectRrset, rrsetName, ErrRrsetNotFound)
	}

	return rrset, nil
}
