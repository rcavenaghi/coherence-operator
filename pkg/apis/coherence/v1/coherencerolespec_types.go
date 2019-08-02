package v1

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
	// Details of the Docker images used in the role
	// +optional
	Images *Images `json:"images,omitempty"`
	// A boolean flag indicating whether members of this role are storage enabled.
	// This value will set the corresponding coherence.distributed.localstorage System property.
	// If not specified the default value is true.
	// This flag is also used to configure the ScalingPolicy value if a value is not specified. If the
	// StorageEnabled field is not specified or is true the scaling will be safe, if StorageEnabled is
	// set to false scaling will be parallel.
	// +optional
	StorageEnabled *bool `json:"storageEnabled,omitempty"`
	// ScalingPolicy describes how the replicas of the cluster role will be scaled.
	// The default if not specified is based upon the value of the StorageEnabled field.
	// If StorageEnabled field is not specified or is true the default scaling will be safe, if StorageEnabled is
	// set to false the default scaling will be parallel.
	// +optional
	ScalingPolicy *ScalingPolicy `json:"scalingPolicy,omitempty"`
	// The readiness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
	// The extra labels to add to the all of the Pods in this roles.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels *map[string]string `json:"labels,omitempty"`
	// CacheConfig is the name of the cache configuration file to use
	CacheConfig *string `json:"cacheConfig,omitempty"`
	// PofConfig is the name of the POF configuration file to use when using POF serializer
	PofConfig *string `json:"pofConfig,omitempty"`
	// OverrideConfig is name of the Coherence operational configuration override file,
	// the default is tangosol-coherence-override.xml
	OverrideConfig *string `json:"overrideConfig,omitEmpty"`
	// MaxHeap is the min/max heap value to pass to the JVM.
	// The format should be the same as that used for Java's -Xms and -Xmx JVM options.
	// If not set the JVM defaults are used.
	MaxHeap *string `json:"maxHeap,omitEmpty"`
	// JvmArgs specifies the options to pass to the Coherence JVM. The default is
	// to use the G1 collector.
	JvmArgs *string `json:"jvmArgs,omitEmpty"`
	// JavaOpts is miscellaneous JVM options to pass to the Coherence store container
	// This options will override the system options computed in the start up script.
	JavaOpts *string `json:"javaOpts,omitEmpty"`
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

// Obtain the full name for  a role.
func (in *CoherenceRoleSpec) GetFullRoleName(cluster *CoherenceCluster) string {
	return cluster.Name + "-" + in.GetRoleName()
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

// DeepCopyWithDefaults returns a copy of this CoherenceRoleSpec with any nil or not set values set
// by the corresponding value in the defaults spec.
func (in *CoherenceRoleSpec) DeepCopyWithDefaults(defaults *CoherenceRoleSpec) *CoherenceRoleSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		} else {
			return nil
		}
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := CoherenceRoleSpec{}

	// Copy EVERY field from "in" to the clone.
	// If a field is not set use the value from the default
	// If the field is a struct it should implement DeepCopyWithDefaults so call that method

	if in.Role != "" {
		clone.Role = in.Role
	} else {
		clone.Role = defaults.Role
	}

	if in.Replicas != nil {
		clone.Replicas = in.Replicas
	} else {
		clone.Replicas = defaults.Replicas
	}

	if in.StorageEnabled != nil {
		clone.StorageEnabled = in.StorageEnabled
	} else {
		clone.StorageEnabled = defaults.StorageEnabled
	}

	if in.ScalingPolicy != nil {
		clone.ScalingPolicy = in.ScalingPolicy
	} else {
		clone.ScalingPolicy = defaults.ScalingPolicy
	}

	if in.CacheConfig != nil {
		clone.CacheConfig = in.CacheConfig
	} else {
		clone.CacheConfig = defaults.CacheConfig
	}

	if in.PofConfig != nil {
		clone.PofConfig = in.PofConfig
	} else {
		clone.PofConfig = defaults.PofConfig
	}

	if in.OverrideConfig != nil {
		clone.OverrideConfig = in.OverrideConfig
	} else {
		clone.OverrideConfig = defaults.OverrideConfig
	}

	if in.MaxHeap != nil {
		clone.MaxHeap = in.MaxHeap
	} else {
		clone.MaxHeap = defaults.MaxHeap
	}

	if in.JvmArgs != nil {
		clone.JvmArgs = in.JvmArgs
	} else {
		clone.JvmArgs = defaults.JvmArgs
	}

	if in.JavaOpts != nil {
		clone.JavaOpts = in.JavaOpts
	} else {
		clone.JavaOpts = defaults.JavaOpts
	}

	clone.Labels = in.mergeMap(in.Labels, defaults.Labels)

	clone.Images = in.Images.DeepCopyWithDefaults(defaults.Images)
	clone.ReadinessProbe = in.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)

	return &clone
}

// Return a map that is two maps merged.
// If both maps are nil then nil is returned.
// Where there are duplicate keys those in m1 take precedence.
// Keys that map to "" will not be added to the merged result
func (in *CoherenceRoleSpec) mergeMap(m1, m2 *map[string]string) *map[string]string {
	if m1 == nil && m2 == nil {
		return nil
	}

	merged := make(map[string]string)

	if m2 != nil {
		for k, v := range *m2 {
			if v != "" {
				merged[k] = v
			}
		}
	}

	if m1 != nil {
		for k, v := range *m1 {
			if v != "" {
				merged[k] = v
			} else {
				delete(merged, k)
			}
		}
	}

	return &merged
}