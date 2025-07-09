// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hwtypes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/hwtypes"
)

type MACAddressResourceModel struct {
	MACAddress hwtypes.MACAddress `tfsdk:"mac_address"`
}

func ExampleMACAddress_ValueMACAddress() {
	// For example purposes, typically the data model would be populated automatically by Plugin Framework via Config, Plan or State.
	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
	data := MACAddressResourceModel{
		MACAddress: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
	}

	// Check that the MACAddress data is known and able to be converted to net.HardwareAddr
	if !data.MACAddress.IsNull() && !data.MACAddress.IsUnknown() {
		macAddr, diags := data.MACAddress.ValueMACAddress()
		if diags.HasError() {
			return
		}

		// Output: true, 00:00:5e:00:53:00
		fmt.Printf("%t, %s\n", data.MACAddress.ValueString() == macAddr.String(), macAddr.String())
	}
}
