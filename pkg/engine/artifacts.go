// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/aks-engine/pkg/api"
	"github.com/Azure/aks-engine/pkg/api/common"
)

// kubernetesComponentFileSpec defines a k8s component that we will deliver via file to a master node vm
type kubernetesComponentFileSpec struct {
	sourceFile      string // filename to source spec data from
	base64Data      string // if not "", this base64-encoded string will take precedent over sourceFile
	destinationFile string // the filename to write to disk on the destination OS
	isEnabled       bool   // is this spec enabled?
}

func kubernetesContainerAddonSettingsInit(p *api.Properties) map[string]kubernetesComponentFileSpec {
	if p.OrchestratorProfile == nil {
		p.OrchestratorProfile = &api.OrchestratorProfile{}
	}
	if p.OrchestratorProfile.KubernetesConfig == nil {
		p.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
	}
	o := p.OrchestratorProfile
	k := o.KubernetesConfig
	return map[string]kubernetesComponentFileSpec{
		DefaultHeapsterAddonName: {
			sourceFile:      "kubernetesmasteraddons-heapster-deployment.yaml",
			base64Data:      k.GetAddonScript(HeapsterAddonName),
			destinationFile: "kube-heapster-deployment.yaml",
			isEnabled:       !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.13.0"),
		},
		MetricsServerAddonName: {
			sourceFile:      "kubernetesmasteraddons-metrics-server-deployment.yaml",
			base64Data:      k.GetAddonScript(MetricsServerAddonName),
			destinationFile: "kube-metrics-server-deployment.yaml",
			isEnabled:       o.IsMetricsServerEnabled(),
		},
		TillerAddonName: {
			sourceFile:      "kubernetesmasteraddons-tiller-deployment.yaml",
			base64Data:      k.GetAddonScript(TillerAddonName),
			destinationFile: "kube-tiller-deployment.yaml",
			isEnabled:       k.IsTillerEnabled(),
		},
		AADPodIdentityAddonName: {
			sourceFile:      "kubernetesmasteraddons-aad-pod-identity-deployment.yaml",
			base64Data:      k.GetAddonScript(AADPodIdentityAddonName),
			destinationFile: "aad-pod-identity-deployment.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsAADPodIdentityEnabled(),
		},
		ACIConnectorAddonName: {
			sourceFile:      "kubernetesmasteraddons-aci-connector-deployment.yaml",
			base64Data:      k.GetAddonScript(ACIConnectorAddonName),
			destinationFile: "aci-connector-deployment.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsACIConnectorEnabled(),
		},
		ClusterAutoscalerAddonName: {
			sourceFile:      "kubernetesmasteraddons-cluster-autoscaler-deployment.yaml",
			base64Data:      k.GetAddonScript(ClusterAutoscalerAddonName),
			destinationFile: "cluster-autoscaler-deployment.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsClusterAutoscalerEnabled(),
		},
		BlobfuseFlexVolumeAddonName: {
			sourceFile:      "kubernetesmasteraddons-blobfuse-flexvolume-installer.yaml",
			base64Data:      k.GetAddonScript(BlobfuseFlexVolumeAddonName),
			destinationFile: "blobfuse-flexvolume-installer.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsBlobfuseFlexVolumeEnabled() && !p.HasCoreOS(),
		},

		SMBFlexVolumeAddonName: {
			sourceFile:      "kubernetesmasteraddons-smb-flexvolume-installer.yaml",
			base64Data:      k.GetAddonScript(SMBFlexVolumeAddonName),
			destinationFile: "smb-flexvolume-installer.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsSMBFlexVolumeEnabled() && !p.HasCoreOS(),
		},
		KeyVaultFlexVolumeAddonName: {
			sourceFile:      "kubernetesmasteraddons-keyvault-flexvolume-installer.yaml",
			base64Data:      k.GetAddonScript(KeyVaultFlexVolumeAddonName),
			destinationFile: "keyvault-flexvolume-installer.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsKeyVaultFlexVolumeEnabled() && !p.HasCoreOS(),
		},
		DashboardAddonName: {
			sourceFile:      "kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
			base64Data:      k.GetAddonScript(DashboardAddonName),
			destinationFile: "kubernetes-dashboard-deployment.yaml",
			isEnabled:       k.IsDashboardEnabled(),
		},
		ReschedulerAddonName: {
			sourceFile:      "kubernetesmasteraddons-kube-rescheduler-deployment.yaml",
			base64Data:      k.GetAddonScript(ReschedulerAddonName),
			destinationFile: "kube-rescheduler-deployment.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsReschedulerEnabled(),
		},
		NVIDIADevicePluginAddonName: {
			sourceFile:      "kubernetesmasteraddons-nvidia-device-plugin-daemonset.yaml",
			base64Data:      k.GetAddonScript(NVIDIADevicePluginAddonName),
			destinationFile: "nvidia-device-plugin.yaml",
			isEnabled:       !p.IsAzureStackCloud() && p.IsNvidiaDevicePluginCapable() && p.IsNVIDIADevicePluginEnabled(),
		},
		ContainerMonitoringAddonName: {
			sourceFile:      "kubernetesmasteraddons-omsagent-daemonset.yaml",
			base64Data:      k.GetAddonScript(ContainerMonitoringAddonName),
			destinationFile: "omsagent-daemonset.yaml",
			isEnabled:       !p.IsAzureStackCloud() && k.IsContainerMonitoringEnabled(),
		},
		IPMASQAgentAddonName: {
			sourceFile:      "ip-masq-agent.yaml",
			base64Data:      k.GetAddonScript(IPMASQAgentAddonName),
			destinationFile: "ip-masq-agent.yaml",
			isEnabled:       k.IsIPMasqAgentEnabled(),
		},
		AzureCNINetworkMonitorAddonName: {
			sourceFile:      "azure-cni-networkmonitor.yaml",
			base64Data:      k.GetAddonScript(AzureCNINetworkMonitorAddonName),
			destinationFile: "azure-cni-networkmonitor.yaml",
			isEnabled:       o.IsAzureCNI() && k.IsAzureCNIMonitoringEnabled(),
		},
		DNSAutoscalerAddonName: {
			sourceFile:      "dns-autoscaler.yaml",
			base64Data:      k.GetAddonScript(DNSAutoscalerAddonName),
			destinationFile: "dns-autoscaler.yaml",
			// TODO enable this when it has been smoke tested
			//common.IsKubernetesVersionGe(p.OrchestratorProfile.OrchestratorVersion, "1.12.0"),
			isEnabled: false,
		},
		CalicoAddonName: {
			sourceFile:      "kubernetesmasteraddons-calico-daemonset.yaml",
			base64Data:      k.GetAddonScript(CalicoAddonName),
			destinationFile: "calico-daemonset.yaml",
			isEnabled:       k.NetworkPolicy == NetworkPolicyCalico,
		},
	}
}

func kubernetesAddonSettingsInit(p *api.Properties) []kubernetesComponentFileSpec {
	if p.OrchestratorProfile == nil {
		p.OrchestratorProfile = &api.OrchestratorProfile{}
	}
	if p.OrchestratorProfile.KubernetesConfig == nil {
		p.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
	}
	o := p.OrchestratorProfile
	k := o.KubernetesConfig
	kubernetesComponentFileSpecs := []kubernetesComponentFileSpec{
		{
			sourceFile:      "kubernetesmasteraddons-kube-dns-deployment.yaml",
			base64Data:      k.GetAddonScript(KubeDNSAddonName),
			destinationFile: "kube-dns-deployment.yaml",
			isEnabled:       !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.12.0"),
		},
		{
			sourceFile:      "coredns.yaml",
			base64Data:      k.GetAddonScript(CoreDNSAddonName),
			destinationFile: "coredns.yaml",
			isEnabled:       common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.12.0"),
		},
		{
			sourceFile:      "kubernetesmasteraddons-kube-proxy-daemonset.yaml",
			base64Data:      k.GetAddonScript(KubeProxyAddonName),
			destinationFile: "kube-proxy-daemonset.yaml",
			isEnabled:       true,
		},
		{
			sourceFile:      "kubernetesmasteraddons-azure-npm-daemonset.yaml",
			base64Data:      k.GetAddonScript(AzureNetworkPolicyAddonName),
			destinationFile: "azure-npm-daemonset.yaml",
			isEnabled:       k.NetworkPolicy == NetworkPolicyAzure && p.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginAzure,
		},
		{
			sourceFile:      "kubernetesmasteraddons-cilium-daemonset.yaml",
			base64Data:      k.GetAddonScript(CiliumAddonName),
			destinationFile: "cilium-daemonset.yaml",
			isEnabled:       k.NetworkPolicy == NetworkPolicyCilium,
		},
		{
			sourceFile:      "kubernetesmasteraddons-flannel-daemonset.yaml",
			base64Data:      k.GetAddonScript(FlannelAddonName),
			destinationFile: "flannel-daemonset.yaml",
			isEnabled:       k.NetworkPlugin == NetworkPluginFlannel,
		},
		{
			sourceFile:      "kubernetesmasteraddons-aad-default-admin-group-rbac.yaml",
			base64Data:      k.GetAddonScript(AADAdminGroupAddonName),
			destinationFile: "aad-default-admin-group-rbac.yaml",
			isEnabled:       p.AADProfile != nil && p.AADProfile.AdminGroupID != "",
		},
		{
			sourceFile:      "kubernetesmasteraddons-azure-cloud-provider-deployment.yaml",
			base64Data:      k.GetAddonScript(AzureCloudProviderAddonName),
			destinationFile: "azure-cloud-provider-deployment.yaml",
			isEnabled:       true,
		},
		{
			sourceFile:      "kubernetesmaster-audit-policy.yaml",
			base64Data:      k.GetAddonScript(AuditPolicyAddonName),
			destinationFile: "audit-policy.yaml",
			isEnabled:       common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0"),
		},
		{
			sourceFile:      "kubernetesmasteraddons-elb-svc.yaml",
			base64Data:      k.GetAddonScript(ELBServiceAddonName),
			destinationFile: "elb-svc.yaml",
			isEnabled:       k.LoadBalancerSku == api.StandardLoadBalancerSku && !p.OrchestratorProfile.IsPrivateCluster(),
		},
		{
			sourceFile:      "kubernetesmasteraddons-pod-security-policy.yaml",
			base64Data:      p.OrchestratorProfile.KubernetesConfig.PodSecurityPolicyConfig["data"],
			destinationFile: "pod-security-policy.yaml",
			isEnabled:       to.Bool(p.OrchestratorProfile.KubernetesConfig.EnablePodSecurityPolicy),
		},
	}

	unmanagedStorageClassesSourceYaml := "kubernetesmasteraddons-unmanaged-azure-storage-classes.yaml"
	managedStorageClassesSourceYaml := "kubernetesmasteraddons-managed-azure-storage-classes.yaml"

	if p.IsAzureStackCloud() {
		unmanagedStorageClassesSourceYaml = "kubernetesmasteraddons-unmanaged-azure-storage-classes-custom.yaml"
		managedStorageClassesSourceYaml = "kubernetesmasteraddons-managed-azure-storage-classes-custom.yaml"
	}

	if len(p.AgentPoolProfiles) > 0 {
		kubernetesComponentFileSpecs = append(kubernetesComponentFileSpecs,
			kubernetesComponentFileSpec{
				sourceFile:      unmanagedStorageClassesSourceYaml,
				base64Data:      p.OrchestratorProfile.KubernetesConfig.GetAddonScript(AzureStorageClassesAddonName),
				destinationFile: "azure-storage-classes.yaml",
				isEnabled:       p.AgentPoolProfiles[0].StorageProfile == api.StorageAccount,
			})
		kubernetesComponentFileSpecs = append(kubernetesComponentFileSpecs,
			kubernetesComponentFileSpec{
				sourceFile:      managedStorageClassesSourceYaml,
				base64Data:      p.OrchestratorProfile.KubernetesConfig.GetAddonScript(AzureStorageClassesAddonName),
				destinationFile: "azure-storage-classes.yaml",
				isEnabled:       p.AgentPoolProfiles[0].StorageProfile == api.ManagedDisks,
			})
	}

	return kubernetesComponentFileSpecs
}

func kubernetesManifestSettingsInit(p *api.Properties) []kubernetesComponentFileSpec {
	if p.OrchestratorProfile == nil {
		p.OrchestratorProfile = &api.OrchestratorProfile{}
	}
	if p.OrchestratorProfile.KubernetesConfig == nil {
		p.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
	}
	o := p.OrchestratorProfile
	k := o.KubernetesConfig
	kubeControllerManagerYaml := "kubernetesmaster-kube-controller-manager.yaml"

	if p.IsAzureStackCloud() {
		kubeControllerManagerYaml = "kubernetesmaster-kube-controller-manager-custom.yaml"
	}

	return []kubernetesComponentFileSpec{
		{
			sourceFile:      "kubernetesmaster-kube-scheduler.yaml",
			base64Data:      k.SchedulerConfig["data"],
			destinationFile: "kube-scheduler.yaml",
			isEnabled:       true,
		},
		{
			sourceFile:      kubeControllerManagerYaml,
			base64Data:      k.ControllerManagerConfig["data"],
			destinationFile: "kube-controller-manager.yaml",
			isEnabled:       true,
		},
		{
			sourceFile:      "kubernetesmaster-cloud-controller-manager.yaml",
			base64Data:      k.CloudControllerManagerConfig["data"],
			destinationFile: "cloud-controller-manager.yaml",
			isEnabled:       k.UseCloudControllerManager != nil && *p.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager,
		},
		{
			sourceFile:      "kubernetesmaster-kube-apiserver.yaml",
			base64Data:      k.APIServerConfig["data"],
			destinationFile: "kube-apiserver.yaml",
			isEnabled:       true,
		},
		{
			sourceFile:      "kubernetesmaster-kube-addon-manager.yaml",
			base64Data:      "", // arbitrary user-provided data not enabled for kube-addon-manager spec
			destinationFile: "kube-addon-manager.yaml",
			isEnabled:       true,
		},
	}
}

func getAddonString(input, destinationPath, destinationFile string) string {
	addonString := getBase64EncodedGzippedCustomScriptFromStr(input)
	return buildConfigString(addonString, destinationFile, destinationPath)
}

func substituteConfigString(input string, kubernetesFeatureSettings []kubernetesComponentFileSpec, sourcePath string, destinationPath string, placeholder string, orchestratorVersion string) string {
	var config string

	versions := strings.Split(orchestratorVersion, ".")
	for _, setting := range kubernetesFeatureSettings {
		if setting.isEnabled {
			var cscript string
			if setting.base64Data != "" {
				var err error
				cscript, err = getStringFromBase64(setting.base64Data)
				if err != nil {
					return ""
				}
				config += getAddonString(cscript, setting.destinationFile, destinationPath)
			} else {
				cscript = getCustomScriptFromFile(setting.sourceFile,
					sourcePath,
					versions[0]+"."+versions[1])
				config += buildConfigString(
					cscript,
					setting.destinationFile,
					destinationPath)
			}
		}
	}

	return strings.Replace(input, placeholder, config, -1)
}

func buildConfigString(configString, destinationFile, destinationPath string) string {
	contents := []string{
		fmt.Sprintf("- path: %s/%s", destinationPath, destinationFile),
		"  permissions: \\\"0644\\\"",
		"  encoding: gzip",
		"  owner: \\\"root\\\"",
		"  content: !!binary |",
		fmt.Sprintf("    %s\\n\\n", configString),
	}

	return strings.Join(contents, "\\n")
}

func getCustomScriptFromFile(sourceFile, sourcePath, version string) string {
	customDataFilePath := getCustomDataFilePath(sourceFile, sourcePath, version)
	return getBase64EncodedGzippedCustomScript(customDataFilePath)
}

func getCustomDataFilePath(sourceFile, sourcePath, version string) string {
	sourceFileFullPath := sourcePath + "/" + sourceFile
	sourceFileFullPathVersioned := sourcePath + "/" + version + "/" + sourceFile

	// Test to check if the versioned file can be read.
	_, err := Asset(sourceFileFullPathVersioned)
	if err == nil {
		sourceFileFullPath = sourceFileFullPathVersioned
	}
	return sourceFileFullPath
}
