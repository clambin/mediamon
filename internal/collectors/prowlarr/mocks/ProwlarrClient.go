// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	prowlarr "github.com/clambin/mediaclients/prowlarr"
	mock "github.com/stretchr/testify/mock"
)

// ProwlarrClient is an autogenerated mock type for the ProwlarrClient type
type ProwlarrClient struct {
	mock.Mock
}

type ProwlarrClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ProwlarrClient) EXPECT() *ProwlarrClient_Expecter {
	return &ProwlarrClient_Expecter{mock: &_m.Mock}
}

// GetApiV1IndexerstatsWithResponse provides a mock function with given fields: ctx, params, reqEditors
func (_m *ProwlarrClient) GetApiV1IndexerstatsWithResponse(ctx context.Context, params *prowlarr.GetApiV1IndexerstatsParams, reqEditors ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error) {
	_va := make([]interface{}, len(reqEditors))
	for _i := range reqEditors {
		_va[_i] = reqEditors[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetApiV1IndexerstatsWithResponse")
	}

	var r0 *prowlarr.GetApiV1IndexerstatsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *prowlarr.GetApiV1IndexerstatsParams, ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error)); ok {
		return rf(ctx, params, reqEditors...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *prowlarr.GetApiV1IndexerstatsParams, ...prowlarr.RequestEditorFn) *prowlarr.GetApiV1IndexerstatsResponse); ok {
		r0 = rf(ctx, params, reqEditors...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*prowlarr.GetApiV1IndexerstatsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *prowlarr.GetApiV1IndexerstatsParams, ...prowlarr.RequestEditorFn) error); ok {
		r1 = rf(ctx, params, reqEditors...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApiV1IndexerstatsWithResponse'
type ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call struct {
	*mock.Call
}

// GetApiV1IndexerstatsWithResponse is a helper method to define mock.On call
//   - ctx context.Context
//   - params *prowlarr.GetApiV1IndexerstatsParams
//   - reqEditors ...prowlarr.RequestEditorFn
func (_e *ProwlarrClient_Expecter) GetApiV1IndexerstatsWithResponse(ctx interface{}, params interface{}, reqEditors ...interface{}) *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call {
	return &ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call{Call: _e.mock.On("GetApiV1IndexerstatsWithResponse",
		append([]interface{}{ctx, params}, reqEditors...)...)}
}

func (_c *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call) Run(run func(ctx context.Context, params *prowlarr.GetApiV1IndexerstatsParams, reqEditors ...prowlarr.RequestEditorFn)) *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]prowlarr.RequestEditorFn, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(prowlarr.RequestEditorFn)
			}
		}
		run(args[0].(context.Context), args[1].(*prowlarr.GetApiV1IndexerstatsParams), variadicArgs...)
	})
	return _c
}

func (_c *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call) Return(_a0 *prowlarr.GetApiV1IndexerstatsResponse, _a1 error) *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call) RunAndReturn(run func(context.Context, *prowlarr.GetApiV1IndexerstatsParams, ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error)) *ProwlarrClient_GetApiV1IndexerstatsWithResponse_Call {
	_c.Call.Return(run)
	return _c
}

// NewProwlarrClient creates a new instance of ProwlarrClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProwlarrClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProwlarrClient {
	mock := &ProwlarrClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
