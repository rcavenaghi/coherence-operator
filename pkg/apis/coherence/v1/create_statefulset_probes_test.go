/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestCreateStatefulSetFromRoleWithEmptyReadinessProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{}

	role := coh.CoherenceRoleSpec{
		ReadinessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create expected probe
	probeExpected := role.UpdateDefaultReadinessProbeAction(role.CreateDefaultReadinessProbe())
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = probeExpected

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithReadinessProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		ReadinessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: nil,
			HTTPGet: &corev1.HTTPGetAction{
				Path:   coh.DefaultReadinessPath,
				Port:   intstr.FromInt(int(coh.DefaultHealthPort)),
				Scheme: "HTTP",
			},
			TCPSocket: nil,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithReadinessProbeSpecWithHttpGet(t *testing.T) {

	handler := &corev1.HTTPGetAction{
		Path: "/test/ready",
		Port: intstr.FromInt(1234),
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec:      nil,
			HTTPGet:   handler,
			TCPSocket: nil,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		ReadinessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithReadinessProbeSpecWithTCPSocket(t *testing.T) {

	handler := &corev1.TCPSocketAction{
		Port: intstr.FromInt(1234),
		Host: "foo.com",
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		ReadinessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithReadinessProbeSpecWithExec(t *testing.T) {

	handler := &corev1.ExecAction{
		Command: []string{"exec", "something"},
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		ReadinessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithEmptyLivenessProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{}

	role := coh.CoherenceRoleSpec{
		LivenessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)

	// Create expected probe
	probeExpected := role.UpdateDefaultLivenessProbeAction(role.CreateDefaultLivenessProbe())
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = probeExpected

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLivenessProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		LivenessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: nil,
			HTTPGet: &corev1.HTTPGetAction{
				Path:   coh.DefaultLivenessPath,
				Port:   intstr.FromInt(int(coh.DefaultHealthPort)),
				Scheme: "HTTP",
			},
			TCPSocket: nil,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithLivenessProbeSpecWithHttpGet(t *testing.T) {

	handler := &corev1.HTTPGetAction{
		Path: "/test/ready",
		Port: intstr.FromInt(1234),
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec:      nil,
			HTTPGet:   handler,
			TCPSocket: nil,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		LivenessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithLivenessProbeSpecWithTCPSocket(t *testing.T) {

	handler := &corev1.TCPSocketAction{
		Port: intstr.FromInt(1234),
		Host: "foo.com",
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		LivenessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}

func TestCreateStatefulSetFromRoleWithLivenessProbeSpecWithExec(t *testing.T) {

	handler := &corev1.ExecAction{
		Command: []string{"exec", "something"},
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	role := coh.CoherenceRoleSpec{
		LivenessProbe: &probe,
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreationWithHelmVerify(t, role, cluster, stsExpected, false)
}
