/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

var _ = Describe("Testing ApplicationSpec struct", func() {

	Context("Copying an ApplicationSpec using DeepCopyWithDefaults", func() {
		var original *coh.ApplicationSpec
		var defaults *coh.ApplicationSpec
		var clone *coh.ApplicationSpec

		var always = corev1.PullAlways
		var never = corev1.PullNever

		var appOne = &coh.ApplicationSpec{
			Type: stringPtr("java"),
			Main: stringPtr("TestMainOne"),
			Args: []string{},
			ImageSpec: coh.ImageSpec{
				Image:           stringPtr("app:1.0"),
				ImagePullPolicy: &always,
			},
			LibDir:    stringPtr("/test/libOne"),
			ConfigDir: stringPtr("/test/confOne"),
		}

		var appTwo = &coh.ApplicationSpec{
			Type: stringPtr("node"),
			Main: stringPtr("TestMainTwo"),
			Args: []string{},
			ImageSpec: coh.ImageSpec{
				Image:           stringPtr("app:2.0"),
				ImagePullPolicy: &never,
			},
			LibDir:    stringPtr("/test/libTwo"),
			ConfigDir: stringPtr("/test/confTwo"),
		}

		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = appOne
				defaults = nil
			})

			It("should copy the original Args", func() {
				Expect(clone.Args).To(Equal(original.Args))
			})

			It("should copy the original Type", func() {
				Expect(*clone.Type).To(Equal(*original.Type))
			})

			It("should copy the original Type", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = appOne
				defaults = &coh.ApplicationSpec{}
			})

			It("should copy the original Args", func() {
				Expect(clone.Args).To(Equal(original.Args))
			})

			It("should copy the original Type", func() {
				Expect(*clone.Type).To(Equal(*original.Type))
			})

			It("should copy the original Type", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = appTwo
			})

			It("should copy the defaults Args", func() {
				Expect(clone.Args).To(Equal(defaults.Args))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.Type).To(Equal(*defaults.Type))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coh.ApplicationSpec{}
				defaults = appTwo
			})

			It("should copy the defaults Args", func() {
				Expect(clone.Args).To(Equal(defaults.Args))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.Type).To(Equal(*defaults.Type))
			})

			It("should copy the defaults LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the defaults ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})

		When("original Args is nil", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = nil
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"one", "two"}
			})

			It("should copy the defaults Args", func() {
				expected := original.DeepCopy()
				expected.Args = defaults.Args
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Args is empty", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"one", "two"}
			})

			It("should copy the original Args", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults Args is nil", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = nil
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("defaults Args is empty", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{}
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("original and defaults Args is populated", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"three", "four"}
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})
	})
})

func TestApplicationSpec_CreateApplicationContainer_NilApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	var app coh.ApplicationSpec
	ok, _ := app.CreateApplicationContainer()
	g.Expect(ok).To(BeFalse())
}

func TestApplicationSpec_CreateApplicationContainer_NoImage(t *testing.T) {
	g := NewGomegaWithT(t)

	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image: nil,
		},
	}

	ok, _ := app.CreateApplicationContainer()
	g.Expect(ok).To(BeFalse())
}

func TestApplicationSpec_CreateApplicationContainer_Default(t *testing.T) {
	g := NewGomegaWithT(t)

	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image: stringPtr("my-test-image:1.0"),
		},
	}

	ok, c := app.CreateApplicationContainer()
	g.Expect(ok).To(BeTrue())
	g.Expect(c.Name).To(Equal(coh.ContainerNameApplication))
	g.Expect(c.Image).To(Equal("my-test-image:1.0"))
	g.Expect(c.Command).To(Equal([]string{coh.DefaultCommandApplication}))
	g.Expect(c.ImagePullPolicy).To(BeEmpty())

	env := []corev1.EnvVar{
		{Name: "EXTERNAL_APP_DIR", Value: coh.ExternalAppDir},
		{Name: "APP_DIR", Value: coh.AppDir},
		{Name: "EXTERNAL_LIB_DIR", Value: coh.ExternalLibDir},
		{Name: "LIB_DIR", Value: coh.LibDir},
		{Name: "EXTERNAL_CONF_DIR", Value: coh.ExternalConfDir},
		{Name: "CONF_DIR", Value: coh.ConfDir},
	}
	g.Expect(c.Env).To(Equal(env))

	mounts := []corev1.VolumeMount{
		{Name: coh.VolumeNameUtils, MountPath: coh.VolumeMountPathUtils},
		{Name: coh.VolumeNameApplication, MountPath: coh.ExternalAppDir},
	}
	g.Expect(c.VolumeMounts).To(Equal(mounts))
}

func TestApplicationSpec_CreateApplicationContainer_WithPullPolicy(t *testing.T) {
	g := NewGomegaWithT(t)

	pullPolicy := corev1.PullAlways
	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image:           stringPtr("my-test-image:1.0"),
			ImagePullPolicy: &pullPolicy,
		},
	}

	ok, c := app.CreateApplicationContainer()
	g.Expect(ok).To(BeTrue())
	g.Expect(c.Name).To(Equal(coh.ContainerNameApplication))
	g.Expect(c.Image).To(Equal("my-test-image:1.0"))
	g.Expect(c.Command).To(Equal([]string{coh.DefaultCommandApplication}))
	g.Expect(c.ImagePullPolicy).To(Equal(pullPolicy))

	env := []corev1.EnvVar{
		{Name: "EXTERNAL_APP_DIR", Value: coh.ExternalAppDir},
		{Name: "APP_DIR", Value: coh.AppDir},
		{Name: "EXTERNAL_LIB_DIR", Value: coh.ExternalLibDir},
		{Name: "LIB_DIR", Value: coh.LibDir},
		{Name: "EXTERNAL_CONF_DIR", Value: coh.ExternalConfDir},
		{Name: "CONF_DIR", Value: coh.ConfDir},
	}
	g.Expect(c.Env).To(Equal(env))

	mounts := []corev1.VolumeMount{
		{Name: coh.VolumeNameUtils, MountPath: coh.VolumeMountPathUtils},
		{Name: coh.VolumeNameApplication, MountPath: coh.ExternalAppDir},
	}
	g.Expect(c.VolumeMounts).To(Equal(mounts))
}

func TestApplicationSpec_CreateApplicationContainer_WithAppDir(t *testing.T) {
	g := NewGomegaWithT(t)

	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image: stringPtr("my-test-image:1.0"),
		},
		AppDir: stringPtr("/home/test/app"),
	}

	ok, c := app.CreateApplicationContainer()
	g.Expect(ok).To(BeTrue())
	g.Expect(c.Name).To(Equal(coh.ContainerNameApplication))
	g.Expect(c.Image).To(Equal("my-test-image:1.0"))
	g.Expect(c.Command).To(Equal([]string{coh.DefaultCommandApplication}))
	g.Expect(c.ImagePullPolicy).To(BeEmpty())

	env := []corev1.EnvVar{
		{Name: "EXTERNAL_APP_DIR", Value: coh.ExternalAppDir},
		{Name: "APP_DIR", Value: "/home/test/app"},
		{Name: "EXTERNAL_LIB_DIR", Value: coh.ExternalLibDir},
		{Name: "LIB_DIR", Value: coh.LibDir},
		{Name: "EXTERNAL_CONF_DIR", Value: coh.ExternalConfDir},
		{Name: "CONF_DIR", Value: coh.ConfDir},
	}
	g.Expect(c.Env).To(Equal(env))

	mounts := []corev1.VolumeMount{
		{Name: coh.VolumeNameUtils, MountPath: coh.VolumeMountPathUtils},
		{Name: coh.VolumeNameApplication, MountPath: coh.ExternalAppDir},
	}
	g.Expect(c.VolumeMounts).To(Equal(mounts))
}

func TestApplicationSpec_CreateApplicationContainer_WithLibDir(t *testing.T) {
	g := NewGomegaWithT(t)

	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image: stringPtr("my-test-image:1.0"),
		},
		LibDir: stringPtr("/home/test/lib"),
	}

	ok, c := app.CreateApplicationContainer()
	g.Expect(ok).To(BeTrue())
	g.Expect(c.Name).To(Equal(coh.ContainerNameApplication))
	g.Expect(c.Image).To(Equal("my-test-image:1.0"))
	g.Expect(c.Command).To(Equal([]string{coh.DefaultCommandApplication}))
	g.Expect(c.ImagePullPolicy).To(BeEmpty())

	env := []corev1.EnvVar{
		{Name: "EXTERNAL_APP_DIR", Value: coh.ExternalAppDir},
		{Name: "APP_DIR", Value: coh.AppDir},
		{Name: "EXTERNAL_LIB_DIR", Value: coh.ExternalLibDir},
		{Name: "LIB_DIR", Value: "/home/test/lib"},
		{Name: "EXTERNAL_CONF_DIR", Value: coh.ExternalConfDir},
		{Name: "CONF_DIR", Value: coh.ConfDir},
	}
	g.Expect(c.Env).To(Equal(env))

	mounts := []corev1.VolumeMount{
		{Name: coh.VolumeNameUtils, MountPath: coh.VolumeMountPathUtils},
		{Name: coh.VolumeNameApplication, MountPath: coh.ExternalAppDir},
	}
	g.Expect(c.VolumeMounts).To(Equal(mounts))
}

func TestApplicationSpec_CreateApplicationContainer_WithConfDir(t *testing.T) {
	g := NewGomegaWithT(t)

	app := coh.ApplicationSpec{
		ImageSpec: coh.ImageSpec{
			Image: stringPtr("my-test-image:1.0"),
		},
		ConfigDir: stringPtr("/home/test/conf"),
	}

	ok, c := app.CreateApplicationContainer()
	g.Expect(ok).To(BeTrue())
	g.Expect(c.Name).To(Equal(coh.ContainerNameApplication))
	g.Expect(c.Image).To(Equal("my-test-image:1.0"))
	g.Expect(c.Command).To(Equal([]string{coh.DefaultCommandApplication}))
	g.Expect(c.ImagePullPolicy).To(BeEmpty())

	env := []corev1.EnvVar{
		{Name: "EXTERNAL_APP_DIR", Value: coh.ExternalAppDir},
		{Name: "APP_DIR", Value: coh.AppDir},
		{Name: "EXTERNAL_LIB_DIR", Value: coh.ExternalLibDir},
		{Name: "LIB_DIR", Value: coh.LibDir},
		{Name: "EXTERNAL_CONF_DIR", Value: coh.ExternalConfDir},
		{Name: "CONF_DIR", Value: "/home/test/conf"},
	}
	g.Expect(c.Env).To(Equal(env))

	mounts := []corev1.VolumeMount{
		{Name: coh.VolumeNameUtils, MountPath: coh.VolumeMountPathUtils},
		{Name: coh.VolumeNameApplication, MountPath: coh.ExternalAppDir},
	}
	g.Expect(c.VolumeMounts).To(Equal(mounts))
}
