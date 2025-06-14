// Copyright (c) HashiCorp, Inc.
// Copyright (c) Christopher Marget
// SPDX-License-Identifier: MPL-2.0

package hwtypes

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*MACAddress)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*MACAddress)(nil)
	_ xattr.ValidateableAttribute                = (*MACAddress)(nil)
	_ function.ValidateableParameter             = (*MACAddress)(nil)
)

// MACAddress is an attribute type that represents a valid IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP
// over InfiniBand link-layer address. Semantic equality logic is defined for MACAddressType, so that addresses
// expressed with varying case and notation are considered equal.
//
// All of the following are semantically equal:
//   - 00:00:5e:00:53:01
//   - 00:00:5E:00:53:01
//   - 00-00-5e-00-53-01
//   - 00-00-5E-00-53-01
//   - 0000.5e00.5301
//   - 0000.5E00.5301
type MACAddress struct {
	basetypes.StringValue
}

// Type returns an MACAddressType.
func (v MACAddress) Type(_ context.Context) attr.Type {
	return MACAddressType{}
}

// Equal returns true if the given value is equivalent.
func (v MACAddress) Equal(o attr.Value) bool {
	other, ok := o.(MACAddress)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the given MAC address string value is semantically equal to the current MAC address string value.
// This comparison utilizes net.ParseMAC and then compares the resulting net.HardwareAddr representations. This means that addresses
// expressed with varying case and notation are considered equal.
//
// All of the following are semantically equal:
//   - 00:00:5e:00:53:01
//   - 00:00:5E:00:53:01
//   - 00-00-5e-00-53-01
//   - 00-00-5E-00-53-01
//   - 0000.5e00.5301
//   - 0000.5E00.5301
func (v MACAddress) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MACAddress)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	// MAC addresses are already validated at this point, ignoring errors
	newMacAddr, _ := net.ParseMAC(newValue.ValueString())
	currentMacAddr, _ := net.ParseMAC(v.ValueString())

	return bytes.Equal(currentMacAddr, newMacAddr), diags
}

// ValidateAttribute implements attribute value validation. This type requires the value provided to be a String
// value that is a valid MAC address. This utilizes the Go `net` library for parsing.
func (v MACAddress) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	_, err := net.ParseMAC(v.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MAC Address String Value",
			"A string value was provided that is not valid MAC string format.\n\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)

		return // leaving this redundant return in case additional validations are added later
	}
}

// ValidateParameter implements provider-defined function parameter value validation. This type requires the value
// provided to be a String value that is a valid MAC address. This utilizes the Go `net` library for parsing.
func (v MACAddress) ValidateParameter(_ context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	_, err := net.ParseMAC(v.ValueString())
	if err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid MAC Address String Value: "+
				"A string value was provided that is not valid MAC string format.\n\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)

		return // leaving this redundant return in case additional validations are added later
	}
}

// ValueMACAddress calls net.ParseMAC with the MACAddress StringValue. A null or unknown value will produce an error diagnostic.
func (v MACAddress) ValueMACAddress() (net.HardwareAddr, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("MACAddress ValueMACAddress Error", "MAC address string value is null"))
		return nil, diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("MACAddress ValueMACAddress Error", "MAC address string value is unknown"))
		return nil, diags
	}

	macAddr, err := net.ParseMAC(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("MACAddress ValueMACAddress Error", err.Error()))
		return nil, diags
	}

	return macAddr, nil
}

// NewMACAddressNull creates a MACAddress with a null value. Determine whether the value is null via IsNull method.
func NewMACAddressNull() MACAddress {
	return MACAddress{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewMACAddressUnknown creates a MACAddress with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewMACAddressUnknown() MACAddress {
	return MACAddress{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewMACAddressValue creates a MACAddress with a known value. Access the value via ValueString method.
func NewMACAddressValue(value string) MACAddress {
	return MACAddress{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewMACAddressPointerValue creates a MACAddress with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewMACAddressPointerValue(value *string) MACAddress {
	return MACAddress{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
