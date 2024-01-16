package selectel

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDomainsZoneV2Basic(t *testing.T) {
	testZoneName := fmt.Sprintf("%s.xyz.", acctest.RandomWithPrefix("tf-acc"))
	resourceZoneName := "zone_tf_acc_test_1"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheckWithProjectID(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDomainsV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsZoneV2Basic(resourceZoneName, testZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainsZoneV2Exists(fmt.Sprintf("selectel_domains_zone_v2.%[1]s", resourceZoneName)),
					resource.TestCheckResourceAttr(fmt.Sprintf("selectel_domains_zone_v2.%[1]s", resourceZoneName), "name", testZoneName),
				),
			},
		},
	})
}
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
