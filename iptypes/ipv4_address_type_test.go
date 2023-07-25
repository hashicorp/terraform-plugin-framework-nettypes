// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestIPv4AddressTypeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            tftypes.Value
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			in: tftypes.Value{},
		},
		"null": {
			in: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"valid IPv4 address - localhost": {
			in: tftypes.NewValue(tftypes.String, "127.0.0.1"),
		},
		"valid IPv4 address - private": {
			in: tftypes.NewValue(tftypes.String, "192.168.255.255"),
		},
		"invalid IPv4 address - no dots": {
			in: tftypes.NewValue(tftypes.String, "192168255255"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"A string value was provided that is not valid IPv4 string format.\n\n"+
						"Given Value: 192168255255\n"+
						"Error: ParseAddr(\"192168255255\"): unable to parse IP",
				),
			},
		},
		"invalid IPv4 address - leading zeroes": {
			in: tftypes.NewValue(tftypes.String, "127.000.000.001"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"A string value was provided that is not valid IPv4 string format.\n\n"+
						"Given Value: 127.000.000.001\n"+
						"Error: ParseAddr(\"127.000.000.001\"): IPv4 field has octet with leading zero",
				),
			},
		},
		"invalid IPv4 address - invalid characters": {
			in: tftypes.NewValue(tftypes.String, "127.0.A.1"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"A string value was provided that is not valid IPv4 string format.\n\n"+
						"Given Value: 127.0.A.1\n"+
						"Error: ParseAddr(\"127.0.A.1\"): unexpected character (at \"A.1\")",
				),
			},
		},
		"invalid IPv4 address - invalid length": {
			in: tftypes.NewValue(tftypes.String, "127.0.1"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"A string value was provided that is not valid IPv4 string format.\n\n"+
						"Given Value: 127.0.1\n"+
						"Error: ParseAddr(\"127.0.1\"): IPv4 address too short",
				),
			},
		},
		"invalid IPv4 address - IPv6 address": {
			in: tftypes.NewValue(tftypes.String, "::"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"An IPv6 string format was provided, string value must be IPv4 format.\n\n"+
						"Given Value: ::\n",
				),
			},
		},
		"wrong-value-type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"IPv4 Address Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
						"expected String value, received tftypes.Value with value: tftypes.Number<\"123\">",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := iptypes.IPv4AddressType{}.Validate(context.Background(), testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
			}
		})
	}
}

func TestIPv4AddressTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"true": {
			in:          tftypes.NewValue(tftypes.String, "127.0.0.1"),
			expectation: iptypes.NewIPv4AddressValue("127.0.0.1"),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: iptypes.NewIPv4AddressUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: iptypes.NewIPv4AddressNull(),
		},
		"wrongType": {
			in:          tftypes.NewValue(tftypes.Number, 123),
			expectedErr: "can't unmarshal tftypes.Number into *string, expected string",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			got, err := iptypes.IPv4AddressType{}.ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectedErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectedErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectedErr, err.Error())
				}
				return
			}
			if err == nil && testCase.expectedErr != "" {
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
