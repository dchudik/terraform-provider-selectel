package selectel

import (
	"context"
	"fmt"
	"net/http"

	domainsV2 "github.com/selectel/domains-go/pkg/v2"
)

func getDomainsV2Client(meta interface{}) (domainsV2.DNSClient[domainsV2.Zone, domainsV2.RRSet], error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, fmt.Errorf("SEL_PROJECT_ID must be set for the dns v2")
	}

	selvpcClient, err := config.GetSelVPCClientWithProjectScope(config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("can't get selvpc client for domains: %w", err)
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
		"filter": zoneName,
	}
	zones, err := client.ListZones(ctx, optsForSearchZone)
	if err != nil {
		return nil, err
	}
	if zones.GetCount() == 0 {
		return nil, ErrZoneNotFound
	}
	if zones.GetCount() > 1 {
		return nil, ErrFoundMultipleZones
	}
	zone := zones.GetItems()[0]
	return zone, err
}
