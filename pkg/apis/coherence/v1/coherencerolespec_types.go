/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoherenceRoleSpec defines a role in a Coherence cluster. A role is one or
// more Pods that perform the same functionality, for example storage members.
// +k8s:openapi-gen=true
type CoherenceRoleSpec struct {
	// The name of this role.
	// This value will be used to set the Coherence role property for all members of this role
	// +optional
	Role string `json:"role,omitempty"`
	// The desired number of cluster members of this role.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Default value is 3.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// The optional application definition
	// +optional
	Application *ApplicationSpec `json:"application,omitempty"`
	// The optional application definition
	// +optional
	Coherence *CoherenceSpec `json:"coherence,omitempty"`
	// The configuration for the Coherence utils image
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// Logging allows configuration of Coherence and java util logging.
	// +optional
	Logging *LoggingSpec `json:"logging,omitempty"`
	// The JVM specific options
	// +optional
	JVM *JVMSpec `json:"jvm,omitempty"`
	// Ports specifies additional port mappings for the Pod and additional Services for those ports
	// +listType=map
	// +listMapKey=name
	// +optional
	Ports []NamedPortSpec `json:"ports,omitempty"`
	// Env is additional environment variable mappings that will be passed to
	// the Coherence container in the Pod
	// To specify extra variables add them as name value pairs the same as they
	// would be added to a Pod containers spec, for example these values:
	//
	// env:
	//   - name "FOO"
	//     value: "foo-value"
	//   - name: "BAR"
	//     value "bar-value"
	//
	// will add the environment variable mappings FOO="foo-value" and BAR="bar-value"
	// +listType=map
	// +listMapKey=name
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// The port that the health check endpoint will bind to.
	// +optional
	HealthPort *int32 `json:"healthPort,omitempty"`
	// The readiness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
	// The liveness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	LivenessProbe *ReadinessProbeSpec `json:"livenessProbe,omitempty"`
	// The configuration to control safe scaling.
	// +optional
	Scaling *ScalingSpec `json:"scaling,omitempty"`
	// Resources is the optional resource requests and limits for the containers
	//  ref: http://kubernetes.io/docs/user-guide/compute-resources/
	//
	// By default the cpu requests is set to zero and the cpu limit set to 32. This
	// is because it appears that K8s defaults cpu to one and since Java 10 the JVM
	// now correctly picks up cgroup cpu limits then the JVM will only see one cpu.
	// By setting resources.requests.cpu=0 and resources.limits.cpu=32 it ensures that
	// the JVM will see the either the number of cpus on the host if this is <= 32 or
	// the JVM will see 32 cpus if the host has > 32 cpus. The limit is set to zero
	// so that there is no hard-limit applied.
	//
	// No default memory limits are applied.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Annotations are free-form yaml that will be added to the store release as annotations
	// Any annotations should be placed BELOW this annotations: key. For example if we wanted to
	// include annotations for Prometheus it would look like this:
	//
	// annotations:
	//   prometheus.io/scrape: "true"
	//   prometheus.io/port: "2408"
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// The extra labels to add to the all of the Pods in this roles.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Volumes defines extra volume mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumes section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/volumes/
	// +listType=map
	// +listMapKey=name
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumeClaimTemplates defines extra PVC mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumeClaimTemplates section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/persistent-volumes/
	// +listType=map
	// +listMapKey=metaData.name
	// +optional
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// VolumeMounts defines extra volume mounts to map to the additional volumes or PVCs declared above
	//   in store.volumes and store.volumeClaimTemplates
	// +listType=map
	// +listMapKey=name
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// Affinity controls Pod scheduling preferences.
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// NodeSelector is the Node labels for pod assignment
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Tolerations is for nodes that have taints on them.
	//   Useful if you want to dedicate nodes to just run the coherence container
	// For example:
	//   tolerations:
	//   - key: "key"
	//     operator: "Equal"
	//     value: "value"
	//     effect: "NoSchedule"
	//
	//   ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
	// +listType=map
	// +listMapKey=key
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext is the PodSecurityContext that will be added to all of the Pods in this role.
	// See: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	// Share a single process namespace between all of the containers in a pod. When this is set containers will
	// be able to view and signal processes from other containers in the same pod, and the first process in each
	// container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set.
	// Optional: Default to false.
	// +optional
	ShareProcessNamespace *bool `json:"shareProcessNamespace,omitempty"`
	// Use the host's ipc namespace. Optional: Default to false.
	// +optional
	HostIPC *bool `json:"hostIPC,omitempty"`
	// Configure various networks and DNS settings for Pods in this rolw.
	// +optional
	Network *NetworkSpec `json:"network,omitempty"`
	// The roles that must be started before this role can start.
	// +listType=map
	// +listMapKey=role
	// +optional
	StartQuorum []StartQuorum `json:"startQuorum,omitempty"`
}

// Obtain the number of replicas required for a role.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replica value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceRoleSpec) GetReplicas() int32 {
	if in == nil {
		return 0
	}
	if in.Replicas == nil {
		return DefaultReplicas
	}
	return *in.Replicas
}

// Set the number of replicas required for a role.
func (in *CoherenceRoleSpec) SetReplicas(replicas int32) {
	if in != nil {
		in.Replicas = &replicas
	}
}

// Obtain the full name for  a role.
func (in *CoherenceRoleSpec) GetFullRoleName(cluster *CoherenceCluster) string {
	if in == nil {
		return ""
	}

	return cluster.GetFullRoleName(in.GetRoleName())
}

// Obtain the name for a role.
// If the Role field is not set the default name is returned.
func (in *CoherenceRoleSpec) GetRoleName() string {
	if in == nil {
		return DefaultRoleName
	}
	if in.Role == "" {
		return DefaultRoleName
	}
	return in.Role
}

func (in *CoherenceRoleSpec) GetCoherenceImage() *string {
	if in != nil && in.Coherence != nil {
		return in.Coherence.Image
	}
	return nil
}

// Ensure that the Coherence image is set for the role.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceRoleSpec) EnsureCoherenceImage(coherenceImage *string) bool {
	if in.Coherence == nil {
		in.Coherence = &CoherenceSpec{}
	}

	return in.Coherence.EnsureImage(coherenceImage)
}

func (in *CoherenceRoleSpec) GetCoherenceUtilsImage() *string {
	if in != nil && in.CoherenceUtils != nil {
		return in.CoherenceUtils.Image
	}
	return nil
}

// Ensure that the Coherence Utils image is set for the role.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceRoleSpec) EnsureCoherenceUtilsImage(utilsImage *string) bool {
	if in.CoherenceUtils == nil {
		in.CoherenceUtils = &ImageSpec{}
	}

	return in.CoherenceUtils.EnsureImage(utilsImage)
}

func (in *CoherenceRoleSpec) GetEffectiveScalingPolicy() ScalingPolicy {
	if in == nil {
		return SafeScaling
	}

	var policy ScalingPolicy

	if in.Scaling == nil || in.Scaling.Policy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if in.Coherence == nil || in.Coherence.StorageEnabled == nil || *in.Coherence.StorageEnabled {
			// storage enabled is either not set or is true so do safe scaling
			policy = ParallelUpSafeDownScaling
		} else {
			// storage enabled is false so do parallel scaling
			policy = ParallelScaling
		}
	} else {
		// scaling policy is set so use it
		policy = *in.Scaling.Policy
	}

	return policy
}

// Returns the port that the health check endpoint will bind to.
func (in *CoherenceRoleSpec) GetHealthPort() int32 {
	if in == nil || in.HealthPort == nil || *in.HealthPort <= 0 {
		return DefaultHealthPort
	}
	return *in.HealthPort
}

// Returns the ScalingProbe to use for checking Status HA for the role.
// This method will not return nil.
func (in *CoherenceRoleSpec) GetScalingProbe() *ScalingProbe {
	if in == nil || in.Scaling == nil || in.Scaling.Probe == nil {
		return in.GetDefaultScalingProbe()
	}
	return in.Scaling.Probe
}

// Obtain a default ScalingProbe
func (in *CoherenceRoleSpec) GetDefaultScalingProbe() *ScalingProbe {
	timeout := 10

	defaultStatusHA := ScalingProbe{
		TimeoutSeconds: &timeout,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/ha",
				Port: intstr.FromString(PortNameHealth),
			},
		},
	}

	return defaultStatusHA.DeepCopy()
}

// DeepCopyWithDefaults returns a copy of this CoherenceRoleSpec with any nil or not set values set
// by the corresponding value in the defaults spec.
func (in *CoherenceRoleSpec) DeepCopyWithDefaults(defaults *CoherenceRoleSpec) *CoherenceRoleSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := CoherenceRoleSpec{}

	// Copy EVERY field from "in" to the clone.
	// If a field is not set use the value from the default
	// If the field is a struct it should implement DeepCopyWithDefaults so call that method

	// Affinity is NOT merged
	if in.Affinity != nil {
		clone.Affinity = in.Affinity
	} else {
		clone.Affinity = defaults.Affinity
	}

	// Annotations are a map and are merged
	clone.Annotations = MergeMap(in.Annotations, defaults.Annotations)
	// Application is merged
	clone.Application = in.Application.DeepCopyWithDefaults(defaults.Application)
	clone.Coherence = in.Coherence.DeepCopyWithDefaults(defaults.Coherence)
	clone.CoherenceUtils = in.CoherenceUtils.DeepCopyWithDefaults(defaults.CoherenceUtils)
	// Environment variables are merged
	clone.Env = in.mergeEnvVar(in.Env, defaults.Env)
	clone.JVM = in.JVM.DeepCopyWithDefaults(defaults.JVM)
	// Labels are a map and are merged
	clone.Labels = MergeMap(in.Labels, defaults.Labels)
	clone.Logging = in.Logging.DeepCopyWithDefaults(defaults.Logging)
	// Network configuration is merged
	clone.Network = in.Network.DeepCopyWithDefaults(defaults.Network)

	// The quorum is NEVER taken from defaults
	clone.StartQuorum = in.StartQuorum

	// NodeSelector is a map and is NOT merged
	clone.NodeSelector = MergeMap(in.NodeSelector, defaults.NodeSelector)
	if in.NodeSelector != nil {
		clone.NodeSelector = in.NodeSelector
	} else {
		clone.NodeSelector = defaults.NodeSelector
	}

	// Ports are named ports in an array and are merged
	if in.Ports != nil {
		clone.Ports = MergeNamedPortSpecs(in.Ports, defaults.Ports)
	} else {
		clone.Ports = defaults.Ports
	}

	// ReadinessProbe is merged
	clone.ReadinessProbe = in.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)

	// Application is NOT merged
	if in.Replicas != nil {
		clone.Replicas = in.Replicas
	} else {
		clone.Replicas = defaults.Replicas
	}

	// Resources is NOT merged
	if in.Resources != nil {
		clone.Resources = in.Resources
	} else {
		clone.Resources = defaults.Resources
	}

	// Role is NOT merged
	if in.Role != "" {
		clone.Role = in.Role
	} else {
		clone.Role = defaults.Role
	}

	// Tolerations is an array but is NOT merged
	if in.Tolerations != nil {
		clone.Tolerations = make([]corev1.Toleration, len(in.Tolerations))
		for i := 0; i < len(in.Tolerations); i++ {
			clone.Tolerations[i] = *in.Tolerations[i].DeepCopy()
		}
	} else if defaults.Tolerations != nil {
		clone.Tolerations = make([]corev1.Toleration, len(defaults.Tolerations))
		for i := 0; i < len(defaults.Tolerations); i++ {
			clone.Tolerations[i] = *defaults.Tolerations[i].DeepCopy()
		}
	}

	// SecurityContext is NOT merged
	if in.SecurityContext != nil {
		clone.SecurityContext = in.SecurityContext
	} else {
		clone.SecurityContext = defaults.SecurityContext
	}

	if in.ShareProcessNamespace != nil {
		clone.ShareProcessNamespace = in.ShareProcessNamespace
	} else {
		clone.ShareProcessNamespace = defaults.ShareProcessNamespace
	}

	if in.HostIPC != nil {
		clone.HostIPC = in.HostIPC
	} else {
		clone.HostIPC = defaults.HostIPC
	}

	// VolumeClaimTemplates is an array of named PersistentVolumeClaims and is merged
	clone.VolumeClaimTemplates = in.mergePersistentVolumeClaims(in.VolumeClaimTemplates, defaults.VolumeClaimTemplates)
	// VolumeMounts is an array of named VolumeMounts and is merged
	clone.VolumeMounts = in.mergeVolumeMounts(in.VolumeMounts, defaults.VolumeMounts)
	// Volumes is an array of named VolumeMounts and is merged
	clone.Volumes = in.mergeVolumes(in.Volumes, defaults.Volumes)

	return &clone
}

func (in *CoherenceRoleSpec) mergeEnvVar(primary, secondary []corev1.EnvVar) []corev1.EnvVar {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.EnvVar{}
	}

	var merged []corev1.EnvVar
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergePersistentVolumeClaims(primary, secondary []corev1.PersistentVolumeClaim) []corev1.PersistentVolumeClaim {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.PersistentVolumeClaim{}
	}

	var merged []corev1.PersistentVolumeClaim
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergeVolumeMounts(primary, secondary []corev1.VolumeMount) []corev1.VolumeMount {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.VolumeMount{}
	}

	var merged []corev1.VolumeMount
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergeVolumes(primary, secondary []corev1.Volume) []corev1.Volume {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.Volume{}
	}

	var merged []corev1.Volume
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

// Create the Kubernetes resources that should be deployed for this role.
// The order of the resources in the returned array is the order that they should be
// created or updated in Kubernetes.
func (in *CoherenceRoleSpec) CreateKubernetesResources(cluster *CoherenceCluster) ([]metav1.Object, error) {
	var res []metav1.Object

	// Create the fluentd ConfigMap if required
	if in.Logging.IsFluentdEnabled() {
		cm, err := in.Logging.CreateConfigMap(cluster, in)
		if err != nil {
			return res, err
		}
		res = append(res, cm)
	}

	// Create the headless Service
	res = append(res, in.CreateHeadlessService(cluster))

	// Create the StatefulSet
	res = append(res, in.CreateStatefulSet(cluster))

	// Create the Services for each port
	res = append(res, in.CreateServicesForPort(cluster)...)

	return res, nil
}

// Create the role's common label set.
func (in *CoherenceRoleSpec) CreateServicesForPort(cluster *CoherenceCluster) []metav1.Object {
	var services []metav1.Object

	if in == nil || in.Ports == nil || len(in.Ports) == 0 {
		return services
	}

	// Create the Services for each port
	for _, p := range in.Ports {
		service := p.CreateService(cluster, in)
		if service != nil {
			services = append(services, service)
		}
	}

	return services
}

// Create the role's common label set.
func (in *CoherenceRoleSpec) CreateCommonLabels(cluster *CoherenceCluster) map[string]string {
	labels := make(map[string]string)
	labels[LabelCoherenceDeployment] = in.GetFullRoleName(cluster)
	labels[LabelCoherenceCluster] = cluster.GetName()
	labels[LabelCoherenceRole] = in.GetRoleName()
	return labels
}

// Create the selector that can be used to match this roles Pods, for example by Services or StatefulSets.
func (in *CoherenceRoleSpec) CreatePodSelectorLabels(cluster *CoherenceCluster) map[string]string {
	selector := in.CreateCommonLabels(cluster)
	selector[LabelComponent] = LabelComponentCoherencePod
	return selector
}

// Create the headless Service for the role's StatefulSet.
func (in *CoherenceRoleSpec) CreateHeadlessService(cluster *CoherenceCluster) *corev1.Service {
	// The labels for the service
	svcLabels := in.CreateCommonLabels(cluster)
	svcLabels[LabelComponent] = LabelComponentCoherenceHeadless

	// The selector for the service
	selector := in.CreatePodSelectorLabels(cluster)

	// Create the Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cluster.GetHeadlessServiceNameForRole(in),
			Labels: svcLabels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP:                "None",
			PublishNotReadyAddresses: true,
			Selector:                 selector,
			Ports: []corev1.ServicePort{
				{
					Name:       PortNameCoherence,
					Protocol:   corev1.ProtocolTCP,
					Port:       7,
					TargetPort: intstr.FromInt(7),
				},
			},
		},
	}

	return svc
}

// Create the role's StatefulSet.
func (in *CoherenceRoleSpec) CreateStatefulSet(cluster *CoherenceCluster) *appsv1.StatefulSet {
	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   in.GetFullRoleName(cluster),
			Labels: in.CreateCommonLabels(cluster),
		},
	}

	// Create the PodSpec labels
	podLabels := in.CreatePodSelectorLabels(cluster)
	// Add the WKA member label
	podLabels[LabelCoherenceWKAMember] = strconv.FormatBool(in.Coherence.IsWKAMember())
	// Add any labels specified for the role
	for k, v := range in.Labels {
		podLabels[k] = v
	}

	replicas := in.GetReplicas()
	var volumeMode int32 = 0777

	cohContainer := in.CreateCoherenceContainer(cluster)

	// Add additional ports
	for _, p := range in.Ports {
		cohContainer.Ports = append(cohContainer.Ports, p.CreatePort())
	}

	// append any additional VolumeMounts
	cohContainer.VolumeMounts = append(cohContainer.VolumeMounts, in.VolumeMounts...)

	// Add the component label
	sts.Labels[LabelComponent] = LabelComponentCoherenceStatefulSet
	sts.Spec = appsv1.StatefulSetSpec{
		Replicas:            &replicas,
		PodManagementPolicy: appsv1.ParallelPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
		RevisionHistoryLimit: pointer.Int32Ptr(5),
		ServiceName:          cluster.GetHeadlessServiceNameForRole(in),
		Selector: &metav1.LabelSelector{
			MatchLabels: in.CreatePodSelectorLabels(cluster),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      podLabels,
				Annotations: in.Annotations,
			},
			Spec: corev1.PodSpec{
				ImagePullSecrets:             cluster.GetImagePullSecrets(),
				ServiceAccountName:           cluster.GetServiceAccountName(),
				AutomountServiceAccountToken: cluster.Spec.AutomountServiceAccountToken,
				SecurityContext:              in.SecurityContext,
				ShareProcessNamespace:        in.ShareProcessNamespace,
				HostIPC:                      notNilBool(in.HostIPC),
				Tolerations:                  in.Tolerations,
				Affinity:                     in.EnsurePodAffinity(cluster),
				NodeSelector:                 in.NodeSelector,
				InitContainers: []corev1.Container{
					in.CreateUtilsContainer(cluster),
				},
				Containers: []corev1.Container{cohContainer},
				Volumes: []corev1.Volume{
					{Name: VolumeNameLogs, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
					{Name: VolumeNameUtils, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
					{Name: VolumeNameApplication, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
					{
						Name: VolumeNameScripts,
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: ConfigMapNameScripts},
								DefaultMode:          &volumeMode,
							},
						},
					},
				},
			},
		},
	}

	// Add the application init-container if required
	hasApp, c := in.Application.CreateApplicationContainer()
	if hasApp {
		sts.Spec.Template.Spec.InitContainers = append(sts.Spec.Template.Spec.InitContainers, c)
	}

	// Add any network settings
	in.Network.UpdateStatefulSet(&sts)
	// Add any JVM settings
	in.JVM.UpdateStatefulSet(&sts)
	// Add any Coherence settings
	in.Coherence.UpdateStatefulSet(cluster, in, &sts)
	// Add any logging settings
	in.Logging.UpdateStatefulSet(&sts, hasApp)

	// append any additional Volumes
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, in.Volumes...)
	// append any additional PVCs
	sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, in.VolumeClaimTemplates...)

	return &sts
}

// Create the Coherence container spec.
func (in *CoherenceRoleSpec) CreateCoherenceContainer(cluster *CoherenceCluster) corev1.Container {
	c := corev1.Container{
		Name:    ContainerNameCoherence,
		Image:   *in.Coherence.Image,
		Command: []string{"/bin/sh", "-x", "/scripts/startCoherence.sh", "server"},
		Env:     in.Env,
		Ports: []corev1.ContainerPort{
			{
				Name:          PortNameCoherence,
				ContainerPort: 7,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          PortNameHealth,
				ContainerPort: notNilInt32OrDefault(in.HealthPort, cluster.GetHealthPort()),
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: VolumeNameLogs, MountPath: VolumeMountPathLogs},
			{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
			{Name: VolumeNameApplication, MountPath: ExternalAppDir},
			{Name: VolumeNameJVM, MountPath: VolumeMountPathJVM},
			{Name: VolumeNameScripts, MountPath: VolumeMountPathScripts},
		},
	}

	if in.Coherence.ImagePullPolicy != nil {
		c.ImagePullPolicy = *in.Coherence.ImagePullPolicy
	}

	healthPort := cluster.GetHealthPort()

	c.Env = append(c.Env, in.CreateDefaultEnv(cluster)...)

	in.Application.UpdateCoherenceContainer(&c)

	if in.Resources != nil {
		// set the container resources if specified
		c.Resources = *in.Resources
	} else {
		// No resources specified so default to 32 cores
		c.Resources = in.CreateDefaultResources()
	}

	c.ReadinessProbe = in.CreateDefaultReadinessProbe()
	in.ReadinessProbe.UpdateProbeSpec(healthPort, DefaultReadinessPath, c.ReadinessProbe)

	c.LivenessProbe = in.CreateDefaultLivenessProbe()
	in.LivenessProbe.UpdateProbeSpec(healthPort, DefaultLivenessPath, c.LivenessProbe)

	return c
}

// Create the default environment variables.
func (in *CoherenceRoleSpec) CreateDefaultEnv(cluster *CoherenceCluster) []corev1.EnvVar {
	healthPort := cluster.GetHealthPort()

	return []corev1.EnvVar{
		{Name: "COH_WKA", Value: cluster.GetWkaServiceName()},
		{Name: "COH_APP_DIR", Value: ExternalAppDir},
		{Name: "COH_EXTRA_CLASSPATH", Value: fmt.Sprintf("%s/*:%s", ExternalLibDir, ExternalConfDir)},
		{
			Name: "COH_MACHINE_NAME", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
		{
			Name: "COH_MEMBER_NAME", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "COH_POD_UID", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.uid",
				},
			},
		},
		{
			Name: "OPERATOR_HOST", ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: ConfigMapNameOperatorConfig},
					Key:                  OperatorConfigKeyHost,
					Optional:             pointer.BoolPtr(true),
				},
			},
		},
		{Name: "COH_SITE_INFO_LOCATION", Value: "http://$(OPERATOR_HOST)/site/$(COH_MACHINE_NAME)"},
		{Name: "COH_RACK_INFO_LOCATION", Value: "http://$(OPERATOR_HOST)/rack/$(COH_MACHINE_NAME)"},
		{Name: "COH_CLUSTER_NAME", Value: cluster.Name},
		{Name: "COH_ROLE", Value: in.GetRoleName()},
		{Name: "COH_UTIL_DIR", Value: VolumeMountPathUtils},
		{Name: "OPERATOR_REQUEST_TIMEOUT", Value: Int32PtrToStringWithDefault(cluster.Spec.OperatorRequestTimeout, 120)},
		{Name: "COH_HEALTH_PORT", Value: Int32ToString(healthPort)},
	}
}

// Create the default Container resources.
func (in *CoherenceRoleSpec) CreateDefaultResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("32"),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("0"),
		},
	}
}

// Create the default readiness probe.
func (in *CoherenceRoleSpec) CreateDefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 30,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    50,
	}
}

// Update the probe with the default readiness probe action.
func (in *CoherenceRoleSpec) UpdateDefaultReadinessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultReadinessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Create the default liveness probe.
func (in *CoherenceRoleSpec) CreateDefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 60,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    5,
	}
}

// Update the probe with the default liveness probe action.
func (in *CoherenceRoleSpec) UpdateDefaultLivenessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultLivenessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Get the Utils init-container spec.
func (in *CoherenceRoleSpec) CreateUtilsContainer(cluster *CoherenceCluster) corev1.Container {
	c := corev1.Container{
		Name:    ContainerNameUtils,
		Image:   *in.CoherenceUtils.Image,
		Command: []string{"/files/utils-init"},
		Env: []corev1.EnvVar{
			{Name: "COH_UTIL_DIR", Value: VolumeMountPathUtils},
			{Name: "COH_CLUSTER_NAME", Value: cluster.Name},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
		},
	}

	// set the image pull policy if set for the role
	if in.CoherenceUtils != nil && in.CoherenceUtils.ImagePullPolicy != nil {
		c.ImagePullPolicy = *in.CoherenceUtils.ImagePullPolicy
	}

	// set the persistence volumes if required
	if in.Coherence != nil {
		in.Coherence.AddPersistenceVolumeMounts(&c)
	}
	return c
}

// Get the Pod Affinity either from that configured for the cluster or the default affinity.
func (in *CoherenceRoleSpec) EnsurePodAffinity(cluster *CoherenceCluster) *corev1.Affinity {
	if in != nil && in.Affinity != nil {
		return in.Affinity
	}
	// return the default affinity which attempts to spread the Pods for a role across fault domains
	return in.CreateDefaultPodAffinity(cluster)
}

// Create the default Pod Affinity to use in a role's StatefulSet.
func (in *CoherenceRoleSpec) CreateDefaultPodAffinity(cluster *CoherenceCluster) *corev1.Affinity {
	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 1,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey: AffinityTopologyKey,
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      LabelCoherenceCluster,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{cluster.GetName()},
								},
								{
									Key:      LabelCoherenceRole,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{in.GetRoleName()},
								},
							},
						},
					},
				},
			},
		},
	}
}
