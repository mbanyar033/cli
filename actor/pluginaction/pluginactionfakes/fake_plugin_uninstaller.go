// This file was generated by counterfeiter
package pluginactionfakes

import (
	"sync"

	"code.cloudfoundry.org/cli/actor/pluginaction"
)

type FakePluginUninstaller struct {
	UninstallStub        func(pluginPath string) error
	uninstallMutex       sync.RWMutex
	uninstallArgsForCall []struct {
		pluginPath string
	}
	uninstallReturns struct {
		result1 error
	}
	uninstallReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePluginUninstaller) Uninstall(pluginPath string) error {
	fake.uninstallMutex.Lock()
	ret, specificReturn := fake.uninstallReturnsOnCall[len(fake.uninstallArgsForCall)]
	fake.uninstallArgsForCall = append(fake.uninstallArgsForCall, struct {
		pluginPath string
	}{pluginPath})
	fake.recordInvocation("Uninstall", []interface{}{pluginPath})
	fake.uninstallMutex.Unlock()
	if fake.UninstallStub != nil {
		return fake.UninstallStub(pluginPath)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.uninstallReturns.result1
}

func (fake *FakePluginUninstaller) UninstallCallCount() int {
	fake.uninstallMutex.RLock()
	defer fake.uninstallMutex.RUnlock()
	return len(fake.uninstallArgsForCall)
}

func (fake *FakePluginUninstaller) UninstallArgsForCall(i int) string {
	fake.uninstallMutex.RLock()
	defer fake.uninstallMutex.RUnlock()
	return fake.uninstallArgsForCall[i].pluginPath
}

func (fake *FakePluginUninstaller) UninstallReturns(result1 error) {
	fake.UninstallStub = nil
	fake.uninstallReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePluginUninstaller) UninstallReturnsOnCall(i int, result1 error) {
	fake.UninstallStub = nil
	if fake.uninstallReturnsOnCall == nil {
		fake.uninstallReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.uninstallReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePluginUninstaller) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.uninstallMutex.RLock()
	defer fake.uninstallMutex.RUnlock()
	return fake.invocations
}

func (fake *FakePluginUninstaller) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ pluginaction.PluginUninstaller = new(FakePluginUninstaller)
