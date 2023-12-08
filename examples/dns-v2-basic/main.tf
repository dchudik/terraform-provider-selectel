terraform {
  required_providers {
    selectel = {
      source  = "selectel/selectel"
      version = "4.0.3"
    }
  }
}

# Edit domain name to your unique domain name
variable "domain_name" {
  default = "tf-provider-test-basic3.ru."
}

resource "selectel_domains_zone_v2" "tf_basic_ru" {
  name = var.domain_name
}

resource "selectel_domains_rrset_v2" "a_tf_basic_ru" {
  zone_id = selectel_domains_zone_v2.tf_basic_ru.id
  name    = format("a.%s",var.domain_name)
  type    = "A"
  ttl     = 60
  records {
    content  = "127.0.0.1"
    disabled = false
  }
}

resource "selectel_domains_rrset_v2" "txt_tf_basic_ru" {
  zone_id = selectel_domains_zone_v2.tf_basic_ru.id
  name    = format("txt.%s",var.domain_name)
  type    = "TXT"
  ttl     = 60
  records {
    content  = "\"Hello terraform\""
    disabled = false
  }
}


data "selectel_domains_zone_v2" "data_tf_selectel_basic_ru" {
  name = "tf-provider-test-basic.ru."
}

data "selectel_domains_rrset_v2" "data_txt_tf_selectel_basic_ru" {
  name    = "txt.tf-provider-test-basic.ru."
  zone_id = data.selectel_domains_zone_v2.data_tf_selectel_basic_ru.id
  type    = "TXT"
}



