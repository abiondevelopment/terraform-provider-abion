// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsCNameRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
			{
				Config: providerConfig + `
			resource "abion_dns_cname_record" "test" {
			  zone  = "pmapitest3.com"
			  name = "@"
			  record = {
				  cname = "www.test.com."
				  ttl = "3600"
				  comments = "test comment"
			  }
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "zone", "pmapitest3.com"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "record.cname", "www.test.com."),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "record.ttl", "3600"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "record.comments", "test comment"),
				),
			},
			// Create and Read testing on subdomain level, move record from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_cname_record" "test" {
			  zone  = "pmapitest3.com"
			  name = "test"
			  record = {
				  cname = "www.test.com."
			  }
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "zone", "pmapitest3.com"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "record.cname", "www.test.com."),
					resource.TestCheckNoResourceAttr("abion_dns_cname_record.test", "record.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_cname_record.test", "record.comments"),
				),
			},

			// Update and Read testing on subdomain level, change cname
			{
				Config: providerConfig + `
			resource "abion_dns_cname_record" "test" {
			  zone  = "pmapitest3.com"
			  name = "test"
			  record = {
				  cname = "www.test3.com."
			  }
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "zone", "pmapitest3.com"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_cname_record.test", "record.cname", "www.test3.com."),
					resource.TestCheckNoResourceAttr("abion_dns_cname_record.test", "record.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_cname_record.test", "record.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsCNameRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_cname_record" "test2" {
			  zone  = "non_existing.com"
			  name = "@"
			  record = {
				  cname = "still_non_existing.com"
			  }
			}
			`,
				ExpectError: regexp.MustCompile(`(?s)Error patching zone.*Zone not found`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
