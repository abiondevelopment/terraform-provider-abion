// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsNSRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify datasource
			{
				Config: providerConfig + `
			resource "abion_dns_ns_record" "test" {
 			  zone  = "pmapitest.com"
			  name = "www"
			  records = [
				{
				  nameserver = "ns1.testabiondns.se."
				  ttl = "3600"
				  comments = "test comment"
				}
			  ]
			}

			data "abion_dns_ns_record" "test_data" {
              zone  = abion_dns_ns_record.test.zone
			  name = "www"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "name", "www"),
					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "records.#", "1"),

					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "records.0.nameserver", "ns1.testabiondns.se."),
					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "records.0.ttl", "3600"),
					resource.TestCheckResourceAttr("data.abion_dns_ns_record.test_data", "records.0.comments", "test comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsNSRecordNonExistingZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			data "abion_dns_ns_record" "non_existing" {
              zone  = "non_existing.com"
			  name = "www"
			}
			`,
				ExpectError: regexp.MustCompile("Unable to Read Zone from Abion API"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsNSRecordNoNSRecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create another record type on subdomain and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_a_record" "test3" {
 			  zone  = "pmapitest.com"
			  name = "test3"
			  records = [
				{
				  ip_address = "203.0.113.0"
				}
			  ]
			}

			data "abion_dns_ns_record" "test_data" {
              zone  = abion_dns_a_record.test3.zone
  			  name = "test3"
			}
			`,
				ExpectError: regexp.MustCompile("No NS records exist on test3 level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
