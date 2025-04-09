// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	context "context"

	kyc "github.com/raystack/frontier/core/kyc"
	mock "github.com/stretchr/testify/mock"
)

// KycService is an autogenerated mock type for the KycService type
type KycService struct {
	mock.Mock
}

type KycService_Expecter struct {
	mock *mock.Mock
}

func (_m *KycService) EXPECT() *KycService_Expecter {
	return &KycService_Expecter{mock: &_m.Mock}
}

// GetKyc provides a mock function with given fields: _a0, _a1
func (_m *KycService) GetKyc(_a0 context.Context, _a1 string) (kyc.KYC, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetKyc")
	}

	var r0 kyc.KYC
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (kyc.KYC, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) kyc.KYC); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(kyc.KYC)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KycService_GetKyc_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKyc'
type KycService_GetKyc_Call struct {
	*mock.Call
}

// GetKyc is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *KycService_Expecter) GetKyc(_a0 interface{}, _a1 interface{}) *KycService_GetKyc_Call {
	return &KycService_GetKyc_Call{Call: _e.mock.On("GetKyc", _a0, _a1)}
}

func (_c *KycService_GetKyc_Call) Run(run func(_a0 context.Context, _a1 string)) *KycService_GetKyc_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *KycService_GetKyc_Call) Return(_a0 kyc.KYC, _a1 error) *KycService_GetKyc_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KycService_GetKyc_Call) RunAndReturn(run func(context.Context, string) (kyc.KYC, error)) *KycService_GetKyc_Call {
	_c.Call.Return(run)
	return _c
}

// ListKycs provides a mock function with given fields: _a0
func (_m *KycService) ListKycs(_a0 context.Context) ([]kyc.KYC, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ListKycs")
	}

	var r0 []kyc.KYC
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]kyc.KYC, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []kyc.KYC); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]kyc.KYC)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KycService_ListKycs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListKycs'
type KycService_ListKycs_Call struct {
	*mock.Call
}

// ListKycs is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *KycService_Expecter) ListKycs(_a0 interface{}) *KycService_ListKycs_Call {
	return &KycService_ListKycs_Call{Call: _e.mock.On("ListKycs", _a0)}
}

func (_c *KycService_ListKycs_Call) Run(run func(_a0 context.Context)) *KycService_ListKycs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *KycService_ListKycs_Call) Return(_a0 []kyc.KYC, _a1 error) *KycService_ListKycs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KycService_ListKycs_Call) RunAndReturn(run func(context.Context) ([]kyc.KYC, error)) *KycService_ListKycs_Call {
	_c.Call.Return(run)
	return _c
}

// SetKyc provides a mock function with given fields: _a0, _a1
func (_m *KycService) SetKyc(_a0 context.Context, _a1 kyc.KYC) (kyc.KYC, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SetKyc")
	}

	var r0 kyc.KYC
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, kyc.KYC) (kyc.KYC, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, kyc.KYC) kyc.KYC); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(kyc.KYC)
	}

	if rf, ok := ret.Get(1).(func(context.Context, kyc.KYC) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KycService_SetKyc_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetKyc'
type KycService_SetKyc_Call struct {
	*mock.Call
}

// SetKyc is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 kyc.KYC
func (_e *KycService_Expecter) SetKyc(_a0 interface{}, _a1 interface{}) *KycService_SetKyc_Call {
	return &KycService_SetKyc_Call{Call: _e.mock.On("SetKyc", _a0, _a1)}
}

func (_c *KycService_SetKyc_Call) Run(run func(_a0 context.Context, _a1 kyc.KYC)) *KycService_SetKyc_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(kyc.KYC))
	})
	return _c
}

func (_c *KycService_SetKyc_Call) Return(_a0 kyc.KYC, _a1 error) *KycService_SetKyc_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KycService_SetKyc_Call) RunAndReturn(run func(context.Context, kyc.KYC) (kyc.KYC, error)) *KycService_SetKyc_Call {
	_c.Call.Return(run)
	return _c
}

// NewKycService creates a new instance of KycService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKycService(t interface {
	mock.TestingT
	Cleanup(func())
}) *KycService {
	mock := &KycService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
