package selectel

import (
	"fmt"
)

func testAccDomainsRrsetV2Basic(resourceRrsetName, rrsetName, rrsetType, rrsetContent string, ttl int, resourceZoneName string) string {
	return fmt.Sprintf(`
		resource "selectel_domains_rrset_v2" %[1]q {
		name = %[2]q
		type = %[3]q
		ttl = %[4]d
		zone_id = selectel_domains_zone_v2.%[5]s.id
		records {
			content = %[6]q
			disabled = false
		}
	}`, resourceRrsetName, rrsetName, rrsetType, ttl, resourceZoneName, rrsetContent)
}
