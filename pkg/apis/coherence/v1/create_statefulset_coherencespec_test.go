/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetFromRoleWithCoherenceSpecEmpty(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithImage(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			ImageSpec: coh.ImageSpec{
				Image: stringPtr("coherence:1.0"),
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].Image = "coherence:1.0"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			ImageSpec: coh.ImageSpec{
				ImagePullPolicy: &policy,
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ImagePullPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithStorageEnabledTrue(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(true),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_STORAGE_ENABLED", Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithStorageEnabledFalse(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(false),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_STORAGE_ENABLED", Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithCacheConfig(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			CacheConfig: stringPtr("test-config.xml"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_CACHE_CONFIG", Value: "test-config.xml"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithOverrideConfig(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			OverrideConfig: stringPtr("test-override.xml"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_OVERRIDE_CONFIG", Value: "test-override.xml"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithLogLevel(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			LogLevel: int32Ptr(9),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_LOG_LEVEL", Value: "9"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithExcludeFromWKATrue(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(true),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithCoherenceSpecWithExcludeFromWKAFalse(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(false),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "true"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
