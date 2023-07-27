// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes

import (
	"context"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable = (*IPv4Prefix)(nil)
)

// IPv4Prefix represents a valid IPv4 CIDR string (RFC 4632). No semantic equality logic is defined for IPv4Prefix,
// so it will follow Terraform's data-consistency rules for strings, which must match byte-for-byte.
type IPv4Prefix struct {
	basetypes.StringValue
}

// Type returns an IPv4PrefixType.
func (v IPv4Prefix) Type(_ context.Context) attr.Type {
	return IPv4PrefixType{}
}

// Equal returns true if the given value is equivalent.
func (v IPv4Prefix) Equal(o attr.Value) bool {
	other, ok := o.(IPv4Prefix)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// ValueIPv4Prefix calls netip.ParsePrefix with the IPv4Prefix StringValue. A null or unknown value will produce an error diagnostic.
func (v IPv4Prefix) ValueIPv4Prefix() (netip.Prefix, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("IPv4Prefix ValueIPv4Prefix Error", "IPv4 CIDR string value is null"))
		return netip.Prefix{}, diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("IPv4Prefix ValueIPv4Prefix Error", "IPv4 CIDR string value is unknown"))
		return netip.Prefix{}, diags
	}

	ipv4Prefix, err := netip.ParsePrefix(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("IPv4Prefix ValueIPv4Prefix Error", err.Error()))
		return netip.Prefix{}, diags
	}

	return ipv4Prefix, nil
}

// NewIPv4PrefixNull creates an IPv4Prefix with a null value. Determine whether the value is null via IsNull method.
func NewIPv4PrefixNull() IPv4Prefix {
	return IPv4Prefix{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewIPv4PrefixUnknown creates an IPv4Prefix with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewIPv4PrefixUnknown() IPv4Prefix {
	return IPv4Prefix{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewIPv4PrefixValue creates an IPv4Prefix with a known value. Access the value via ValueString method.
func NewIPv4PrefixValue(value string) IPv4Prefix {
	return IPv4Prefix{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewIPv4PrefixPointerValue creates an IPv4Prefix with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewIPv4PrefixPointerValue(value *string) IPv4Prefix {
	return IPv4Prefix{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
