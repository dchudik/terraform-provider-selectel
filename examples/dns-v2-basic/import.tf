# Edit zone name for import from your project
variable "zone_name_for_import" {
  default = "tf-provider-test-import-basic.ru."
}

# Edit rrset name for import from your project
variable "rrset_name_for_import" {
  default = "a.tf-provider-test-import-basic.ru."
}

# Edit rrset type for import from your project
variable "rrset_type_for_import" {
  default = "A"
}

import {
    id = var.zone_name_for_import
    to = selectel_domains_zone_v2.tf_basic_import_ru
}

resource "selectel_domains_zone_v2" "tf_basic_import_ru" {
  name = var.zone_name_for_import
}

# For import rrset use zone_name/rrset_name/rrset_type as resource id 
import {
    id = format("%s/%s/%s", var.zone_name_for_import, 
      var.rrset_name_for_import, var.rrset_type_for_import)
    to = selectel_domains_rrset_v2.a_tf_basic_import_ru
}

resource "selectel_domains_rrset_v2" "a_tf_basic_import_ru" {
  name = var.rrset_name_for_import
  # Edit zone id if you don't want import zone  
  zone_id = selectel_domains_zone_v2.tf_basic_import_ru.id
  type = var.rrset_type_for_import
  ttl = 60
  # If you set not all records
  # records wich not set will be removed
  records {
    content = "1.2.3.5"
    disabled = false
  }
}