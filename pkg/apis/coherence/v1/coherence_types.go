/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"bytes"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Common Coherence API structs

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// ----- constants ----------------------------------------------------------

const (
	// The default number of replicas that will be created for a role if no value is specified in the spec
	DefaultReplicas int32 = 3

	// The default health check port.
	DefaultHealthPort int32 = 6676

	// The defaultrole name that will be used for a role if no value is specified in the spec
	DefaultRoleName = "storage"

	// The suffix appended to a cluster name to give the WKA service name
	WKAServiceNameSuffix = "-wka"

	// The key of the label used to hold the Coherence cluster name
	CoherenceClusterLabel string = "coherenceCluster"

	// The key of the label used to hold the Coherence role name
	CoherenceRoleLabel string = "coherenceRole"

	// The key of the label used to hold the component name
	CoherenceComponentLabel string = "component"

	// The key of the label used to hold the Coherence Operator version name
	CoherenceOperatorVersionLabel string = "coherenceOperatorVersion"
)

// ----- helper functions ---------------------------------------------------

// Return a map that is two maps merged.
// If both maps are nil then nil is returned.
// Where there are duplicate keys those in m1 take precedence.
// Keys that map to "" will not be added to the merged result
func MergeMap(m1, m2 map[string]string) map[string]string {
	if m1 == nil && m2 == nil {
		return nil
	}

	merged := make(map[string]string)

	for k, v := range m2 {
		if v != "" {
			merged[k] = v
		}
	}

	for k, v := range m1 {
		if v != "" {
			merged[k] = v
		} else {
			delete(merged, k)
		}
	}

	return merged
}

func notNilBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func notNilString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func notNilStringOrDefault(s *string, dflt string) string {
	if s == nil {
		return dflt
	}
	return *s
}

func notNilInt32(i *int32) int32 {
	return notNilInt32OrDefault(i, 0)
}

func notNilInt32OrDefault(i *int32, dflt int32) int32 {
	if i == nil {
		return dflt
	}
	return *i
}

// Ensure that the StatefulSet has a container with the specified name
func EnsureContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	c := FindContainer(name, sts)
	if c == nil {
		c = &corev1.Container{Name: name}
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *c)
	}
	return c
}

// Ensure that the StatefulSet has a container with the specified name
func ReplaceContainer(sts *appsv1.StatefulSet, cNew *corev1.Container) {
	for i, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == cNew.Name {
			sts.Spec.Template.Spec.Containers[i] = *cNew
			return
		}
	}
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *cNew)
}

// Find the StatefulSet container with the specified name.
func FindContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	for _, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Find the StatefulSet init-container with the specified name.
func FindInitContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Ensure that the StatefulSet has a volume with the specified name
func ReplaceVolume(sts *appsv1.StatefulSet, volNew corev1.Volume) {
	for i, v := range sts.Spec.Template.Spec.Volumes {
		if v.Name == volNew.Name {
			sts.Spec.Template.Spec.Volumes[i] = volNew
			return
		}
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volNew)
}

// ----- ApplicationSpec struct ---------------------------------------------

// The specification of the application deployed into the Coherence
// role members.
// +k8s:openapi-gen=true
type ApplicationSpec struct {
	// The application type to execute.
	// This field would be set if using the Coherence Graal image and running a none-Java
	// application. For example if the application was a Node application this field
	// would be set to "node". The default is to run a plain Java application.
	// +optional
	Type *string `json:"type,omitempty"`
	// Class is the Coherence container main class.  The default value is
	// com.tangosol.net.DefaultCacheServer.
	// If the application type is non-Java this would be the name of the corresponding language specific
	// runnable, for example if the application type is "node" the main may be a Javascript file.
	// +optional
	Main *string `json:"main,omitempty"`
	// Args is the optional arguments to pass to the main class.
	// +listType=atomic
	// +optional
	Args []string `json:"args,omitempty"`
	// The inlined application image definition
	ImageSpec `json:",inline"`
	// The application folder in the custom artifacts Docker image containing
	// application artifacts.
	// This will effectively become the working directory of the Coherence container.
	// If not set the application directory default value is "/app".
	// +optional
	AppDir *string `json:"appDir,omitempty"`
	// The folder in the custom artifacts Docker image containing jar
	// files to be added to the classpath of the Coherence container.
	// If not set the lib directory default value is "/app/lib".
	// +optional
	LibDir *string `json:"libDir,omitempty"`
	// The folder in the custom artifacts Docker image containing
	// configuration files to be added to the classpath of the Coherence container.
	// If not set the config directory default value is "/app/conf".
	// +optional
	ConfigDir *string `json:"configDir,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this ApplicationSpec struct with any nil or not set
// values set by the corresponding value in the defaults Images struct.
func (in *ApplicationSpec) DeepCopyWithDefaults(defaults *ApplicationSpec) *ApplicationSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ApplicationSpec{}
	clone.ImageSpec = *in.ImageSpec.DeepCopyWithDefaults(&defaults.ImageSpec)

	if in.Type != nil {
		clone.Type = in.Type
	} else {
		clone.Type = defaults.Type
	}

	if in.Main != nil {
		clone.Main = in.Main
	} else {
		clone.Main = defaults.Main
	}

	if in.Args != nil {
		clone.Args = in.Args
	} else {
		clone.Args = defaults.Args
	}

	if in.AppDir != nil {
		clone.AppDir = in.AppDir
	} else {
		clone.AppDir = defaults.AppDir
	}

	if in.LibDir != nil {
		clone.LibDir = in.LibDir
	} else {
		clone.LibDir = defaults.LibDir
	}

	if in.ConfigDir != nil {
		clone.ConfigDir = in.ConfigDir
	} else {
		clone.ConfigDir = defaults.ConfigDir
	}

	return &clone
}

// Create the application init-container if enabled
func (in *ApplicationSpec) CreateApplicationContainer() (bool, corev1.Container) {
	if in == nil || in.Image == nil {
		return false, corev1.Container{}
	}

	appDir := notNilStringOrDefault(in.AppDir, AppDir)
	libDir := notNilStringOrDefault(in.LibDir, LibDir)
	confDir := notNilStringOrDefault(in.ConfigDir, ConfDir)

	c := corev1.Container{
		Name:    ContainerNameApplication,
		Image:   *in.Image,
		Command: []string{DefaultCommandApplication},
		Env: []corev1.EnvVar{
			{Name: "EXTERNAL_APP_DIR", Value: ExternalAppDir},
			{Name: "APP_DIR", Value: appDir},
			{Name: "EXTERNAL_LIB_DIR", Value: ExternalLibDir},
			{Name: "LIB_DIR", Value: libDir},
			{Name: "EXTERNAL_CONF_DIR", Value: ExternalConfDir},
			{Name: "CONF_DIR", Value: confDir},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
			{Name: VolumeNameApplication, MountPath: ExternalAppDir},
		},
	}

	if in.ImagePullPolicy != nil {
		c.ImagePullPolicy = *in.ImagePullPolicy
	}

	return true, c
}

func (in *ApplicationSpec) UpdateCoherenceContainer(c *corev1.Container) {
	if in == nil {
		return
	}

	if in.Type != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: "APP_TYPE", Value: *in.Type})
	}
	if in.Main != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_MAIN_CLASS", Value: *in.Main})
	}
	if in.Args != nil && len(in.Args) > 0 {
		args := strings.Join(in.Args, " ")
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_MAIN_ARGS", Value: args})
	}
}

// ----- CoherenceSpec struct -----------------------------------------------

// The Coherence specific configuration.
// +k8s:openapi-gen=true
type CoherenceSpec struct {
	// The Coherence images configuration.
	ImageSpec `json:",inline"`
	// A boolean flag indicating whether members of this role are storage enabled.
	// This value will set the corresponding coherence.distributed.localstorage System property.
	// If not specified the default value is true.
	// This flag is also used to configure the ScalingPolicy value if a value is not specified. If the
	// StorageEnabled field is not specified or is true the scaling will be safe, if StorageEnabled is
	// set to false scaling will be parallel.
	// +optional
	StorageEnabled *bool `json:"storageEnabled,omitempty"`
	// CacheConfig is the name of the cache configuration file to use
	// +optional
	CacheConfig *string `json:"cacheConfig,omitempty"`
	// OverrideConfig is name of the Coherence operational configuration override file,
	// the default is tangosol-coherence-override.xml
	// +optional
	OverrideConfig *string `json:"overrideConfig,omitempty"`
	// The Coherence log level, default being 5 (info level).
	// +optional
	LogLevel *int32 `json:"logLevel,omitempty"`
	// Persistence values configure the on-disc data persistence settings.
	// The bool Enabled enables or disabled on disc persistence of data.
	// +optional
	Persistence *PersistentStorageSpec `json:"persistence,omitempty"`
	// Snapshot values configure the on-disc persistence data snapshot (backup) settings.
	// The bool Enabled enables or disabled a different location for
	// persistence snapshot data. If set to false then snapshot files will be written
	// to the same volume configured for persistence data in the Persistence section.
	// +optional
	Snapshot *PersistentStorageSpec `json:"snapshot,omitempty"`
	// Management configures Coherence management over REST
	//   Note: Coherence management over REST will be available in 12.2.1.4.
	// +optional
	Management *PortSpecWithSSL `json:"management,omitempty"`
	// Metrics configures Coherence metrics publishing
	//   Note: Coherence metrics publishing will be available in 12.2.1.4.
	// +optional
	Metrics *PortSpecWithSSL `json:"metrics,omitempty"`
	// Exclude members of this role from being part of the cluster's WKA list.
	ExcludeFromWKA *bool `json:"excludeFromWKA,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this CoherenceSpec struct with any nil or not set
// values set by the corresponding value in the defaults CoherenceSpec struct.
func (in *CoherenceSpec) DeepCopyWithDefaults(defaults *CoherenceSpec) *CoherenceSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := CoherenceSpec{}
	clone.ImageSpec = *in.ImageSpec.DeepCopyWithDefaults(&defaults.ImageSpec)
	clone.Persistence = in.Persistence.DeepCopyWithDefaults(defaults.Persistence)
	clone.Snapshot = in.Snapshot.DeepCopyWithDefaults(defaults.Snapshot)
	clone.Management = in.Management.DeepCopyWithDefaults(defaults.Management)
	clone.Metrics = in.Metrics.DeepCopyWithDefaults(defaults.Metrics)

	if in.StorageEnabled != nil {
		clone.StorageEnabled = in.StorageEnabled
	} else {
		clone.StorageEnabled = defaults.StorageEnabled
	}

	if in.CacheConfig != nil {
		clone.CacheConfig = in.CacheConfig
	} else {
		clone.CacheConfig = defaults.CacheConfig
	}

	if in.OverrideConfig != nil {
		clone.OverrideConfig = in.OverrideConfig
	} else {
		clone.OverrideConfig = defaults.OverrideConfig
	}

	if in.LogLevel != nil {
		clone.LogLevel = in.LogLevel
	} else {
		clone.LogLevel = defaults.LogLevel
	}

	if in.ExcludeFromWKA != nil {
		clone.ExcludeFromWKA = in.ExcludeFromWKA
	} else {
		clone.ExcludeFromWKA = defaults.ExcludeFromWKA
	}

	return &clone
}

// IsWKAMember returns true if this role is a WKA list member.
func (in *CoherenceSpec) IsWKAMember() bool {
	return in != nil && (in.ExcludeFromWKA == nil || !*in.ExcludeFromWKA)
}

// Determine whether persistence is enabled
func (in *CoherenceSpec) IsPersistenceEnabled() bool {
	return in.Persistence != nil && in.Persistence.Enabled != nil && *in.Persistence.Enabled
}

// Determine whether snapshots is enabled
func (in *CoherenceSpec) IsSnapshotsEnabled() bool {
	return in.Snapshot != nil && in.Snapshot.Enabled != nil && *in.Snapshot.Enabled
}

// Add the persistence and snapshot volume mounts to the specified container
func (in *CoherenceSpec) AddPersistenceVolumeMounts(c *corev1.Container) {
	if in == nil {
		// nothing to update
		return
	}

	// Add the persistence volume mount if required
	if in.IsPersistenceEnabled() {
		// add the persistence volume mount
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      VolumeNamePersistence,
			MountPath: VolumeMountPathPersistence,
		})
	}

	// Add the snapshot volume mount if required
	if in.IsSnapshotsEnabled() {
		// Snapshots is enabled so use the snapshot mount
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{Name: VolumeNameSnapshots, MountPath: VolumeMountPathSnapshots})
		//} else {
		//	// no specific snapshot spec set so use the root mount point
		//	rootSnapshots := corev1.VolumeMount{Name: VolumeNameSnapshots, MountPath: VolumeMountPathRootSnapshots}
		//	c.VolumeMounts = append(c.VolumeMounts, rootSnapshots)
	}
}

// Add the persistence and snapshot persistent volume claims
func (in *CoherenceSpec) AddPersistencePVCs(cluster *CoherenceCluster, role *CoherenceRoleSpec, sts *appsv1.StatefulSet) {
	// Add the persistence PVC if required
	if required, pvc := in.Persistence.CreatePersistentVolumeClaim(cluster, role, VolumeNamePersistence); required {
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *pvc)
	}

	// Add the snapshot PVC if required
	if required, pvc := in.Snapshot.CreatePersistentVolumeClaim(cluster, role, VolumeNameSnapshots); required {
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *pvc)
	}
}

// Add the persistence and snapshot volumes
func (in *CoherenceSpec) AddPersistenceVolumes(sts *appsv1.StatefulSet) {
	// Add the persistence volume if required
	if in.IsPersistenceEnabled() && in.Persistence.Volume != nil {
		source := corev1.VolumeSource{}
		in.Persistence.Volume.DeepCopyInto(&source)
		vol := corev1.Volume{
			Name:         VolumeNamePersistence,
			VolumeSource: source,
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
	}

	// Add the snapshot volume if required
	if in.IsSnapshotsEnabled() && in.Snapshot.Volume != nil {
		source := corev1.VolumeSource{}
		in.Snapshot.Volume.DeepCopyInto(&source)
		vol := corev1.Volume{
			Name:         VolumeNameSnapshots,
			VolumeSource: source,
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
	}
}

// Apply Coherence settings to the StatefulSet.
func (in *CoherenceSpec) UpdateStatefulSet(cluster *CoherenceCluster, role *CoherenceRoleSpec, sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	// Get the Coherence container
	c := EnsureContainer(ContainerNameCoherence, sts)
	defer ReplaceContainer(sts, c)

	if in.CacheConfig != nil && *in.CacheConfig != "" {
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_CACHE_CONFIG", Value: *in.CacheConfig})
	}

	if in.OverrideConfig != nil && *in.OverrideConfig != "" {
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_OVERRIDE_CONFIG", Value: *in.OverrideConfig})
	}

	if in.LogLevel != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOG_LEVEL", Value: Int32PtrToString(in.LogLevel)})
	}

	if in.StorageEnabled != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_STORAGE_ENABLED", Value: BoolPtrToString(in.StorageEnabled)})
	}

	if in.IsPersistenceEnabled() {
		// enable persistence environment variable
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"})
	}

	if in.IsSnapshotsEnabled() {
		// enable snapshot environment variable
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"})
	}

	in.Management.AddSSLVolumes(sts, c, VolumeNameManagementSSL, VolumeMountPathManagementCerts)
	c.Env = append(c.Env, in.Management.CreateEnvVars("COH_MGMT", VolumeMountPathManagementCerts, DefaultManagementPort)...)

	in.Metrics.AddSSLVolumes(sts, c, VolumeNameMetricsSSL, VolumeMountPathMetricsCerts)
	c.Env = append(c.Env, in.Metrics.CreateEnvVars("COH_METRICS", VolumeMountPathMetricsCerts, DefaultMetricsPort)...)

	in.AddPersistenceVolumeMounts(c)
	in.AddPersistenceVolumes(sts)
	in.AddPersistencePVCs(cluster, role, sts)
}

// ----- JVMSpec struct -----------------------------------------------------

// The JVM configuration.
// +k8s:openapi-gen=true
type JVMSpec struct {
	// Args specifies the options (System properties, -XX: args etc) to pass to the JVM.
	// +listType=atomic
	// +optional
	Args []string `json:"args,omitempty"`
	// The settings for enabling debug mode in the JVM.
	// +optional
	Debug *JvmDebugSpec `json:"debug,omitempty"`
	// If set to true Adds the  -XX:+UseContainerSupport JVM option to ensure that the JVM
	// respects any container resource limits.
	// The default value is true
	// +optional
	UseContainerLimits *bool `json:"useContainerLimits,omitempty"`
	// If set to true, enabled continuour flight recorder recordings.
	// This will add the JVM options -XX:+UnlockCommercialFeatures -XX:+FlightRecorder
	// -XX:FlightRecorderOptions=defaultrecording=true,dumponexit=true,dumponexitpath=/dumps
	// +optional
	FlightRecorder *bool `json:"flightRecorder,omitempty"`
	// Set JVM garbage collector options.
	// +optional
	Gc *JvmGarbageCollectorSpec `json:"gc,omitempty"`
	// +optional
	DiagnosticsVolume *corev1.VolumeSource `json:"diagnosticsVolume,omitempty"`
	// Configure the JVM memory options.
	// +optional
	Memory *JvmMemorySpec `json:"memory,omitempty"`
	// Configure JMX using JMXMP.
	// +optional
	Jmxmp *JvmJmxmpSpec `json:"jmxmp,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JVMSpec struct with any nil or not set
// values set by the corresponding value in the defaults JVMSpec struct.
func (in *JVMSpec) DeepCopyWithDefaults(defaults *JVMSpec) *JVMSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JVMSpec{}
	clone.Debug = in.Debug.DeepCopyWithDefaults(defaults.Debug)
	clone.Gc = in.Gc.DeepCopyWithDefaults(defaults.Gc)
	clone.Memory = in.Memory.DeepCopyWithDefaults(defaults.Memory)
	clone.Jmxmp = in.Jmxmp.DeepCopyWithDefaults(defaults.Jmxmp)

	if in.UseContainerLimits != nil {
		clone.UseContainerLimits = in.UseContainerLimits
	} else {
		clone.UseContainerLimits = defaults.UseContainerLimits
	}

	if in.FlightRecorder != nil {
		clone.FlightRecorder = in.FlightRecorder
	} else {
		clone.FlightRecorder = defaults.FlightRecorder
	}

	if in.DiagnosticsVolume != nil {
		clone.DiagnosticsVolume = in.DiagnosticsVolume
	} else {
		clone.DiagnosticsVolume = defaults.DiagnosticsVolume
	}

	// Merge Args
	if in.Args != nil {
		clone.Args = []string{}
		clone.Args = append(clone.Args, defaults.Args...)
		clone.Args = append(clone.Args, in.Args...)
	} else if defaults.Args != nil {
		clone.Args = []string{}
		clone.Args = append(clone.Args, defaults.Args...)
	}

	return &clone
}

// Update the StatefulSet with any JVM specific settings
func (in *JVMSpec) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	c := EnsureContainer(ContainerNameCoherence, sts)
	defer ReplaceContainer(sts, c)

	var gc *JvmGarbageCollectorSpec

	if in != nil {
		// Add debug settings
		in.Debug.UpdateCoherenceContainer(c)

		// Add environment variables to the Coherence container
		if in.Args != nil && len(in.Args) > 0 {
			args := strings.Join(in.Args, " ")
			c.Env = append(c.Env, corev1.EnvVar{Name: "JVM_ARGS", Value: args})
		}

		if in.Memory != nil {
			c.Env = append(c.Env, in.Memory.CreateEnvVars()...)
		}

		if in.Jmxmp != nil {
			c.Env = append(c.Env, in.Jmxmp.CreateEnvVars()...)
		}

		if in.Gc != nil {
			gc = in.Gc
		}
	}

	c.Env = append(c.Env, gc.CreateEnvVars()...)

	// Configure the JVM to use container limits (true by default)
	useContainerLimits := in == nil || in.UseContainerLimits == nil || *in.UseContainerLimits
	c.Env = append(c.Env, corev1.EnvVar{Name: "JVM_USE_CONTAINER_LIMITS", Value: strconv.FormatBool(useContainerLimits)})

	// Configure the JVM to use Flight Recorder (true by default)
	useFlightRecorder := in == nil || in.FlightRecorder == nil || *in.FlightRecorder
	c.Env = append(c.Env, corev1.EnvVar{Name: "JVM_FLIGHT_RECORDER", Value: strconv.FormatBool(useFlightRecorder)})

	// Add diagnostic volume if specified otherwise use an empty-volume
	if in != nil && in.DiagnosticsVolume != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         VolumeNameJVM,
			VolumeSource: *in.DiagnosticsVolume,
		})
	} else {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         VolumeNameJVM,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}
}

// ----- ImageSpec struct ---------------------------------------------------

// CoherenceInternalImageSpec defines the settings for a Docker image
// +k8s:openapi-gen=true
type ImageSpec struct {
	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image *string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// Ensure that the image value is set.
func (in *ImageSpec) EnsureImage(image *string) bool {
	if in != nil && in.Image == nil {
		in.Image = image
		return true
	}
	return false
}

// DeepCopyWithDefaults returns a copy of this ImageSpec struct with any nil or not set values set
// by the corresponding value in the defaults ImageSpec struct.
func (in *ImageSpec) DeepCopyWithDefaults(defaults *ImageSpec) *ImageSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ImageSpec{}

	if in.Image != nil {
		clone.Image = in.Image
	} else {
		clone.Image = defaults.Image
	}

	if in.ImagePullPolicy != nil {
		clone.ImagePullPolicy = in.ImagePullPolicy
	} else {
		clone.ImagePullPolicy = defaults.ImagePullPolicy
	}

	return &clone
}

// ----- LoggingSpec struct -------------------------------------------------

// LoggingSpec defines the settings for the Coherence Pod logging
// +k8s:openapi-gen=true
type LoggingSpec struct {
	// ConfigFile allows the location of the Java util logging configuration file to be overridden.
	//  If this value is not set the logging.properties file embedded in this chart will be used.
	//  If this value is set the configuration will be located by trying the following locations in order:
	//    1. If store.logging.configMapName is set then the config map will be mounted as a volume and the logging
	//         properties file will be located as a file location relative to the ConfigMap volume mount point.
	//    2. If userArtifacts.imageName is set then using this value as a file name relative to the location of the
	//         configuration files directory in the user artifacts image.
	//    3. Using this value as an absolute file name.
	// +optional
	ConfigFile *string `json:"configFile,omitempty"`
	// ConfigMapName allows a config map to be mounted as a volume containing the logging
	//  configuration file to use.
	// +optional
	ConfigMapName *string `json:"configMapName,omitempty"`
	// Configures whether Fluentd is enabled and the configuration
	// of the Fluentd side-car container
	// +optional
	Fluentd *FluentdSpec `json:"fluentd,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this LoggingSpec struct with any nil or not set values set
// by the corresponding value in the defaults LoggingSpec struct.
func (in *LoggingSpec) DeepCopyWithDefaults(defaults *LoggingSpec) *LoggingSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := LoggingSpec{}
	clone.Fluentd = in.Fluentd.DeepCopyWithDefaults(defaults.Fluentd)

	if in.ConfigFile != nil {
		clone.ConfigFile = in.ConfigFile
	} else {
		clone.ConfigFile = defaults.ConfigFile
	}

	if in.ConfigMapName != nil {
		clone.ConfigMapName = in.ConfigMapName
	} else {
		clone.ConfigMapName = defaults.ConfigMapName
	}

	return &clone
}

func (in *LoggingSpec) CreateConfigMap(cluster *CoherenceCluster, role *CoherenceRoleSpec) (*corev1.ConfigMap, error) {
	fluentdConfig, err := in.GetFluentdConfig(cluster, role)
	if err != nil {
		return nil, err
	}

	labels := role.CreateCommonLabels(cluster)
	labels[LabelComponent] = LabelComponentEfkConfig

	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf(EfkConfigMapNameTemplate, role.GetFullRoleName(cluster)),
			Labels: labels,
		},
		Data: map[string]string{"fluentd-coherence.conf": fluentdConfig},
	}

	return &cm, nil
}

func (in *LoggingSpec) IsFluentdEnabled() bool {
	return in != nil && in.Fluentd != nil && in.Fluentd.Enabled != nil && *in.Fluentd.Enabled
}

func (in *LoggingSpec) GetFluentdConfig(cluster *CoherenceCluster, role *CoherenceRoleSpec) (string, error) {
	l := LoggingConfigTemplate{
		ClusterName: cluster.Name,
		RoleName:    role.GetRoleName(),
		Logging:     in,
	}

	return l.Parse()
}

// Apply the logging configuration to the StatefulSet
func (in *LoggingSpec) UpdateStatefulSet(sts *appsv1.StatefulSet, hasApp bool) {
	c := EnsureContainer(ContainerNameCoherence, sts)
	defer ReplaceContainer(sts, c)

	if in == nil {
		// just set the default logging config
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: DefaultLoggingConfig})
		return
	}

	if in.ConfigFile != nil && *in.ConfigFile != "" {
		switch {
		case in.ConfigMapName != nil && *in.ConfigMapName != "":
			// Logging config should come from the ConfigMap
			c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: VolumeMountPathLoggingConfig + "/" + *in.ConfigFile})
		case hasApp:
			// Logging config should come from the external config directory
			c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: ExternalConfDir + "/" + *in.ConfigFile})
		default:
			// Logging config is as set
			c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: *in.ConfigFile})
		}
	} else {
		// Logging config is the default
		c.Env = append(c.Env, corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: DefaultLoggingConfig})
	}

	if in.ConfigMapName != nil && *in.ConfigMapName != "" {
		// Append the ConfigMap volume mount
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      VolumeNameLoggingConfig,
			MountPath: VolumeMountPathLoggingConfig,
		})

		// Append the ConfigMap volume
		vol := corev1.Volume{
			Name: VolumeNameLoggingConfig,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "test-logging-configmap",
					},
					DefaultMode: pointer.Int32Ptr(int32(0777)),
				},
			},
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
	}

	// Apply any fluentd configuration
	in.Fluentd.UpdateStatefulSet(sts)
}

// ----- LoggingConfigTemplate struct ---------------------------------------

// A struct used when converting the fluentd config template to a string using go templating.
type LoggingConfigTemplate struct {
	ClusterName string
	RoleName    string
	Logging     *LoggingSpec
}

// Parse the fluentd configuration.
func (in LoggingConfigTemplate) Parse() (string, error) {
	t, err := template.New("efk").Parse(EfkConfig)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, in); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ----- PersistentStorageSpec struct ---------------------------------------

// PersistenceStorageSpec defines the persistence settings for the Coherence
// +k8s:openapi-gen=true
type PersistentStorageSpec struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// PersistentVolumeClaim allows the configuration of a normal k8s persistent volume claim
	// for persistence data.
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"` // from k8s.io/api/core/v1
	// Volume allows the configuration of a normal k8s volume mapping
	// for persistence data instead of a persistent volume claim. If a value is defined
	// for store.persistence.volume then no PVC will be created and persistence data
	// will instead be written to this volume. It is up to the deployer to understand
	// the consequences of this and how the guarantees given when using PVCs differ
	// to the storage guarantees for the particular volume type configured here.
	// +optional
	Volume *corev1.VolumeSource `json:"volume,omitempty"` // from k8s.io/api/core/v1
}

// DeepCopyWithDefaults returns a copy of this PersistentStorageSpec struct with any nil or not set values set
// by the corresponding value in the defaults PersistentStorageSpec struct.
func (in *PersistentStorageSpec) DeepCopyWithDefaults(defaults *PersistentStorageSpec) *PersistentStorageSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := PersistentStorageSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.PersistentVolumeClaim != nil {
		clone.PersistentVolumeClaim = in.PersistentVolumeClaim
	} else {
		clone.PersistentVolumeClaim = defaults.PersistentVolumeClaim
	}

	if in.Volume != nil {
		clone.Volume = in.Volume
	} else {
		clone.Volume = defaults.Volume
	}

	return &clone
}

// Create a PersistentVolumeClaim
func (in *PersistentStorageSpec) CreatePersistentVolumeClaim(cluster *CoherenceCluster, role *CoherenceRoleSpec, name string) (bool, *corev1.PersistentVolumeClaim) {
	if in == nil || in.Enabled == nil || !*in.Enabled || in.Volume != nil {
		// Either persistence is disabled or we're using a normal Volume
		return false, nil
	}

	spec := corev1.PersistentVolumeClaimSpec{}
	if in.PersistentVolumeClaim != nil {
		in.PersistentVolumeClaim.DeepCopyInto(&spec)
	}

	labels := role.CreateCommonLabels(cluster)
	labels[LabelComponent] = LabelComponentPVC

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: spec,
	}

	return true, pvc
}

// ----- SSLSpec struct -----------------------------------------------------

// SSLSpec defines the SSL settings for a Coherence component over REST endpoint.
// +k8s:openapi-gen=true
type SSLSpec struct {
	// Enabled is a boolean flag indicating whether enables or disables SSL on the Coherence management
	// over REST endpoint, the default is false (disabled).
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Secrets is the name of the k8s secrets containing the Java key stores and password files.
	//   This value MUST be provided if SSL is enabled on the Coherence management over ReST endpoint.
	// +optional
	Secrets *string `json:"secrets,omitempty"`
	// Keystore is the name of the Java key store file in the k8s secret to use as the SSL keystore
	//   when configuring component over REST to use SSL.
	// +optional
	KeyStore *string `json:"keyStore,omitempty"`
	// KeyStorePasswordFile is the name of the file in the k8s secret containing the keystore
	//   password when configuring component over REST to use SSL.
	// +optional
	KeyStorePasswordFile *string `json:"keyStorePasswordFile,omitempty"`
	// KeyStorePasswordFile is the name of the file in the k8s secret containing the key
	//   password when configuring component over REST to use SSL.
	// +optional
	KeyPasswordFile *string `json:"keyPasswordFile,omitempty"`
	// KeyStoreAlgorithm is the name of the keystore algorithm for the keystore in the k8s secret
	//   used when configuring component over REST to use SSL. If not set the default is SunX509
	// +optional
	KeyStoreAlgorithm *string `json:"keyStoreAlgorithm,omitempty"`
	// KeyStoreProvider is the name of the keystore provider for the keystore in the k8s secret
	//   used when configuring component over REST to use SSL.
	// +optional
	KeyStoreProvider *string `json:"keyStoreProvider,omitempty"`
	// KeyStoreType is the name of the Java keystore type for the keystore in the k8s secret used
	//   when configuring component over REST to use SSL. If not set the default is JKS.
	// +optional
	KeyStoreType *string `json:"keyStoreType,omitempty"`
	// TrustStore is the name of the Java trust store file in the k8s secret to use as the SSL
	//   trust store when configuring component over REST to use SSL.
	// +optional
	TrustStore *string `json:"trustStore,omitempty"`
	// TrustStorePasswordFile is the name of the file in the k8s secret containing the trust store
	//   password when configuring component over REST to use SSL.
	// +optional
	TrustStorePasswordFile *string `json:"trustStorePasswordFile,omitempty"`
	// TrustStoreAlgorithm is the name of the keystore algorithm for the trust store in the k8s
	//   secret used when configuring component over REST to use SSL.  If not set the default is SunX509.
	// +optional
	TrustStoreAlgorithm *string `json:"trustStoreAlgorithm,omitempty"`
	// TrustStoreProvider is the name of the keystore provider for the trust store in the k8s
	//   secret used when configuring component over REST to use SSL.
	// +optional
	TrustStoreProvider *string `json:"trustStoreProvider,omitempty"`
	// TrustStoreType is the name of the Java keystore type for the trust store in the k8s secret
	//   used when configuring component over REST to use SSL. If not set the default is JKS.
	// +optional
	TrustStoreType *string `json:"trustStoreType,omitempty"`
	// RequireClientCert is a boolean flag indicating whether the client certificate will be
	//   authenticated by the server (two-way SSL) when configuring component over REST to use SSL.
	//   If not set the default is false
	// +optional
	RequireClientCert *bool `json:"requireClientCert,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this SSLSpec struct with any nil or not set values set
// by the corresponding value in the defaults SSLSpec struct.
func (in *SSLSpec) DeepCopyWithDefaults(defaults *SSLSpec) *SSLSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := SSLSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Secrets != nil {
		clone.Secrets = in.Secrets
	} else {
		clone.Secrets = defaults.Secrets
	}

	if in.KeyStore != nil {
		clone.KeyStore = in.KeyStore
	} else {
		clone.KeyStore = defaults.KeyStore
	}

	if in.KeyStorePasswordFile != nil {
		clone.KeyStorePasswordFile = in.KeyStorePasswordFile
	} else {
		clone.KeyStorePasswordFile = defaults.KeyStorePasswordFile
	}

	if in.KeyPasswordFile != nil {
		clone.KeyPasswordFile = in.KeyPasswordFile
	} else {
		clone.KeyPasswordFile = defaults.KeyPasswordFile
	}

	if in.KeyStoreAlgorithm != nil {
		clone.KeyStoreAlgorithm = in.KeyStoreAlgorithm
	} else {
		clone.KeyStoreAlgorithm = defaults.KeyStoreAlgorithm
	}

	if in.KeyStoreProvider != nil {
		clone.KeyStoreProvider = in.KeyStoreProvider
	} else {
		clone.KeyStoreProvider = defaults.KeyStoreProvider
	}

	if in.KeyStoreType != nil {
		clone.KeyStoreType = in.KeyStoreType
	} else {
		clone.KeyStoreType = defaults.KeyStoreType
	}

	if in.TrustStore != nil {
		clone.TrustStore = in.TrustStore
	} else {
		clone.TrustStore = defaults.TrustStore
	}

	if in.TrustStorePasswordFile != nil {
		clone.TrustStorePasswordFile = in.TrustStorePasswordFile
	} else {
		clone.TrustStorePasswordFile = defaults.TrustStorePasswordFile
	}

	if in.TrustStoreAlgorithm != nil {
		clone.TrustStoreAlgorithm = in.TrustStoreAlgorithm
	} else {
		clone.TrustStoreAlgorithm = defaults.TrustStoreAlgorithm
	}

	if in.TrustStoreProvider != nil {
		clone.TrustStoreProvider = in.TrustStoreProvider
	} else {
		clone.TrustStoreProvider = defaults.TrustStoreProvider
	}

	if in.TrustStoreType != nil {
		clone.TrustStoreType = in.TrustStoreType
	} else {
		clone.TrustStoreType = defaults.TrustStoreType
	}

	if in.RequireClientCert != nil {
		clone.RequireClientCert = in.RequireClientCert
	} else {
		clone.RequireClientCert = defaults.RequireClientCert
	}

	return &clone
}

// Create the SSL environment variables
func (in *SSLSpec) CreateEnvVars(prefix, secretMount string) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil {
		return envVars
	}

	if in.Enabled != nil && *in.Enabled {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_ENABLED", Value: "true"})
	}

	if in.Secrets != nil && *in.Secrets != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_CERTS", Value: secretMount})
	}

	if in.KeyStore != nil && *in.KeyStore != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEYSTORE", Value: *in.KeyStore})
	}

	if in.KeyStorePasswordFile != nil && *in.KeyStorePasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEYSTORE_PASSWORD_FILE", Value: *in.KeyStorePasswordFile})
	}

	if in.KeyPasswordFile != nil && *in.KeyPasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEY_PASSWORD_FILE", Value: *in.KeyPasswordFile})
	}

	if in.KeyStoreAlgorithm != nil && *in.KeyStoreAlgorithm != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEYSTORE_ALGORITHM", Value: *in.KeyStoreAlgorithm})
	}

	if in.KeyStoreProvider != nil && *in.KeyStoreProvider != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEYSTORE_PROVIDER", Value: *in.KeyStoreProvider})
	}

	if in.KeyStoreType != nil && *in.KeyStoreType != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_KEYSTORE_TYPE", Value: *in.KeyStoreType})
	}

	if in.TrustStore != nil && *in.TrustStore != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_TRUSTSTORE", Value: *in.TrustStore})
	}

	if in.TrustStorePasswordFile != nil && *in.TrustStorePasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_TRUSTSTORE_PASSWORD_FILE", Value: *in.TrustStorePasswordFile})
	}

	if in.TrustStoreAlgorithm != nil && *in.TrustStoreAlgorithm != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_TRUSTSTORE_ALGORITHM", Value: *in.TrustStoreAlgorithm})
	}

	if in.TrustStoreProvider != nil && *in.TrustStoreProvider != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_TRUSTSTORE_PROVIDER", Value: *in.TrustStoreProvider})
	}

	if in.TrustStoreType != nil && *in.TrustStoreType != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_TRUSTSTORE_TYPE", Value: *in.TrustStoreType})
	}

	if in.RequireClientCert != nil && *in.RequireClientCert {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + "_SSL_REQUIRE_CLIENT_CERT", Value: "true"})
	}

	return envVars
}

// ----- PortSpec struct ----------------------------------------------------
// PortSpec defines the port settings for a Coherence component
// +k8s:openapi-gen=true
type PortSpec struct {
	// Port specifies the port used.
	// +optional
	Port int32 `json:"port,omitempty"`
	// Protocol for container port. Must be UDP or TCP. Defaults to "TCP"
	// +optional
	Protocol *corev1.Protocol `json:"protocol,omitempty"`
	// Service specifies the service used to expose the port.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`
	// The port on each node on which this service is exposed when type=NodePort or LoadBalancer.
	// Usually assigned by the system. If specified, it will be allocated to the service
	// if unused or else creation of the service will fail.
	// Default is to auto-allocate a port if the ServiceType of this Service requires one.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +optional
	NodePort *int32 `json:"nodePort,omitempty"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	HostPort *int32 `json:"hostPort,omitempty"`
	// What host IP to bind the external port to.
	// +optional
	HostIP *string `json:"hostIP,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this PortSpec struct with any nil or not set values set
// by the corresponding value in the defaults PortSpec struct.
func (in *PortSpec) DeepCopyWithDefaults(defaults *PortSpec) *PortSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := PortSpec{}

	if in.Port != 0 {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.Protocol != nil {
		clone.Protocol = in.Protocol
	} else {
		clone.Protocol = defaults.Protocol
	}

	if in.Service != nil {
		clone.Service = in.Service
	} else {
		clone.Service = defaults.Service
	}

	if in.NodePort != nil {
		clone.NodePort = in.HostPort
	} else {
		clone.NodePort = defaults.HostPort
	}

	if in.HostPort != nil {
		clone.HostPort = in.HostPort
	} else {
		clone.HostPort = defaults.HostPort
	}

	if in.HostIP != nil {
		clone.HostIP = in.HostIP
	} else {
		clone.HostIP = defaults.HostIP
	}

	return &clone
}

// ----- NamedPortSpec struct ----------------------------------------------------
// NamedPortSpec defines a named port for a Coherence component
// +k8s:openapi-gen=true
type NamedPortSpec struct {
	// Name specifies the name of th port.
	// +optional
	Name     string `json:"name,omitempty"`
	PortSpec `json:",inline"`
}

// DeepCopyWithDefaults returns a copy of this NamedPortSpec struct with any nil or not set values set
// by the corresponding value in the defaults NamedPortSpec struct.
func (in *NamedPortSpec) DeepCopyWithDefaults(defaults *NamedPortSpec) *NamedPortSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := NamedPortSpec{}

	if in.Name != "" {
		clone.Name = in.Name
	} else {
		clone.Name = defaults.Name
	}

	if in.Port != 0 {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.Protocol != nil {
		clone.Protocol = in.Protocol
	} else {
		clone.Protocol = defaults.Protocol
	}

	if in.Service != nil {
		clone.Service = in.Service.DeepCopyWithDefaults(defaults.Service)
	} else {
		clone.Service = defaults.Service
	}

	return &clone
}

// Create the Kubernetes service to expose this port.
func (in *NamedPortSpec) CreateService(cluster *CoherenceCluster, role *CoherenceRoleSpec) *corev1.Service {
	if in == nil || !in.IsEnabled() {
		return nil
	}

	var name string
	if in.Service != nil && in.Service.Name != nil {
		name = in.Service.GetName()
	} else {
		name = fmt.Sprintf("%s-%s", role.GetFullRoleName(cluster), in.Name)
	}

	// The labels for the service
	svcLabels := role.CreateCommonLabels(cluster)
	svcLabels[LabelComponent] = LabelComponentPortService
	if in.Service != nil {
		for k, v := range in.Service.Labels {
			svcLabels[k] = v
		}
	}

	// The service annotations
	var ann map[string]string
	if in.Service != nil && in.Service.Annotations != nil {
		ann = in.Service.Annotations
	}

	// Create the Service spec
	spec := in.Service.createServiceSpec()

	// Add the port
	spec.Ports = []corev1.ServicePort{
		{
			Name:       in.Name,
			Protocol:   in.GetProtocol(),
			Port:       in.GetPort(),
			TargetPort: intstr.FromString(in.Name),
			NodePort:   in.GetNodePort(),
		},
	}

	// Add the service selector
	spec.Selector = role.CreatePodSelectorLabels(cluster)

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      svcLabels,
			Annotations: ann,
		},
		Spec: spec,
	}

	return &svc
}

func (in *NamedPortSpec) IsEnabled() bool {
	return in != nil && in.Service.IsEnabled()
}

func (in *NamedPortSpec) GetProtocol() corev1.Protocol {
	if in == nil || in.Protocol == nil {
		return corev1.ProtocolTCP
	}
	return *in.Protocol
}

func (in *NamedPortSpec) GetPort() int32 {
	switch {
	case in == nil:
		return 0
	case in != nil && in.Service != nil && in.Service.Port != nil:
		return *in.Service.Port
	default:
		return in.Port
	}
}

func (in *NamedPortSpec) GetNodePort() int32 {
	if in == nil || in.NodePort == nil {
		return 0
	}
	return *in.NodePort
}

func (in *NamedPortSpec) CreatePort() corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          in.Name,
		ContainerPort: in.GetPort(),
		Protocol:      in.GetProtocol(),
		HostPort:      notNilInt32(in.HostPort),
		HostIP:        notNilString(in.HostIP),
	}
}

// Merge merges two arrays of NamedPortSpec structs.
// Any NamedPortSpec instances in both arrays that share the same name will be merged,
// the field set in the primary NamedPortSpec will take precedence over those in the
// secondary NamedPortSpec.
func MergeNamedPortSpecs(primary, secondary []NamedPortSpec) []NamedPortSpec {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []NamedPortSpec{}
	}

	var mr []NamedPortSpec
	mr = append(mr, primary...)

	for _, p := range secondary {
		found := false
		for i, pp := range primary {
			if pp.Name == p.Name {
				clone := pp.DeepCopyWithDefaults(&p)
				mr[i] = *clone
				found = true
				break
			}
		}

		if !found {
			mr = append(mr, p)
		}
	}

	return mr
}

// ----- JvmDebugSpec struct ---------------------------------------------------

// The JVM Debug specific configuration.
// See:
// +k8s:openapi-gen=true
type JvmDebugSpec struct {
	// Enabled is a flag to enable or disable running the JVM in debug mode. Default is disabled.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// A boolean true if the target VM is to be suspended immediately before the main class is loaded;
	// false otherwise. The default value is false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`
	// Attach specifies the address of the debugger that the JVM should attempt to connect back to
	// instead of listening on a port.
	// +optional
	Attach *string `json:"attach,omitempty"`
	// The port that the debugger will listen on; the default is 5005.
	// +optional
	Port *int32 `json:"port,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JvmDebugSpec struct with any nil or not set values set
// by the corresponding value in the defaults JvmDebugSpec struct.
func (in *JvmDebugSpec) DeepCopyWithDefaults(defaults *JvmDebugSpec) *JvmDebugSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JvmDebugSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Suspend != nil {
		clone.Suspend = in.Suspend
	} else {
		clone.Suspend = defaults.Suspend
	}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.Attach != nil {
		clone.Attach = in.Attach
	} else {
		clone.Attach = defaults.Attach
	}

	return &clone
}

// Update the Coherence Container with any JVM specific settings
func (in *JvmDebugSpec) UpdateCoherenceContainer(c *corev1.Container) {
	if in == nil || in.Enabled == nil || !*in.Enabled {
		// nothing to do, debug is either nil or disabled
		return
	}

	c.Ports = append(c.Ports, corev1.ContainerPort{
		Name:          PortNameDebug,
		ContainerPort: notNilInt32OrDefault(in.Port, DefaultDebugPort),
	})

	c.Env = append(c.Env, in.CreateEnvVars()...)
}

// Create the JVM debugger environment variables for the Coherence container.
func (in *JvmDebugSpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil || in.Enabled == nil || !*in.Enabled {
		return envVars
	}

	envVars = append(envVars,
		corev1.EnvVar{Name: "JVM_DEBUG_ENABLED", Value: "true"},
		corev1.EnvVar{Name: "JVM_DEBUG_PORT", Value: Int32PtrToStringWithDefault(in.Port, DefaultDebugPort)},
	)

	if in != nil && in.Suspend != nil && *in.Suspend {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_DEBUG_SUSPEND", Value: "true"})
	}

	if in != nil && in.Attach != nil {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_DEBUG_ATTACH", Value: *in.Attach})
	}

	return envVars
}

// ----- JVM GC struct ------------------------------------------------------

// Options for managing the JVM garbage collector.
type JvmGarbageCollectorSpec struct {
	// The name of the JVM garbage collector to use.
	// G1 - adds the -XX:+UseG1GC option
	// CMS - adds the -XX:+UseConcMarkSweepGC option
	// Parallel - adds the -XX:+UseParallelGC
	// Default - use the JVMs default collector
	// The field value is case insensitive
	// If not set G1 is used.
	// If set to a value other than those above then
	// the default collector for the JVM will be used.
	// +optional
	Collector *string `json:"collector,omitempty"`
	// Args specifies the GC options to pass to the JVM.
	// +optional
	Args []string `json:"args,omitempty"`
	// Enable the following GC logging args  -verbose:gc -XX:+PrintGCDetails -XX:+PrintGCTimeStamps
	// -XX:+PrintHeapAtGC -XX:+PrintTenuringDistribution -XX:+PrintGCApplicationStoppedTime
	// -XX:+PrintGCApplicationConcurrentTime
	// Default is true
	// +optional
	Logging *bool `json:"logging,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JvmGarbageCollectorSpec struct with any nil or not set values set
// by the corresponding value in the defaults JvmGarbageCollectorSpec struct.
func (in *JvmGarbageCollectorSpec) DeepCopyWithDefaults(defaults *JvmGarbageCollectorSpec) *JvmGarbageCollectorSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JvmGarbageCollectorSpec{}

	if in.Collector != nil {
		clone.Collector = in.Collector
	} else {
		clone.Collector = defaults.Collector
	}

	if in.Args != nil {
		clone.Args = in.Args
	} else {
		clone.Args = defaults.Args
	}

	if in.Logging != nil {
		clone.Logging = in.Logging
	} else {
		clone.Logging = defaults.Logging
	}

	return &clone
}

// Create the GC environment variables for the Coherence container.
func (in *JvmGarbageCollectorSpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	// Add any GC args
	if in != nil && in.Args != nil && len(in.Args) > 0 {
		args := strings.Join(in.Args, " ")
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_GC_ARGS", Value: args})
	}

	// Set the collector to use
	if in != nil && in.Collector != nil && *in.Collector != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_GC_COLLECTOR", Value: *in.Collector})
	}

	// Enable or disable GC logging
	if in != nil && in.Logging != nil {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_GC_LOGGING", Value: BoolPtrToString(in.Logging)})
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_GC_LOGGING", Value: "true"})
	}

	return envVars
}

// ----- JVM MemoryGC struct ------------------------------------------------

// Options for managing the JVM memory.
type JvmMemorySpec struct {
	// HeapSize is the min/max heap value to pass to the JVM.
	// The format should be the same as that used for Java's -Xms and -Xmx JVM options.
	// If not set the JVM defaults are used.
	// +optional
	HeapSize *string `json:"heapSize,omitempty"`
	// StackSize is the stack sixe value to pass to the JVM.
	// The format should be the same as that used for Java's -Xss JVM option.
	// If not set the JVM defaults are used.
	// +optional
	StackSize *string `json:"stackSize,omitempty"`
	// MetaspaceSize is the min/max metaspace size to pass to the JVM.
	// This sets the -XX:MetaspaceSize and -XX:MaxMetaspaceSize=size JVM options.
	// If not set the JVM defaults are used.
	// +optional
	MetaspaceSize *string `json:"metaspaceSize,omitempty"`
	// DirectMemorySize sets the maximum total size (in bytes) of the New I/O (the java.nio package) direct-buffer
	// allocations. This value sets the -XX:MaxDirectMemorySize JVM option.
	// If not set the JVM defaults are used.
	// +optional
	DirectMemorySize *string `json:"directMemorySize,omitempty"`
	// Adds the -XX:NativeMemoryTracking=mode  JVM options
	// where mode is on of "off", "summary" or "detail", the default is "summary"
	// If not set to "off" also add -XX:+PrintNMTStatistics
	// +optional
	NativeMemoryTracking *string `json:"nativeMemoryTracking,omitempty"`
	// Configure the JVM behaviour when an OutOfMemoryError occurs.
	// +optional
	OnOutOfMemory *JvmOutOfMemorySpec `json:"onOutOfMemory,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JvmMemorySpec struct with any nil or not set values set
// by the corresponding value in the defaults JvmMemorySpec struct.
func (in *JvmMemorySpec) DeepCopyWithDefaults(defaults *JvmMemorySpec) *JvmMemorySpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JvmMemorySpec{}
	clone.OnOutOfMemory = in.OnOutOfMemory.DeepCopyWithDefaults(defaults.OnOutOfMemory)

	if in.HeapSize != nil {
		clone.HeapSize = in.HeapSize
	} else {
		clone.HeapSize = defaults.HeapSize
	}

	if in.StackSize != nil {
		clone.StackSize = in.StackSize
	} else {
		clone.StackSize = defaults.StackSize
	}

	if in.MetaspaceSize != nil {
		clone.MetaspaceSize = in.MetaspaceSize
	} else {
		clone.MetaspaceSize = defaults.MetaspaceSize
	}

	if in.DirectMemorySize != nil {
		clone.DirectMemorySize = in.DirectMemorySize
	} else {
		clone.DirectMemorySize = defaults.DirectMemorySize
	}

	if in.NativeMemoryTracking != nil {
		clone.NativeMemoryTracking = in.NativeMemoryTracking
	} else {
		clone.NativeMemoryTracking = defaults.NativeMemoryTracking
	}

	return &clone
}

// Create the environment variables to add to the Coherence container
func (in *JvmMemorySpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil {
		return envVars
	}

	if in.HeapSize != nil && *in.HeapSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_HEAP_SIZE", Value: *in.HeapSize})
	}

	if in.DirectMemorySize != nil && *in.DirectMemorySize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_DIRECT_MEMORY_SIZE", Value: *in.DirectMemorySize})
	}

	if in.StackSize != nil && *in.StackSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_STACK_SIZE", Value: *in.StackSize})
	}

	if in.MetaspaceSize != nil && *in.MetaspaceSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_METASPACE_SIZE", Value: *in.MetaspaceSize})
	}

	if in.NativeMemoryTracking != nil && *in.NativeMemoryTracking != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "JVM_NATIVE_MEMORY_TRACKING", Value: *in.NativeMemoryTracking})
	}

	envVars = append(envVars, in.OnOutOfMemory.CreateEnvVars()...)

	return envVars
}

// ----- JVM Out Of Memory struct -------------------------------------------

// Options for managing the JVM behaviour when an OutOfMemoryError occurs.
type JvmOutOfMemorySpec struct {
	// If set to true the JVM will exit when an OOM error occurs.
	// Default is true
	// +optional
	Exit *bool `json:"exit,omitempty"`
	// If set to true adds the -XX:+HeapDumpOnOutOfMemoryError JVM option to cause a heap dump
	// to be created when an OOM error occurs.
	// Default is true
	// +optional
	HeapDump *bool `json:"heapDump,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JvmOutOfMemorySpec struct with any nil or not set values set
// by the corresponding value in the defaults JvmOutOfMemorySpec struct.
func (in *JvmOutOfMemorySpec) DeepCopyWithDefaults(defaults *JvmOutOfMemorySpec) *JvmOutOfMemorySpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JvmOutOfMemorySpec{}

	if in.Exit != nil {
		clone.Exit = in.Exit
	} else {
		clone.Exit = defaults.Exit
	}

	if in.HeapDump != nil {
		clone.HeapDump = in.HeapDump
	} else {
		clone.HeapDump = defaults.HeapDump
	}

	return &clone
}

// Create the environment variables for the out of memory actions
func (in *JvmOutOfMemorySpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in != nil {
		if in.Exit != nil {
			envVars = append(envVars, corev1.EnvVar{Name: "JVM_OOM_EXIT", Value: BoolPtrToString(in.Exit)})
		}
		if in.HeapDump != nil {
			envVars = append(envVars, corev1.EnvVar{Name: "JVM_OOM_HEAP_DUMP", Value: BoolPtrToString(in.HeapDump)})
		}
	}

	return envVars
}

// ----- JvmJmxmpSpec struct -------------------------------------------------------

// Options for configuring JMX using JMXMP.
type JvmJmxmpSpec struct {
	// If set to true the JMXMP support will be enabled.
	// Default is false
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// The port tht the JMXMP MBeanServer should bind to.
	// If not set the default port is 9099
	// +optional
	Port *int32 `json:"port,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JvmJmxmpSpec struct with any nil or not set values set
// by the corresponding value in the defaults JvmJmxmpSpec struct.
func (in *JvmJmxmpSpec) DeepCopyWithDefaults(defaults *JvmJmxmpSpec) *JvmJmxmpSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := JvmJmxmpSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	return &clone
}

// Create any required environment variables for the Coherence container
func (in *JvmJmxmpSpec) CreateEnvVars() []corev1.EnvVar {
	enabled := in != nil && in.Enabled != nil && *in.Enabled

	envVars := []corev1.EnvVar{{Name: "JVM_JMXMP_ENABLED", Value: strconv.FormatBool(enabled)}}
	envVars = append(envVars, corev1.EnvVar{Name: "JVM_JMXMP_PORT", Value: Int32PtrToStringWithDefault(in.Port, DefaultJmxmpPort)})

	return envVars
}

// ----- PortSpecWithSSL struct ----------------------------------------------------

// PortSpecWithSSL defines a port with SSL settings for a Coherence component
// +k8s:openapi-gen=true
type PortSpecWithSSL struct {
	// Enable or disable flag.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// The port to bind to.
	// +optional
	Port *int32 `json:"port,omitempty"`
	// SSL configures SSL settings for a Coherence component
	// +optional
	SSL *SSLSpec `json:"ssl,omitempty"`
}

// IsSSLEnabled returns true if this port is SSL enabled
func (in *PortSpecWithSSL) IsSSLEnabled() bool {
	return in != nil && in.Enabled != nil && *in.Enabled
}

// DeepCopyWithDefaults returns a copy of this PortSpecWithSSL struct with any nil or not set values set
// by the corresponding value in the defaults PortSpecWithSSL struct.
func (in *PortSpecWithSSL) DeepCopyWithDefaults(defaults *PortSpecWithSSL) *PortSpecWithSSL {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := PortSpecWithSSL{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.SSL != nil {
		clone.SSL = in.SSL
	} else {
		clone.SSL = defaults.SSL
	}

	return &clone
}

// Create environment variables for the Coherence container
func (in *PortSpecWithSSL) CreateEnvVars(prefix, secretMount string, defaultPort int32) []corev1.EnvVar {
	if in == nil || !notNilBool(in.Enabled) {
		// disabled
		return []corev1.EnvVar{{Name: prefix + "_ENABLED", Value: "false"}}
	}

	envVars := []corev1.EnvVar{{Name: prefix + "_ENABLED", Value: "true"}}
	envVars = append(envVars, in.SSL.CreateEnvVars(prefix, secretMount)...)

	// add the port environment variable
	port := notNilInt32OrDefault(in.Port, defaultPort)
	envVars = append(envVars, corev1.EnvVar{Name: prefix + "_PORT", Value: Int32ToString(port)})

	return envVars
}

// Add the SSL secret volume and volume mount if required
func (in *PortSpecWithSSL) AddSSLVolumes(sts *appsv1.StatefulSet, c *corev1.Container, volName, path string) {
	if in == nil || !notNilBool(in.Enabled) || in.SSL == nil || !notNilBool(in.SSL.Enabled) {
		// the port spec is nil or disabled or SSL is nil or disabled
		return
	}

	if in.SSL.Secrets != nil && *in.SSL.Secrets != "" {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volName,
			ReadOnly:  true,
			MountPath: path,
		})

		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  *in.SSL.Secrets,
					DefaultMode: pointer.Int32Ptr(int32(0777)),
				},
			},
		})
	}

}

// ----- ServiceSpec struct -------------------------------------------------
// ServiceSpec defines the settings for a Service
// +k8s:openapi-gen=true
type ServiceSpec struct {
	// Enabled controls whether to create the service yaml or not
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// An optional name to use to override the generated service name.
	// +optional
	Name *string `json:"name,omitempty"`
	// The service port value
	// +optional
	Port *int32 `json:"port,omitempty"`
	// Type is the K8s service type (typically ClusterIP or LoadBalancer)
	// The default is "ClusterIP".
	// +optional
	Type *corev1.ServiceType `json:"type,omitempty"`
	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +optional
	ExternalIPs []string `json:"externalIPs,omitempty"`
	// clusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service; otherwise, creation
	// of the service will fail. This field can not be changed through updates.
	// Valid values are "None", empty string (""), or a valid IP address. "None"
	// can be specified for headless services when proxying is not required.
	// Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if
	// type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP *string `json:"clusterIP,omitempty"`
	// LoadBalancerIP is the IP address of the load balancer
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`
	// The extra labels to add to the service.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is free form yaml that will be added to the service annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Supports "ClientIP" and "None". Used to maintain session affinity.
	// Enable client IP based session affinity.
	// Must be ClientIP or None.
	// Defaults to None.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	SessionAffinity *corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`
	// If specified and supported by the platform, this will restrict traffic through the cloud-provider
	// load-balancer will be restricted to the specified client IPs. This field will be ignored if the
	// cloud-provider does not support the feature."
	// More info: https://kubernetes.io/docs/tasks/access-application-cluster/configure-cloud-provider-firewall/
	// +listType=atomic
	// +optional
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`
	// externalName is the external reference that kubedns or equivalent will
	// return as a CNAME record for this service. No proxying will be involved.
	// Must be a valid RFC-1123 hostname (https://tools.ietf.org/html/rfc1123)
	// and requires Type to be ExternalName.
	// +optional
	ExternalName *string `json:"externalName,omitempty"`
	// externalTrafficPolicy denotes if this Service desires to route external
	// traffic to node-local or cluster-wide endpoints. "Local" preserves the
	// client source IP and avoids a second hop for LoadBalancer and Nodeport
	// type services, but risks potentially imbalanced traffic spreading.
	// "Cluster" obscures the client source IP and may cause a second hop to
	// another node, but should have good overall load-spreading.
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	// healthCheckNodePort specifies the healthcheck nodePort for the service.
	// If not specified, HealthCheckNodePort is created by the service api
	// backend with the allocated nodePort. Will use user-specified nodePort value
	// if specified by the client. Only effects when Type is set to LoadBalancer
	// and ExternalTrafficPolicy is set to Local.
	// +optional
	HealthCheckNodePort *int32 `json:"healthCheckNodePort,omitempty"`
	// publishNotReadyAddresses, when set to true, indicates that DNS implementations
	// must publish the notReadyAddresses of subsets for the Endpoints associated with
	// the Service. The default value is false.
	// The primary use case for setting this field is to use a StatefulSet's Headless Service
	// to propagate SRV records for its Pods without respect to their readiness for purpose
	// of peer discovery.
	// +optional
	PublishNotReadyAddresses *bool `json:"publishNotReadyAddresses,omitempty"`
	// sessionAffinityConfig contains the configurations of session affinity.
	// +optional
	SessionAffinityConfig *corev1.SessionAffinityConfig `json:"sessionAffinityConfig,omitempty"`
	// ipFamily specifies whether this Service has a preference for a particular IP family (e.g. IPv4 vs.
	// IPv6).  If a specific IP family is requested, the clusterIP field will be allocated from that family, if it is
	// available in the cluster.  If no IP family is requested, the cluster's primary IP family will be used.
	// Other IP fields (loadBalancerIP, loadBalancerSourceRanges, externalIPs) and controllers which
	// allocate external load-balancers should use the same IP family.  Endpoints for this Service will be of
	// this family.  This field is immutable after creation. Assigning a ServiceIPFamily not available in the
	// cluster (e.g. IPv6 in IPv4 only cluster) is an error condition and will fail during clusterIP assignment.
	// +optional
	IPFamily *corev1.IPFamily `json:"ipFamily,omitempty"`
}

// Set the Type of the service.
func (in *ServiceSpec) GetName() string {
	if in == nil || in.Name == nil {
		return ""
	}
	return *in.Name
}

// Set the Type of the service.
func (in *ServiceSpec) IsEnabled() bool {
	if in == nil || in.Enabled == nil {
		return true
	}
	return *in.Enabled
}

// Set the Type of the service.
func (in *ServiceSpec) SetServiceType(t corev1.ServiceType) {
	if in != nil {
		in.Type = &t
	}
}

// DeepCopyWithDefaults returns a copy of this ServiceSpec struct with any nil or not set values set
// by the corresponding value in the defaults ServiceSpec struct.
func (in *ServiceSpec) DeepCopyWithDefaults(defaults *ServiceSpec) *ServiceSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ServiceSpec{}
	// Annotations are a map and are merged
	clone.Annotations = MergeMap(in.Annotations, defaults.Annotations)
	// Labels are a map and are merged
	clone.Labels = MergeMap(in.Labels, defaults.Labels)

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Type != nil {
		clone.Type = in.Type
	} else {
		clone.Type = defaults.Type
	}

	if in.Name != nil {
		clone.Name = in.Name
	} else {
		clone.Name = defaults.Name
	}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.LoadBalancerIP != nil {
		clone.LoadBalancerIP = in.LoadBalancerIP
	} else {
		clone.LoadBalancerIP = defaults.LoadBalancerIP
	}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.SessionAffinity != nil {
		clone.SessionAffinity = in.SessionAffinity
	} else {
		clone.SessionAffinity = defaults.SessionAffinity
	}

	if in.LoadBalancerSourceRanges != nil {
		clone.LoadBalancerSourceRanges = in.LoadBalancerSourceRanges
	} else {
		clone.LoadBalancerSourceRanges = defaults.LoadBalancerSourceRanges
	}

	if in.ExternalName != nil {
		clone.ExternalName = in.ExternalName
	} else {
		clone.ExternalName = defaults.ExternalName
	}

	if in.ExternalTrafficPolicy != nil {
		clone.ExternalTrafficPolicy = in.ExternalTrafficPolicy
	} else {
		clone.ExternalTrafficPolicy = defaults.ExternalTrafficPolicy
	}

	if in.HealthCheckNodePort != nil {
		clone.HealthCheckNodePort = in.HealthCheckNodePort
	} else {
		clone.HealthCheckNodePort = defaults.HealthCheckNodePort
	}

	if in.PublishNotReadyAddresses != nil {
		clone.PublishNotReadyAddresses = in.PublishNotReadyAddresses
	} else {
		clone.PublishNotReadyAddresses = defaults.PublishNotReadyAddresses
	}

	if in.SessionAffinityConfig != nil {
		clone.SessionAffinityConfig = in.SessionAffinityConfig
	} else {
		clone.SessionAffinityConfig = defaults.SessionAffinityConfig
	}

	if in.ClusterIP != nil {
		clone.ClusterIP = in.ClusterIP
	} else {
		clone.ClusterIP = defaults.ClusterIP
	}

	if in.IPFamily != nil {
		clone.IPFamily = in.IPFamily
	} else {
		clone.IPFamily = defaults.IPFamily
	}

	clone.ExternalIPs = MergeStringSlice(in.ExternalIPs, defaults.ExternalIPs)

	return &clone
}

// Create the service spec for the port.
func (in *ServiceSpec) createServiceSpec() corev1.ServiceSpec {
	spec := corev1.ServiceSpec{}
	if in != nil {
		if in.Type != nil {
			spec.Type = *in.Type
		}
		if in.LoadBalancerIP != nil {
			spec.LoadBalancerIP = *in.LoadBalancerIP
		}
		if in.SessionAffinity != nil {
			spec.SessionAffinity = *in.SessionAffinity
		}
		spec.LoadBalancerSourceRanges = in.LoadBalancerSourceRanges
		if in.ExternalName != nil {
			spec.ExternalName = *in.ExternalName
		}
		if in.ExternalTrafficPolicy != nil {
			spec.ExternalTrafficPolicy = *in.ExternalTrafficPolicy
		}
		if in.HealthCheckNodePort != nil {
			spec.HealthCheckNodePort = *in.HealthCheckNodePort
		}
		if in.PublishNotReadyAddresses != nil {
			spec.PublishNotReadyAddresses = *in.PublishNotReadyAddresses
		}
		if in.ClusterIP != nil {
			spec.ClusterIP = *in.ClusterIP
		}
		spec.SessionAffinityConfig = in.SessionAffinityConfig
		spec.IPFamily = in.IPFamily
		spec.ExternalIPs = in.ExternalIPs
	}
	return spec
}

// ----- ScalingSpec -----------------------------------------------------

// The configuration to control safe scaling.
type ScalingSpec struct {
	// ScalingPolicy describes how the replicas of the cluster role will be scaled.
	// The default if not specified is based upon the value of the StorageEnabled field.
	// If StorageEnabled field is not specified or is true the default scaling will be safe, if StorageEnabled is
	// set to false the default scaling will be parallel.
	// +optional
	Policy *ScalingPolicy `json:"policy,omitempty"`
	// The probe to use to determine whether a role is Status HA.
	// If not set the default handler will be used.
	// In most use-cases the default handler would suffice but in
	// advanced use-cases where the application code has a different
	// concept of Status HA to just checking Coherence services then
	// a different handler may be specified.
	// +optional
	Probe *ScalingProbe `json:"probe,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this ScalingSpec struct with any nil or not set values set
// by the corresponding value in the defaults ScalingSpec struct.
func (in *ScalingSpec) DeepCopyWithDefaults(defaults *ScalingSpec) *ScalingSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ScalingSpec{}
	clone.Probe = in.Probe.DeepCopyWithDefaults(defaults.Probe)

	if in.Policy != nil {
		clone.Policy = in.Policy
	} else {
		clone.Policy = defaults.Policy
	}

	return &clone
}

// ----- ScalingProbe ----------------------------------------------------

// ScalingProbe is the handler that will be used to determine how to check for StatusHA in a CoherenceRole.
// StatusHA checking is primarily used during scaling of a role, a role must be in a safe Status HA state
// before scaling takes place. If StatusHA handler is disabled for a role (by specifically setting Enabled
// to false then no check will take place and a role will be assumed to be safe).
// +k8s:openapi-gen=true
type ScalingProbe struct {
	corev1.Handler `json:",inline"`
	// Number of seconds after which the handler times out (only applies to http and tcp handlers).
	// Defaults to 1 second. Minimum value is 1.
	// +optional
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// Returns the timeout value in seconds.
func (in *ScalingProbe) GetTimeout() time.Duration {
	if in == nil || in.TimeoutSeconds == nil || *in.TimeoutSeconds <= 0 {
		return time.Second
	}

	return time.Second * time.Duration(*in.TimeoutSeconds)
}

// DeepCopyWithDefaults returns a copy of this ReadinessProbeSpec struct with any nil or not set values set
// by the corresponding value in the defaults ReadinessProbeSpec struct.
func (in *ScalingProbe) DeepCopyWithDefaults(defaults *ScalingProbe) *ScalingProbe {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ScalingProbe{}

	if in.TimeoutSeconds != nil {
		clone.TimeoutSeconds = in.TimeoutSeconds
	} else {
		clone.TimeoutSeconds = defaults.TimeoutSeconds
	}

	if in.Handler.HTTPGet != nil {
		clone.Handler.HTTPGet = in.Handler.HTTPGet
	} else {
		clone.Handler.HTTPGet = defaults.Handler.HTTPGet
	}

	if in.Handler.TCPSocket != nil {
		clone.Handler.TCPSocket = in.Handler.TCPSocket
	} else {
		clone.Handler.TCPSocket = defaults.Handler.TCPSocket
	}

	if in.Handler.Exec != nil {
		clone.Handler.Exec = in.Handler.Exec
	} else {
		clone.Handler.Exec = defaults.Handler.Exec
	}

	return &clone
}

// ----- ReadinessProbeSpec struct ------------------------------------------

// ReadinessProbeSpec defines the settings for the Coherence Pod readiness probe
// +k8s:openapi-gen=true
type ReadinessProbeSpec struct {
	// The action taken to determine the health of a container
	ProbeHandler `json:",inline"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// +optional
	PeriodSeconds *int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// +optional
	SuccessThreshold *int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// +optional
	FailureThreshold *int32 `json:"failureThreshold,omitempty"`
}

type ProbeHandler struct {
	// One and only one of the following should be specified.
	// Exec specifies the action to take.
	// +optional
	Exec *corev1.ExecAction `json:"exec,omitempty"`
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *corev1.HTTPGetAction `json:"httpGet,omitempty"`
	// TCPSocket specifies an action involving a TCP port.
	// TCP hooks not yet supported
	// +optional
	TCPSocket *corev1.TCPSocketAction `json:"tcpSocket,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this ReadinessProbeSpec struct with any nil or not set values set
// by the corresponding value in the defaults ReadinessProbeSpec struct.
func (in *ReadinessProbeSpec) DeepCopyWithDefaults(defaults *ReadinessProbeSpec) *ReadinessProbeSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := ReadinessProbeSpec{}

	if in.InitialDelaySeconds != nil {
		clone.InitialDelaySeconds = in.InitialDelaySeconds
	} else {
		clone.InitialDelaySeconds = defaults.InitialDelaySeconds
	}

	if in.TimeoutSeconds != nil {
		clone.TimeoutSeconds = in.TimeoutSeconds
	} else {
		clone.TimeoutSeconds = defaults.TimeoutSeconds
	}

	if in.PeriodSeconds != nil {
		clone.PeriodSeconds = in.PeriodSeconds
	} else {
		clone.PeriodSeconds = defaults.PeriodSeconds
	}

	if in.SuccessThreshold != nil {
		clone.SuccessThreshold = in.SuccessThreshold
	} else {
		clone.SuccessThreshold = defaults.SuccessThreshold
	}

	if in.FailureThreshold != nil {
		clone.FailureThreshold = in.FailureThreshold
	} else {
		clone.FailureThreshold = defaults.FailureThreshold
	}

	return &clone
}

// Update the specified probe spec with the required configuration
func (in *ReadinessProbeSpec) UpdateProbeSpec(port int32, path string, probe *corev1.Probe) {
	switch {
	case in != nil && in.Exec != nil:
		probe.Exec = in.Exec
	case in != nil && in.HTTPGet != nil:
		probe.HTTPGet = in.HTTPGet
	case in != nil && in.TCPSocket != nil:
		probe.TCPSocket = in.TCPSocket
	default:
		probe.HTTPGet = &corev1.HTTPGetAction{
			Path:   path,
			Port:   intstr.FromInt(int(port)),
			Scheme: corev1.URISchemeHTTP,
		}
	}

	if in != nil {
		if in.InitialDelaySeconds != nil {
			probe.InitialDelaySeconds = *in.InitialDelaySeconds
		}
		if in.PeriodSeconds != nil {
			probe.PeriodSeconds = *in.PeriodSeconds
		}
		if in.FailureThreshold != nil {
			probe.FailureThreshold = *in.FailureThreshold
		}
		if in.SuccessThreshold != nil {
			probe.SuccessThreshold = *in.SuccessThreshold
		}
		if in.TimeoutSeconds != nil {
			probe.TimeoutSeconds = *in.TimeoutSeconds
		}
	}
}

// ----- FluentdSpec struct -------------------------------------------------

// FluentdSpec defines the settings for the fluentd image
// +k8s:openapi-gen=true
type FluentdSpec struct {
	ImageSpec `json:",inline"`
	// Controls whether or not log capture via a Fluentd sidecar container to an EFK stack is enabled.
	// If this flag i set to true it is expected that the coherence-monitoring-config secret exists in
	// the namespace that the cluster is being deployed to. This secret is either created by the
	// Coherence Operator Helm chart if it was installed with the correct parameters or it should
	// have already been created manually.
	Enabled *bool `json:"enabled,omitempty"`
	// The Fluentd configuration file configuring source for application log.
	// +optional
	ConfigFile *string `json:"configFile,omitempty"`
	// This value should be source.tag from fluentd.application.configFile.
	// +optional
	Tag *string `json:"tag,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this FluentdSpec struct with any nil or not set values set
// by the corresponding value in the defaults FluentdSpec struct.
func (in *FluentdSpec) DeepCopyWithDefaults(defaults *FluentdSpec) *FluentdSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := FluentdSpec{}
	clone.ImageSpec = *in.ImageSpec.DeepCopyWithDefaults(&defaults.ImageSpec)

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.ConfigFile != nil {
		clone.ConfigFile = in.ConfigFile
	} else {
		clone.ConfigFile = defaults.ConfigFile
	}

	if in.Tag != nil {
		clone.Tag = in.Tag
	} else {
		clone.Tag = defaults.Tag
	}

	return &clone
}

func (in *FluentdSpec) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	if in == nil || in.Enabled == nil || !*in.Enabled {
		// either the fluentd spec is nil or disabled
		return
	}

	// add the fluentd container
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, in.CreateFluentdContainer())

	// add the fluentd configuration ConfigMap volume
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: VolumeNameFluentdConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf(EfkConfigMapNameTemplate, sts.Name)},
				DefaultMode:          pointer.Int32Ptr(420),
			},
		},
	})
}

func (in *FluentdSpec) CreateFluentdContainer() corev1.Container {
	var pullPolicy corev1.PullPolicy
	if in.ImagePullPolicy == nil {
		pullPolicy = corev1.PullIfNotPresent
	} else {
		pullPolicy = *in.ImagePullPolicy
	}

	var imageName string
	if in.Image == nil {
		imageName = DefaultFluentdImage
	} else {
		imageName = *in.Image
	}

	return corev1.Container{
		Name:            ContainerNameFluentd,
		Image:           imageName,
		ImagePullPolicy: pullPolicy,
		Args:            []string{"-c", "/etc/fluent.conf"},
		Env: []corev1.EnvVar{
			{
				Name: "COHERENCE_POD_ID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name:  "FLUENTD_CONF",
				Value: "fluentd-coherence.conf",
			},
			{
				Name:  "FLUENT_ELASTICSEARCH_SED_DISABLE",
				Value: "true",
			},
			{
				Name: "ELASTICSEARCH_HOST",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: CoherenceMonitoringConfigName,
						},
						Key: LoggingConfigKeyElasticSearchHost,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_PORT",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: CoherenceMonitoringConfigName,
						},
						Key: LoggingConfigElasticSearchPort,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: CoherenceMonitoringConfigName,
						},
						Key: LoggingConfigElasticSearchUser,
					},
				},
			},
			{
				Name: "ELASTICSEARCH_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: CoherenceMonitoringConfigName,
						},
						Key: LoggingConfigElasticSearchPassword,
					},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      VolumeNameFluentdConfig,
				MountPath: VolumeMountPathFluentdConfig,
				SubPath:   VolumeMountSubPathFluentdConfig,
			},
			{
				Name:      VolumeNameLogs,
				MountPath: VolumeMountPathLogs,
			},
		},
	}
}

// ----- ScalingPolicy type -------------------------------------------------

// ScalingPolicy describes a policy for scaling a cluster role
type ScalingPolicy string

// Scaling policy constants
const (
	// Safe means that a role will be scaled up or down in a safe manner to ensure no data loss.
	SafeScaling ScalingPolicy = "Safe"
	// Parallel means that a role will be scaled up or down by adding or removing members in parallel.
	// If the members of the role are storage enabled then this could cause data loss
	ParallelScaling ScalingPolicy = "Parallel"
	// ParallelUpSafeDown means that a role will be scaled up by adding or removing members in parallel
	// but will be scaled down in a safe manner to ensure no data loss.
	ParallelUpSafeDownScaling ScalingPolicy = "ParallelUpSafeDown"
)

// ----- LocalObjectReference -----------------------------------------------

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
}

// ----- NetworkSpec --------------------------------------------------------

// NetworkSpec configures various networking and DNS settings for Pods in a role.
// +k8s:openapi-gen=true
type NetworkSpec struct {
	// Specifies the DNS parameters of a pod. Parameters specified here will be merged to the
	// generated DNS configuration based on DNSPolicy.
	// +optional
	DNSConfig *PodDNSConfig `json:"dnsConfig,omitempty"`
	// Set DNS policy for the pod. Defaults to "ClusterFirst". Valid values are 'ClusterFirstWithHostNet',
	// 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy
	// selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS
	// policy explicitly to 'ClusterFirstWithHostNet'.
	// +optional
	DNSPolicy *corev1.DNSPolicy `json:"dnsPolicy,omitempty"`
	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified.
	// This is only valid for non-hostNetwork pods.
	// +listType=map
	// +listMapKey=ip
	// +optional
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`
	// Host networking requested for this pod. Use the host's network namespace. If this option is set,
	// the ports that will be used must be specified. Default to false.
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty"`
	// Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.
	// +optional
	Hostname *string `json:"hostname,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this NetworkSpec struct with any nil or not set values set
// by the corresponding value in the defaults NetworkSpec struct.
func (in *NetworkSpec) DeepCopyWithDefaults(defaults *NetworkSpec) *NetworkSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := NetworkSpec{}
	clone.DNSConfig = in.DNSConfig.DeepCopyWithDefaults(defaults.DNSConfig)

	if in.DNSPolicy != nil {
		clone.DNSPolicy = in.DNSPolicy
	} else {
		clone.DNSPolicy = defaults.DNSPolicy
	}

	// merge HostAlias list
	m := make(map[string]corev1.HostAlias)
	if defaults.HostAliases != nil {
		for _, h := range defaults.HostAliases {
			m[h.IP] = h
		}
	}
	if in.HostAliases != nil {
		for _, h := range in.HostAliases {
			m[h.IP] = h
		}
	}
	if len(m) > 0 {
		i := 0
		clone.HostAliases = make([]corev1.HostAlias, len(m))
		for _, h := range m {
			clone.HostAliases[i] = h
			i++
		}
	}

	if in.HostNetwork != nil {
		clone.HostNetwork = in.HostNetwork
	} else {
		clone.HostNetwork = defaults.HostNetwork
	}

	if in.Hostname != nil {
		clone.Hostname = in.Hostname
	} else {
		clone.Hostname = defaults.Hostname
	}

	return &clone
}

// Update the specified StatefulSet's network settings.
func (in *NetworkSpec) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	in.DNSConfig.UpdateStatefulSet(sts)

	if in.DNSPolicy != nil {
		sts.Spec.Template.Spec.DNSPolicy = *in.DNSPolicy
	}

	sts.Spec.Template.Spec.HostAliases = in.HostAliases
	sts.Spec.Template.Spec.HostNetwork = notNilBool(in.HostNetwork)
	sts.Spec.Template.Spec.Hostname = notNilString(in.Hostname)
}

// ----- PodDNSConfig -------------------------------------------------------

// PodDNSConfig defines the DNS parameters of a pod in addition to
// those generated from DNSPolicy.
// +k8s:openapi-gen=true
type PodDNSConfig struct {
	// A list of DNS name server IP addresses.
	// This will be appended to the base nameservers generated from DNSPolicy.
	// Duplicated nameservers will be removed.
	// +listType=atomic
	// +optional
	Nameservers []string `json:"nameservers,omitempty"`
	// A list of DNS search domains for host-name lookup.
	// This will be appended to the base search paths generated from DNSPolicy.
	// Duplicated search paths will be removed.
	// +listType=atomic
	// +optional
	Searches []string `json:"searches,omitempty"`
	// A list of DNS resolver options.
	// This will be merged with the base options generated from DNSPolicy.
	// Duplicated entries will be removed. Resolution options given in Options
	// will override those that appear in the base DNSPolicy.
	// +listType=map
	// +listMapKey=name
	// +optional
	Options []corev1.PodDNSConfigOption `json:"options,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this PodDNSConfig struct with any nil or not set values set
// by the corresponding value in the defaults PodDNSConfig struct.
func (in *PodDNSConfig) DeepCopyWithDefaults(defaults *PodDNSConfig) *PodDNSConfig {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := PodDNSConfig{}

	// merge Options list
	m := make(map[string]corev1.PodDNSConfigOption)
	if defaults.Options != nil {
		for _, opt := range defaults.Options {
			m[opt.Name] = opt
		}
	}
	if in.Options != nil {
		for _, opt := range in.Options {
			m[opt.Name] = opt
		}
	}
	if len(m) > 0 {
		i := 0
		clone.Options = make([]corev1.PodDNSConfigOption, len(m))
		for _, opt := range m {
			clone.Options[i] = opt
			i++
		}
	}

	if in.Nameservers != nil {
		clone.Nameservers = []string{}
		clone.Nameservers = append(clone.Nameservers, defaults.Nameservers...)
		clone.Nameservers = append(clone.Nameservers, in.Nameservers...)
	} else if defaults.Nameservers != nil {
		clone.Nameservers = []string{}
		clone.Nameservers = append(clone.Nameservers, defaults.Nameservers...)
	}

	if in.Searches != nil {
		clone.Searches = []string{}
		clone.Searches = append(clone.Searches, defaults.Searches...)
		clone.Searches = append(clone.Searches, in.Searches...)
	} else if defaults.Searches != nil {
		clone.Searches = []string{}
		clone.Searches = append(clone.Searches, defaults.Searches...)
	}

	return &clone
}

func (in *PodDNSConfig) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	cfg := corev1.PodDNSConfig{}

	if in.Nameservers != nil && len(in.Nameservers) > 0 {
		cfg.Nameservers = in.Nameservers
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}

	if in.Searches != nil && len(in.Searches) > 0 {
		cfg.Searches = in.Searches
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}

	if in.Options != nil && len(in.Options) > 0 {
		cfg.Options = in.Options
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}
}

// ----- StartQuorum --------------------------------------------------------

// StartQuorum defines the order that roles will be created when initially
// creating a new cluster.
// +k8s:openapi-gen=true
type StartQuorum struct {
	// The list of roles to start first.
	// +optional
	Role string `json:"role"`
	// The number of the dependency Pods that should have been started
	// before this roles will be started.
	// +optional
	PodCount int32 `json:"podCount,omitempty"`
}

// ----- StartStatus --------------------------------------------------------

// StartQuorumStatus tracks the state of a role's start quorums.
type StartQuorumStatus struct {
	// The inlined start quorum.
	StartQuorum `json:",inline"`
	// Whether this quorum's condition has been met
	Ready bool `json:"ready"`
}

func MergeStringSlice(s1, s2 []string) []string {
	m := make(map[string]int)
	for _, eip := range s2 {
		m[eip] = 0
	}

	for _, eip := range s1 {
		m[eip] = 0
	}

	var merged []string
	for k := range m {
		merged = append(merged, k)
	}
	return merged
}

// Convert an int32 pointer to a string using the default if the pointer is nil.
func Int32PtrToStringWithDefault(i *int32, d int32) string {
	if i == nil {
		return Int32ToString(d)
	}
	return Int32ToString(*i)
}

// Convert an int32 pointer to a string.
func Int32PtrToString(i *int32) string {
	return Int32ToString(*i)
}

// Convert an int32 to a string.
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// Convert a bool pointer to a string.
func BoolPtrToString(b *bool) string {
	if b != nil && *b {
		return "true"
	}
	return "false"
}
