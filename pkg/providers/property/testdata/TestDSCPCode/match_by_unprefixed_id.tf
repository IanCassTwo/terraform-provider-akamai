provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cp_code" "test" {
  name     = "test2"
  contract = "ctr_test"
  group    = "grp_test"
}
