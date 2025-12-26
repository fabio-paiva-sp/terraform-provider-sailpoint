terraform {
  required_providers {
    sailpoint = {
      source = "hashicorp.com/sailpoint/sailpoint"
    }
  }
}

provider "sailpoint" {}

data "sailpoint_coffees" "example" {}
