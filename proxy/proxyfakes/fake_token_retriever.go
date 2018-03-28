// Code generated by counterfeiter. DO NOT EDIT.
package proxyfakes

import (
	"sync"

	"code.cloudfoundry.org/gcp-broker-proxy/proxy"
	"golang.org/x/oauth2"
)

type FakeTokenRetriever struct {
	GetTokenStub        func() (*oauth2.Token, error)
	getTokenMutex       sync.RWMutex
	getTokenArgsForCall []struct{}
	getTokenReturns     struct {
		result1 *oauth2.Token
		result2 error
	}
	getTokenReturnsOnCall map[int]struct {
		result1 *oauth2.Token
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTokenRetriever) GetToken() (*oauth2.Token, error) {
	fake.getTokenMutex.Lock()
	ret, specificReturn := fake.getTokenReturnsOnCall[len(fake.getTokenArgsForCall)]
	fake.getTokenArgsForCall = append(fake.getTokenArgsForCall, struct{}{})
	fake.recordInvocation("GetToken", []interface{}{})
	fake.getTokenMutex.Unlock()
	if fake.GetTokenStub != nil {
		return fake.GetTokenStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTokenReturns.result1, fake.getTokenReturns.result2
}

func (fake *FakeTokenRetriever) GetTokenCallCount() int {
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	return len(fake.getTokenArgsForCall)
}

func (fake *FakeTokenRetriever) GetTokenReturns(result1 *oauth2.Token, result2 error) {
	fake.GetTokenStub = nil
	fake.getTokenReturns = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeTokenRetriever) GetTokenReturnsOnCall(i int, result1 *oauth2.Token, result2 error) {
	fake.GetTokenStub = nil
	if fake.getTokenReturnsOnCall == nil {
		fake.getTokenReturnsOnCall = make(map[int]struct {
			result1 *oauth2.Token
			result2 error
		})
	}
	fake.getTokenReturnsOnCall[i] = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeTokenRetriever) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getTokenMutex.RLock()
	defer fake.getTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTokenRetriever) recordInvocation(key string, args []interface{}) {
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

var _ proxy.TokenRetriever = new(FakeTokenRetriever)