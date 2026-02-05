// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
)

type IPv4PrefixResourceModel struct {
	IPv4CIDR cidrtypes.IPv4Prefix `tfsdk:"ipv4_cidr"`
}

func ExampleIPv4Prefix_ValueIPv4Prefix() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPv4PrefixResourceModel{
		IPv4CIDR: cidrtypes.NewIPv4PrefixValue("127.0.0.0/8"),
	}

	// Check that the IPv4CIDR data is known and able to be converted to netip.Prefix
	if !data.IPv4CIDR.IsNull() && !data.IPv4CIDR.IsUnknown() {
		ipPrefix, diags := data.IPv4CIDR.ValueIPv4Prefix()
		if diags.HasError() {
			return
		}

		loopback := netip.MustParseAddr("127.0.0.1")
		// Output: true, 127.0.0.0/8
		fmt.Printf("%t, %s\n", ipPrefix.Contains(loopback), ipPrefix)
	}
}
