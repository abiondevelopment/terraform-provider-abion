// Copyright (c) Abion AB
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IntPointerToInt32 convert *int to types.Int32.
func IntPointerToInt32(input *int) types.Int32 {
	if input == nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(*input))
}

// Int32ToIntPointer convert types.Int32 to *int.
func Int32ToIntPointer(input types.Int32) *int {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}
	value := int(input.ValueInt32())
	return &value
}

func StringPointerToTerraformString(input *string) types.String {
	if input == nil {
		return types.StringNull()
	}
	return types.StringPointerValue(input)
}
