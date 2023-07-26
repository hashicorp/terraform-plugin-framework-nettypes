// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

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
