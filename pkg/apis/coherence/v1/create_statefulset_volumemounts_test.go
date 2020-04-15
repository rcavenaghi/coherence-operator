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

func TestCreateStatefulSetFromRoleWithEmptyVolumeMounts(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		VolumeMounts: []corev1.VolumeMount{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithOneVolumeMount(t *testing.T) {

	mountOne := corev1.VolumeMount{
		Name:      "volume-one",
		ReadOnly:  true,
		MountPath: "/home/root/one",
	}

	role := coh.CoherenceRoleSpec{
		VolumeMounts: []corev1.VolumeMount{mountOne},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mountOne)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
func TestCreateStatefulSetFromRoleWithTwoVolumeMounts(t *testing.T) {

	mountOne := corev1.VolumeMount{
		Name:      "volume-one",
		ReadOnly:  true,
		MountPath: "/home/root/one",
	}

	mountTwo := corev1.VolumeMount{
		Name:      "volume-two",
		ReadOnly:  true,
		MountPath: "/home/root/two",
	}

	role := coh.CoherenceRoleSpec{
		VolumeMounts: []corev1.VolumeMount{mountOne, mountTwo},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mountOne, mountTwo)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
