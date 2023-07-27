// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iptypes

import (
	"context"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable = (*IPv4Address)(nil)
)

// IPv4Address represents a valid IPv4 address string (dotted decimal, no leading zeroes). No semantic equality
// logic is defined for IPv4Address, so it will follow Terraform's data-consistency rules for strings, which must match byte-for-byte.
type IPv4Address struct {
	basetypes.StringValue
}

// Type returns an IPv4AddressType.
func (v IPv4Address) Type(_ context.Context) attr.Type {
	return IPv4AddressType{}
}

// Equal returns true if the given value is equivalent.
func (v IPv4Address) Equal(o attr.Value) bool {
	other, ok := o.(IPv4Address)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// ValueIPv4Address calls netip.ParseAddr with the IPv4Address StringValue. A null or unknown value will produce an error diagnostic.
func (v IPv4Address) ValueIPv4Address() (netip.Addr, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("IPv4Address ValueIPv4Address Error", "IPv4 address string value is null"))
		return netip.Addr{}, diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("IPv4Address ValueIPv4Address Error", "IPv4 address string value is unknown"))
		return netip.Addr{}, diags
	}

	ipv4Addr, err := netip.ParseAddr(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("IPv4Address ValueIPv4Address Error", err.Error()))
		return netip.Addr{}, diags
	}

	return ipv4Addr, nil
}

// NewIPv4AddressNull creates an IPv4Address with a null value. Determine whether the value is null via IsNull method.
func NewIPv4AddressNull() IPv4Address {
	return IPv4Address{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewIPv4AddressUnknown creates an IPv4Address with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewIPv4AddressUnknown() IPv4Address {
	return IPv4Address{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewIPv4AddressValue creates an IPv4Address with a known value. Access the value via ValueString method.
func NewIPv4AddressValue(value string) IPv4Address {
	return IPv4Address{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewIPv4AddressPointerValue creates an IPv4Address with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewIPv4AddressPointerValue(value *string) IPv4Address {
	return IPv4Address{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
