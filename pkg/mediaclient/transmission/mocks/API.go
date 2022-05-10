// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	transmission "github.com/clambin/mediamon/pkg/mediaclient/transmission"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// GetSessionParameters provides a mock function with given fields: ctx
func (_m *API) GetSessionParameters(ctx context.Context) (transmission.SessionParameters, error) {
	ret := _m.Called(ctx)

	var r0 transmission.SessionParameters
	if rf, ok := ret.Get(0).(func(context.Context) transmission.SessionParameters); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmission.SessionParameters)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSessionStatistics provides a mock function with given fields: ctx
func (_m *API) GetSessionStatistics(ctx context.Context) (transmission.SessionStats, error) {
	ret := _m.Called(ctx)

	var r0 transmission.SessionStats
	if rf, ok := ret.Get(0).(func(context.Context) transmission.SessionStats); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmission.SessionStats)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAPI creates a new instance of API. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t testing.TB) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
