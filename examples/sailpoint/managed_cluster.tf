resource "sailpoint_managed_cluster" "mycluster" {
  name        = "Testing cluster v1.2"
  type        = "standard"
  description = "My cluster created via Terraform"
  configuration = {
    gmtOffset = "-3"
    debug = false
  }
}