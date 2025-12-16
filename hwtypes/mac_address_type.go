// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package hwtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*MACAddressType)(nil)

// MACAddressType is an attribute type that represents a valid IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP
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
type MACAddressType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t MACAddressType) String() string {
	return "hwtypes.MACAddressType"
}

// ValueType returns the Value type.
func (t MACAddressType) ValueType(ctx context.Context) attr.Value {
	return MACAddress{}
}

// Equal returns true if the given type is equivalent.
func (t MACAddressType) Equal(o attr.Type) bool {
	other, ok := o.(MACAddressType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t MACAddressType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return MACAddress{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t MACAddressType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
