// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsSRVRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify datasource
			{
				Config: providerConfig + `
			resource "abion_dns_srv_record" "test" {
 			  zone  = "pmapitest.com"
			  name = "@"
			  records = [
				{
				  target   = "server1.pmapitest.com."
				  port     = "443"
				  priority = "1"
				  weight   = "100"
				  comments = "test comment"
				  ttl      = "3600"
				},
			  ]
			}

			data "abion_dns_srv_record" "test_data" {
              zone  = abion_dns_srv_record.test.zone
              name  = abion_dns_srv_record.test.name
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "name", "@"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.#", "1"),

					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.target", "server1.pmapitest.com."),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.port", "443"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.priority", "1"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.weight", "100"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.ttl", "3600"),
					resource.TestCheckResourceAttr("data.abion_dns_srv_record.test_data", "records.0.comments", "test comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsSRVRecordNonExistingZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			data "abion_dns_srv_record" "non_existing" {
              zone  = "non_existing.com"
			  name = "@"
			}
			`,
				ExpectError: regexp.MustCompile("Unable to Read Zone from Abion API"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsSRVRecordNoRecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_srv_record" "test2" {
 			  zone  = "pmapitest.com"
			  name = "@"
			  records = [
				{
				  target   = "server1.pmapitest.com."
				  port     = "443"
				  priority = "1"
				  weight   = "100"
				},
			  ]
			}

			data "abion_dns_srv_record" "test_data" {
              zone  = abion_dns_srv_record.test2.zone
  			  name = "xxx"
			}
			`,
				ExpectError: regexp.MustCompile("No records exist on xxx level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsSRVRecordNoARecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create another record type on subdomain and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_aaaa_record" "test3" {
 			  zone  = "pmapitest.com"
			  name = "test3"
			  records = [
				{
				  ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff0"
				}
			  ]
			}

			data "abion_dns_srv_record" "test_data" {
              zone  = abion_dns_aaaa_record.test3.zone
  			  name = "test3"
			}
			`,
				ExpectError: regexp.MustCompile("No SRV records exist on test3 level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
