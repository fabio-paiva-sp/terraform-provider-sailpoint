terraform {
  required_providers {
    sailpoint = {
      source = "hashicorp.com/sailpoint/sailpoint"
    }
  }
}

provider "sailpoint" {

}

data "sailpoint_managed_clusters" "example" {
   filters = "type eq \"iai\""
}

output "clusters" {
  value = data.sailpoint_managed_clusters.example
}

data "sailpoint_managed_cluster" "example" {
  id = "85cef106836d44df9f757751815d1ea9"
}

output "cluster_by_id" {
  value = data.sailpoint_managed_cluster.example
}
