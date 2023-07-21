// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"context"
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestIPv6AddressStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentIpAddr iptypes.IPv6Address
		givenIpAddr   basetypes.StringValuable
		expectedMatch bool
		expectedDiags diag.Diagnostics
	}{
		"not equal - IPv6 address mismatch": {
			currentIpAddr: iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:1"),
			expectedMatch: false,
		},
		"not equal - IPv6 address compressed mismatch": {
			currentIpAddr: iptypes.NewIPv6AddressValue("FF01::"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("FF01::1"),
			expectedMatch: false,
		},
		"not equal - IPv4-Mapped IPv6 address mismatch": {
			currentIpAddr: iptypes.NewIPv6AddressValue("::FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("::FFFF:192.168.255.254"),
			expectedMatch: false,
		},
		"semantically equal - byte-for-byte match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:0"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped byte-for-byte match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("::FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("::FFFF:192.168.255.255"),
			expectedMatch: true,
		},
		"semantically equal - compressed all zeroes match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("::"),
			expectedMatch: true,
		},
		"semantically equal - compressed all leading zeroes match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("2001:0DB8:0000:0000:0008:0800:0200C:417A"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("2001:DB8::8:800:200C:417A"),
			expectedMatch: true,
		},
		"semantically equal - start compressed match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("::101"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:101"),
			expectedMatch: true,
		},
		"semantically equal - middle compressed match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("2001:DB8::8:800:200C:417A"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("2001:DB8:0:0:8:800:200C:417A"),
			expectedMatch: true,
		},
		"semantically equal - end compressed match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("FF01:0:0:0:0:0:0:0"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("FF01::"),
			expectedMatch: true,
		},
		"semantically equal - IPv4-Mapped compressed match": {
			currentIpAddr: iptypes.NewIPv6AddressValue("0:0:0:0:0:FFFF:192.168.255.255"),
			givenIpAddr:   iptypes.NewIPv6AddressValue("::FFFF:192.168.255.255"),
			expectedMatch: true,
		},
		"error - not given IPv6Address IPv6 value": {
			currentIpAddr: iptypes.NewIPv6AddressValue("0:0:0:0:0:0:0:0"),
			givenIpAddr:   basetypes.NewStringValue("0:0:0:0:0:0:0:0"),
			expectedMatch: false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: iptypes.IPv6Address\n"+
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
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
			}
		})
	}
}

func TestIPv6AddressValueIPv6Address(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipValue        iptypes.IPv6Address
		expectedIpAddr netip.Addr
		expectedDiags  diag.Diagnostics
	}{
		"IPv6 address value is null ": {
			ipValue: iptypes.NewIPv6AddressNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv6Address ValueIPv6Address Error",
					"IPv6 address string value is null",
				),
			},
		},
		"IPv6 address value is unknown ": {
			ipValue: iptypes.NewIPv6AddressUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"IPv6Address ValueIPv6Address Error",
					"IPv6 address string value is unknown",
				),
			},
		},
		"valid IPv6 address ": {
			ipValue:        iptypes.NewIPv6AddressValue("2001:DB8::8:800:200C:417A"),
			expectedIpAddr: netip.MustParseAddr("2001:DB8::8:800:200C:417A"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipAddr, diags := testCase.ipValue.ValueIPv6Address()

			if ipAddr != testCase.expectedIpAddr {
				t.Errorf("Unexpected difference in netip.Addr, got: %s, expected: %s", ipAddr, testCase.expectedIpAddr)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
			}
		})
	}
}
