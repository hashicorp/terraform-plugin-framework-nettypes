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

func TestIPv6PrefixTypeValidate(t *testing.T) {
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
		"valid IPv6 prefix - unspecified": {
			in: tftypes.NewValue(tftypes.String, "::/128"),
		},
		"valid IPv6 prefix - full": {
			in: tftypes.NewValue(tftypes.String, "2001:0DB8:0:0:0:0:0:0CD3/60"),
		},
		"valid IPv6 prefix - trailing double colon": {
			in: tftypes.NewValue(tftypes.String, "FF00::/8"),
		},
		"valid IPv6 prefix - leading double colon": {
			in: tftypes.NewValue(tftypes.String, "::8:800:200C:417A/32"),
		},
		"valid IPv6 prefix - middle double colon": {
			in: tftypes.NewValue(tftypes.String, "2001:0DB8::CD30/60"),
		},
		"valid IPv6 prefix - lowercase": {
			in: tftypes.NewValue(tftypes.String, "2001:0db8::cd30/60"),
		},
		"valid IPv6 prefix - IPv4-Mapped": {
			in: tftypes.NewValue(tftypes.String, "::FFFF:1.2.3.0/112"),
		},
		"valid IPv6 prefix - IPv4-Compatible": {
			in: tftypes.NewValue(tftypes.String, "::1.2.3.0/112"),
		},
		"invalid IPv6 prefix - invalid address colon end": {
			in: tftypes.NewValue(tftypes.String, "0:0:0:0:0:0:0:/128"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: 0:0:0:0:0:0:0:/128\n"+
						"Error: netip.ParsePrefix(\"0:0:0:0:0:0:0:/128\"): ParseAddr(\"0:0:0:0:0:0:0:\"): colon must be followed by more characters (at \":\")",
				),
			},
		},
		"invalid IPv6 prefix - address too many colons": {
			in: tftypes.NewValue(tftypes.String, "0:0::1::/112"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: 0:0::1::/112\n"+
						"Error: netip.ParsePrefix(\"0:0::1::/112\"): ParseAddr(\"0:0::1::\"): multiple :: in address (at \":\")",
				),
			},
		},
		"invalid IPv6 prefix - address trailing numbers": {
			in: tftypes.NewValue(tftypes.String, "0:0:0:0:0:0:0:1:99/112"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: 0:0:0:0:0:0:0:1:99/112\n"+
						"Error: netip.ParsePrefix(\"0:0:0:0:0:0:0:1:99/112\"): ParseAddr(\"0:0:0:0:0:0:0:1:99\"): trailing garbage after address (at \"99\")",
				),
			},
		},
		"invalid IPv6 prefix - invalid prefix length": {
			in: tftypes.NewValue(tftypes.String, "FF00::/999"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: FF00::/999\n"+
						"Error: netip.ParsePrefix(\"FF00::/999\"): prefix length out of range",
				),
			},
		},
		"invalid IPv6 prefix - invalid prefix characters": {
			in: tftypes.NewValue(tftypes.String, "FF00::/notcorrect"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: FF00::/notcorrect\n"+
						"Error: netip.ParsePrefix(\"FF00::/notcorrect\"): bad bits after slash: \"notcorrect\"",
				),
			},
		},
		"invalid IPv6 prefix - IPv4 prefix": {
			in: tftypes.NewValue(tftypes.String, "172.16.0.0/12"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"An IPv4 CIDR string format was provided, string value must be IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: 172.16.0.0/12\n",
				),
			},
		},
		"wrong-value-type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"IPv6 Prefix Type Validation Error",
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

			diags := cidrtypes.IPv6PrefixType{}.Validate(context.Background(), testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPv6PrefixTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"true": {
			in:          tftypes.NewValue(tftypes.String, "FF00::/8"),
			expectation: cidrtypes.NewIPv6PrefixValue("FF00::/8"),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: cidrtypes.NewIPv6PrefixUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: cidrtypes.NewIPv6PrefixNull(),
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

			got, err := cidrtypes.IPv6PrefixType{}.ValueFromTerraform(ctx, testCase.in)
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
