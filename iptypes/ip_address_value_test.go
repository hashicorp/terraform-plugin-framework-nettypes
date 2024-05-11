// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

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

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

func TestIPAddressStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentIpAddr iptypes.IPAddress
		givenIpAddr   basetypes.StringValuable
		expectedMatch bool
		expectedDiags diag.Diagnostics
	}{
		"not equal - IPv6 address mismatch": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPAddressValue("0:0:0:0:0:0:0:1"),
			expectedMatch: false,
		},
		"not equal - IPv6 address compressed mismatch": {
			currentIpAddr: iptypes.NewIPAddressValue("FF01::"),
			givenIpAddr:   iptypes.NewIPAddressValue("FF01::1"),
			expectedMatch: false,
		},
		"not equal - IPv4-Mapped IPv6 address mismatch": {
			currentIpAddr: iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPAddressValue("::FFFF:192.168.255.254"),
			expectedMatch: false,
		},
		"semantically equal - byte-for-byte match": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPAddressValue("0:0:0:0:0:0:0:0"),
			expectedMatch: true,
		},
		"semantically equal - case insensitive": {
			currentIpAddr: iptypes.NewIPAddressValue("2001:0DB8:0000:0000:0008:0800:0200C:417A"),
			givenIpAddr:   iptypes.NewIPAddressValue("2001:0db8:0000:0000:0008:0800:0200c:417a"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped byte-for-byte match": {
			currentIpAddr: iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
			expectedMatch: true,
		},
		"semantically equal - compressed all zeroes match": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPAddressValue("::"),
			expectedMatch: true,
		},
		"semantically equal - compressed all leading zeroes match": {
			currentIpAddr: iptypes.NewIPAddressValue("2001:0DB8:0000:0000:0008:0800:0200C:417A"),
			givenIpAddr:   iptypes.NewIPAddressValue("2001:DB8::8:800:200C:417A"),
			expectedMatch: true,
		},
		"semantically equal - start compressed match": {
			currentIpAddr: iptypes.NewIPAddressValue("::101"),
			givenIpAddr:   iptypes.NewIPAddressValue("0:0:0:0:0:0:0:101"),
			expectedMatch: true,
		},
		"semantically equal - middle compressed match": {
			currentIpAddr: iptypes.NewIPAddressValue("2001:DB8::8:800:200C:417A"),
			givenIpAddr:   iptypes.NewIPAddressValue("2001:DB8:0:0:8:800:200C:417A"),
			expectedMatch: true,
		},
		"semantically equal - end compressed match": {
			currentIpAddr: iptypes.NewIPAddressValue("FF01:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPAddressValue("FF01::"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped compressed match": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped match": {
			currentIpAddr: iptypes.NewIPAddressValue("FF01:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPAddressValue("FF01::"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped mismatch": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
			expectedMatch: true,
		},
		"semantically equal - IPv4 match": {
			currentIpAddr: iptypes.NewIPAddressValue("104.28.204.175"),
			givenIpAddr:   iptypes.NewIPAddressValue("104.28.204.175"),
			expectedMatch: true,
		},
		"semantically equal - IPv4 mismatch": {
			currentIpAddr: iptypes.NewIPAddressValue("104.28.204.175"),
			givenIpAddr:   iptypes.NewIPAddressValue("104.28.204.174"),
			expectedMatch: false,
		},
		"error - not given IPAddress IPv6 value": {
			currentIpAddr: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   basetypes.NewStringValue("0:0:0:0:0:0:0:0"),
			expectedMatch: false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: iptypes.IPAddress\n"+
						"Got Value Type: basetypes.StringValue",
				),
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentIpAddr.StringSemanticEquals(context.Background(), testCase.givenIpAddr)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPAddressValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue  iptypes.IPAddress
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			addressValue: iptypes.IPAddress{},
		},
		"null": {
			addressValue: iptypes.NewIPAddressNull(),
		},
		"unknown": {
			addressValue: iptypes.NewIPAddressUnknown(),
		},
		"valid IPv4 address - localhost": {
			addressValue: iptypes.NewIPAddressValue("127.0.0.1"),
		},
		"valid IPv4 address - private": {
			addressValue: iptypes.NewIPAddressValue("192.168.255.255"),
		},
		"valid IPv6 address - unspecified": {
			addressValue: iptypes.NewIPAddressValue("::"),
		},
		"valid IPv6 address - full": {
			addressValue: iptypes.NewIPAddressValue("1:2:3:4:5:6:7:8"),
		},
		"valid IPv6 address - trailing double colon": {
			addressValue: iptypes.NewIPAddressValue("FF01::"),
		},
		"valid IPv6 address - leading double colon": {
			addressValue: iptypes.NewIPAddressValue("::8:800:200C:417A"),
		},
		"valid IPv6 address - middle double colon": {
			addressValue: iptypes.NewIPAddressValue("2001:DB8::8:800:200C:417A"),
		},
		"valid IPv6 address - lowercase": {
			addressValue: iptypes.NewIPAddressValue("2001:db8::8:800:200c:417a"),
		},
		"valid IPv6 address - IPv4-Mapped": {
			addressValue: iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
		},
		"valid IPv6 address - IPv4-Compatible": {
			addressValue: iptypes.NewIPAddressValue("::127.0.0.1"),
		},
		"invalid IPv4 address - no dots": {
			addressValue: iptypes.NewIPAddressValue("192168255255"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 192168255255\n"+
						"Error: ParseAddr(\"192168255255\"): unable to parse IP",
				),
			},
		},
		"invalid IPv4 address - leading zeroes": {
			addressValue: iptypes.NewIPAddressValue("127.000.000.001"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 127.000.000.001\n"+
						"Error: ParseAddr(\"127.000.000.001\"): IPv4 field has octet with leading zero",
				),
			},
		},
		"invalid IPv4 address - invalid characters": {
			addressValue: iptypes.NewIPAddressValue("127.0.A.1"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 127.0.A.1\n"+
						"Error: ParseAddr(\"127.0.A.1\"): unexpected character (at \"A.1\")",
				),
			},
		},
		"invalid IPv4 address - invalid length": {
			addressValue: iptypes.NewIPAddressValue("127.0.1"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 127.0.1\n"+
						"Error: ParseAddr(\"127.0.1\"): IPv4 address too short",
				),
			},
		},
		"invalid IPv6 address - invalid colon end": {
			addressValue: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 0:0:0:0:0:0:0:\n"+
						"Error: ParseAddr(\"0:0:0:0:0:0:0:\"): colon must be followed by more characters (at \":\")",
				),
			},
		},
		"invalid IPv6 address - too many colons": {
			addressValue: iptypes.NewIPAddressValue("0:0::1::"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 0:0::1::\n"+
						"Error: ParseAddr(\"0:0::1::\"): multiple :: in address (at \":\")",
				),
			},
		},
		"invalid IPv6 address - trailing numbers": {
			addressValue: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:1:99"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IP Address String Value",
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
						"Given Value: 0:0:0:0:0:0:0:1:99\n"+
						"Error: ParseAddr(\"0:0:0:0:0:0:0:1:99\"): trailing garbage after address (at \"99\")",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.addressValue.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{Path: path.Root("test")},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestIPAddressValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue    iptypes.IPAddress
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			addressValue: iptypes.IPAddress{},
		},
		"null": {
			addressValue: iptypes.NewIPAddressNull(),
		},
		"unknown": {
			addressValue: iptypes.NewIPAddressUnknown(),
		},
		"valid IPv4 address - localhost": {
			addressValue: iptypes.NewIPAddressValue("127.0.0.1"),
		},
		"valid IPv4 address - private": {
			addressValue: iptypes.NewIPAddressValue("192.168.255.255"),
		},
		"valid IPv6 address - unspecified": {
			addressValue: iptypes.NewIPAddressValue("::"),
		},
		"valid IPv6 address - full": {
			addressValue: iptypes.NewIPAddressValue("1:2:3:4:5:6:7:8"),
		},
		"valid IPv6 address - trailing double colon": {
			addressValue: iptypes.NewIPAddressValue("FF01::"),
		},
		"valid IPv6 address - leading double colon": {
			addressValue: iptypes.NewIPAddressValue("::8:800:200C:417A"),
		},
		"valid IPv6 address - middle double colon": {
			addressValue: iptypes.NewIPAddressValue("2001:DB8::8:800:200C:417A"),
		},
		"valid IPv6 address - lowercase": {
			addressValue: iptypes.NewIPAddressValue("2001:db8::8:800:200c:417a"),
		},
		"valid IPv6 address - IPv4-Mapped": {
			addressValue: iptypes.NewIPAddressValue("::FFFF:192.168.255.255"),
		},
		"valid IPv6 address - IPv4-Compatible": {
			addressValue: iptypes.NewIPAddressValue("::127.0.0.1"),
		},
		"invalid IPv4 address - no dots": {
			addressValue: iptypes.NewIPAddressValue("192168255255"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 192168255255\n"+
					"Error: ParseAddr(\"192168255255\"): unable to parse IP",
			),
		},
		"invalid IPv4 address - leading zeroes": {
			addressValue: iptypes.NewIPAddressValue("127.000.000.001"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 127.000.000.001\n"+
					"Error: ParseAddr(\"127.000.000.001\"): IPv4 field has octet with leading zero",
			),
		},
		"invalid IPv4 address - invalid characters": {
			addressValue: iptypes.NewIPAddressValue("127.0.A.1"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 127.0.A.1\n"+
					"Error: ParseAddr(\"127.0.A.1\"): unexpected character (at \"A.1\")",
			),
		},
		"invalid IPv4 address - invalid length": {
			addressValue: iptypes.NewIPAddressValue("127.0.1"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 127.0.1\n"+
					"Error: ParseAddr(\"127.0.1\"): IPv4 address too short",
			),
		},
		"invalid IPv6 address - invalid colon end": {
			addressValue: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 0:0:0:0:0:0:0:\n"+
					"Error: ParseAddr(\"0:0:0:0:0:0:0:\"): colon must be followed by more characters (at \":\")",
			),
		},
		"invalid IPv6 address - too many colons": {
			addressValue: iptypes.NewIPAddressValue("0:0::1::"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 0:0::1::\n"+
					"Error: ParseAddr(\"0:0::1::\"): multiple :: in address (at \":\")",
			),
		},
		"invalid IPv6 address - trailing numbers": {
			addressValue: iptypes.NewIPAddressValue("0:0:0:0:0:0:0:1:99"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IP Address String Value: "+
					"A string value was provided that is not valid IPv4 or IPv6 string format (RFC 791, RFC 4291).\n\n"+
					"Given Value: 0:0:0:0:0:0:0:1:99\n"+
					"Error: ParseAddr(\"0:0:0:0:0:0:0:1:99\"): trailing garbage after address (at \"99\")",
			),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.addressValue.ValidateParameter(
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

func TestIPAddressValueIPAddress_ipv4(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipValue        iptypes.IPAddress
		expectedIpAddr netip.Addr
		expectedDiags  diag.Diagnostics
	}{
		"IP address value is null ": {
			ipValue: iptypes.NewIPAddressNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPAddress ValueIPAddress Error",
					"IP address string value is null",
				),
			},
		},
		"IP address value is unknown ": {
			ipValue: iptypes.NewIPAddressUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPAddress ValueIPAddress Error",
					"IP address string value is unknown",
				),
			},
		},
		"valid IPv6 address ": {
			ipValue:        iptypes.NewIPAddressValue("2001:DB8::8:800:200C:417A"),
			expectedIpAddr: netip.MustParseAddr("2001:DB8::8:800:200C:417A"),
		},
		"valid IPv4 address ": {
			ipValue:        iptypes.NewIPAddressValue("127.0.0.1"),
			expectedIpAddr: netip.MustParseAddr("127.0.0.1"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipAddr, diags := testCase.ipValue.ValueIPAddress()

			if ipAddr != testCase.expectedIpAddr {
				t.Errorf("Unexpected difference in netip.Addr, got: %s, expected: %s", ipAddr, testCase.expectedIpAddr)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
