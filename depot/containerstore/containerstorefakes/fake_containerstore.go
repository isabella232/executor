// This file was generated by counterfeiter
package containerstorefakes

import (
	"io"
	"sync"

	"code.cloudfoundry.org/executor"
	"code.cloudfoundry.org/executor/depot/containerstore"
	"code.cloudfoundry.org/lager"
	"github.com/tedsuo/ifrit"
)

type FakeContainerStore struct {
	ReserveStub        func(logger lager.Logger, req *executor.AllocationRequest) (executor.Container, error)
	reserveMutex       sync.RWMutex
	reserveArgsForCall []struct {
		logger lager.Logger
		req    *executor.AllocationRequest
	}
	reserveReturns struct {
		result1 executor.Container
		result2 error
	}
	DestroyStub        func(logger lager.Logger, guid string) error
	destroyMutex       sync.RWMutex
	destroyArgsForCall []struct {
		logger lager.Logger
		guid   string
	}
	destroyReturns struct {
		result1 error
	}
	InitializeStub        func(logger lager.Logger, req *executor.RunRequest) error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		logger lager.Logger
		req    *executor.RunRequest
	}
	initializeReturns struct {
		result1 error
	}
	CreateStub        func(logger lager.Logger, guid string) (executor.Container, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		logger lager.Logger
		guid   string
	}
	createReturns struct {
		result1 executor.Container
		result2 error
	}
	RunStub        func(logger lager.Logger, guid string) error
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		logger lager.Logger
		guid   string
	}
	runReturns struct {
		result1 error
	}
	StopStub        func(logger lager.Logger, guid string) error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
		logger lager.Logger
		guid   string
	}
	stopReturns struct {
		result1 error
	}
	GetStub        func(logger lager.Logger, guid string) (executor.Container, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		logger lager.Logger
		guid   string
	}
	getReturns struct {
		result1 executor.Container
		result2 error
	}
	ListStub        func(logger lager.Logger) []executor.Container
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		logger lager.Logger
	}
	listReturns struct {
		result1 []executor.Container
	}
	MetricsStub        func(logger lager.Logger) (map[string]executor.ContainerMetrics, error)
	metricsMutex       sync.RWMutex
	metricsArgsForCall []struct {
		logger lager.Logger
	}
	metricsReturns struct {
		result1 map[string]executor.ContainerMetrics
		result2 error
	}
	RemainingResourcesStub        func(logger lager.Logger) executor.ExecutorResources
	remainingResourcesMutex       sync.RWMutex
	remainingResourcesArgsForCall []struct {
		logger lager.Logger
	}
	remainingResourcesReturns struct {
		result1 executor.ExecutorResources
	}
	GetFilesStub        func(logger lager.Logger, guid, sourcePath string) (io.ReadCloser, error)
	getFilesMutex       sync.RWMutex
	getFilesArgsForCall []struct {
		logger     lager.Logger
		guid       string
		sourcePath string
	}
	getFilesReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	NewRegistryPrunerStub        func(logger lager.Logger) ifrit.Runner
	newRegistryPrunerMutex       sync.RWMutex
	newRegistryPrunerArgsForCall []struct {
		logger lager.Logger
	}
	newRegistryPrunerReturns struct {
		result1 ifrit.Runner
	}
	NewContainerReaperStub        func(logger lager.Logger) ifrit.Runner
	newContainerReaperMutex       sync.RWMutex
	newContainerReaperArgsForCall []struct {
		logger lager.Logger
	}
	newContainerReaperReturns struct {
		result1 ifrit.Runner
	}
}

func (fake *FakeContainerStore) Reserve(logger lager.Logger, req *executor.AllocationRequest) (executor.Container, error) {
	fake.reserveMutex.Lock()
	fake.reserveArgsForCall = append(fake.reserveArgsForCall, struct {
		logger lager.Logger
		req    *executor.AllocationRequest
	}{logger, req})
	fake.reserveMutex.Unlock()
	if fake.ReserveStub != nil {
		return fake.ReserveStub(logger, req)
	} else {
		return fake.reserveReturns.result1, fake.reserveReturns.result2
	}
}

func (fake *FakeContainerStore) ReserveCallCount() int {
	fake.reserveMutex.RLock()
	defer fake.reserveMutex.RUnlock()
	return len(fake.reserveArgsForCall)
}

func (fake *FakeContainerStore) ReserveArgsForCall(i int) (lager.Logger, *executor.AllocationRequest) {
	fake.reserveMutex.RLock()
	defer fake.reserveMutex.RUnlock()
	return fake.reserveArgsForCall[i].logger, fake.reserveArgsForCall[i].req
}

func (fake *FakeContainerStore) ReserveReturns(result1 executor.Container, result2 error) {
	fake.ReserveStub = nil
	fake.reserveReturns = struct {
		result1 executor.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerStore) Destroy(logger lager.Logger, guid string) error {
	fake.destroyMutex.Lock()
	fake.destroyArgsForCall = append(fake.destroyArgsForCall, struct {
		logger lager.Logger
		guid   string
	}{logger, guid})
	fake.destroyMutex.Unlock()
	if fake.DestroyStub != nil {
		return fake.DestroyStub(logger, guid)
	} else {
		return fake.destroyReturns.result1
	}
}

func (fake *FakeContainerStore) DestroyCallCount() int {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return len(fake.destroyArgsForCall)
}

func (fake *FakeContainerStore) DestroyArgsForCall(i int) (lager.Logger, string) {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.destroyArgsForCall[i].logger, fake.destroyArgsForCall[i].guid
}

func (fake *FakeContainerStore) DestroyReturns(result1 error) {
	fake.DestroyStub = nil
	fake.destroyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerStore) Initialize(logger lager.Logger, req *executor.RunRequest) error {
	fake.initializeMutex.Lock()
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
		logger lager.Logger
		req    *executor.RunRequest
	}{logger, req})
	fake.initializeMutex.Unlock()
	if fake.InitializeStub != nil {
		return fake.InitializeStub(logger, req)
	} else {
		return fake.initializeReturns.result1
	}
}

func (fake *FakeContainerStore) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *FakeContainerStore) InitializeArgsForCall(i int) (lager.Logger, *executor.RunRequest) {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return fake.initializeArgsForCall[i].logger, fake.initializeArgsForCall[i].req
}

func (fake *FakeContainerStore) InitializeReturns(result1 error) {
	fake.InitializeStub = nil
	fake.initializeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerStore) Create(logger lager.Logger, guid string) (executor.Container, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		logger lager.Logger
		guid   string
	}{logger, guid})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(logger, guid)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2
	}
}

func (fake *FakeContainerStore) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeContainerStore) CreateArgsForCall(i int) (lager.Logger, string) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].logger, fake.createArgsForCall[i].guid
}

func (fake *FakeContainerStore) CreateReturns(result1 executor.Container, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 executor.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerStore) Run(logger lager.Logger, guid string) error {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		logger lager.Logger
		guid   string
	}{logger, guid})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		return fake.RunStub(logger, guid)
	} else {
		return fake.runReturns.result1
	}
}

func (fake *FakeContainerStore) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeContainerStore) RunArgsForCall(i int) (lager.Logger, string) {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.runArgsForCall[i].logger, fake.runArgsForCall[i].guid
}

func (fake *FakeContainerStore) RunReturns(result1 error) {
	fake.RunStub = nil
	fake.runReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerStore) Stop(logger lager.Logger, guid string) error {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
		logger lager.Logger
		guid   string
	}{logger, guid})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub(logger, guid)
	} else {
		return fake.stopReturns.result1
	}
}

func (fake *FakeContainerStore) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *FakeContainerStore) StopArgsForCall(i int) (lager.Logger, string) {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return fake.stopArgsForCall[i].logger, fake.stopArgsForCall[i].guid
}

func (fake *FakeContainerStore) StopReturns(result1 error) {
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerStore) Get(logger lager.Logger, guid string) (executor.Container, error) {
	fake.getMutex.Lock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		logger lager.Logger
		guid   string
	}{logger, guid})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(logger, guid)
	} else {
		return fake.getReturns.result1, fake.getReturns.result2
	}
}

func (fake *FakeContainerStore) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeContainerStore) GetArgsForCall(i int) (lager.Logger, string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].logger, fake.getArgsForCall[i].guid
}

func (fake *FakeContainerStore) GetReturns(result1 executor.Container, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 executor.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerStore) List(logger lager.Logger) []executor.Container {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(logger)
	} else {
		return fake.listReturns.result1
	}
}

func (fake *FakeContainerStore) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeContainerStore) ListArgsForCall(i int) lager.Logger {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return fake.listArgsForCall[i].logger
}

func (fake *FakeContainerStore) ListReturns(result1 []executor.Container) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []executor.Container
	}{result1}
}

func (fake *FakeContainerStore) Metrics(logger lager.Logger) (map[string]executor.ContainerMetrics, error) {
	fake.metricsMutex.Lock()
	fake.metricsArgsForCall = append(fake.metricsArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.metricsMutex.Unlock()
	if fake.MetricsStub != nil {
		return fake.MetricsStub(logger)
	} else {
		return fake.metricsReturns.result1, fake.metricsReturns.result2
	}
}

func (fake *FakeContainerStore) MetricsCallCount() int {
	fake.metricsMutex.RLock()
	defer fake.metricsMutex.RUnlock()
	return len(fake.metricsArgsForCall)
}

func (fake *FakeContainerStore) MetricsArgsForCall(i int) lager.Logger {
	fake.metricsMutex.RLock()
	defer fake.metricsMutex.RUnlock()
	return fake.metricsArgsForCall[i].logger
}

func (fake *FakeContainerStore) MetricsReturns(result1 map[string]executor.ContainerMetrics, result2 error) {
	fake.MetricsStub = nil
	fake.metricsReturns = struct {
		result1 map[string]executor.ContainerMetrics
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerStore) RemainingResources(logger lager.Logger) executor.ExecutorResources {
	fake.remainingResourcesMutex.Lock()
	fake.remainingResourcesArgsForCall = append(fake.remainingResourcesArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.remainingResourcesMutex.Unlock()
	if fake.RemainingResourcesStub != nil {
		return fake.RemainingResourcesStub(logger)
	} else {
		return fake.remainingResourcesReturns.result1
	}
}

func (fake *FakeContainerStore) RemainingResourcesCallCount() int {
	fake.remainingResourcesMutex.RLock()
	defer fake.remainingResourcesMutex.RUnlock()
	return len(fake.remainingResourcesArgsForCall)
}

func (fake *FakeContainerStore) RemainingResourcesArgsForCall(i int) lager.Logger {
	fake.remainingResourcesMutex.RLock()
	defer fake.remainingResourcesMutex.RUnlock()
	return fake.remainingResourcesArgsForCall[i].logger
}

func (fake *FakeContainerStore) RemainingResourcesReturns(result1 executor.ExecutorResources) {
	fake.RemainingResourcesStub = nil
	fake.remainingResourcesReturns = struct {
		result1 executor.ExecutorResources
	}{result1}
}

func (fake *FakeContainerStore) GetFiles(logger lager.Logger, guid string, sourcePath string) (io.ReadCloser, error) {
	fake.getFilesMutex.Lock()
	fake.getFilesArgsForCall = append(fake.getFilesArgsForCall, struct {
		logger     lager.Logger
		guid       string
		sourcePath string
	}{logger, guid, sourcePath})
	fake.getFilesMutex.Unlock()
	if fake.GetFilesStub != nil {
		return fake.GetFilesStub(logger, guid, sourcePath)
	} else {
		return fake.getFilesReturns.result1, fake.getFilesReturns.result2
	}
}

func (fake *FakeContainerStore) GetFilesCallCount() int {
	fake.getFilesMutex.RLock()
	defer fake.getFilesMutex.RUnlock()
	return len(fake.getFilesArgsForCall)
}

func (fake *FakeContainerStore) GetFilesArgsForCall(i int) (lager.Logger, string, string) {
	fake.getFilesMutex.RLock()
	defer fake.getFilesMutex.RUnlock()
	return fake.getFilesArgsForCall[i].logger, fake.getFilesArgsForCall[i].guid, fake.getFilesArgsForCall[i].sourcePath
}

func (fake *FakeContainerStore) GetFilesReturns(result1 io.ReadCloser, result2 error) {
	fake.GetFilesStub = nil
	fake.getFilesReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerStore) NewRegistryPruner(logger lager.Logger) ifrit.Runner {
	fake.newRegistryPrunerMutex.Lock()
	fake.newRegistryPrunerArgsForCall = append(fake.newRegistryPrunerArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.newRegistryPrunerMutex.Unlock()
	if fake.NewRegistryPrunerStub != nil {
		return fake.NewRegistryPrunerStub(logger)
	} else {
		return fake.newRegistryPrunerReturns.result1
	}
}

func (fake *FakeContainerStore) NewRegistryPrunerCallCount() int {
	fake.newRegistryPrunerMutex.RLock()
	defer fake.newRegistryPrunerMutex.RUnlock()
	return len(fake.newRegistryPrunerArgsForCall)
}

func (fake *FakeContainerStore) NewRegistryPrunerArgsForCall(i int) lager.Logger {
	fake.newRegistryPrunerMutex.RLock()
	defer fake.newRegistryPrunerMutex.RUnlock()
	return fake.newRegistryPrunerArgsForCall[i].logger
}

func (fake *FakeContainerStore) NewRegistryPrunerReturns(result1 ifrit.Runner) {
	fake.NewRegistryPrunerStub = nil
	fake.newRegistryPrunerReturns = struct {
		result1 ifrit.Runner
	}{result1}
}

func (fake *FakeContainerStore) NewContainerReaper(logger lager.Logger) ifrit.Runner {
	fake.newContainerReaperMutex.Lock()
	fake.newContainerReaperArgsForCall = append(fake.newContainerReaperArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.newContainerReaperMutex.Unlock()
	if fake.NewContainerReaperStub != nil {
		return fake.NewContainerReaperStub(logger)
	} else {
		return fake.newContainerReaperReturns.result1
	}
}

func (fake *FakeContainerStore) NewContainerReaperCallCount() int {
	fake.newContainerReaperMutex.RLock()
	defer fake.newContainerReaperMutex.RUnlock()
	return len(fake.newContainerReaperArgsForCall)
}

func (fake *FakeContainerStore) NewContainerReaperArgsForCall(i int) lager.Logger {
	fake.newContainerReaperMutex.RLock()
	defer fake.newContainerReaperMutex.RUnlock()
	return fake.newContainerReaperArgsForCall[i].logger
}

func (fake *FakeContainerStore) NewContainerReaperReturns(result1 ifrit.Runner) {
	fake.NewContainerReaperStub = nil
	fake.newContainerReaperReturns = struct {
		result1 ifrit.Runner
	}{result1}
}

var _ containerstore.ContainerStore = new(FakeContainerStore)
