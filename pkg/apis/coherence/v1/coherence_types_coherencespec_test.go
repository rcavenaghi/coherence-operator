/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"testing"
)

// Tests for the CoherenceSpec struct

func TestCoherenceSpec_IsPersistenceEnabled_WhenNil(t *testing.T) {
	g := NewGomegaWithT(t)
	var c coh.CoherenceSpec
	g.Expect(c.IsPersistenceEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsPersistenceEnabled_WhenPersistenceIsNil(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{}
	g.Expect(c.IsPersistenceEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsPersistenceEnabled_WhenPersistenceEnabledIsNil(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Persistence: &coh.PersistentStorageSpec{
			Enabled: nil,
		},
	}
	g.Expect(c.IsPersistenceEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsPersistenceEnabled_WhenPersistenceEnabledIsFalse(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Persistence: &coh.PersistentStorageSpec{
			Enabled: boolPtr(false),
		},
	}
	g.Expect(c.IsPersistenceEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsPersistenceEnabled_WhenPersistenceEnabledIsTrue(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Persistence: &coh.PersistentStorageSpec{
			Enabled: boolPtr(true),
		},
	}
	g.Expect(c.IsPersistenceEnabled()).To(BeTrue())
}

func TestCoherenceSpec_IsSnapshotsEnabled_WhenNil(t *testing.T) {
	g := NewGomegaWithT(t)
	var c coh.CoherenceSpec
	g.Expect(c.IsSnapshotsEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsSnapshotsEnabled_WhenPersistenceIsNil(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{}
	g.Expect(c.IsSnapshotsEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsSnapshotsEnabled_WhenPersistenceEnabledIsNil(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Snapshot: &coh.PersistentStorageSpec{
			Enabled: nil,
		},
	}
	g.Expect(c.IsSnapshotsEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsSnapshotsEnabled_WhenPersistenceEnabledIsFalse(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Snapshot: &coh.PersistentStorageSpec{
			Enabled: boolPtr(false),
		},
	}
	g.Expect(c.IsSnapshotsEnabled()).To(BeFalse())
}

func TestCoherenceSpec_IsSnapshotsEnabled_WhenPersistenceEnabledIsTrue(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceSpec{
		Snapshot: &coh.PersistentStorageSpec{
			Enabled: boolPtr(true),
		},
	}
	g.Expect(c.IsSnapshotsEnabled()).To(BeTrue())
}
