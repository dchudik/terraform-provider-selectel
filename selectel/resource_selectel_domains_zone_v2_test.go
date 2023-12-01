package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccDomainsZoneV2Basic(resourceName, zoneName string) string {
	return fmt.Sprintf(`
resource "selectel_domains_zone_v2" %[1]q {
  name = %[2]q
}`, resourceName, zoneName)
}

func testAccCheckDomainsV2ZoneDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	client, err := getDomainsV2Client(meta)
	if err != nil {
		return err
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		log.Printf("RT: %s", rs.Type)
		if rs.Type != "selectel_domains_zone_v2" {
			continue
		}

		zoneID := rs.Primary.ID

		_, err = client.GetZone(ctx, zoneID, nil)
		if err == nil {
			return errors.New("domain still exists")
		}
	}

	return nil
}
