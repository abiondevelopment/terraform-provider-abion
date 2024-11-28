// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsPTRRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing on root level
			{
				Config: providerConfig + `
			resource "abion_dns_ptr_record" "test" {
			  zone  = "pmapitest6.com"
			  name = "@"
			  records = [
				{
				  ptr = "www.example0.com."
				},
				{
				  ptr = "www.example1.com."
				  ttl = "3600"
				},
				{
				  ptr = "www.example2.com."
				  comments = "test comment"
				}
			  ]
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "zone", "pmapitest6.com"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "name", "@"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.0.ptr", "www.example0.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.1.ptr", "www.example1.com."),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.1.ttl", "3600"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.2.ptr", "www.example2.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.2.ttl"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.2.comments", "test comment"),
				),
			},
			// Create and Read testing on subdomain level, move records from root to subdomain
			{
				Config: providerConfig + `
			resource "abion_dns_ptr_record" "test" {
			  zone  = "pmapitest6.com"
			  name = "test"
			  records = [
				{
				  ptr = "www.example0.com."
				},
				{
				  ptr = "www.example1.com."
				},
				{
				  ptr = "www.example2.com."
				}
			  ]	
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "zone", "pmapitest6.com"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.#", "3"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.0.ptr", "www.example0.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.1.ptr", "www.example1.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.1.comments"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.2.ptr", "www.example2.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.2.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.2.comments"),
				),
			},

			// Update and Read testing on subdomain level, remove ip address
			{
				Config: providerConfig + `
			resource "abion_dns_ptr_record" "test" {
			  zone  = "pmapitest6.com"
			  name = "test"
			  records = [
				{
				  ptr = "www.example0.com."
				},
				{
				  ptr = "www.example2.com."
				}
			  ]	
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "zone", "pmapitest6.com"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "name", "test"),
					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.#", "2"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.0.ptr", "www.example0.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.0.comments"),

					resource.TestCheckResourceAttr("abion_dns_ptr_record.test", "records.1.ptr", "www.example2.com."),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.1.ttl"),
					resource.TestCheckNoResourceAttr("abion_dns_ptr_record.test", "records.1.comments"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsPTRRecordNonExistingZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify error non existing zone
			{
				Config: providerConfig + `
			resource "abion_dns_ptr_record" "test2" {
			  zone  = "non_existing.com"
			  name = "@"
			  records = [
				{
				  ptr = "www.example0.com."
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
