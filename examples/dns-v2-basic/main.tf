terraform {
  required_providers {
    selectel = {
      source  = "selectel/selectel"
      version = "4.0.3"
    }
  }
}

# Edit zone name to your unique zone name
variable "zone_name" {
  default = "tf-provider-test-basic.ru."
}

resource "selectel_domains_zone_v2" "tf_basic_ru" {
  name = var.zone_name
}

resource "selectel_domains_rrset_v2" "a_tf_basic_ru" {
  zone_id = selectel_domains_zone_v2.tf_basic_ru.id
  name    = format("a.%s",var.zone_name)
  type    = "A"
  ttl     = 60
  records {
    content  = "127.0.0.1"
    disabled = false
  }
}

resource "selectel_domains_rrset_v2" "txt_tf_basic_ru" {
  zone_id = selectel_domains_zone_v2.tf_basic_ru.id
  name    = format("txt.%s",var.zone_name)
  type    = "TXT"
  ttl     = 60
  records {
    content  = "\"Hello terraform\""
    disabled = false
  }
}

# Edit zone name for get information about zone from your project
variable "zone_name_for_data" {
  default = "tf-provider-test-data-basic.ru."
}

data "selectel_domains_zone_v2" "data_tf_selectel_basic_ru" {
  name = var.zone_name_for_data
}

data "selectel_domains_rrset_v2" "data_txt_tf_selectel_basic_ru" {
  name    = format("txt.%s", var.zone_name_for_data)
  zone_id = data.selectel_domains_zone_v2.data_tf_selectel_basic_ru.id
  type    = "TXT"
}



