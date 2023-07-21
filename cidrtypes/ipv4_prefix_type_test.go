// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestIPv4PrefixTypeValidate(t *testing.T) {
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
		"valid IPv4 prefix - loopback": {
			in: tftypes.NewValue(tftypes.String, "127.0.0.0/8"),
		},
		"valid IPv4 prefix - private": {
			in: tftypes.NewValue(tftypes.String, "172.16.0.0/12"),
		},
		"invalid IPv4 prefix - invalid address with leading zeroes": {
			in: tftypes.NewValue(tftypes.String, "127.0.0.000/8"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: 127.0.0.000/8\n"+
						"Error: netip.ParsePrefix(\"127.0.0.000/8\"): ParseAddr(\"127.0.0.000\"): IPv4 field has octet with leading zero",
				),
			},
		},
		"invalid IPv4 prefix - invalid address length": {
			in: tftypes.NewValue(tftypes.String, "127.0.1/8"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: 127.0.1/8\n"+
						"Error: netip.ParsePrefix(\"127.0.1/8\"): ParseAddr(\"127.0.1\"): IPv4 address too short",
				),
			},
		},
		"invalid IPv4 prefix - invalid prefix length": {
			in: tftypes.NewValue(tftypes.String, "127.0.0.0/999"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: 127.0.0.0/999\n"+
						"Error: netip.ParsePrefix(\"127.0.0.0/999\"): prefix length out of range",
				),
			},
		},
		"invalid IPv4 prefix - invalid prefix characters": {
			in: tftypes.NewValue(tftypes.String, "127.0.0.0/notcorrect"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: 127.0.0.0/notcorrect\n"+
						"Error: netip.ParsePrefix(\"127.0.0.0/notcorrect\"): bad bits after slash: \"notcorrect\"",
				),
			},
		},
		"invalid IPv4 prefix - IPv6 prefix": {
			in: tftypes.NewValue(tftypes.String, "::/128"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"An IPv6 CIDR string format was provided, string value must be IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: ::/128\n",
				),
			},
		},
		"wrong-value-type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"IPv4 Prefix Type Validation Error",
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

			diags := cidrtypes.IPv4PrefixType{}.Validate(context.Background(), testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
			}
		})
	}
}

func TestIPv4PrefixTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"true": {
			in:          tftypes.NewValue(tftypes.String, "172.16.0.0/12"),
			expectation: cidrtypes.NewIPv4PrefixValue("172.16.0.0/12"),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: cidrtypes.NewIPv4PrefixUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: cidrtypes.NewIPv4PrefixNull(),
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

			got, err := cidrtypes.IPv4PrefixType{}.ValueFromTerraform(ctx, testCase.in)
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
