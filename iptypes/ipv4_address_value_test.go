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

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

func TestIPv4AddressValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue  iptypes.IPv4Address
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			addressValue: iptypes.IPv4Address{},
		},
		"null": {
			addressValue: iptypes.NewIPv4AddressNull(),
		},
		"unknown": {
			addressValue: iptypes.NewIPv4AddressUnknown(),
		},
		"valid IPv4 address - localhost": {
			addressValue: iptypes.NewIPv4AddressValue("127.0.0.1"),
		},
		"valid IPv4 address - private": {
			addressValue: iptypes.NewIPv4AddressValue("192.168.255.255"),
		},
		"invalid IPv4 address - no dots": {
			addressValue: iptypes.NewIPv4AddressValue("192168255255"),
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
			addressValue: iptypes.NewIPv4AddressValue("127.000.000.001"),
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
			addressValue: iptypes.NewIPv4AddressValue("127.0.A.1"),
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
			addressValue: iptypes.NewIPv4AddressValue("127.0.1"),
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
			addressValue: iptypes.NewIPv4AddressValue("::"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 Address String Value",
					"An IPv6 string format was provided, string value must be IPv4 format.\n\n"+
						"Given Value: ::\n",
				),
			},
		},
	}

	for name, testCase := range testCases {

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

func TestIPv4AddressValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue    iptypes.IPv4Address
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			addressValue: iptypes.IPv4Address{},
		},
		"null": {
			addressValue: iptypes.NewIPv4AddressNull(),
		},
		"unknown": {
			addressValue: iptypes.NewIPv4AddressUnknown(),
		},
		"valid IPv4 address - localhost": {
			addressValue: iptypes.NewIPv4AddressValue("127.0.0.1"),
		},
		"valid IPv4 address - private": {
			addressValue: iptypes.NewIPv4AddressValue("192.168.255.255"),
		},
		"invalid IPv4 address - no dots": {
			addressValue: iptypes.NewIPv4AddressValue("192168255255"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 Address String Value: "+
					"A string value was provided that is not valid IPv4 string format.\n\n"+
					"Given Value: 192168255255\n"+
					"Error: ParseAddr(\"192168255255\"): unable to parse IP",
			),
		},
		"invalid IPv4 address - leading zeroes": {
			addressValue: iptypes.NewIPv4AddressValue("127.000.000.001"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 Address String Value: "+
					"A string value was provided that is not valid IPv4 string format.\n\n"+
					"Given Value: 127.000.000.001\n"+
					"Error: ParseAddr(\"127.000.000.001\"): IPv4 field has octet with leading zero",
			),
		},
		"invalid IPv4 address - invalid characters": {
			addressValue: iptypes.NewIPv4AddressValue("127.0.A.1"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 Address String Value: "+
					"A string value was provided that is not valid IPv4 string format.\n\n"+
					"Given Value: 127.0.A.1\n"+
					"Error: ParseAddr(\"127.0.A.1\"): unexpected character (at \"A.1\")",
			),
		},
		"invalid IPv4 address - invalid length": {
			addressValue: iptypes.NewIPv4AddressValue("127.0.1"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 Address String Value: "+
					"A string value was provided that is not valid IPv4 string format.\n\n"+
					"Given Value: 127.0.1\n"+
					"Error: ParseAddr(\"127.0.1\"): IPv4 address too short",
			),
		},
		"invalid IPv4 address - IPv6 address": {
			addressValue: iptypes.NewIPv4AddressValue("::"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 Address String Value: "+
					"An IPv6 string format was provided, string value must be IPv4 format.\n\n"+
					"Given Value: ::\n",
			),
		},
	}

	for name, testCase := range testCases {

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

func TestIPv4AddressValueIPv4Address(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipValue        iptypes.IPv4Address
		expectedIpAddr netip.Addr
		expectedDiags  diag.Diagnostics
	}{
		"IPv4 address value is null ": {
			ipValue: iptypes.NewIPv4AddressNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv4Address ValueIPv4Address Error",
					"IPv4 address string value is null",
				),
			},
		},
		"IPv4 address value is unknown ": {
			ipValue: iptypes.NewIPv4AddressUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv4Address ValueIPv4Address Error",
					"IPv4 address string value is unknown",
				),
			},
		},
		"valid IPv4 address ": {
			ipValue:        iptypes.NewIPv4AddressValue("127.0.0.1"),
			expectedIpAddr: netip.MustParseAddr("127.0.0.1"),
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipAddr, diags := testCase.ipValue.ValueIPv4Address()

			if ipAddr != testCase.expectedIpAddr {
				t.Errorf("Unexpected difference in netip.Addr, got: %s, expected: %s", ipAddr, testCase.expectedIpAddr)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
