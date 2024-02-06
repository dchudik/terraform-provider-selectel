package selectel

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	domainsV2 "github.com/selectel/domains-go/pkg/v2"
)

const resourceRRSetName = "rrset_tf_acc_test_1"

func TestAccDomainsRRSetV2DataSourceBasic(t *testing.T) {
	testZoneName := fmt.Sprintf("%s.ru.", acctest.RandomWithPrefix("tf-acc"))
	testRRSetName := fmt.Sprintf("%[1]s.%[2]s", acctest.RandomWithPrefix("tf-acc"), testZoneName)
	testRRSetType := domainsV2.TXT
	testRRSetTTL := 60
	testRRSetContent := fmt.Sprintf("\"%[1]s\"", acctest.RandString(16))
	dataSourceRRSetName := fmt.Sprintf("data.selectel_domains_rrset_v2.%[1]s", resourceRRSetName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheckWithProjectID(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDomainsV2RRSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsRRSetV2DataSourceBasic(resourceRRSetName, testRRSetName, string(testRRSetType), testRRSetContent, testRRSetTTL, resourceZoneName, testZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainsRRSetV2ID(dataSourceRRSetName),
					resource.TestCheckResourceAttr(dataSourceRRSetName, "name", testRRSetName),
					resource.TestCheckResourceAttr(dataSourceRRSetName, "type", string(testRRSetType)),
					resource.TestCheckResourceAttrSet(dataSourceRRSetName, "zone_id"),
				),
			},
		},
	})
}

func testAccDomainsRRSetV2ID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("can't find rrset: %s", name)
		}

		if rs.Primary.ID == "" {
			return errors.New("rrset data source ID not set")
		}

		return nil
	}
}

func testAccDomainsRRSetV2DataSourceBasic(resourceRRSetName, rrsetName, rrsetType, rrsetContent string, ttl int, resourceZoneName, zoneName string) string {
	return fmt.Sprintf(`
	%[1]s

	%[2]s

	data "selectel_domains_rrset_v2" %[3]q {
	  name = selectel_domains_rrset_v2.%[3]s.name
	  type = selectel_domains_rrset_v2.%[3]s.type
	  zone_id = selectel_domains_zone_v2.%[4]s.id
	}
`, testAccDomainsZoneV2Basic(resourceZoneName, zoneName), testAccDomainsRRSetV2Basic(resourceRRSetName, rrsetName, rrsetType, rrsetContent, ttl, resourceZoneName), resourceRRSetName, resourceZoneName)
}
