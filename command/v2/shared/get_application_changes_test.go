package shared_test

import (
	"fmt"

	"code.cloudfoundry.org/cli/actor/pushaction"
	"code.cloudfoundry.org/cli/actor/v2action"
	. "code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetApplicationChanges", func() {
	var (
		appName string

		appConfig pushaction.ApplicationConfig
		changes   []ui.Change
	)

	BeforeEach(func() {
		appName = "steve"

		appConfig = pushaction.ApplicationConfig{
			CurrentApplication: pushaction.Application{
				Application: v2action.Application{
					Name:      appName,
					StackGUID: "some-old-stack-guid",
				}},
			DesiredApplication: pushaction.Application{
				Application: v2action.Application{
					Name:      appName,
					StackGUID: "some-new-stack-guid",
				}},
			Path: "/foo/bar",
			CurrentRoutes: []v2action.Route{
				{Host: "route1", Domain: v2action.Domain{Name: "example.com"}},
				{Host: "route2", Domain: v2action.Domain{Name: "example.com"}},
			},
			DesiredRoutes: []v2action.Route{
				{Host: "route3", Domain: v2action.Domain{Name: "example.com"}},
				{Host: "route4", Domain: v2action.Domain{Name: "example.com"}},
			},
		}
	})

	JustBeforeEach(func() {
		changes = GetApplicationChanges(appConfig)
	})

	Context("name", func() {
		It("sets the first change to name", func() {
			Expect(changes[0]).To(Equal(ui.Change{
				Header:       "name:",
				CurrentValue: appName,
				NewValue:     appName,
			}))
		})
	})

	Context("docker image", func() {
		BeforeEach(func() {
			appConfig.CurrentApplication.DockerImage = "some-path"
			appConfig.DesiredApplication.DockerImage = "some-new-path"
		})

		It("set the second change to docker image", func() {
			Expect(changes[1]).To(Equal(ui.Change{
				Header:       "docker image:",
				CurrentValue: "some-path",
				NewValue:     "some-new-path",
			}))
		})
	})

	Context("path", func() {
		It("sets the second change to path", func() {
			Expect(changes[1]).To(Equal(ui.Change{
				Header:       "path:",
				CurrentValue: "/foo/bar",
				NewValue:     "/foo/bar",
			}))
		})
	})

	Context("buildpack", func() {
		Describe("new app with no specified buildpack", func() {
			It("does not provide a buildpack change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("buildpack:"), fmt.Sprintf("entry %d should not be a buildpack", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(
				currentBuildpack string, currentDetectedBuildpack string,
				desiredBuildpack string, desiredDetectedBuildpack string,
				currentValue string, newValue string,
			) {
				appConfig.CurrentApplication.Buildpack = currentBuildpack
				appConfig.CurrentApplication.DetectedBuildpack = currentDetectedBuildpack
				appConfig.DesiredApplication.Buildpack = desiredBuildpack
				appConfig.DesiredApplication.DetectedBuildpack = desiredDetectedBuildpack

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "buildpack:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},

			Entry("new app with buildpack specified",
				"", "",
				"some-new-buildpack", "",
				"", "some-new-buildpack",
			),
			Entry("existing buildpack with new buildpack specified",
				"some-old-buildpack", "",
				"some-new-buildpack", "",
				"some-old-buildpack", "some-new-buildpack",
			),
			Entry("existing detected buildpack with new buildpack specified",
				"", "some-detected-buildpack",
				"some-new-buildpack", "",
				"some-detected-buildpack", "some-new-buildpack",
			),
			Entry("existing detected buildpack with new detected buildpack",
				"", "some-detected-buildpack",
				"", "some-detected-buildpack",
				"some-detected-buildpack", "some-detected-buildpack",
			),

			// Can never happen because desired starts as a copy of current
			Entry("existing buildpack with no new buildpack specified",
				"some-old-buildpack", "",
				"", "",
				"some-old-buildpack", "",
			),
			// Can never happen because desired starts as a copy of current
			Entry("existing detected buildpack with no new buildpack specified",
				"", "some-detected-buildpack",
				"", "",
				"some-detected-buildpack", "",
			),
		)
	})

	Context("command", func() {
		Describe("new app with no specified command", func() {
			It("does not provide a command change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("command:"), fmt.Sprintf("entry %d should not be command", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(
				currentCommand string, currentDetectedCommand string,
				desiredCommand string, desiredDetectedCommand string,
				currentValue string, newValue string,
			) {
				appConfig.CurrentApplication.Command = currentCommand
				appConfig.CurrentApplication.DetectedStartCommand = currentDetectedCommand
				appConfig.DesiredApplication.Command = desiredCommand
				appConfig.DesiredApplication.DetectedStartCommand = desiredDetectedCommand

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "command:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with command specified",
				"", "",
				"some-new-command", "",
				"", "some-new-command",
			),
			Entry("existing command with new command specified",
				"some-old-command", "",
				"some-new-command", "",
				"some-old-command", "some-new-command",
			),
			Entry("existing detected command with new command specified",
				"", "some-detected-command",
				"some-new-command", "",
				"some-detected-command", "some-new-command",
			),
			Entry("existing detected command with new detected command",
				"", "some-detected-command",
				"", "some-detected-command",
				"some-detected-command", "some-detected-command",
			),

			// Can never happen because desired starts as a copy of current
			Entry("existing command with no new command specified",
				"some-old-command", "",
				"", "",
				"some-old-command", "",
			),
			// Can never happen because desired starts as a copy of current
			Entry("existing detected command with no new command specified",
				"", "some-detected-command",
				"", "",
				"some-detected-command", "",
			),
		)
	})

	Context("disk_quota", func() {
		Describe("new app with no specified disk_quota", func() {
			It("does not provide a disk_quota change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("disk quota:"), fmt.Sprintf("entry %d should not be disk quota", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingDiskQuota int, newDiskQuota int, currentValue string, newValue string) {
				appConfig.CurrentApplication.DiskQuota = uint64(existingDiskQuota)
				appConfig.DesiredApplication.DiskQuota = uint64(newDiskQuota)

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "disk quota:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with disk_quota specified", 0, 200, "", "200M"),
			Entry("existing disk_quota with no disk_quota specified", 100, 0, "100M", "0"),
			Entry("existing disk_quota with new disk_quota specified", 100, 200, "100M", "200M"),
		)
	})

	Context("health-check-http-endpoint", func() {
		Describe("new app with no specified health check http endpoint", func() {
			It("does not provide an http endpoint check type change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("health check http endpoint:"), fmt.Sprintf("entry %d should not be health check http endpoint", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingType string, newType string, currentValue string, newValue string) {
				appConfig.CurrentApplication.HealthCheckHTTPEndpoint = existingType
				appConfig.DesiredApplication.HealthCheckHTTPEndpoint = newType

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "health check http endpoint:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with http-endpoint specified", "", "some-new-http-endpoint", "", "some-new-http-endpoint"),
			Entry("existing http-endpoint with no http-endpoint specified", "some-old-http-endpoint", "", "some-old-http-endpoint", ""),
			Entry("existing http-endpoint with new http-endpoint specified", "some-old-http-endpoint", "some-new-http-endpoint", "some-old-http-endpoint", "some-new-http-endpoint"),
		)
	})

	Context("health-check-timeout", func() {
		Describe("new app with no specified health check timeout", func() {
			It("does not provide an health check timeout change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("health check http endpoint:"), fmt.Sprintf("entry %d should not be health check http endpoint", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingType int, newType int, currentValue int, newValue int) {
				appConfig.CurrentApplication.HealthCheckTimeout = existingType
				appConfig.DesiredApplication.HealthCheckTimeout = newType

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "health check timeout:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with health-check-timeout specified", 0, 200, 0, 200),
			Entry("existing health-check-timeout with no health-check-timeout specified", 100, 0, 100, 0),
			Entry("existing health-check-timeout with new health-check-timeout specified", 100, 200, 100, 200),
		)
	})

	Context("health-check-type", func() {
		Describe("new app with no specified health-check-type", func() {
			It("does not provide a health check type change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("health check type:"), fmt.Sprintf("entry %d should not be health check type", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingType string, newType string, currentValue string, newValue string) {
				appConfig.CurrentApplication.HealthCheckType = existingType
				appConfig.DesiredApplication.HealthCheckType = newType

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "health check type:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with health-check-type specified", "", "some-new-health-check-type", "", "some-new-health-check-type"),
			Entry("existing health-check-type with no health-check-type specified", "some-old-health-check-type", "", "some-old-health-check-type", ""),
			Entry("existing health-check-type with new health-check-type specified", "some-old-health-check-type", "some-new-health-check-type", "some-old-health-check-type", "some-new-health-check-type"),
		)
	})

	Context("instances", func() {
		Describe("new app with no specified instances", func() {
			It("does not provide an instances change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("instances:"), fmt.Sprintf("entry %d should not be instances", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingInstances int, newInstances int, currentValue int, newValue int) {
				appConfig.CurrentApplication.Instances = existingInstances
				appConfig.DesiredApplication.Instances = newInstances

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "instances:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with instances specified", 0, 200, 0, 200),
			Entry("existing instances with no instances specified", 100, 0, 100, 0),
			Entry("existing instances with new instances specified", 100, 200, 100, 200),
		)
	})

	Context("memory", func() {
		Describe("new app with no specified memory", func() {
			It("does not provide a memory change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("memory:"), fmt.Sprintf("entry %d should not be memory", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingMemory int, newMemory int, currentValue string, newValue string) {
				appConfig.CurrentApplication.Memory = uint64(existingMemory)
				appConfig.DesiredApplication.Memory = uint64(newMemory)

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "memory:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with memory specified", 0, 200, "", "200M"),
			Entry("existing memory with no memory specified", 100, 0, "100M", "0"),
			Entry("existing memory with new memory specified", 100, 200, "100M", "200M"),
		)
	})

	Context("stack", func() {
		Describe("new app with no specified stack", func() {
			It("does not provide an stack change", func() {
				for i, change := range changes {
					Expect(change.Header).ToNot(Equal("stack:"), fmt.Sprintf("entry %d should not be stack", i))
				}
			})
		})

		DescribeTable("non-empty values",
			func(existingStack string, newStack string, currentValue string, newValue string) {
				appConfig.CurrentApplication.Stack.Name = existingStack
				appConfig.DesiredApplication.Stack.Name = newStack

				changes = GetApplicationChanges(appConfig)

				Expect(changes[2]).To(Equal(ui.Change{
					Header:       "stack:",
					CurrentValue: currentValue,
					NewValue:     newValue,
				}))
			},
			Entry("new app with stack specified", "", "some-new-stack", "", "some-new-stack"),
			Entry("existing stack with no stack specified", "some-old-stack", "", "some-old-stack", ""),
			Entry("existing stack with new stack specified", "some-old-stack", "some-new-stack", "some-old-stack", "some-new-stack"),
		)
	})

	Context("services", func() {
		BeforeEach(func() {
			appConfig.CurrentServices = map[string]v2action.ServiceInstance{"service-1": {}, "service-2": {}}
			appConfig.DesiredServices = map[string]v2action.ServiceInstance{"service-3": {}, "service-4": {}}
		})

		It("sets the third change to services", func() {
			Expect(len(changes)).To(BeNumerically(">=", 2))
			change := changes[2]
			Expect(change.Header).To(Equal("services:"))
			Expect(change.CurrentValue).To(ConsistOf([]string{"service-1", "service-2"}))
			Expect(change.NewValue).To(ConsistOf([]string{"service-3", "service-4"}))
		})
	})

	Context("user provided environment variables", func() {
		var oldMap, newMap map[string]string

		BeforeEach(func() {
			oldMap = map[string]string{"a": "b"}
			newMap = map[string]string{"1": "2"}
			appConfig.CurrentApplication.EnvironmentVariables = oldMap
			appConfig.DesiredApplication.EnvironmentVariables = newMap
		})

		It("sets the fourth change to routes", func() {
			Expect(changes[3]).To(Equal(ui.Change{
				Header:       "env:",
				CurrentValue: oldMap,
				NewValue:     newMap,
			}))
		})
	})

	Context("routes", func() {
		It("sets the fifth change to routes", func() {
			Expect(changes[4]).To(Equal(ui.Change{
				Header:       "routes:",
				CurrentValue: []string{"route1.example.com", "route2.example.com"},
				NewValue:     []string{"route3.example.com", "route4.example.com"},
			}))
		})
	})
})
