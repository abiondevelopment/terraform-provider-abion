// Copyright Abion AB 2024, 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"terraform-provider-abion/internal/utils"

	abionclient "github.com/abiondevelopment/abion-go/abion"
)

func createRecordPatchRequest(zoneName string, subDomainOrRoot string, recordType utils.RecordType, data []abionclient.Record) abionclient.ZoneRequest {
	return abionclient.ZoneRequest{
		Data: abionclient.Zone{
			Type: "zone",
			ID:   zoneName,
			Attributes: abionclient.Attributes{
				Records: map[string]map[string][]abionclient.Record{
					subDomainOrRoot: {recordType.String(): data},
				},
			},
		},
	}
}
