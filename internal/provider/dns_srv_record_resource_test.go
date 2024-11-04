// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsSRVRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
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
				},
				{
				  target   = "server2.pmapitest.com."
				  port     = "5103"
				  priority = "100"
				  weight   = "10"
				  comments = "test comment"
				},
				{
				  target   = "server3.pmapitest.com."
				  port     = "443"
				  priority = "100"
				  weight   = "30"
				  ttl      = "3600"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.target", "server1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.priority", "1"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.weight", "100"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.target", "server2.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.port", "5103"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.priority", "100"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.weight", "10"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.comments", "test comment"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.1.ttl"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.target", "server3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.priority", "100"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.weight", "30"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.2.comments"),
				),
			},
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_srv_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  target   = "server1.pmapitest.com."
				  port     = "443"
				  priority = "1"
				  weight   = "100"
				},
				{
				  target   = "server2.pmapitest.com."
				  port     = "5103"
				  priority = "100"
				  weight   = "10"
				},
				{
				  target   = "server3.pmapitest.com."
				  port     = "443"
				  priority = "100"
				  weight   = "30"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.target", "server1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.priority", "1"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.weight", "100"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.target", "server2.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.port", "5103"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.priority", "100"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.weight", "10"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.target", "server3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.priority", "100"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.2.weight", "30"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.2.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.2.comments"),
				),
			},

			// Update and Read testing on subdomain level, remove srv record
			{
				Config: providerConfig + `
			resource "abion_dns_srv_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  target   = "server1.pmapitest.com."
				  port     = "443"
				  priority = "1"
				  weight   = "100"
				},
				{
				  target   = "server3.pmapitest.com."
				  port     = "443"
				  priority = "100"
				  weight   = "30"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.target", "server1.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.priority", "1"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.0.weight", "100"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.target", "server3.pmapitest.com."),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.port", "443"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.priority", "100"),
					resource.TestCheckResourceAttr("abion_dns_srv_record.test", "records.1.weight", "30"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_srv_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsSRVRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_srv_record" "test2" {
			  zone  = "non_existing.com"
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
			`,
				ExpectError: regexp.MustCompile(`(?s)Error patching zone.*Zone not found`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
