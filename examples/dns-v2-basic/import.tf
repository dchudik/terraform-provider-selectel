# import {
#     id = "tf-provider-test-import-basic.ru."
#     provider = selectel
#     to = tf_basic_import_ru
# }

# resource "selectel_domains_zone_v2" "tf_basic_import_ru" {
#   name = "tf-provider-test-import-basic.ru."
# }


# import {
#     id = "tf-provider-test-import-basic.ru./a.tf-provider-test-import-basic.ru./A"
#     provider = selectel
#     to = a_tf_basic_import_ru
# }

# resource "selectel_domains_rrset_v2" "a_tf_basic_import_ru" {
#   name = "tf-provider-test-import-basic.ru."
#   zone_id = "890caaH7-jy5s-441e-93ad-0a7b75402tas26"
#   type = "A"
#   ttl = 60
#   records {
#     content = "Hello tf"
#     disabled = false
#   }
# }