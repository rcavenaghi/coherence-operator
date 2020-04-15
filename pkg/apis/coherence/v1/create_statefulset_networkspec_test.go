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

func TestCreateStatefulSetFromRoleWithNetworkSpec(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigNameServers(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Nameservers: []string{"one", "two"},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Nameservers: []string{"one", "two"},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigEmptyNameServers(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Nameservers: []string{},
				Searches:    nil,
				Options:     nil,
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigSearches(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{"one", "two"},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Searches: []string{"one", "two"},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigEmptySearches(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{},
				Options:  nil,
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigOptions(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Options: []corev1.PodDNSConfigOption{
					{
						Name:  "Foo",
						Value: stringPtr("Bar"),
					},
				},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Options: []corev1.PodDNSConfigOption{
			{
				Name:  "Foo",
				Value: stringPtr("Bar"),
			},
		},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSConfigEmptyOptions(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Options: []corev1.PodDNSConfigOption{},
			},
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithDNSPolicy(t *testing.T) {

	policy := corev1.DNSClusterFirst

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			DNSPolicy: &policy,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.DNSPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithHostAliases(t *testing.T) {

	aliases := []corev1.HostAlias{
		{
			IP:        "10.10.10.10",
			Hostnames: []string{"foo.com"},
		},
	}

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			HostAliases: aliases,
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.HostAliases = aliases

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithHostNetworkFalse(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(false),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.HostNetwork = false

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithHostNetworkTrue(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(true),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.HostNetwork = true

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}

func TestCreateStatefulSetFromRoleWithNetworkSpecWithHostname(t *testing.T) {

	role := coh.CoherenceRoleSpec{
		Network: &coh.NetworkSpec{
			Hostname: stringPtr("foo.com"),
		},
	}

	// Create the test cluster
	cluster := createTestCluster(role)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(cluster, role)
	stsExpected.Spec.Template.Spec.Hostname = "foo.com"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, role, cluster, stsExpected)
}
