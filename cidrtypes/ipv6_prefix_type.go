// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cidrtypes

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*IPv6PrefixType)(nil)
	_ xattr.TypeWithValidate  = (*IPv6PrefixType)(nil)
)

// IPv6PrefixType is an attribute type that represents a valid IPv6 CIDR string (RFC 4291). Semantic equality logic is defined for IPv6PrefixType
// such that a CIDR string with the address zero bits `compressed` will be considered equivalent to the `non-compressed` string.
//
// Examples:
//   - `0:0:0:0:0:0:0:0/128` is semantically equal to `::/128`
//   - `2001:0DB8:0:0:0:0:0:0CD3/60` is semantically equal to `2001:0DB8::CD30/60`
//   - `FF00:0:0:0:0:0:0:0/8` is semantically equal to `FF00::/8`
type IPv6PrefixType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t IPv6PrefixType) String() string {
	return "cidrtypes.IPv6PrefixType"
}

// ValueType returns the Value type.
func (t IPv6PrefixType) ValueType(ctx context.Context) attr.Value {
	return IPv6Prefix{}
}

// Equal returns true if the given type is equivalent.
func (t IPv6PrefixType) Equal(o attr.Type) bool {
	other, ok := o.(IPv6PrefixType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value that is a valid IPv6 CIDR.
func (t IPv6PrefixType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"IPv6 Prefix Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string

	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"IPv6 Prefix Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	ipPrefix, err := netip.ParsePrefix(valueString)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Invalid IPv6 CIDR String Value",
			"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	if ipPrefix.Addr().Is4() {
		diags.AddAttributeError(
			path,
			"Invalid IPv6 CIDR String Value",
			"An IPv4 CIDR string format was provided, string value must be IPv6 CIDR string format (RFC 4291).\n\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	if !ipPrefix.IsValid() || !ipPrefix.Addr().Is6() {
		diags.AddAttributeError(
			path,
			"Invalid IPv6 CIDR String Value",
			"A string value was provided that is not valid IPv6 CIDR string format (RFC 4291).\n\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	return diags
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t IPv6PrefixType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return IPv6Prefix{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t IPv6PrefixType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
