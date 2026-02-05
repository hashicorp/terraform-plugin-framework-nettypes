// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"context"
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
)

func TestIPv6PrefixStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentIpPrefix cidrtypes.IPv6Prefix
		givenIpPrefix   basetypes.StringValuable
		expectedMatch   bool
		expectedDiags   diag.Diagnostics
	}{
		"not equal - IPv6 prefix mismatch": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:1/128"),
			expectedMatch:   false,
		},
		"not equal - IPv6 prefix compressed mismatch": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("FF01::/8"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("FF01::1/8"),
			expectedMatch:   false,
		},
		"not equal - IPv4-Mapped IPv6 prefix mismatch": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.0/112"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.4/112"),
			expectedMatch:   false,
		},
		"semantically equal - byte-for-byte match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
			expectedMatch:   true,
		},
		"semantically equal - case insensitive": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0000:0000:0008:0800:200C:417A/60"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("2001:0db8:0000:0000:0008:0800:200c:417a/60"),
			expectedMatch:   true,
		},
		"semantically equal - IPv4-Mapped byte-for-byte match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.4/112"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.4/112"),
			expectedMatch:   true,
		},
		"semantically equal - compressed all zeroes match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("::/128"),
			expectedMatch:   true,
		},
		"semantically equal - compressed all leading zeroes match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0000:0000:0008:0800:200C:417A/60"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("2001:DB8::8:800:200C:417A/60"),
			expectedMatch:   true,
		},
		"semantically equal - start compressed match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("::101/128"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:101/128"),
			expectedMatch:   true,
		},
		"semantically equal - middle compressed match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("2001:DB8::8:800:200C:417A/60"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("2001:DB8:0:0:8:800:200C:417A/60"),
			expectedMatch:   true,
		},
		"semantically equal - end compressed match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("FF01:0:0:0:0:0:0:0/8"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("FF01::/8"),
			expectedMatch:   true,
		},
		"semantically equal - IPv4-Mapped compressed match": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:FFFF:1.2.3.4/112"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.4/112"),
			expectedMatch:   true,
		},
		"error - not given IPv6Prefix IPv6 value": {
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
			givenIpPrefix:   basetypes.NewStringValue("0:0:0:0:0:0:0:0/128"),
			expectedMatch:   false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: cidrtypes.IPv6Prefix\n"+
						"Got Value Type: basetypes.StringValue",
				),
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentIpPrefix.StringSemanticEquals(context.Background(), testCase.givenIpPrefix)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPv6PrefixValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefixValue   cidrtypes.IPv6Prefix
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			prefixValue: cidrtypes.IPv6Prefix{},
		},
		"null": {
			prefixValue: cidrtypes.NewIPv6PrefixNull(),
		},
		"unknown": {
			prefixValue: cidrtypes.NewIPv6PrefixUnknown(),
		},
		"valid IPv6 prefix - unspecified": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::/128"),
		},
		"valid IPv6 prefix - full": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0:0:0:0:0:0CD3/60"),
		},
		"valid IPv6 prefix - trailing double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00:0:0:0:0:0:0:0/8"),
		},
		"valid IPv6 prefix - leading double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
		},
		"valid IPv6 prefix - middle double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0:0:0:0:0:0CD3/60"),
		},
		"valid IPv6 prefix - lowercase": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0db8:0:0:0:0:0:0cd3/60"),
		},
		"valid IPv6 prefix - IPv4-Mapped": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.0/112"),
		},
		"valid IPv6 prefix - IPv4-Compatible": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::1.2.3.0/112"),
		},
		"invalid IPv6 prefix - invalid address colon end": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:/128"),
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
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0::1::/112"),
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
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:1:99/112"),
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
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00::/999"),
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
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00::/notcorrect"),
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
			prefixValue: cidrtypes.NewIPv6PrefixValue("172.16.0.0/12"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv6 CIDR String Value",
					"An IPv4 CIDR string format was provided, string value must be IPv6 CIDR string format (RFC 4291).\n\n"+
						"Given Value: 172.16.0.0/12\n",
				),
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.prefixValue.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{
					Path: path.Root("test"),
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPv6PrefixValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefixValue     cidrtypes.IPv6Prefix
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			prefixValue: cidrtypes.IPv6Prefix{},
		},
		"null": {
			prefixValue: cidrtypes.NewIPv6PrefixNull(),
		},
		"unknown": {
			prefixValue: cidrtypes.NewIPv6PrefixUnknown(),
		},
		"valid IPv6 prefix - unspecified": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::/128"),
		},
		"valid IPv6 prefix - full": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0:0:0:0:0:0CD3/60"),
		},
		"valid IPv6 prefix - trailing double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00:0:0:0:0:0:0:0/8"),
		},
		"valid IPv6 prefix - leading double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:0/128"),
		},
		"valid IPv6 prefix - middle double colon": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0:0:0:0:0:0CD3/60"),
		},
		"valid IPv6 prefix - lowercase": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("2001:0db8:0:0:0:0:0:0cd3/60"),
		},
		"valid IPv6 prefix - IPv4-Mapped": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::FFFF:1.2.3.0/112"),
		},
		"valid IPv6 prefix - IPv4-Compatible": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("::1.2.3.0/112"),
		},
		"invalid IPv6 prefix - invalid address colon end": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:/128"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: 0:0:0:0:0:0:0:/128\n"+
					"Error: netip.ParsePrefix(\"0:0:0:0:0:0:0:/128\"): ParseAddr(\"0:0:0:0:0:0:0:\"): colon must be followed by more characters (at \":\")",
			),
		},
		"invalid IPv6 prefix - address too many colons": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0::1::/112"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: 0:0::1::/112\n"+
					"Error: netip.ParsePrefix(\"0:0::1::/112\"): ParseAddr(\"0:0::1::\"): multiple :: in address (at \":\")",
			),
		},
		"invalid IPv6 prefix - address trailing numbers": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("0:0:0:0:0:0:0:1:99/112"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: 0:0:0:0:0:0:0:1:99/112\n"+
					"Error: netip.ParsePrefix(\"0:0:0:0:0:0:0:1:99/112\"): ParseAddr(\"0:0:0:0:0:0:0:1:99\"): trailing garbage after address (at \"99\")",
			),
		},
		"invalid IPv6 prefix - invalid prefix length": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00::/999"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: FF00::/999\n"+
					"Error: netip.ParsePrefix(\"FF00::/999\"): prefix length out of range",
			),
		},
		"invalid IPv6 prefix - invalid prefix characters": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("FF00::/notcorrect"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: FF00::/notcorrect\n"+
					"Error: netip.ParsePrefix(\"FF00::/notcorrect\"): bad bits after slash: \"notcorrect\"",
			),
		},
		"invalid IPv6 prefix - IPv4 prefix": {
			prefixValue: cidrtypes.NewIPv6PrefixValue("172.16.0.0/12"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv6 CIDR String Value: "+
					"An IPv4 CIDR string format was provided, string value must be IPv6 CIDR string format (RFC 4291).\n\n"+
					"Given Value: 172.16.0.0/12\n",
			),
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.prefixValue.ValidateParameter(
				context.Background(),
				function.ValidateParameterRequest{
					Position: 0,
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Error, testCase.expectedFuncErr); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPv6PrefixValueIPv6Prefix(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipValue          cidrtypes.IPv6Prefix
		expectedIpPrefix netip.Prefix
		expectedDiags    diag.Diagnostics
	}{
		"IPv6 prefix value is null ": {
			ipValue: cidrtypes.NewIPv6PrefixNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv6Prefix ValueIPv6Prefix Error",
					"IPv6 CIDR string value is null",
				),
			},
		},
		"IPv6 prefix value is unknown ": {
			ipValue: cidrtypes.NewIPv6PrefixUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv6Prefix ValueIPv6Prefix Error",
					"IPv6 CIDR string value is unknown",
				),
			},
		},
		"valid IPv6 prefix ": {
			ipValue:          cidrtypes.NewIPv6PrefixValue("2001:0DB8::CD30/60"),
			expectedIpPrefix: netip.MustParsePrefix("2001:0DB8::CD30/60"),
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipAddr, diags := testCase.ipValue.ValueIPv6Prefix()

			if ipAddr != testCase.expectedIpPrefix {
				t.Errorf("Unexpected difference in netip.Prefix, got: %s, expected: %s", ipAddr, testCase.expectedIpPrefix)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
