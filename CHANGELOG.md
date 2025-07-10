## 0.3.0 (July 10, 2025)

NOTES:

* all: This Go module has been updated to Go 1.23 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.23 release notes](https://go.dev/doc/go1.23) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#109](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/109))

FEATURES:

* New Custom Types: IPAddress and IPPrefix (Union of IP v4 and v6) ([#6](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/6))
* hwtypes/MACAddress: Add new MACAddressType custom type implementation, representing a MAC address string ([#127](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/127))

## 0.2.0 (September 09, 2024)

BREAKING CHANGES:

* cidrtypes: Removed `Validate()` method from `IPv4PrefixType` following deprecation of `xattr.TypeWithValidate` ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* cidrtypes: Removed `Validate()` method from `IPv6PrefixType` following deprecation of `xattr.TypeWithValidate` ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Removed `Validate()` method from `IPv4AddressType` following deprecation of `xattr.TypeWithValidate` ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Removed `Validate()` method from `IPv6AddressType` following deprecation of `xattr.TypeWithValidate` ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))

NOTES:

* all: This Go module has been updated to Go 1.22 per the [Go support policy](https://go.dev/doc/devel/release#policy). It is recommended to review the [Go 1.22 release notes](https://go.dev/doc/go1.22) before upgrading. Any consumers building on earlier Go versions may experience errors. ([#77](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/77))

ENHANCEMENTS:

* cidrtypes: Added `ValidateAttribute()` method to `IPv4Prefix` type, which supports validating an attribute value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* cidrtypes: Added `ValidateAttribute()` method to `IPv6Prefix` type, which supports validating an attribute value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Added `ValidateAttribute()` method to `IPv4Address` type, which supports validating an attribute value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Added `ValidateAttribute()` method to `IPv6Address` type, which supports validating an attribute value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* cidrtypes: Added `ValidateParameter()` method to `IPv4Prefix` type, which supports validating a provider-defined function parameter value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* cidrtypes: Added `ValidateParameter()` method to `IPv6Prefix` type, which supports validating a provider-defined function parameter value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Added `ValidateParameter()` method to `IPv4Address` type, which supports validating a provider-defined function parameter value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))
* iptypes: Added `ValidateParameter()` method to `IPv6Address` type, which supports validating a provider-defined function parameter value ([#55](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/55))

## 0.1.0 (July 28, 2023)

FEATURES:

* nettypes/iptypes: Add new IPv4Address custom type implementation, representing an IPv4 address string ([#2](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/2))
* nettypes/iptypes: Add new IPv6Address custom type implementation, representing an IPv6 address string ([#2](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/2))
* nettypes/cidrtypes: Add new IPv4Prefix custom type implementation, representing an IPv4 CIDR string ([#2](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/2))
* nettypes/cidrtypes: Add new IPv6Prefix custom type implementation, representing an IPv6 CIDR string ([#2](https://github.com/hashicorp/terraform-plugin-framework-nettypes/issues/2))

