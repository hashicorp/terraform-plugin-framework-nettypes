// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package cidrtypes_test

import (
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
)

type IPPrefixResourceModel struct {
	IPCIDR cidrtypes.IPPrefix `tfsdk:"ip_cidr"`
}

func ExampleIPPrefix_ValueIPPrefix_v4() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPPrefixResourceModel{
		IPCIDR: cidrtypes.NewIPPrefixValue("127.0.0.0/8"),
	}

	// Check that the IPv4CIDR data is known and able to be converted to netip.Prefix
	if !data.IPCIDR.IsNull() && !data.IPCIDR.IsUnknown() {
		ipPrefix, diags := data.IPCIDR.ValueIPPrefix()
		if diags.HasError() {
			return
		}

		loopback := netip.MustParseAddr("127.0.0.1")
		// Output: true, 127.0.0.0/8
		fmt.Printf("%t, %s\n", ipPrefix.Contains(loopback), ipPrefix)
	}
}

func ExampleIPPrefix_ValueIPPrefix_v6() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPPrefixResourceModel{
		IPCIDR: cidrtypes.NewIPPrefixValue("::1/128"),
	}

	// Check that the IPv6CIDR data is known and able to be converted to netip.Prefix
	if !data.IPCIDR.IsNull() && !data.IPCIDR.IsUnknown() {
		ipPrefix, diags := data.IPCIDR.ValueIPPrefix()
		if diags.HasError() {
			return
		}

		loopback := netip.MustParseAddr("::1")
		// Output: true, ::1/128
		fmt.Printf("%t, %s\n", ipPrefix.Contains(loopback), ipPrefix)
	}
}
