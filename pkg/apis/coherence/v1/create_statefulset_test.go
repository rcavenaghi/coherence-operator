/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateStatefulSetFromMinimalRoleSpec(t *testing.T) {
	// Create minimal role spec
	role := coh.CoherenceRoleSpec{}
	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithName(t *testing.T) {
	// create a role with a name
	role := coh.CoherenceRoleSpec{
		Role: "data",
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// The sts name should be the full role name
	stsExpected.Name = role.GetFullRoleName(cluster)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithReplicas(t *testing.T) {
	// create a role with a name
	role := coh.CoherenceRoleSpec{
		Replicas: pointer.Int32Ptr(50),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithEnvVars(t *testing.T) {
	// create a role with environment variables
	ev := []corev1.EnvVar{
		{Name: "FOO", Value: "FOO_VAL"},
	}
	role := coh.CoherenceRoleSpec{
		Env: ev,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	// add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, ev...)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithEmptyEnvVars(t *testing.T) {
	// create a role with empty environment variables
	role := coh.CoherenceRoleSpec{
		Env: []corev1.EnvVar{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithHealthPort(t *testing.T) {
	// create a role with a custom health port
	role := coh.CoherenceRoleSpec{
		HealthPort: int32Ptr(210),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Port = intstr.FromInt(210)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe.HTTPGet.Port = intstr.FromInt(210)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLabels(t *testing.T) {
	// create a role with empty environment variables
	labels := make(map[string]string)
	labels["foo"] = "foo-label"
	labels["bar"] = "bar-label"

	role := coh.CoherenceRoleSpec{
		Labels: labels,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	for k, v := range labels {
		stsExpected.Spec.Template.Labels[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithAnnotations(t *testing.T) {
	// create a role with empty environment variables
	annotations := make(map[string]string)
	annotations["foo"] = "foo-annotation"
	annotations["bar"] = "bar-annotation"

	role := coh.CoherenceRoleSpec{
		Annotations: annotations,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	if stsExpected.Spec.Template.Annotations == nil {
		stsExpected.Spec.Template.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		stsExpected.Spec.Template.Annotations[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithResources(t *testing.T) {
	res := corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("8"),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("4"),
		},
	}

	role := coh.CoherenceRoleSpec{
		Resources: &res,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].Resources = res

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithAffinity(t *testing.T) {
	// Create a test affinity spec
	sel := metav1.LabelSelector{
		MatchLabels: map[string]string{"Foo": "Bar"},
	}
	affinity := corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &sel,
				},
			},
		},
	}

	// Create the role with the affinity spec
	role := coh.CoherenceRoleSpec{
		Affinity: &affinity,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Affinity = &affinity

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNodeSelector(t *testing.T) {
	selector := make(map[string]string)
	selector["foo"] = "foo-label"
	selector["bar"] = "bar-label"

	// Create the role with the node selector
	role := coh.CoherenceRoleSpec{
		NodeSelector: selector,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.NodeSelector = selector

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithTolerations(t *testing.T) {
	tolerations := []corev1.Toleration{
		{
			Key:      "Foo",
			Operator: corev1.TolerationOpEqual,
			Value:    "Bar",
		},
	}

	role := coh.CoherenceRoleSpec{
		Tolerations: tolerations,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Tolerations = tolerations

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithSecurityContext(t *testing.T) {
	ctx := corev1.PodSecurityContext{
		RunAsUser:    pointer.Int64Ptr(1000),
		RunAsNonRoot: boolPtr(true),
	}

	role := coh.CoherenceRoleSpec{
		SecurityContext: &ctx,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.SecurityContext = &ctx

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithShareProcessNamespaceFalse(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		ShareProcessNamespace: boolPtr(false),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.ShareProcessNamespace = boolPtr(false)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithShareProcessNamespaceTrue(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		ShareProcessNamespace: boolPtr(true),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.ShareProcessNamespace = boolPtr(true)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithHostIPCFalse(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		HostIPC: boolPtr(false),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.HostIPC = false

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithHostIPCNamespaceTrue(t *testing.T) {
	role := coh.CoherenceRoleSpec{
		HostIPC: boolPtr(true),
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.HostIPC = true

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
