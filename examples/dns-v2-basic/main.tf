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
  records {
    content  = "\"Hello also terraform\""
    disabled = false
  }
}

