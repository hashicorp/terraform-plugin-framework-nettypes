// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package hwtypes_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/hwtypes"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMACAddressStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentMacAddr hwtypes.MACAddress
		givenMacAddr   basetypes.StringValuable
		expectedMatch  bool
		expectedDiags  diag.Diagnostics
	}{
		"not equal - MAC address value mismatch": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00:00:5e:00:53:01"),
			expectedMatch:  false,
		},
		"not equal - MAC address case mismatch": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00:00:5E:00:53:01"),
			expectedMatch:  false,
		},
		"not equal - MAC address delimiter mismatch": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00-00:5e-00-53-01"),
			expectedMatch:  false,
		},
		"semantically equal - colon-delimited byte-for-byte match": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			expectedMatch:  true,
		},
		"semantically equal - colon-delimited case insensitive": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00:00:5E:00:53:00"),
			expectedMatch:  true,
		},
		"semantically equal - hyphen-delimited byte-for-byte match": {
			currentMacAddr: hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
			expectedMatch:  true,
		},
		"semantically equal - hyphen-delimited case insensitive": {
			currentMacAddr: hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00-00-5E-00-53-00"),
			expectedMatch:  true,
		},
		"semantically equal - dot-delimited byte-for-byte match": {
			currentMacAddr: hwtypes.NewMACAddressValue("0000.5e00.5300"),
			givenMacAddr:   hwtypes.NewMACAddressValue("0000.5e00.5300"),
			expectedMatch:  true,
		},
		"semantically equal - dot-delimited case insensitive": {
			currentMacAddr: hwtypes.NewMACAddressValue("0000.5e00.5300"),
			givenMacAddr:   hwtypes.NewMACAddressValue("0000.5E00.5300"),
			expectedMatch:  true,
		},
		"semantically equal - dot vs colon delimited": {
			currentMacAddr: hwtypes.NewMACAddressValue("0000.5e00.5300"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			expectedMatch:  true,
		},
		"semantically equal - dot vs hyphen delimited": {
			currentMacAddr: hwtypes.NewMACAddressValue("0000.5e00.5300"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
			expectedMatch:  true,
		},
		"semantically equal - colon vs hyphen delimited": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
			expectedMatch:  true,
		},
		"error - not given MACAddress MAC value": {
			currentMacAddr: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			givenMacAddr:   basetypes.NewStringValue("00:00:5e:00:53:00"),
			expectedMatch:  false,
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Semantic Equality Check Error",
					"An unexpected value type was received while performing semantic equality checks. "+
						"Please report this to the provider developers.\n\n"+
						"Expected Value Type: hwtypes.MACAddress\n"+
						"Got Value Type: basetypes.StringValue",
				),
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentMacAddr.StringSemanticEquals(context.Background(), testCase.givenMacAddr)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestMACAddressValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue  hwtypes.MACAddress
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			addressValue: hwtypes.MACAddress{},
		},
		"null": {
			addressValue: hwtypes.NewMACAddressNull(),
		},
		"unknown": {
			addressValue: hwtypes.NewMACAddressUnknown(),
		},
		"valid MAC address - lowercase - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
		},
		"valid MAC address - uppercase - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5E:00:53:00"),
		},
		"valid MAC address - lowercase - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
		},
		"valid MAC address - uppercase - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-5E-00-53-00"),
		},
		"valid MAC address - lowercase - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.5e00.5300"),
		},
		"valid MAC address - uppercase - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.5E00.5300"),
		},
		"valid MAC address - lowercase - IPoIB - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01"),
		},
		"valid MAC address - uppercase - IPoIB - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:00:00:FE:80:00:00:00:00:00:00:02:00:5E:10:00:00:00:01"),
		},
		"valid MAC address - lowercase - IPoIB - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01"),
		},
		"valid MAC address - uppercase - IPoIB - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-00-00-FE-80-00-00-00-00-00-00-02-00-5E-10-00-00-00-01"),
		},
		"valid MAC address - lowercase - IPoIB - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001"),
		},
		"valid MAC address - uppercase - IPoIB - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.0000.FE80.0000.0000.0000.0200.5E10.0000.0001"),
		},
		"invalid MAC address - 7 bytes": {
			addressValue: hwtypes.NewMACAddressValue("0:0:0:0:0:0:0"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid MAC Address String Value",
					"A string value was provided that is not valid MAC string format.\n\n"+
						"Given Value: 0:0:0:0:0:0:0\n"+
						"Error: address 0:0:0:0:0:0:0: invalid MAC address",
				),
			},
		},
		"invalid MAC address - bogus digit": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5G:00:53:00"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid MAC Address String Value",
					"A string value was provided that is not valid MAC string format.\n\n"+
						"Given Value: 00:00:5G:00:53:00\n"+
						"Error: address 00:00:5G:00:53:00: invalid MAC address",
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.addressValue.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{Path: path.Root("test")},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestMACAddressValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		addressValue    hwtypes.MACAddress
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			addressValue: hwtypes.MACAddress{},
		},
		"null": {
			addressValue: hwtypes.NewMACAddressNull(),
		},
		"unknown": {
			addressValue: hwtypes.NewMACAddressUnknown(),
		},
		"valid MAC address - lowercase - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
		},
		"valid MAC address - uppercase - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5E:00:53:00"),
		},
		"valid MAC address - lowercase - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-5e-00-53-00"),
		},
		"valid MAC address - uppercase - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-5E-00-53-00"),
		},
		"valid MAC address - lowercase - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.5e00.5300"),
		},
		"valid MAC address - uppercase - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.5E00.5300"),
		},
		"valid MAC address - lowercase - IPoIB - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01"),
		},
		"valid MAC address - uppercase - IPoIB - colon-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00:00:00:00:FE:80:00:00:00:00:00:00:02:00:5E:10:00:00:00:01"),
		},
		"valid MAC address - lowercase - IPoIB - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01"),
		},
		"valid MAC address - uppercase - IPoIB - hyphen-delimited": {
			addressValue: hwtypes.NewMACAddressValue("00-00-00-00-FE-80-00-00-00-00-00-00-02-00-5E-10-00-00-00-01"),
		},
		"valid MAC address - lowercase - IPoIB - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001"),
		},
		"valid MAC address - uppercase - IPoIB - dot-delimited": {
			addressValue: hwtypes.NewMACAddressValue("0000.0000.FE80.0000.0000.0000.0200.5E10.0000.0001"),
		},
		"invalid MAC address - 7 bytes": {
			addressValue: hwtypes.NewMACAddressValue("0:0:0:0:0:0:0"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid MAC Address String Value: "+
					"A string value was provided that is not valid MAC string format.\n\n"+
					"Given Value: 0:0:0:0:0:0:0\n"+
					"Error: address 0:0:0:0:0:0:0: invalid MAC address",
			),
		},
		"invalid MAC address - bogus digit": {
			addressValue: hwtypes.NewMACAddressValue("00:00:5G:00:53:00"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				"Invalid MAC Address String Value: "+
					"A string value was provided that is not valid MAC string format.\n\n"+
					"Given Value: 00:00:5G:00:53:00\n"+
					"Error: address 00:00:5G:00:53:00: invalid MAC address",
			),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.addressValue.ValidateParameter(
				context.Background(),
				function.ValidateParameterRequest{
					Position: 0,
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Error, testCase.expectedFuncErr); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestMACAddressValueMACAddress(t *testing.T) {
	t.Parallel()

	mustParseMac := func(s string) net.HardwareAddr {
		mac, err := net.ParseMAC(s)
		if err != nil {
			panic(err)
		}
		return mac
	}

	testCases := map[string]struct {
		macValue        hwtypes.MACAddress
		expectedMacAddr net.HardwareAddr
		expectedDiags   diag.Diagnostics
	}{
		"MAC address value is null ": {
			macValue: hwtypes.NewMACAddressNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"MACAddress ValueMACAddress Error",
					"MAC address string value is null",
				),
			},
		},
		"MAC address value is unknown ": {
			macValue: hwtypes.NewMACAddressUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"MACAddress ValueMACAddress Error",
					"MAC address string value is unknown",
				),
			},
		},
		"valid MAC address ": {
			macValue:        hwtypes.NewMACAddressValue("00:00:5e:00:53:00"),
			expectedMacAddr: mustParseMac("00:00:5e:00:53:00"),
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			macAddr, diags := testCase.macValue.ValueMACAddress()

			if !bytes.Equal(macAddr, testCase.expectedMacAddr) {
				t.Errorf("Unexpected difference in net.HardwareAddr, got: %s, expected: %s", macAddr, testCase.expectedMacAddr)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
