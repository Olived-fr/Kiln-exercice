// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	list "kiln-exercice/internal/usecase/delegation/list"

	mock "github.com/stretchr/testify/mock"
)

// DelegationUseCase is an autogenerated mock type for the DelegationUseCase type
type DelegationUseCase struct {
	mock.Mock
}

type DelegationUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *DelegationUseCase) EXPECT() *DelegationUseCase_Expecter {
	return &DelegationUseCase_Expecter{mock: &_m.Mock}
}

// ListDelegations provides a mock function with given fields: ctx, input
func (_m *DelegationUseCase) ListDelegations(ctx context.Context, input list.Input) ([]list.DelegationData, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for ListDelegations")
	}

	var r0 []list.DelegationData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, list.Input) ([]list.DelegationData, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, list.Input) []list.DelegationData); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]list.DelegationData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, list.Input) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DelegationUseCase_ListDelegations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListDelegations'
type DelegationUseCase_ListDelegations_Call struct {
	*mock.Call
}

// ListDelegations is a helper method to define mock.On call
//   - ctx context.Context
//   - input list.Input
func (_e *DelegationUseCase_Expecter) ListDelegations(ctx interface{}, input interface{}) *DelegationUseCase_ListDelegations_Call {
	return &DelegationUseCase_ListDelegations_Call{Call: _e.mock.On("ListDelegations", ctx, input)}
}

func (_c *DelegationUseCase_ListDelegations_Call) Run(run func(ctx context.Context, input list.Input)) *DelegationUseCase_ListDelegations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(list.Input))
	})
	return _c
}

func (_c *DelegationUseCase_ListDelegations_Call) Return(_a0 []list.DelegationData, _a1 error) *DelegationUseCase_ListDelegations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DelegationUseCase_ListDelegations_Call) RunAndReturn(run func(context.Context, list.Input) ([]list.DelegationData, error)) *DelegationUseCase_ListDelegations_Call {
	_c.Call.Return(run)
	return _c
}

// NewDelegationUseCase creates a new instance of DelegationUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDelegationUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *DelegationUseCase {
	mock := &DelegationUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
