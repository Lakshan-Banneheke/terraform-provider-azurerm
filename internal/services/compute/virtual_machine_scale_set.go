// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func VirtualMachineScaleSetAdditionalCapabilitiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// NOTE: requires registration to use:
				// $ az feature show --namespace Microsoft.Compute --name UltraSSDWithVMSS
				// $ az provider register -n Microsoft.Compute
				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ForceNew: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAdditionalCapabilities(input []interface{}) *virtualmachinescalesets.AdditionalCapabilities {
	capabilities := virtualmachinescalesets.AdditionalCapabilities{}

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})

		capabilities.UltraSSDEnabled = utils.Bool(raw["ultra_ssd_enabled"].(bool))
	}

	return &capabilities
}

func FlattenVirtualMachineScaleSetAdditionalCapabilities(input *virtualmachinescalesets.AdditionalCapabilities) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ultraSsdEnabled := false
	if input.UltraSSDEnabled != nil {
		ultraSsdEnabled = *input.UltraSSDEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"ultra_ssd_enabled": ultraSsdEnabled,
		},
	}
}

func VirtualMachineScaleSetNetworkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"ip_configuration": virtualMachineScaleSetIPConfigurationSchema(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networkValidate.NetworkSecurityGroupID,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func VirtualMachineScaleSetGalleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 100,
		Computed: !features.FourPointOhBeta(),
		ConflictsWith: func() []string {
			if !features.FourPointOhBeta() {
				return []string{"gallery_applications"}
			}
			return []string{}
		}(),
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				// Example: https://mystorageaccount.blob.core.windows.net/configurations/settings.config
				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				// NOTE: Per the service team, "this is a pass through value that we just add to the model but don't depend on. It can be any string."
				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func VirtualMachineScaleSetGalleryApplicationsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      100,
		Computed:      !features.FourPointOhBeta(),
		ConflictsWith: []string{"gallery_application"},
		Deprecated:    "`gallery_applications` has been renamed to `gallery_application` and will be deprecated in 4.0",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"package_reference_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
					Deprecated:   "`package_reference_id` has been renamed to `version_id` and will be deprecated in 4.0",
				},

				// Example: https://mystorageaccount.blob.core.windows.net/configurations/settings.config
				"configuration_reference_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					Deprecated:   "`configuration_reference_blob_uri` has been renamed to `configuration_blob_uri` and will be deprecated in 4.0",
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				// NOTE: Per the service team, "this is a pass through value that we just add to the model but don't depend on. It can be any string."
				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func expandVirtualMachineScaleSetGalleryApplication(input []interface{}) *[]virtualmachinescalesets.VMGalleryApplication {
	if len(input) == 0 {
		return nil
	}

	out := make([]virtualmachinescalesets.VMGalleryApplication, 0)

	for _, v := range input {
		packageReferenceId := v.(map[string]interface{})["version_id"].(string)
		configurationReference := v.(map[string]interface{})["configuration_blob_uri"].(string)
		order := v.(map[string]interface{})["order"].(int)
		tag := v.(map[string]interface{})["tag"].(string)

		app := &virtualmachinescalesets.VMGalleryApplication{
			PackageReferenceId:     packageReferenceId,
			ConfigurationReference: utils.String(configurationReference),
			Order:                  utils.Int64(int64(order)),
			Tags:                   utils.String(tag),
		}

		out = append(out, *app)
	}

	return &out
}

func flattenVirtualMachineScaleSetGalleryApplication(input *[]virtualmachinescalesets.VMGalleryApplication) []interface{} {
	if len(*input) == 0 {
		return nil
	}

	out := make([]interface{}, 0)

	for _, v := range *input {
		var configurationReference, tag string
		var order int

		if v.ConfigurationReference != nil {
			configurationReference = *v.ConfigurationReference
		}

		if v.Order != nil {
			order = int(*v.Order)
		}

		if v.Tags != nil {
			tag = *v.Tags
		}

		app := map[string]interface{}{
			"version_id":             v.PackageReferenceId,
			"configuration_blob_uri": configurationReference,
			"order":                  order,
			"tag":                    tag,
		}

		out = append(out, app)
	}

	return out
}

func expandVirtualMachineScaleSetGalleryApplications(input []interface{}) *[]virtualmachinescalesets.VMGalleryApplication {
	if len(input) == 0 {
		return nil
	}

	out := make([]virtualmachinescalesets.VMGalleryApplication, 0)

	for _, v := range input {
		packageReferenceId := v.(map[string]interface{})["package_reference_id"].(string)
		configurationReference := v.(map[string]interface{})["configuration_reference_blob_uri"].(string)
		order := v.(map[string]interface{})["order"].(int)
		tag := v.(map[string]interface{})["tag"].(string)

		app := &virtualmachinescalesets.VMGalleryApplication{
			PackageReferenceId:     packageReferenceId,
			ConfigurationReference: utils.String(configurationReference),
			Order:                  utils.Int64(int64(order)),
			Tags:                   utils.String(tag),
		}

		out = append(out, *app)
	}

	return &out
}

func flattenVirtualMachineScaleSetGalleryApplications(input *[]virtualmachinescalesets.VMGalleryApplication) []interface{} {
	if len(*input) == 0 {
		return nil
	}

	out := make([]interface{}, 0)

	for _, v := range *input {
		var configurationReference, tag string
		var order int

		if v.ConfigurationReference != nil {
			configurationReference = *v.ConfigurationReference
		}

		if v.Order != nil {
			order = int(*v.Order)
		}

		if v.Tags != nil {
			tag = *v.Tags
		}

		app := map[string]interface{}{
			"package_reference_id":             v.PackageReferenceId,
			"configuration_reference_blob_uri": configurationReference,
			"order":                            order,
			"tag":                              tag,
		}

		out = append(out, app)
	}

	return out
}

func VirtualMachineScaleSetScaleInPolicySchema() *pluginsdk.Schema {
	if !features.FourPointOhBeta() {
		return &pluginsdk.Schema{
			Type:          pluginsdk.TypeList,
			Optional:      true,
			Computed:      !features.FourPointOhBeta(),
			MaxItems:      1,
			ConflictsWith: []string{"scale_in_policy"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"rule": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(compute.VirtualMachineScaleSetScaleInRulesDefault),
						ValidateFunc: validation.StringInSlice([]string{
							string(compute.VirtualMachineScaleSetScaleInRulesDefault),
							string(compute.VirtualMachineScaleSetScaleInRulesNewestVM),
							string(compute.VirtualMachineScaleSetScaleInRulesOldestVM),
						}, false),
					},

					"force_deletion_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		}
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"rule": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.VirtualMachineScaleSetScaleInRulesDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.VirtualMachineScaleSetScaleInRulesDefault),
						string(compute.VirtualMachineScaleSetScaleInRulesNewestVM),
						string(compute.VirtualMachineScaleSetScaleInRulesOldestVM),
					}, false),
				},

				"force_deletion_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetScaleInPolicy(input []interface{}) *virtualmachinescalesets.ScaleInPolicy {
	if len(input) == 0 {
		return nil
	}

	rule := input[0].(map[string]interface{})["rule"].(string)
	forceDeletion := input[0].(map[string]interface{})["force_deletion_enabled"].(bool)

	return &virtualmachinescalesets.ScaleInPolicy{
		Rules:         &[]virtualmachinescalesets.VirtualMachineScaleSetScaleInRules{virtualmachinescalesets.VirtualMachineScaleSetScaleInRules(rule)},
		ForceDeletion: utils.Bool(forceDeletion),
	}
}

func FlattenVirtualMachineScaleSetScaleInPolicy(input *virtualmachinescalesets.ScaleInPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	rule := string(virtualmachinescalesets.VirtualMachineScaleSetScaleInRulesDefault)
	var forceDeletion bool
	if rules := input.Rules; rules != nil && len(*rules) > 0 {
		rule = string((*rules)[0])
	}

	if input.ForceDeletion != nil {
		forceDeletion = *input.ForceDeletion
	}

	return []interface{}{
		map[string]interface{}{
			"rule":                   rule,
			"force_deletion_enabled": forceDeletion,
		},
	}
}

func VirtualMachineScaleSetSpotRestorePolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ForceNew: true,
				},

				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "PT1H",
					ForceNew:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetSpotRestorePolicy(input []interface{}) *virtualmachinescalesets.SpotRestorePolicy {
	if len(input) == 0 {
		return nil
	}

	enabled := input[0].(map[string]interface{})["enabled"].(bool)
	timeout := input[0].(map[string]interface{})["timeout"].(string)

	return &virtualmachinescalesets.SpotRestorePolicy{
		Enabled:        utils.Bool(enabled),
		RestoreTimeout: utils.String(timeout),
	}
}

func FlattenVirtualMachineScaleSetSpotRestorePolicy(input *virtualmachinescalesets.SpotRestorePolicy) []interface{} {
	if input == nil {
		return nil
	}

	var enabled bool
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	var restore string
	if input.RestoreTimeout != nil {
		restore = *input.RestoreTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
			"timeout": restore,
		},
	}
}

func VirtualMachineScaleSetNetworkInterfaceSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ip_configuration": virtualMachineScaleSetIPConfigurationSchemaForDataSource(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				"network_security_group_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func virtualMachineScaleSetIPConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
					},
					Set:      pluginsdk.HashString,
					MaxItems: 20,
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"load_balancer_inbound_nat_rules_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"public_ip_address": virtualMachineScaleSetPublicIPAddressSchema(),

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.IPVersionIPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPVersionIPv4),
						string(compute.IPVersionIPv6),
					}, false),
				},
			},
		},
	}
}

func virtualMachineScaleSetIPConfigurationSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"load_balancer_inbound_nat_rules_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"public_ip_address": virtualMachineScaleSetPublicIPAddressSchemaForDataSource(),

				"subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func virtualMachineScaleSetPublicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},
				"ip_tag": {
					// TODO: does this want to be a Set?
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(compute.IPVersionIPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPVersionIPv4),
						string(compute.IPVersionIPv6),
					}, false),
				},
				// TODO: preview feature
				// $ az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress
				// $ az provider register -n Microsoft.Network
				"public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: networkValidate.PublicIpPrefixID,
				},
			},
		},
	}
}

func virtualMachineScaleSetPublicIPAddressSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"domain_name_label": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"idle_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"ip_tag": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"type": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"public_ip_prefix_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetNetworkInterface(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandVirtualMachineScaleSetIPConfiguration(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration{
			Name: raw["name"].(string),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandVirtualMachineScaleSetIPConfiguration(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, error) {
	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))
	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               pointer.To(version),
			ApplicationGatewayBackendAddressPools: expandIDsToSubResources(raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()),
			ApplicationSecurityGroups:             expandIDsToSubResources(raw["application_security_group_ids"].(*pluginsdk.Set).List()),
			LoadBalancerBackendAddressPools:       expandIDsToSubResources(raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()),
			LoadBalancerInboundNatPools:           expandIDsToSubResources(raw["load_balancer_inbound_nat_rules_ids"].(*pluginsdk.Set).List()),
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandVirtualMachineScaleSetPublicIPAddress(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandVirtualMachineScaleSetPublicIPAddress(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	ipTags := make([]virtualmachinescalesets.VirtualMachineScaleSetIPTag, 0)
	for _, ipTagV := range raw["ip_tag"].([]interface{}) {
		ipTagRaw := ipTagV.(map[string]interface{})
		ipTags = append(ipTags, virtualmachinescalesets.VirtualMachineScaleSetIPTag{
			Tag:       utils.String(ipTagRaw["tag"].(string)),
			IPTagType: utils.String(ipTagRaw["type"].(string)),
		})
	}

	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			IPTags:                 pointer.To(ipTags),
			PublicIPAddressVersion: pointer.To(virtualmachinescalesets.IPVersion(raw["version"].(string))),
		},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = utils.Int64(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	if publicIPPrefixID := raw["public_ip_prefix_id"].(string); publicIPPrefixID != "" {
		publicIPAddressConfig.Properties.PublicIPPrefix = &virtualmachinescalesets.SubResource{
			Id: utils.String(publicIPPrefixID),
		}
	}

	return &publicIPAddressConfig
}

func ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandVirtualMachineScaleSetIPConfigurationUpdate(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration{
			Name: utils.String(raw["name"].(string)),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            &ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandVirtualMachineScaleSetIPConfigurationUpdate(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	loadBalancerInboundNatPoolIdsRaw := raw["load_balancer_inbound_nat_rules_ids"].(*pluginsdk.Set).List()
	loadBalancerInboundNatPoolIds := expandIDsToSubResources(loadBalancerInboundNatPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))

	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration{
		Name: utils.String(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               pointer.To(version),
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			LoadBalancerInboundNatPools:           loadBalancerInboundNatPoolIds,
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandVirtualMachineScaleSetPublicIPAddressUpdate(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandVirtualMachineScaleSetPublicIPAddressUpdate(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration {
	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration{
		Name:       utils.String(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties{},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = utils.Int64(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	return &publicIPAddressConfig
}

func FlattenVirtualMachineScaleSetNetworkInterface(input *[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		if v.Properties == nil {
			continue
		}
		var networkSecurityGroupId string
		if v.Properties.NetworkSecurityGroup != nil && v.Properties.NetworkSecurityGroup.Id != nil {
			networkSecurityGroupId = *v.Properties.NetworkSecurityGroup.Id
		}
		var enableAcceleratedNetworking, enableIPForwarding, primary bool
		if v.Properties.EnableAcceleratedNetworking != nil {
			enableAcceleratedNetworking = *v.Properties.EnableAcceleratedNetworking
		}
		if v.Properties.EnableIPForwarding != nil {
			enableIPForwarding = *v.Properties.EnableIPForwarding
		}
		if v.Properties.Primary != nil {
			primary = *v.Properties.Primary
		}

		var dnsServers []interface{}
		if settings := v.Properties.DnsSettings; settings != nil {
			dnsServers = utils.FlattenStringSlice(v.Properties.DnsSettings.DnsServers)
		}

		var ipConfigurations []interface{}
		for _, configRaw := range v.Properties.IPConfigurations {
			config := flattenVirtualMachineScaleSetIPConfiguration(configRaw)
			ipConfigurations = append(ipConfigurations, config)
		}

		results = append(results, map[string]interface{}{
			"name":                          v.Name,
			"dns_servers":                   dnsServers,
			"enable_accelerated_networking": enableAcceleratedNetworking,
			"enable_ip_forwarding":          enableIPForwarding,
			"ip_configuration":              ipConfigurations,
			"network_security_group_id":     networkSecurityGroupId,
			"primary":                       primary,
		})
	}

	return results
}

func flattenVirtualMachineScaleSetIPConfiguration(input virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration) map[string]interface{} {
	if input.Properties == nil {
		return map[string]interface{}{}
	}
	var subnetId string

	if input.Properties.Subnet != nil && input.Properties.Subnet.Id != nil {
		subnetId = *input.Properties.Subnet.Id
	}

	var primary bool
	if input.Properties.Primary != nil {
		primary = *input.Properties.Primary
	}

	var publicIPAddresses []interface{}
	if input.Properties.PublicIPAddressConfiguration != nil {
		publicIPAddresses = append(publicIPAddresses, flattenVirtualMachineScaleSetPublicIPAddress(*input.Properties.PublicIPAddressConfiguration))
	}

	applicationGatewayBackendAddressPoolIds := flattenSubResourcesToIDs(input.Properties.ApplicationGatewayBackendAddressPools)
	applicationSecurityGroupIds := flattenSubResourcesToIDs(input.Properties.ApplicationSecurityGroups)
	loadBalancerBackendAddressPoolIds := flattenSubResourcesToIDs(input.Properties.LoadBalancerBackendAddressPools)
	loadBalancerInboundNatRuleIds := flattenSubResourcesToIDs(input.Properties.LoadBalancerInboundNatPools)

	return map[string]interface{}{
		"name":              input.Name,
		"primary":           primary,
		"public_ip_address": publicIPAddresses,
		"subnet_id":         subnetId,
		"version":           string(pointer.From(input.Properties.PrivateIPAddressVersion)),
		"application_gateway_backend_address_pool_ids": applicationGatewayBackendAddressPoolIds,
		"application_security_group_ids":               applicationSecurityGroupIds,
		"load_balancer_backend_address_pool_ids":       loadBalancerBackendAddressPoolIds,
		"load_balancer_inbound_nat_rules_ids":          loadBalancerInboundNatRuleIds,
	}
}

func flattenVirtualMachineScaleSetPublicIPAddress(input virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration) map[string]interface{} {
	if input.Properties == nil {
		return map[string]interface{}{}
	}
	ipTags := make([]interface{}, 0)
	if input.Properties.IPTags != nil {
		for _, rawTag := range *input.Properties.IPTags {
			var tag, tagType string

			if rawTag.IPTagType != nil {
				tagType = *rawTag.IPTagType
			}

			if rawTag.Tag != nil {
				tag = *rawTag.Tag
			}

			ipTags = append(ipTags, map[string]interface{}{
				"tag":  tag,
				"type": tagType,
			})
		}
	}

	var domainNameLabel, publicIPPrefixId, version string
	if input.Properties.DnsSettings != nil {
		domainNameLabel = input.Properties.DnsSettings.DomainNameLabel
	}

	if input.Properties.PublicIPPrefix != nil && input.Properties.PublicIPPrefix.Id != nil {
		publicIPPrefixId = *input.Properties.PublicIPPrefix.Id
	}

	if pointer.From(input.Properties.PublicIPAddressVersion) != "" {
		version = string(pointer.From(input.Properties.PublicIPAddressVersion))
	}

	var idleTimeoutInMinutes int
	if input.Properties.IdleTimeoutInMinutes != nil {
		idleTimeoutInMinutes = int(*input.Properties.IdleTimeoutInMinutes)
	}

	return map[string]interface{}{
		"name":                    input.Name,
		"domain_name_label":       domainNameLabel,
		"idle_timeout_in_minutes": idleTimeoutInMinutes,
		"ip_tag":                  ipTags,
		"public_ip_prefix_id":     publicIPPrefixId,
		"version":                 version,
	}
}

func VirtualMachineScaleSetDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: does this want to be a Set?
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},

				"create_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.DiskCreateOptionTypesEmpty),
						string(compute.DiskCreateOptionTypesFromImage),
					}, false),
					Default: string(compute.DiskCreateOptionTypesEmpty),
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 2000), // TODO: confirm upper bounds
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesPremiumV2LRS),
						string(compute.StorageAccountTypesPremiumZRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesStandardSSDZRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// TODO rename `ultra_ssd_disk_iops_read_write` to `disk_iops_read_write` in 4.0
				"ultra_ssd_disk_iops_read_write": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Computed:     true,
				},

				// TODO rename `ultra_ssd_disk_mbps_read_write` to `disk_mbps_read_write` in 4.0
				"ultra_ssd_disk_mbps_read_write": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Computed:     true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetDataDisk(input []interface{}, ultraSSDEnabled bool) (*[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, error) {
	disks := make([]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		storageAccountType := virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))
		disk := virtualmachinescalesets.VirtualMachineScaleSetDataDisk{
			Caching:    pointer.To(virtualmachinescalesets.CachingTypes(raw["caching"].(string))),
			DiskSizeGB: utils.Int64(int64(raw["disk_size_gb"].(int))),
			Lun:        int64(raw["lun"].(int)),
			ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: pointer.To(storageAccountType),
			},
			WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
			CreateOption:            virtualmachinescalesets.DiskCreateOptionTypes(raw["create_option"].(string)),
		}

		if name := raw["name"]; name != nil && name.(string) != "" {
			disk.Name = utils.String(name.(string))
		}

		if id := raw["disk_encryption_set_id"].(string); id != "" {
			disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
				Id: utils.String(id),
			}
		}

		var iops int
		if diskIops, ok := raw["ultra_ssd_disk_iops_read_write"]; ok && diskIops.(int) > 0 {
			iops = diskIops.(int)
		}

		if iops > 0 && !ultraSSDEnabled && storageAccountType != virtualmachinescalesets.StorageAccountTypesPremiumVTwoLRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_iops_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		var mbps int
		if diskMbps, ok := raw["ultra_ssd_disk_mbps_read_write"]; ok && diskMbps.(int) > 0 {
			mbps = diskMbps.(int)
		}

		if mbps > 0 && !ultraSSDEnabled && storageAccountType != virtualmachinescalesets.StorageAccountTypesPremiumVTwoLRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_mbps_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if iops > 0 {
			disk.DiskIOPSReadWrite = utils.Int64(int64(iops))
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if mbps > 0 {
			disk.DiskMBpsReadWrite = utils.Int64(int64(mbps))
		}

		disks = append(disks, disk)
	}

	return &disks, nil
}

func FlattenVirtualMachineScaleSetDataDisk(input *[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		var name string
		if v.Name != nil {
			name = *v.Name
		}

		diskSizeGb := 0
		if v.DiskSizeGB != nil && *v.DiskSizeGB != 0 {
			diskSizeGb = int(*v.DiskSizeGB)
		}

		storageAccountType := ""
		diskEncryptionSetId := ""
		if v.ManagedDisk != nil {
			storageAccountType = string(pointer.From(v.ManagedDisk.StorageAccountType))
			if v.ManagedDisk.DiskEncryptionSet != nil && v.ManagedDisk.DiskEncryptionSet.Id != nil {
				diskEncryptionSetId = *v.ManagedDisk.DiskEncryptionSet.Id
			}
		}

		writeAcceleratorEnabled := false
		if v.WriteAcceleratorEnabled != nil {
			writeAcceleratorEnabled = *v.WriteAcceleratorEnabled
		}

		iops := 0
		if v.DiskIOPSReadWrite != nil {
			iops = int(*v.DiskIOPSReadWrite)
		}

		mbps := 0
		if v.DiskMBpsReadWrite != nil {
			mbps = int(*v.DiskMBpsReadWrite)
		}

		dataDisk := map[string]interface{}{
			"name":                           name,
			"caching":                        string(pointer.From(v.Caching)),
			"create_option":                  string(v.CreateOption),
			"lun":                            v.Lun,
			"disk_encryption_set_id":         diskEncryptionSetId,
			"disk_size_gb":                   diskSizeGb,
			"storage_account_type":           storageAccountType,
			"ultra_ssd_disk_iops_read_write": iops,
			"ultra_ssd_disk_mbps_read_write": mbps,
			"write_accelerator_enabled":      writeAcceleratorEnabled,
		}

		output = append(output, dataDisk)
	}

	return output
}

func VirtualMachineScaleSetOSDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// note: OS Disks don't support Ultra SSDs or PremiumV2_LRS
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesPremiumZRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesStandardSSDZRS),
					}, false),
				},

				"diff_disk_settings": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskOptionsLocal),
								}, false),
							},
							"placement": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
								Default:  string(compute.DiffDiskPlacementCacheDisk),
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskPlacementCacheDisk),
									string(compute.DiffDiskPlacementResourceDisk),
								}, false),
							},
						},
					},
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:      true,
					ValidateFunc:  validate.DiskEncryptionSetID,
					ConflictsWith: []string{"os_disk.0.secure_vm_disk_encryption_set_id"},
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},

				"secure_vm_disk_encryption_set_id": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ForceNew:      true,
					ValidateFunc:  validate.DiskEncryptionSetID,
					ConflictsWith: []string{"os_disk.0.disk_encryption_set_id"},
				},

				"security_encryption_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.SecurityEncryptionTypesVMGuestStateOnly),
						string(compute.SecurityEncryptionTypesDiskWithVMGuestState),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetOSDisk(input []interface{}, osType virtualmachinescalesets.OperatingSystemTypes) (*virtualmachinescalesets.VirtualMachineScaleSetOSDisk, error) {
	raw := input[0].(map[string]interface{})
	caching := raw["caching"].(string)
	disk := virtualmachinescalesets.VirtualMachineScaleSetOSDisk{
		Caching: pointer.To(virtualmachinescalesets.CachingTypes(caching)),
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))),
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),

		// these have to be hard-coded so there's no point exposing them
		CreateOption: virtualmachinescalesets.DiskCreateOptionTypesFromImage,
		OsType:       pointer.To(osType),
	}

	securityEncryptionType := raw["security_encryption_type"].(string)
	if securityEncryptionType != "" {
		disk.ManagedDisk.SecurityProfile = &virtualmachinescalesets.VMDiskSecurityProfile{
			SecurityEncryptionType: pointer.To(virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType)),
		}
	}
	if secureVMDiskEncryptionId := raw["secure_vm_disk_encryption_set_id"].(string); secureVMDiskEncryptionId != "" {
		if virtualmachinescalesets.SecurityEncryptionTypesDiskWithVMGuestState != virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType) {
			return nil, fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		disk.ManagedDisk.SecurityProfile.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(secureVMDiskEncryptionId),
		}
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int64(int64(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		if caching != string(compute.CachingTypesReadOnly) {
			// Restriction per https://docs.microsoft.com/azure/virtual-machines/ephemeral-os-disks-deploy#vm-template-deployment
			return nil, fmt.Errorf("`diff_disk_settings` can only be set when `caching` is set to `ReadOnly`")
		}

		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		disk.DiffDiskSettings = &virtualmachinescalesets.DiffDiskSettings{
			Option:    pointer.To(virtualmachinescalesets.DiffDiskOptions(diffDiskRaw["option"].(string))),
			Placement: pointer.To(virtualmachinescalesets.DiffDiskPlacement(diffDiskRaw["placement"].(string))),
		}
	}

	return &disk, nil
}

func ExpandVirtualMachineScaleSetOSDiskUpdate(input []interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk {
	raw := input[0].(map[string]interface{})
	disk := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk{
		Caching: pointer.To(virtualmachinescalesets.CachingTypes(raw["caching"].(string))),
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))),
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int64(int64(osDiskSize))
	}

	return &disk
}

func FlattenVirtualMachineScaleSetOSDisk(input *virtualmachinescalesets.VirtualMachineScaleSetOSDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option":    string(pointer.From(input.DiffDiskSettings.Option)),
			"placement": string(pointer.From(input.DiffDiskSettings.Placement)),
		})
	}

	diskSizeGb := 0
	if input.DiskSizeGB != nil && *input.DiskSizeGB != 0 {
		diskSizeGb = int(*input.DiskSizeGB)
	}

	storageAccountType := ""
	diskEncryptionSetId := ""
	secureVMDiskEncryptionSetId := ""
	securityEncryptionType := ""
	if input.ManagedDisk != nil {
		storageAccountType = string(pointer.From(input.ManagedDisk.StorageAccountType))
		if input.ManagedDisk.DiskEncryptionSet != nil && input.ManagedDisk.DiskEncryptionSet.Id != nil {
			diskEncryptionSetId = *input.ManagedDisk.DiskEncryptionSet.Id
		}

		if securityProfile := input.ManagedDisk.SecurityProfile; securityProfile != nil {
			securityEncryptionType = string(pointer.From(securityProfile.SecurityEncryptionType))
			if securityProfile.DiskEncryptionSet != nil && securityProfile.DiskEncryptionSet.Id != nil {
				secureVMDiskEncryptionSetId = *securityProfile.DiskEncryptionSet.Id
			}
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"caching":                          string(pointer.From(input.Caching)),
			"disk_size_gb":                     diskSizeGb,
			"diff_disk_settings":               diffDiskSettings,
			"storage_account_type":             storageAccountType,
			"write_accelerator_enabled":        writeAcceleratorEnabled,
			"disk_encryption_set_id":           diskEncryptionSetId,
			"secure_vm_disk_encryption_set_id": secureVMDiskEncryptionSetId,
			"security_encryption_type":         securityEncryptionType,
		},
	}
}

func VirtualMachineScaleSetAutomatedOSUpgradePolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO: should these be optional + defaulted?
				"disable_automatic_rollback": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_automatic_os_upgrade": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(input []interface{}) *virtualmachinescalesets.AutomaticOSUpgradePolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &virtualmachinescalesets.AutomaticOSUpgradePolicy{
		DisableAutomaticRollback: utils.Bool(raw["disable_automatic_rollback"].(bool)),
		EnableAutomaticOSUpgrade: utils.Bool(raw["enable_automatic_os_upgrade"].(bool)),
	}
}

func FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(input *virtualmachinescalesets.AutomaticOSUpgradePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	disableAutomaticRollback := false
	if input.DisableAutomaticRollback != nil {
		disableAutomaticRollback = *input.DisableAutomaticRollback
	}

	enableAutomaticOSUpgrade := false
	if input.EnableAutomaticOSUpgrade != nil {
		enableAutomaticOSUpgrade = *input.EnableAutomaticOSUpgrade
	}

	return []interface{}{
		map[string]interface{}{
			"disable_automatic_rollback":  disableAutomaticRollback,
			"enable_automatic_os_upgrade": enableAutomaticOSUpgrade,
		},
	}
}

func VirtualMachineScaleSetRollingUpgradePolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cross_zone_upgrades_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"max_batch_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"max_unhealthy_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"max_unhealthy_upgraded_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"pause_time_between_batches": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azValidate.ISO8601Duration,
				},
				"prioritize_unhealthy_instances_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetRollingUpgradePolicy(input []interface{}, isZonal bool) (*virtualmachinescalesets.RollingUpgradePolicy, error) {
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	rollingUpgradePolicy := &virtualmachinescalesets.RollingUpgradePolicy{
		MaxBatchInstancePercent:             utils.Int64(int64(raw["max_batch_instance_percent"].(int))),
		MaxUnhealthyInstancePercent:         utils.Int64(int64(raw["max_unhealthy_instance_percent"].(int))),
		MaxUnhealthyUpgradedInstancePercent: utils.Int64(int64(raw["max_unhealthy_upgraded_instance_percent"].(int))),
		PauseTimeBetweenBatches:             utils.String(raw["pause_time_between_batches"].(string)),
		PrioritizeUnhealthyInstances:        utils.Bool(raw["prioritize_unhealthy_instances_enabled"].(bool)),
	}

	enableCrossZoneUpgrade := raw["cross_zone_upgrades_enabled"].(bool)
	if isZonal {
		// EnableCrossZoneUpgrade can only be set when for zonal scale set
		rollingUpgradePolicy.EnableCrossZoneUpgrade = utils.Bool(enableCrossZoneUpgrade)
	} else if enableCrossZoneUpgrade {
		return nil, fmt.Errorf("`rolling_upgrade_policy.0.cross_zone_upgrades_enabled` can only be set to `true` when `zones` is specified")
	}

	return rollingUpgradePolicy, nil
}

func FlattenVirtualMachineScaleSetRollingUpgradePolicy(input *virtualmachinescalesets.RollingUpgradePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enableCrossZoneUpgrade := false
	if input.EnableCrossZoneUpgrade != nil {
		enableCrossZoneUpgrade = *input.EnableCrossZoneUpgrade
	}

	maxBatchInstancePercent := 0
	if input.MaxBatchInstancePercent != nil {
		maxBatchInstancePercent = int(*input.MaxBatchInstancePercent)
	}

	maxUnhealthyInstancePercent := 0
	if input.MaxUnhealthyInstancePercent != nil {
		maxUnhealthyInstancePercent = int(*input.MaxUnhealthyInstancePercent)
	}

	maxUnhealthyUpgradedInstancePercent := 0
	if input.MaxUnhealthyUpgradedInstancePercent != nil {
		maxUnhealthyUpgradedInstancePercent = int(*input.MaxUnhealthyUpgradedInstancePercent)
	}

	pauseTimeBetweenBatches := ""
	if input.PauseTimeBetweenBatches != nil {
		pauseTimeBetweenBatches = *input.PauseTimeBetweenBatches
	}

	prioritizeUnhealthyInstances := false
	if input.PrioritizeUnhealthyInstances != nil {
		prioritizeUnhealthyInstances = *input.PrioritizeUnhealthyInstances
	}

	return []interface{}{
		map[string]interface{}{
			"cross_zone_upgrades_enabled":             enableCrossZoneUpgrade,
			"max_batch_instance_percent":              maxBatchInstancePercent,
			"max_unhealthy_instance_percent":          maxUnhealthyInstancePercent,
			"max_unhealthy_upgraded_instance_percent": maxUnhealthyUpgradedInstancePercent,
			"pause_time_between_batches":              pauseTimeBetweenBatches,
			"prioritize_unhealthy_instances_enabled":  prioritizeUnhealthyInstances,
		},
	}
}

// TODO remove VirtualMachineScaleSetTerminateNotificationSchema in 4.0
func VirtualMachineScaleSetTerminateNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		MaxItems:   1,
		Deprecated: "`terminate_notification` has been renamed to `termination_notification` and will be removed in 4.0.",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
		ConflictsWith: []string{"termination_notification"},
	}
}

func VirtualMachineScaleSetTerminationNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
		ConflictsWith: func() []string {
			if !features.FourPointOhBeta() {
				return []string{"terminate_notification"}
			}
			return []string{}
		}(),
	}
}

func ExpandVirtualMachineScaleSetScheduledEventsProfile(input []interface{}) *virtualmachinescalesets.ScheduledEventsProfile {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	enabled := raw["enabled"].(bool)
	timeout := raw["timeout"].(string)

	return &virtualmachinescalesets.ScheduledEventsProfile{
		TerminateNotificationProfile: &virtualmachinescalesets.TerminateNotificationProfile{
			Enable:           &enabled,
			NotBeforeTimeout: &timeout,
		},
	}
}

func FlattenVirtualMachineScaleSetScheduledEventsProfile(input *virtualmachinescalesets.ScheduledEventsProfile) []interface{} {
	// if enabled is set to false, there will be no ScheduledEventsProfile in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.Enable != nil {
		enabled = *input.TerminateNotificationProfile.Enable
	}

	timeout := "PT5M"
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.NotBeforeTimeout != nil {
		timeout = *input.TerminateNotificationProfile.NotBeforeTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
			"timeout": timeout,
		},
	}
}

func VirtualMachineScaleSetAutomaticRepairsPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"grace_period": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT30M",
					// this field actually has a range from 30m to 90m, is there a function that can do this validation?
					ValidateFunc: azValidate.ISO8601Duration,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(input []interface{}) *virtualmachinescalesets.AutomaticRepairsPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &virtualmachinescalesets.AutomaticRepairsPolicy{
		Enabled:     utils.Bool(raw["enabled"].(bool)),
		GracePeriod: utils.String(raw["grace_period"].(string)),
	}
}

func FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(input *virtualmachinescalesets.AutomaticRepairsPolicy) []interface{} {
	// if enabled is set to false, there will be no AutomaticRepairsPolicy in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.Enabled != nil {
		enabled = *input.Enabled
	}

	gracePeriod := "PT30M"
	if input != nil && input.GracePeriod != nil {
		gracePeriod = *input.GracePeriod
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":      enabled,
			"grace_period": gracePeriod,
		},
	}
}

func VirtualMachineScaleSetExtensionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"auto_upgrade_minor_version": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"force_update_tag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				// Need to check `protected_settings_from_key_vault` conflicting with `protected_settings` in iteration
				"protected_settings_from_key_vault": protectedSettingsFromKeyVaultSchema(false),

				"provision_after_extensions": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"settings": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},
			},
		},
		Set: virtualMachineScaleSetExtensionHash,
	}
}

func virtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["auto_upgrade_minor_version"].(bool)))

		if v, ok = m["force_update_tag"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}

		if v, ok := m["provision_after_extensions"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
			if err == nil {
				serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serializedSettings))
				}
			}
		}

		if v, ok := m["protected_settings"]; ok {
			settings := v.(string)
			if settings != "" {
				expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
				if err == nil {
					serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
					if err == nil {
						buf.WriteString(fmt.Sprintf("%s-", serializedSettings))
					}
				}
			}
		}

		if v, ok := m["protected_settings_from_key_vault"]; ok {
			protectedSettingsFromKeyVault := v.([]interface{})
			if len(protectedSettingsFromKeyVault) > 0 {
				buf.WriteString(fmt.Sprintf("%s-", protectedSettingsFromKeyVault[0].(map[string]interface{})["secret_url"].(string)))
				buf.WriteString(fmt.Sprintf("%s-", protectedSettingsFromKeyVault[0].(map[string]interface{})["source_vault_id"].(string)))
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func expandVirtualMachineScaleSetExtensions(input []interface{}) (extensionProfile *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, hasHealthExtension bool, err error) {
	extensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
	if len(input) == 0 {
		return extensionProfile, false, nil
	}

	extensions := make([]virtualmachinescalesets.VirtualMachineScaleSetExtension, 0)
	for _, v := range input {
		extensionRaw := v.(map[string]interface{})
		extension := virtualmachinescalesets.VirtualMachineScaleSetExtension{
			Name: utils.String(extensionRaw["name"].(string)),
		}
		extensionType := extensionRaw["type"].(string)

		extensionProps := virtualmachinescalesets.VirtualMachineScaleSetExtensionProperties{
			Publisher:                utils.String(extensionRaw["publisher"].(string)),
			Type:                     &extensionType,
			TypeHandlerVersion:       utils.String(extensionRaw["type_handler_version"].(string)),
			AutoUpgradeMinorVersion:  utils.Bool(extensionRaw["auto_upgrade_minor_version"].(bool)),
			EnableAutomaticUpgrade:   utils.Bool(extensionRaw["automatic_upgrade_enabled"].(bool)),
			ProvisionAfterExtensions: utils.ExpandStringSlice(extensionRaw["provision_after_extensions"].([]interface{})),
		}

		if extensionType == "ApplicationHealthLinux" || extensionType == "ApplicationHealthWindows" {
			hasHealthExtension = true
		}

		if forceUpdateTag := extensionRaw["force_update_tag"]; forceUpdateTag != nil {
			extensionProps.ForceUpdateTag = utils.String(forceUpdateTag.(string))
		}

		if val, ok := extensionRaw["settings"]; ok && val.(string) != "" {
			settings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON from `settings`: %+v", err)
			}
			extensionProps.Settings = pointer.To(interface{}(settings))
		}

		protectedSettingsFromKeyVault := expandProtectedSettingsFromKeyVault(extensionRaw["protected_settings_from_key_vault"].([]interface{}))
		extensionProps.ProtectedSettingsFromKeyVault = protectedSettingsFromKeyVault

		if val, ok := extensionRaw["protected_settings"]; ok && val.(string) != "" {
			if protectedSettingsFromKeyVault != nil {
				return nil, false, fmt.Errorf("`protected_settings_from_key_vault` cannot be used with `protected_settings`")
			}

			protectedSettings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON from `protected_settings`: %+v", err)
			}
			extensionProps.ProtectedSettings = pointer.To(interface{}(protectedSettings))
		}

		extension.Properties = &extensionProps
		extensions = append(extensions, extension)
	}
	extensionProfile.Extensions = &extensions

	return extensionProfile, hasHealthExtension, nil
}

func flattenVirtualMachineScaleSetExtensions(input *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, d *pluginsdk.ResourceData) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if input == nil || input.Extensions == nil {
		return result, nil
	}

	// extensionsFromState holds the "extension" block, which is used to retrieve the "protected_settings" to fill it back the state,
	// since it is not returned from the API.
	extensionsFromState := map[string]map[string]interface{}{}
	if extSet, ok := d.GetOk("extension"); ok && extSet != nil {
		extensions := extSet.(*pluginsdk.Set).List()
		for _, ext := range extensions {
			if ext == nil {
				continue
			}
			ext := ext.(map[string]interface{})
			extensionsFromState[ext["name"].(string)] = ext
		}
	}

	for _, v := range *input.Extensions {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		autoUpgradeMinorVersion := false
		enableAutomaticUpgrade := false
		forceUpdateTag := ""
		provisionAfterExtension := make([]interface{}, 0)
		protectedSettings := ""
		var protectedSettingsFromKeyVault *virtualmachinescalesets.KeyVaultSecretReference
		extPublisher := ""
		extSettings := ""
		extType := ""
		extTypeVersion := ""

		if props := v.Properties; props != nil {
			if props.Publisher != nil {
				extPublisher = *props.Publisher
			}

			if props.Type != nil {
				extType = *props.Type
			}

			if props.TypeHandlerVersion != nil {
				extTypeVersion = *props.TypeHandlerVersion
			}

			if props.AutoUpgradeMinorVersion != nil {
				autoUpgradeMinorVersion = *props.AutoUpgradeMinorVersion
			}

			if props.EnableAutomaticUpgrade != nil {
				enableAutomaticUpgrade = *props.EnableAutomaticUpgrade
			}

			if props.ForceUpdateTag != nil {
				forceUpdateTag = *props.ForceUpdateTag
			}

			if props.ProvisionAfterExtensions != nil {
				provisionAfterExtension = utils.FlattenStringSlice(props.ProvisionAfterExtensions)
			}

			if props.Settings != nil {
				settingsRaw := *props.Settings
				if settings, ok := settingsRaw.(map[string]interface{}); ok {
					extSettingsRaw, err := pluginsdk.FlattenJsonToString(settings)
					if err != nil {
						return nil, err
					}
					extSettings = extSettingsRaw
				}
			}

			protectedSettingsFromKeyVault = props.ProtectedSettingsFromKeyVault
		}
		// protected_settings isn't returned, so we attempt to get it from state otherwise set to empty string
		if ext, ok := extensionsFromState[name]; ok {
			if protectedSettingsFromState, ok := ext["protected_settings"]; ok {
				if protectedSettingsFromState.(string) != "" && protectedSettingsFromState.(string) != "{}" {
					protectedSettings = protectedSettingsFromState.(string)
				}
			}
		}

		result = append(result, map[string]interface{}{
			"name":                              name,
			"auto_upgrade_minor_version":        autoUpgradeMinorVersion,
			"automatic_upgrade_enabled":         enableAutomaticUpgrade,
			"force_update_tag":                  forceUpdateTag,
			"provision_after_extensions":        provisionAfterExtension,
			"protected_settings":                protectedSettings,
			"protected_settings_from_key_vault": flattenProtectedSettingsFromKeyVault(protectedSettingsFromKeyVault),
			"publisher":                         extPublisher,
			"settings":                          extSettings,
			"type":                              extType,
			"type_handler_version":              extTypeVersion,
		})
	}
	return result, nil
}
