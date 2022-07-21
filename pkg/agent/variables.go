// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package agent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/agentbaker/pkg/agent/datamodel"
	"github.com/blang/semver"
)

var dockerShimFlags = []string{"--cni-bin-dir", "--cni-cache-dir", "--cni-conf-dir", "--docker-endpoint", "--image-pull-progress-deadline", "--network-plugin", "--network-plugin-mtu"}

// getCustomDataVariables returns cloudinit data used by Linux
func getCustomDataVariables(config *datamodel.NodeBootstrappingConfiguration) paramsMap {
	cs := config.ContainerService
	cloudInitFiles := map[string]interface{}{
		"cloudInitData": paramsMap{
			"provisionStartScript":         getBase64EncodedGzippedCustomScript(kubernetesCSEStartScript, config),
			"provisionScript":              getBase64EncodedGzippedCustomScript(kubernetesCSEMainScript, config),
			"provisionSource":              getBase64EncodedGzippedCustomScript(kubernetesCSEHelpersScript, config),
			"provisionSourceUbuntu":        getBase64EncodedGzippedCustomScript(kubernetesCSEHelpersScriptUbuntu, config),
			"provisionSourceMariner":       getBase64EncodedGzippedCustomScript(kubernetesCSEHelpersScriptMariner, config),
			"provisionInstalls":            getBase64EncodedGzippedCustomScript(kubernetesCSEInstall, config),
			"provisionInstallsUbuntu":      getBase64EncodedGzippedCustomScript(kubernetesCSEInstallUbuntu, config),
			"provisionInstallsMariner":     getBase64EncodedGzippedCustomScript(kubernetesCSEInstallMariner, config),
			"provisionConfigs":             getBase64EncodedGzippedCustomScript(kubernetesCSEConfig, config),
			"provisionSendLogs":            getBase64EncodedGzippedCustomScript(kubernetesCSESendLogs, config),
			"provisionRedactCloudConfig":   getBase64EncodedGzippedCustomScript(kubernetesCSERedactCloudConfig, config),
			"customSearchDomainsScript":    getBase64EncodedGzippedCustomScript(kubernetesCustomSearchDomainsScript, config),
			"dhcpv6SystemdService":         getBase64EncodedGzippedCustomScript(dhcpv6SystemdService, config),
			"dhcpv6ConfigurationScript":    getBase64EncodedGzippedCustomScript(dhcpv6ConfigurationScript, config),
			"kubeletSystemdService":        getBase64EncodedGzippedCustomScript(kubeletSystemdService, config),
			"krustletSystemdService":       getBase64EncodedGzippedCustomScript(krustletSystemdService, config),
			"reconcilePrivateHostsScript":  getBase64EncodedGzippedCustomScript(reconcilePrivateHostsScript, config),
			"reconcilePrivateHostsService": getBase64EncodedGzippedCustomScript(reconcilePrivateHostsService, config),
			"ensureNoDupEbtablesScript":    getBase64EncodedGzippedCustomScript(ensureNoDupEbtablesScript, config),
			"ensureNoDupEbtablesService":   getBase64EncodedGzippedCustomScript(ensureNoDupEbtablesService, config),
			"bindMountScript":              getBase64EncodedGzippedCustomScript(bindMountScript, config),
			"bindMountSystemdService":      getBase64EncodedGzippedCustomScript(bindMountSystemdService, config),
			"migPartitionSystemdService":   getBase64EncodedGzippedCustomScript(migPartitionSystemdService, config),
			"migPartitionScript":           getBase64EncodedGzippedCustomScript(migPartitionScript, config),
			"containerdKubeletDropin":      getBase64EncodedGzippedCustomScript(containerdKubeletDropin, config),
			"cgroupv2KubeletDropin":        getBase64EncodedGzippedCustomScript(cgroupv2KubeletDropin, config),
			"componentConfigDropin":        getBase64EncodedGzippedCustomScript(componentConfigDropin, config),
			"tlsBootstrapDropin":           getBase64EncodedGzippedCustomScript(tlsBootstrapDropin, config),
			"bindMountDropin":              getBase64EncodedGzippedCustomScript(bindMountDropin, config),
			"httpProxyDropin":              getBase64EncodedGzippedCustomScript(httpProxyDropin, config),
			"componentManifestFile":        getBase64EncodedGzippedCustomScript(componentManifestFile, config),
			"syncTunnelLogsScript":         getBase64EncodedGzippedCustomScript(syncTunnelLogsScript, config),
		},
	}

	cloudInitData := cloudInitFiles["cloudInitData"].(paramsMap)
	if cs.IsAKSCustomCloud() {
		if strings.EqualFold(string(config.OSSKU), string("CBLMariner")) {
			cloudInitData["initAKSCustomCloud"] = getBase64EncodedGzippedCustomScript(initAKSCustomCloudMarinerScript, config)
		} else {
			cloudInitData["initAKSCustomCloud"] = getBase64EncodedGzippedCustomScript(initAKSCustomCloudScript, config)
		}
	}

	if !cs.Properties.IsVHDDistroForAllNodes() {
		cloudInitData["provisionCIS"] = getBase64EncodedGzippedCustomScript(kubernetesCISScript, config)
		cloudInitData["kmsSystemdService"] = getBase64EncodedGzippedCustomScript(kmsSystemdService, config)
		cloudInitData["aptPreferences"] = getBase64EncodedGzippedCustomScript(aptPreferences, config)
		cloudInitData["healthMonitorScript"] = getBase64EncodedGzippedCustomScript(kubernetesHealthMonitorScript, config)
		cloudInitData["kubeletMonitorSystemdService"] = getBase64EncodedGzippedCustomScript(kubernetesKubeletMonitorSystemdService, config)
		cloudInitData["dockerMonitorSystemdService"] = getBase64EncodedGzippedCustomScript(kubernetesDockerMonitorSystemdService, config)
		cloudInitData["dockerMonitorSystemdTimer"] = getBase64EncodedGzippedCustomScript(kubernetesDockerMonitorSystemdTimer, config)
		cloudInitData["containerdMonitorSystemdService"] = getBase64EncodedGzippedCustomScript(kubernetesContainerdMonitorSystemdService, config)
		cloudInitData["containerdMonitorSystemdTimer"] = getBase64EncodedGzippedCustomScript(kubernetesContainerdMonitorSystemdTimer, config)
		cloudInitData["dockerClearMountPropagationFlags"] = getBase64EncodedGzippedCustomScript(dockerClearMountPropagationFlags, config)
	}

	return cloudInitFiles
}

// getWindowsCustomDataVariables returns custom data for Windows
// TODO(qinhao): combine this function with `getCSECommandVariables` after we support passing variables from cse command to customdata
func getWindowsCustomDataVariables(config *datamodel.NodeBootstrappingConfiguration) paramsMap {
	cs := config.ContainerService
	// these variables is subet of
	customData := map[string]interface{}{
		"tenantID":                             config.TenantID,
		"subscriptionId":                       config.SubscriptionID,
		"resourceGroup":                        config.ResourceGroupName,
		"location":                             cs.Location,
		"vmType":                               cs.Properties.GetVMType(),
		"subnetName":                           cs.Properties.GetSubnetName(),
		"nsgName":                              cs.Properties.GetNSGName(),
		"virtualNetworkName":                   cs.Properties.GetVirtualNetworkName(),
		"routeTableName":                       cs.Properties.GetRouteTableName(),
		"primaryAvailabilitySetName":           cs.Properties.GetPrimaryAvailabilitySetName(),
		"primaryScaleSetName":                  config.PrimaryScaleSetName,
		"useManagedIdentityExtension":          useManagedIdentity(cs),
		"useInstanceMetadata":                  useInstanceMetadata(cs),
		"loadBalancerSku":                      cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku,
		"excludeMasterFromStandardLB":          true,
		"enableTelemetry":                      false,
		"windowsEnableCSIProxy":                cs.Properties.WindowsProfile.IsCSIProxyEnabled(),
		"windowsCSIProxyURL":                   cs.Properties.WindowsProfile.CSIProxyURL,
		"windowsProvisioningScriptsPackageURL": cs.Properties.WindowsProfile.ProvisioningScriptsPackageURL,
		"windowsPauseImageURL":                 cs.Properties.WindowsProfile.WindowsPauseImageURL,
		"alwaysPullWindowsPauseImage":          strconv.FormatBool(cs.Properties.WindowsProfile.IsAlwaysPullWindowsPauseImage()),
		"windowsCalicoPackageURL":              cs.Properties.WindowsProfile.WindowsCalicoPackageURL,
		"windowsSecureTlsEnabled":              cs.Properties.WindowsProfile.IsWindowsSecureTlsEnabled(),
		"windowsGmsaPackageUrl":                cs.Properties.WindowsProfile.WindowsGmsaPackageUrl,
		"windowsCSEScriptsPackageURL":          cs.Properties.WindowsProfile.CseScriptsPackageURL,
		"isNotRebootWindowsNode":               strconv.FormatBool(config.AgentPoolProfile.IsNotRebootWindowsNode()),
	}

	return customData
}

func getCSECommandVariables(config *datamodel.NodeBootstrappingConfiguration) paramsMap {
	cs := config.ContainerService
	profile := config.AgentPoolProfile
	return map[string]interface{}{
		"tenantID":                        config.TenantID,
		"subscriptionId":                  config.SubscriptionID,
		"resourceGroup":                   config.ResourceGroupName,
		"location":                        cs.Location,
		"vmType":                          cs.Properties.GetVMType(),
		"subnetName":                      cs.Properties.GetSubnetName(),
		"nsgName":                         cs.Properties.GetNSGName(),
		"virtualNetworkName":              cs.Properties.GetVirtualNetworkName(),
		"virtualNetworkResourceGroupName": cs.Properties.GetVNetResourceGroupName(),
		"routeTableName":                  cs.Properties.GetRouteTableName(),
		"primaryAvailabilitySetName":      cs.Properties.GetPrimaryAvailabilitySetName(),
		"primaryScaleSetName":             config.PrimaryScaleSetName,
		"useManagedIdentityExtension":     useManagedIdentity(cs),
		"useInstanceMetadata":             useInstanceMetadata(cs),
		"loadBalancerSku":                 cs.Properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku,
		"excludeMasterFromStandardLB":     true,
		"maximumLoadBalancerRuleCount":    getMaximumLoadBalancerRuleCount(cs),
		"userAssignedIdentityID":          config.UserAssignedIdentityClientID,
		"isVHD":                           isVHD(profile),
		"gpuNode":                         strconv.FormatBool(config.EnableNvidia),
		"sgxNode":                         strconv.FormatBool(datamodel.IsSgxEnabledSKU(profile.VMSize)),
		"configGPUDriverIfNeeded":         config.ConfigGPUDriverIfNeeded,
		"enableGPUDevicePluginIfNeeded":   config.EnableGPUDevicePluginIfNeeded,
		"migNode":                         strconv.FormatBool(datamodel.IsMIGNode(config.GPUInstanceProfile)),
		"gpuInstanceProfile":              config.GPUInstanceProfile,
	}
}

func useManagedIdentity(cs *datamodel.ContainerService) string {
	useManagedIdentity := cs.Properties.OrchestratorProfile.KubernetesConfig != nil &&
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity
	return strconv.FormatBool(useManagedIdentity)
}

func useInstanceMetadata(cs *datamodel.ContainerService) string {
	useInstanceMetadata := cs.Properties.OrchestratorProfile.KubernetesConfig != nil &&
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseInstanceMetadata != nil &&
		*cs.Properties.OrchestratorProfile.KubernetesConfig.UseInstanceMetadata
	return strconv.FormatBool(useInstanceMetadata)
}

func getMaximumLoadBalancerRuleCount(cs *datamodel.ContainerService) int {
	if cs.Properties.OrchestratorProfile.KubernetesConfig != nil {
		return cs.Properties.OrchestratorProfile.KubernetesConfig.MaximumLoadBalancerRuleCount
	}
	return 0
}

func isVHD(profile *datamodel.AgentPoolProfile) string {
	//NOTE: update as new distro is introduced
	return strconv.FormatBool(profile.IsVHDDistro())
}

func getOutBoundCmd(nbc *datamodel.NodeBootstrappingConfiguration, cloudSpecConfig *datamodel.AzureEnvironmentSpecConfig) string {
	cs := nbc.ContainerService
	if cs.Properties.FeatureFlags.IsFeatureEnabled("BlockOutboundInternet") {
		return ""
	}
	registry := ""
	if cloudSpecConfig.CloudName == datamodel.AzureChinaCloud {
		registry = `gcr.azk8s.cn`
	} else if cs.IsAKSCustomCloud() {
		registry = cs.Properties.CustomCloudEnv.McrURL
	} else {
		registry = `mcr.microsoft.com`
	}

	// curl on Ubuntu 16.04 (shipped prior to AKS 1.18) doesn't support proxy TLS
	// so we need to use nc for the connectivity check
	clusterVersion, _ := semver.Make(cs.Properties.OrchestratorProfile.OrchestratorVersion)
	minVersion, _ := semver.Make("1.18.0")

	connectivityCheckCommand := ""
	if clusterVersion.GTE(minVersion) {
		connectivityCheckCommand = `curl -v --insecure --proxy-insecure https://` + registry + `/v2/`
	} else {
		connectivityCheckCommand = `nc -vz ` + registry + ` 443`
	}

	if registry == "" {
		return ""
	}

	// only use https proxy, if user doesn't specify httpsProxy we autofill it with value from httpProxy
	proxyVars := ""
	if nbc.HTTPProxyConfig != nil {
		if nbc.HTTPProxyConfig.HTTPProxy != nil {
			// from https://curl.se/docs/manual.html, curl uses http_proxy but uppercase for others?
			proxyVars = fmt.Sprintf("export http_proxy=\"%s\";", *nbc.HTTPProxyConfig.HTTPProxy)
		}
		if nbc.HTTPProxyConfig.HTTPSProxy != nil {
			proxyVars = fmt.Sprintf("export HTTPS_PROXY=\"%s\"; %s", *nbc.HTTPProxyConfig.HTTPSProxy, proxyVars)
		}
		if nbc.HTTPProxyConfig.NoProxy != nil {
			proxyVars = fmt.Sprintf("export NO_PROXY=\"%s\"; %s", strings.Join(*nbc.HTTPProxyConfig.NoProxy, ","), proxyVars)
		}
	}

	cmd := `retrycmd_if_failure() { r=$1; w=$2; t=$3; shift && shift && shift; for i in $(seq 1 $r); do timeout $t ${@}; [ $? -eq 0  ] && break || if [ $i -eq $r ]; then return 1; else sleep $w; fi; done }; ERR_OUTBOUND_CONN_FAIL=50; retrycmd_if_failure 100 1 10 ` + connectivityCheckCommand + ` >> /var/log/azure/cluster-provision-cse-output.log 2>&1 || time ` + connectivityCheckCommand + ` || exit $ERR_OUTBOUND_CONN_FAIL;`

	if proxyVars != "" {
		cmd = fmt.Sprintf("%s %s", proxyVars, cmd)
	}

	return cmd
}
