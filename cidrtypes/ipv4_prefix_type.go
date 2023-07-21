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
	_ basetypes.StringTypable = (*IPv4PrefixType)(nil)
	_ xattr.TypeWithValidate  = (*IPv4PrefixType)(nil)
)

// IPv4PrefixType is an attribute type that represents a valid IPv4 CIDR string (RFC 4632). No semantic equality
// logic is defined for IPv4PrefixType, so it will follow Terraform's data-consistency rules for strings, which must match byte-for-byte.
type IPv4PrefixType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t IPv4PrefixType) String() string {
	return "cidrtypes.IPv4PrefixType"
}

// ValueType returns the Value type.
func (t IPv4PrefixType) ValueType(ctx context.Context) attr.Value {
	return IPv4Prefix{}
}

// Equal returns true if the given type is equivalent.
func (t IPv4PrefixType) Equal(o attr.Type) bool {
	other, ok := o.(IPv4PrefixType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value that is a valid IPv4 CIDR (RFC 4632).
func (t IPv4PrefixType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"IPv4 Prefix Type Validation Error",
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
			"IPv4 Prefix Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	ipPrefix, err := netip.ParsePrefix(valueString)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Invalid IPv4 CIDR String Value",
			"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	// TODO: is this correct way to determine IPv6 CIDR?
	if ipPrefix.Addr().Is6() {
		diags.AddAttributeError(
			path,
			"Invalid IPv4 CIDR String Value",
			"An IPv6 CIDR string format was provided, string value must be IPv4 CIDR string format (RFC 4632).\n\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	// TODO: is this correct way to determine IPv4 CIDR?
	if !ipPrefix.IsValid() || !ipPrefix.Addr().Is4() {
		diags.AddAttributeError(
			path,
			"Invalid IPv4 CIDR String Value",
			"A string value was provided that is not valid IPv4 CIDR string format (RFC 4632).\n\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	return diags
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t IPv4PrefixType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return IPv4Prefix{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t IPv4PrefixType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
