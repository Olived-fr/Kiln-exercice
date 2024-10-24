// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"

	tzkt "kiln-exercice/pkg/tzkt"
)

// XTZSDK is an autogenerated mock type for the XTZSDK type
type XTZSDK struct {
	mock.Mock
}

type XTZSDK_Expecter struct {
	mock *mock.Mock
}

func (_m *XTZSDK) EXPECT() *XTZSDK_Expecter {
	return &XTZSDK_Expecter{mock: &_m.Mock}
}

// GetDelegations provides a mock function with given fields: ctx, from, to
func (_m *XTZSDK) GetDelegations(ctx context.Context, from time.Time, to time.Time) ([]tzkt.Delegation, error) {
	ret := _m.Called(ctx, from, to)

	if len(ret) == 0 {
		panic("no return value specified for GetDelegations")
	}

	var r0 []tzkt.Delegation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) ([]tzkt.Delegation, error)); ok {
		return rf(ctx, from, to)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) []tzkt.Delegation); ok {
		r0 = rf(ctx, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tzkt.Delegation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XTZSDK_GetDelegations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDelegations'
type XTZSDK_GetDelegations_Call struct {
	*mock.Call
}

// GetDelegations is a helper method to define mock.On call
//   - ctx context.Context
//   - from time.Time
//   - to time.Time
func (_e *XTZSDK_Expecter) GetDelegations(ctx interface{}, from interface{}, to interface{}) *XTZSDK_GetDelegations_Call {
	return &XTZSDK_GetDelegations_Call{Call: _e.mock.On("GetDelegations", ctx, from, to)}
}

func (_c *XTZSDK_GetDelegations_Call) Run(run func(ctx context.Context, from time.Time, to time.Time)) *XTZSDK_GetDelegations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Time), args[2].(time.Time))
	})
	return _c
}

func (_c *XTZSDK_GetDelegations_Call) Return(_a0 []tzkt.Delegation, _a1 error) *XTZSDK_GetDelegations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XTZSDK_GetDelegations_Call) RunAndReturn(run func(context.Context, time.Time, time.Time) ([]tzkt.Delegation, error)) *XTZSDK_GetDelegations_Call {
	_c.Call.Return(run)
	return _c
}

// NewXTZSDK creates a new instance of XTZSDK. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewXTZSDK(t interface {
	mock.TestingT
	Cleanup(func())
}) *XTZSDK {
	mock := &XTZSDK{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
