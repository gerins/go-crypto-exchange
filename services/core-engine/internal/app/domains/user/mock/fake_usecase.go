// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"context"
	"core-engine/internal/app/domains/user"
	"sync"
)

type FakeUsecase struct {
	LoginStub        func(context.Context, user.LoginRequest) (user.LoginResponse, error)
	loginMutex       sync.RWMutex
	loginArgsForCall []struct {
		arg1 context.Context
		arg2 user.LoginRequest
	}
	loginReturns struct {
		result1 user.LoginResponse
		result2 error
	}
	loginReturnsOnCall map[int]struct {
		result1 user.LoginResponse
		result2 error
	}
	RegisterStub        func(context.Context, user.RegisterRequest) (user.User, error)
	registerMutex       sync.RWMutex
	registerArgsForCall []struct {
		arg1 context.Context
		arg2 user.RegisterRequest
	}
	registerReturns struct {
		result1 user.User
		result2 error
	}
	registerReturnsOnCall map[int]struct {
		result1 user.User
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUsecase) Login(arg1 context.Context, arg2 user.LoginRequest) (user.LoginResponse, error) {
	fake.loginMutex.Lock()
	ret, specificReturn := fake.loginReturnsOnCall[len(fake.loginArgsForCall)]
	fake.loginArgsForCall = append(fake.loginArgsForCall, struct {
		arg1 context.Context
		arg2 user.LoginRequest
	}{arg1, arg2})
	stub := fake.LoginStub
	fakeReturns := fake.loginReturns
	fake.recordInvocation("Login", []interface{}{arg1, arg2})
	fake.loginMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsecase) LoginCallCount() int {
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	return len(fake.loginArgsForCall)
}

func (fake *FakeUsecase) LoginCalls(stub func(context.Context, user.LoginRequest) (user.LoginResponse, error)) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = stub
}

func (fake *FakeUsecase) LoginArgsForCall(i int) (context.Context, user.LoginRequest) {
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	argsForCall := fake.loginArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUsecase) LoginReturns(result1 user.LoginResponse, result2 error) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = nil
	fake.loginReturns = struct {
		result1 user.LoginResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsecase) LoginReturnsOnCall(i int, result1 user.LoginResponse, result2 error) {
	fake.loginMutex.Lock()
	defer fake.loginMutex.Unlock()
	fake.LoginStub = nil
	if fake.loginReturnsOnCall == nil {
		fake.loginReturnsOnCall = make(map[int]struct {
			result1 user.LoginResponse
			result2 error
		})
	}
	fake.loginReturnsOnCall[i] = struct {
		result1 user.LoginResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeUsecase) Register(arg1 context.Context, arg2 user.RegisterRequest) (user.User, error) {
	fake.registerMutex.Lock()
	ret, specificReturn := fake.registerReturnsOnCall[len(fake.registerArgsForCall)]
	fake.registerArgsForCall = append(fake.registerArgsForCall, struct {
		arg1 context.Context
		arg2 user.RegisterRequest
	}{arg1, arg2})
	stub := fake.RegisterStub
	fakeReturns := fake.registerReturns
	fake.recordInvocation("Register", []interface{}{arg1, arg2})
	fake.registerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUsecase) RegisterCallCount() int {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	return len(fake.registerArgsForCall)
}

func (fake *FakeUsecase) RegisterCalls(stub func(context.Context, user.RegisterRequest) (user.User, error)) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = stub
}

func (fake *FakeUsecase) RegisterArgsForCall(i int) (context.Context, user.RegisterRequest) {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	argsForCall := fake.registerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUsecase) RegisterReturns(result1 user.User, result2 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	fake.registerReturns = struct {
		result1 user.User
		result2 error
	}{result1, result2}
}

func (fake *FakeUsecase) RegisterReturnsOnCall(i int, result1 user.User, result2 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	if fake.registerReturnsOnCall == nil {
		fake.registerReturnsOnCall = make(map[int]struct {
			result1 user.User
			result2 error
		})
	}
	fake.registerReturnsOnCall[i] = struct {
		result1 user.User
		result2 error
	}{result1, result2}
}

func (fake *FakeUsecase) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUsecase) recordInvocation(key string, args []interface{}) {
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

var _ user.Usecase = new(FakeUsecase)
