// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsTXTRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
			{
				Config: providerConfig + `
			resource "abion_dns_txt_record" "test" {
			  zone  = "pmapitest.com"
			  name = "@"
			  records = [
				{
				  txt_data = "txt data 1"
				},
				{
				  txt_data = "txt data 2"
				  ttl = "3600"
				},
				{
				  txt_data = "txt data 3"
				  comments = "test comment"
				}
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.0.txt_data", "txt data 1"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.1.txt_data", "txt data 2"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.1.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.2.txt_data", "txt data 3"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.2.ttl"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.2.comments", "test comment"),
				),
			},
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_txt_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  txt_data = "txt data 1"
				},
				{
				  txt_data = "txt data 2"
				},
				{
				  txt_data = "txt data 3"
				}
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.0.txt_data", "txt data 1"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.1.txt_data", "txt data 2"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.2.txt_data", "txt data 3"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.2.comments"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.2.ttl"),
				),
			},

			// Update and Read testing on subdomain level, remove ip address
			{
				Config: providerConfig + `
			resource "abion_dns_txt_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  txt_data = "txt data 1"
				},
				{
				  txt_data = "txt data 3"
				}
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.0.txt_data", "txt data 1"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.comments"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.0.ttl"),

					resource.TestCheckResourceAttr("abion_dns_txt_record.test", "records.1.txt_data", "txt data 3"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_txt_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsTXTRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_txt_record" "test2" {
			  zone  = "non_existing.com"
			  name = "@"
			  records = [
				{
				  txt_data = "text data 1"
				},
			  ]
			}
			`,
				ExpectError: regexp.MustCompile(`(?s)Error patching zone.*Zone not found`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
