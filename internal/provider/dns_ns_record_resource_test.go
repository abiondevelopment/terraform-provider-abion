// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsNSRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_ns_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  nameserver = "ns1.testabiondns.se."
				},
				{
				  nameserver = "ns2.testabiondns.se."
				  ttl = "3600"
				},
				{
				  nameserver = "ns3.testabiondns.se."
				  comments = "test comment"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.0.nameserver", "ns1.testabiondns.se."),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.1.nameserver", "ns2.testabiondns.se."),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.1.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.2.nameserver", "ns3.testabiondns.se."),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.2.ttl"),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.2.comments", "test comment"),
				),
			},

			// Update and Read testing on subdomain level, remove ip address
			{
				Config: providerConfig + `
			resource "abion_dns_ns_record" "test" {
			  zone  = "pmapitest.com"
			  name = "test"
			  records = [
				{
				  nameserver = "ns1.testabiondns.se."
				},
				{
				  nameserver = "ns3.testabiondns.se."
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "zone", "pmapitest.com"),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.0.nameserver", "ns1.testabiondns.se."),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_ns_record.test", "records.1.nameserver", "ns3.testabiondns.se."),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ns_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsNSRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_ns_record" "test2" {
			  zone  = "non_existing.com"
			  name = "test"
			  records = [
				{
				  nameserver = "ns1.testabiondns.se."
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
