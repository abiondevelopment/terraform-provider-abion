// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package utils

// RecordType defines a type for DNS record types.
type RecordType string

// Constants for DNS record types.
const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCName RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeTXT   RecordType = "TXT"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypePTR   RecordType = "PTR"
)

// String method to convert the RecordType to a string.
func (r RecordType) String() string {
	return string(r)
}
