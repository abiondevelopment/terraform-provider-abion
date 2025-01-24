// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsCAARecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify datasource
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test" {
 			  zone  = "pmapitest9.com"
			  name = "@"
			  records = [
				{
				  flag     = "0"
				  tag      = "issue"
				  value    = "letsencrypt.test"
				  comments = "test comment"
				  ttl      = "3600"
				},
			  ]
			}

			data "abion_dns_caa_record" "test_data" {
              zone  = abion_dns_caa_record.test.zone
              name  = abion_dns_caa_record.test.name
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "zone", "pmapitest9.com"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "name", "@"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.#", "1"),

					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.0.flag", "0"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.0.tag", "issue"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.0.value", "letsencrypt.test"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.0.ttl", "3600"),
					resource.TestCheckResourceAttr("data.abion_dns_caa_record.test_data", "records.0.comments", "test comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsCAARecordNonExistingZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			data "abion_dns_caa_record" "non_existing" {
              zone  = "non_existing.com"
			  name  = "@"
			}
			`,
				ExpectError: regexp.MustCompile("Unable to Read Zone from Abion API"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsCAARecordNoRecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test2" {
 			  zone  = "pmapitest9.com"
			  name = "@"
			  records = [
				{
				  flag   = "0"
				  tag    = "issue"
				  value  = "letsencrypt.test"
				},
			  ]
			}

			data "abion_dns_caa_record" "test_data" {
              zone  = abion_dns_caa_record.test2.zone
  			  name  = "xxx"
			}
			`,
				ExpectError: regexp.MustCompile("No records exist on xxx level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsCAARecordNoARecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create another record type on subdomain and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_aaaa_record" "test3" {
 			  zone  = "pmapitest9.com"
			  name  = "test3"
			  records = [
				{
				  ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff0"
				}
			  ]
			}

			data "abion_dns_caa_record" "test_data" {
              zone  = abion_dns_aaaa_record.test3.zone
  			  name  = "test3"
			}
			`,
				ExpectError: regexp.MustCompile("No CAA records exist on test3 level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
