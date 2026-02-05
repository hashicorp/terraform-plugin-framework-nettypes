// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

type IPv4AddressResourceModel struct {
	IPv4Address iptypes.IPv4Address `tfsdk:"ipv4_address"`
}

func ExampleIPv4Address_ValueIPv4Address() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPv4AddressResourceModel{
		IPv4Address: iptypes.NewIPv4AddressValue("127.0.0.1"),
	}

	// Check that the IPv4Address data is known and able to be converted to netip.Addr
	if !data.IPv4Address.IsNull() && !data.IPv4Address.IsUnknown() {
		ipAddr, diags := data.IPv4Address.ValueIPv4Address()
		if diags.HasError() {
			return
		}

		// Output: true, 127.0.0.1
		fmt.Printf("%t, %s\n", ipAddr.IsLoopback(), ipAddr)
	}
}
