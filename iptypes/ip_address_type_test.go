// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

func TestIPAddressTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"valid-ipv4": {
			in:          tftypes.NewValue(tftypes.String, "127.0.0.1"),
			expectation: iptypes.NewIPAddressValue("127.0.0.1"),
		},
		"valid-ipv6": {
			in:          tftypes.NewValue(tftypes.String, "FF01::101"),
			expectation: iptypes.NewIPAddressValue("FF01::101"),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: iptypes.NewIPAddressUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: iptypes.NewIPAddressNull(),
		},
		"wrongType": {
			in:          tftypes.NewValue(tftypes.Number, 123),
			expectedErr: "can't unmarshal tftypes.Number into *string, expected string",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			got, err := iptypes.IPAddressType{}.ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectedErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectedErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectedErr, err.Error())
				}
				return
			}
			if testCase.expectedErr != "" {
				t.Fatalf("Expected error to be %q, didn't get an error", testCase.expectedErr)
			}
			if !got.Equal(testCase.expectation) {
				t.Errorf("Expected %+v, got %+v", testCase.expectation, got)
			}
			if testCase.expectation.IsNull() != testCase.in.IsNull() {
				t.Errorf("Expected null-ness match: expected %t, got %t", testCase.expectation.IsNull(), testCase.in.IsNull())
			}
			if testCase.expectation.IsUnknown() != !testCase.in.IsKnown() {
				t.Errorf("Expected unknown-ness match: expected %t, got %t", testCase.expectation.IsUnknown(), !testCase.in.IsKnown())
			}
		})
	}
}
