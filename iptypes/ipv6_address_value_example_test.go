// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

type IPv6AddressResourceModel struct {
	IPv6Address iptypes.IPv6Address `tfsdk:"ipv6_address"`
}

func ExampleIPv6Address_ValueIPv6Address() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPv6AddressResourceModel{
		IPv6Address: iptypes.NewIPv6AddressValue("::1"),
	}

	// Check that the IPv6Address data is known and able to be converted to netip.Addr
	if !data.IPv6Address.IsNull() && !data.IPv6Address.IsUnknown() {
		ipAddr, diags := data.IPv6Address.ValueIPv6Address()
		if diags.HasError() {
			return
		}

		// Output: true, ::1
		fmt.Printf("%t, %s\n", ipAddr.IsLoopback(), ipAddr)
	}
}
