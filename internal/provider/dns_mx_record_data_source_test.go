// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsMXRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify datasource
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test" {
 			  zone  = "pmapitest4.com"
			  name = "@"
			  records = [
				{
			      host     = "mail1.pmapitest4.com."
      			  priority = "10"
				  ttl = "3600"
				  comments = "test comment"
				}
			  ]
			}

			data "abion_dns_mx_record" "test_data" {
              zone  = abion_dns_mx_record.test.zone
              name  = abion_dns_mx_record.test.name
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "zone", "pmapitest4.com"),
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "name", "@"),
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "records.#", "1"),

					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "records.0.host", "mail1.pmapitest4.com."),
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "records.0.priority", "10"),
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "records.0.ttl", "3600"),
					resource.TestCheckResourceAttr("data.abion_dns_mx_record.test_data", "records.0.comments", "test comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsMXRecordNonExistingZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			data "abion_dns_mx_record" "non_existing" {
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

func TestAccDnsMXRecordNoRecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create record and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test2" {
 			  zone  = "pmapitest4.com"
			  name = "@"
			  records = [
				{
			      host     = "mail.pmapitest4.com."
      			  priority = "10"
				}
			  ]
			}

			data "abion_dns_mx_record" "test_data" {
              zone  = abion_dns_mx_record.test2.zone
  			  name = "xxx"
			}
			`,
				ExpectError: regexp.MustCompile("No records exist on xxx level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsMXRecordNoARecordOnSubDomainLevelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create another record type on subdomain and verify error
			{
				Config: providerConfig + `
			resource "abion_dns_aaaa_record" "test3" {
 			  zone  = "pmapitest4.com"
			  name = "test3"
			  records = [
				{
				  ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff0"
				}
			  ]
			}

			data "abion_dns_mx_record" "test_data" {
              zone  = abion_dns_aaaa_record.test3.zone
  			  name = "test3"
			}
			`,
				ExpectError: regexp.MustCompile("No MX records exist on test3 level"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
