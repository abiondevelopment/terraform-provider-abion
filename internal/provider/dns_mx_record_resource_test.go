// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsMXRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test" {
			  zone  = "pmapitest.com"
			  name = "@"
			  records = [
				{
				  host     = "mail1.pmapitest.com."
				  priority = "10"
				},
				{
				  host     = "mail2.pmapitest.com."
				  priority = "20"
				  comments = "test comment"
				},
				{
				  host     = "mail3.pmapitest.com."
				  priority = "30"
				  ttl      = "3600"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.host", "mail1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.priority", "10"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.host", "mail2.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.priority", "20"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.comments", "test comment"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.1.ttl"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.2.host", "mail3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.2.priority", "30"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.2.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.2.comments"),
				),
			},
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  host     = "mail1.pmapitest.com."
				  priority = "10"
				},
				{
				  host     = "mail2.pmapitest.com."
				  priority = "20"
				},
				{
				  host     = "mail3.pmapitest.com."
				  priority = "30"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.host", "mail1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.priority", "10"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.host", "mail2.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.priority", "20"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.2.host", "mail3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.2.priority", "30"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.2.comments"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.2.ttl"),
				),
			},

			// Update and Read testing on subdomain level, remove mx record
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  host     = "mail1.pmapitest.com."
				  priority = "10"
				},
				{
				  host     = "mail3.pmapitest.com."
				  priority = "30"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.host", "mail1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.0.priority", "10"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.comments"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.0.ttl"),

					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.host", "mail3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_mx_record.test", "records.1.priority", "30"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_mx_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsMXRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_mx_record" "test2" {
			  zone  = "non_existing.com"
			  name = "@"
			  records = [
				{
				  host     = "mail1.pmapitest.com."
				  priority = "10"
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
