// Copyright (c) HashiCorp, Inc.
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

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
)

func TestIPv4PrefixValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefixValue   cidrtypes.IPv4Prefix
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			prefixValue: cidrtypes.IPv4Prefix{},
		},
		"null": {
			prefixValue: cidrtypes.NewIPv4PrefixNull(),
		},
		"unknown": {
			prefixValue: cidrtypes.NewIPv4PrefixUnknown(),
		},
		"valid IPv4 prefix - loopback": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/8"),
		},
		"valid IPv4 prefix - private": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("172.16.0.0/12"),
		},
		"invalid IPv4 prefix - invalid address no dots": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("192168255255/8"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: 192168255255/8\n"+
						"Error: netip.ParsePrefix(\"192168255255/8\"): ParseAddr(\"192168255255\"): unable to parse IP",
				),
			},
		},
		"invalid IPv4 prefix - invalid address with leading zeroes": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.000/8"),
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
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.1/8"),
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
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/999"),
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
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/notcorrect"),
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
			prefixValue: cidrtypes.NewIPv4PrefixValue("::/128"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid IPv4 CIDR String Value",
					"An IPv6 CIDR string format was provided, string value must be IPv4 CIDR string format (RFC 4632).\n\n"+
						"Given Value: ::/128\n",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
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

func TestIPv4PrefixValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefixValue     cidrtypes.IPv4Prefix
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			prefixValue: cidrtypes.IPv4Prefix{},
		},
		"null": {
			prefixValue: cidrtypes.NewIPv4PrefixNull(),
		},
		"unknown": {
			prefixValue: cidrtypes.NewIPv4PrefixUnknown(),
		},
		"valid IPv4 prefix - loopback": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/8"),
		},
		"valid IPv4 prefix - private": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("172.16.0.0/12"),
		},
		"invalid IPv4 prefix - invalid address no dots": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("192168255255/8"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: 192168255255/8\n"+
					"Error: netip.ParsePrefix(\"192168255255/8\"): ParseAddr(\"192168255255\"): unable to parse IP",
			),
		},
		"invalid IPv4 prefix - invalid address with leading zeroes": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.000/8"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: 127.0.0.000/8\n"+
					"Error: netip.ParsePrefix(\"127.0.0.000/8\"): ParseAddr(\"127.0.0.000\"): IPv4 field has octet with leading zero",
			),
		},
		"invalid IPv4 prefix - invalid address length": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.1/8"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: 127.0.1/8\n"+
					"Error: netip.ParsePrefix(\"127.0.1/8\"): ParseAddr(\"127.0.1\"): IPv4 address too short",
			),
		},
		"invalid IPv4 prefix - invalid prefix length": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/999"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: 127.0.0.0/999\n"+
					"Error: netip.ParsePrefix(\"127.0.0.0/999\"): prefix length out of range",
			),
		},
		"invalid IPv4 prefix - invalid prefix characters": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("127.0.0.0/notcorrect"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: 127.0.0.0/notcorrect\n"+
					"Error: netip.ParsePrefix(\"127.0.0.0/notcorrect\"): bad bits after slash: \"notcorrect\"",
			),
		},
		"invalid IPv4 prefix - IPv6 prefix": {
			prefixValue: cidrtypes.NewIPv4PrefixValue("::/128"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid IPv4 CIDR String Value: "+
					"An IPv6 CIDR string format was provided, string value must be IPv4 CIDR string format (RFC 4632).\n\n"+
					"Given Value: ::/128\n",
			),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
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

func TestIPv4PrefixValueIPv4Prefix(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefixValue      cidrtypes.IPv4Prefix
		expectedIpPrefix netip.Prefix
		expectedDiags    diag.Diagnostics
	}{
		"IPv4 CIDR value is null ": {
			prefixValue: cidrtypes.NewIPv4PrefixNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv4Prefix ValueIPv4Prefix Error",
					"IPv4 CIDR string value is null",
				),
			},
		},
		"IPv4 CIDR value is unknown ": {
			prefixValue: cidrtypes.NewIPv4PrefixUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv4Prefix ValueIPv4Prefix Error",
					"IPv4 CIDR string value is unknown",
				),
			},
		},
		"valid IPv4 CIDR ": {
			prefixValue:      cidrtypes.NewIPv4PrefixValue("172.16.0.0/12"),
			expectedIpPrefix: netip.MustParsePrefix("172.16.0.0/12"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipPrefix, diags := testCase.prefixValue.ValueIPv4Prefix()

			if ipPrefix != testCase.expectedIpPrefix {
				t.Errorf("Unexpected difference in netip.Prefix, got: %s, expected: %s", ipPrefix, testCase.expectedIpPrefix)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
