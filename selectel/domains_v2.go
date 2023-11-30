package selectel

import (
	"fmt"
	"net/http"

	domainsV2 "github.com/selectel/domains-go/pkg/v2"
)

// TODO: add DNS v2 in examples
func getDomainsV2Client(meta interface{}) (domainsV2.DNSClient[domainsV2.Zone, domainsV2.RRSet], error) {
	config := meta.(*Config)

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
