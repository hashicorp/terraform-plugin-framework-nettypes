// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"context"
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0000:0000:0008:0800:0200C:417A/60"),
			givenIpPrefix:   cidrtypes.NewIPv6PrefixValue("2001:0db8:0000:0000:0008:0800:0200c:417a/60"),
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
			currentIpPrefix: cidrtypes.NewIPv6PrefixValue("2001:0DB8:0000:0000:0008:0800:0200C:417A/60"),
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
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentIpPrefix.StringSemanticEquals(context.Background(), testCase.givenIpPrefix)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
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
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipAddr, diags := testCase.ipValue.ValueIPv6Prefix()

			if ipAddr != testCase.expectedIpPrefix {
				t.Errorf("Unexpected difference in netip.Prefix, got: %s, expected: %s", ipAddr, testCase.expectedIpPrefix)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (+got, -expected): %s", diff)
			}
		})
	}
}
