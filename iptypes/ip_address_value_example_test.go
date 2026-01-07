// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package iptypes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
)

type IPAddressResourceModel struct {
	IPAddress iptypes.IPAddress `tfsdk:"ip_address"`
}

func ExampleIPAddress_ValueIPAddress_v4() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPAddressResourceModel{
		IPAddress: iptypes.NewIPAddressValue("127.0.0.1"),
	}

	// Check that the IPAddress data is known and able to be converted to netip.Addr
	if !data.IPAddress.IsNull() && !data.IPAddress.IsUnknown() {
		ipAddr, diags := data.IPAddress.ValueIPAddress()
		if diags.HasError() {
			return
		}

		// Output: true, 127.0.0.1
		fmt.Printf("%t, %s\n", ipAddr.IsLoopback(), ipAddr)
	}
}

func ExampleIPAddress_ValueIPAddress_v6() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := IPAddressResourceModel{
		IPAddress: iptypes.NewIPAddressValue("::1"),
	}

	// Check that the IPAddress data is known and able to be converted to netip.Addr
	if !data.IPAddress.IsNull() && !data.IPAddress.IsUnknown() {
		ipAddr, diags := data.IPAddress.ValueIPAddress()
		if diags.HasError() {
			return
		}

		// Output: true, ::1
		fmt.Printf("%t, %s\n", ipAddr.IsLoopback(), ipAddr)
	}
}
