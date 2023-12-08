# Edit domain name for import from your project
variable "domain_name_for_import" {
  default = "tf-provider-test-import-basic.ru."
}

import {
    id = var.domain_name_for_import
    provider = selectel
    to = tf_basic_import_ru
}

resource "selectel_domains_zone_v2" "tf_basic_import_ru" {
  name = var.domain_name_for_import
}

# Edit rrset name for import from your project
variable "rrset_name_for_import" {
  default = format("a.%s", var.domain_name_for_import)
}
# Edit rrset type for import from your project
variable "rrset_type_for_import" {
  default = "A"
}

# For import rrset use domain_name/rrset_name/rrset_type as resource id 
import {
    id = format("%s/%s/%s", 
        var.domain_name_for_import, 
        var.rrset_name_for_import, var.rrset_type_for_import)
    provider = selectel
    to = a_tf_basic_import_ru
}

resource "selectel_domains_rrset_v2" "a_tf_basic_import_ru" {
  name = var.rrset_name_for_import
  # Edit zone id if you don't want import zone  
  zone_id = tf_basic_import_ru.id
  type = var.rrset_type_for_import
  ttl = 60
  records {
    content = "Hello tf"
    disabled = false
  }
}