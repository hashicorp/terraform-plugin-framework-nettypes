// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
)

type IPv6PrefixResourceModel struct {
	IPv6CIDR cidrtypes.IPv6Prefix `tfsdk:"ipv6_cidr"`
}

func ExampleIPv6Prefix_ValueIPv6Prefix() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPv6PrefixResourceModel{
		IPv6CIDR: cidrtypes.NewIPv6PrefixValue("::1/128"),
	}

	// Check that the IPv6CIDR data is known and able to be converted to netip.Prefix
	if !data.IPv6CIDR.IsNull() && !data.IPv6CIDR.IsUnknown() {
		ipPrefix, diags := data.IPv6CIDR.ValueIPv6Prefix()
		if diags.HasError() {
			return
		}

		loopback := netip.MustParseAddr("::1")
		// Output: true, ::1/128
		fmt.Printf("%t, %s\n", ipPrefix.Contains(loopback), ipPrefix)
	}
}
