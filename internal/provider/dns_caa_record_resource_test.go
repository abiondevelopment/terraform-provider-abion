// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsCAARecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test" {
			  zone  = "pmapitest9.com"
			  name = "@"
			  records = [
				{
				  flag   = "0"
				  tag    = "iodef"
				  value  = "mailto:webmaster@test.test"
				  ttl    = "3600"
				},
				{
				  flag   = "0"
				  tag    = "issue"
				  value  = "letsencrypt.test"
				},
				{
				  flag     = "0"
				  tag      = "issuewild"
				  value    = ";"
				  comments = "test comment"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "zone", "pmapitest9.com"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.tag", "iodef"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.value", "mailto:webmaster@test.test"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.tag", "issue"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.value", "letsencrypt.test"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.tag", "issuewild"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.value", ";"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.comments", "test comment"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.2.ttl"),
				),
			},
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test" {
			  zone  = "pmapitest9.com"
			  name = "test"
			  records = [
				{
				  flag   = "0"
				  tag    = "iodef"
				  value  = "mailto:webmaster@test.test"
				},
				{
				  flag   = "0"
				  tag    = "issue"
				  value  = "letsencrypt.test"
				},
				{
				  flag   = "0"
				  tag    = "issuewild"
				  value  = ";"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "zone", "pmapitest9.com"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.tag", "iodef"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.value", "mailto:webmaster@test.test"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.tag", "issue"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.value", "letsencrypt.test"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.tag", "issuewild"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.2.value", ";"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.2.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.2.comments"),
				),
			},

			// Update and Read testing on subdomain level, remove caa record
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test" {
			  zone  = "pmapitest9.com"
			  name = "test"
			  records = [
				{
				  flag  = "0"
				  tag   = "iodef"
				  value = "mailto:webmaster@test.test"
				},
				{
				  flag  = "0"
				  tag   = "issue"
				  value = "letsencrypt.test"
				},
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "zone", "pmapitest9.com"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.tag", "iodef"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.0.value", "mailto:webmaster@test.test"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.flag", "0"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.tag", "issue"),
					resource.TestCheckResourceAttr("abion_dns_caa_record.test", "records.1.value", "letsencrypt.test"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_caa_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsCAARecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_caa_record" "test2" {
			  zone  = "non_existing.com"
			  name = "@"
			  records = [
				{
				  flag  = "0"
				  tag   = "issue"
				  value = "letsencrypt.test"
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
