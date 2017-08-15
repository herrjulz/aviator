// Code generated by counterfeiter. DO NOT EDIT.
package processorfakes

import (
	"sync"

	"github.com/JulzDiverse/aviator/processor"
)

type FakeSpruceClient struct {
	MergeWithOptsStub        func()
	mergeWithOptsMutex       sync.RWMutex
	mergeWithOptsArgsForCall []struct{}
	invocations              map[string][][]interface{}
	invocationsMutex         sync.RWMutex
}

func (fake *FakeSpruceClient) MergeWithOpts() {
	fake.mergeWithOptsMutex.Lock()
	fake.mergeWithOptsArgsForCall = append(fake.mergeWithOptsArgsForCall, struct{}{})
	fake.recordInvocation("MergeWithOpts", []interface{}{})
	fake.mergeWithOptsMutex.Unlock()
	if fake.MergeWithOptsStub != nil {
		fake.MergeWithOptsStub()
	}
}

func (fake *FakeSpruceClient) MergeWithOptsCallCount() int {
	fake.mergeWithOptsMutex.RLock()
	defer fake.mergeWithOptsMutex.RUnlock()
	return len(fake.mergeWithOptsArgsForCall)
}

func (fake *FakeSpruceClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mergeWithOptsMutex.RLock()
	defer fake.mergeWithOptsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSpruceClient) recordInvocation(key string, args []interface{}) {
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

var _ processor.SpruceClient = new(FakeSpruceClient)