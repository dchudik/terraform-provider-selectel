package selectel

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDomainsZoneV2DataSourceBasic(t *testing.T) {
	// TODO: fix with dot in end. When somain don't end dot, then tests not pass
	testZoneName := fmt.Sprintf("%s.ru.", acctest.RandomWithPrefix("tf-acc"))
	resourceZoneName := "zone_tf_acc_test_1"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheckWithProjectID(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDomainsV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsZoneV2DataSourceBasic(resourceZoneName, testZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainsZoneV2DataSourceID(fmt.Sprintf("data.selectel_domains_zone_v2.%[1]s", resourceZoneName)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.selectel_domains_zone_v2.%[1]s", resourceZoneName), "name", testZoneName),
				),
			},
		},
	})
}

func testAccDomainsZoneV2DataSourceID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("can't find zone data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return errors.New("zone data source ID not set")
		}

		return nil
	}
}

func testAccDomainsZoneV2DataSourceBasic(resourceName, zoneName string) string {
	return fmt.Sprintf(`
	%[1]s

	data "selectel_domains_zone_v2" %[2]q {
	  name = selectel_domains_zone_v2.%[2]s.name
	}
`, testAccDomainsZoneV2Basic(resourceName, zoneName), resourceName, zoneName)
}
