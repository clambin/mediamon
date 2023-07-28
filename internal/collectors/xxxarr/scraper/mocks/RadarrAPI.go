// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	xxxarr "github.com/clambin/mediaclients/xxxarr"
)

// RadarrAPI is an autogenerated mock type for the RadarrAPI type
type RadarrAPI struct {
	mock.Mock
}

// GetCalendar provides a mock function with given fields: ctx
func (_m *RadarrAPI) GetCalendar(ctx context.Context) ([]xxxarr.RadarrCalendarResponse, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.RadarrCalendarResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.RadarrCalendarResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.RadarrCalendarResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.RadarrCalendarResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHealth provides a mock function with given fields: ctx
func (_m *RadarrAPI) GetHealth(ctx context.Context) ([]xxxarr.RadarrHealthResponse, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.RadarrHealthResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.RadarrHealthResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.RadarrHealthResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.RadarrHealthResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMovieByID provides a mock function with given fields: ctx, movieID
func (_m *RadarrAPI) GetMovieByID(ctx context.Context, movieID int) (xxxarr.RadarrMovieResponse, error) {
	ret := _m.Called(ctx, movieID)

	var r0 xxxarr.RadarrMovieResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (xxxarr.RadarrMovieResponse, error)); ok {
		return rf(ctx, movieID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) xxxarr.RadarrMovieResponse); ok {
		r0 = rf(ctx, movieID)
	} else {
		r0 = ret.Get(0).(xxxarr.RadarrMovieResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, movieID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMovies provides a mock function with given fields: ctx
func (_m *RadarrAPI) GetMovies(ctx context.Context) ([]xxxarr.RadarrMovieResponse, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.RadarrMovieResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.RadarrMovieResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.RadarrMovieResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.RadarrMovieResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetQueue provides a mock function with given fields: ctx
func (_m *RadarrAPI) GetQueue(ctx context.Context) (xxxarr.RadarrQueueResponse, error) {
	ret := _m.Called(ctx)

	var r0 xxxarr.RadarrQueueResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (xxxarr.RadarrQueueResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) xxxarr.RadarrQueueResponse); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(xxxarr.RadarrQueueResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetQueuePage provides a mock function with given fields: ctx, pageNr
func (_m *RadarrAPI) GetQueuePage(ctx context.Context, pageNr int) (xxxarr.RadarrQueueResponse, error) {
	ret := _m.Called(ctx, pageNr)

	var r0 xxxarr.RadarrQueueResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (xxxarr.RadarrQueueResponse, error)); ok {
		return rf(ctx, pageNr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) xxxarr.RadarrQueueResponse); ok {
		r0 = rf(ctx, pageNr)
	} else {
		r0 = ret.Get(0).(xxxarr.RadarrQueueResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, pageNr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSystemStatus provides a mock function with given fields: ctx
func (_m *RadarrAPI) GetSystemStatus(ctx context.Context) (xxxarr.RadarrSystemStatusResponse, error) {
	ret := _m.Called(ctx)

	var r0 xxxarr.RadarrSystemStatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (xxxarr.RadarrSystemStatusResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) xxxarr.RadarrSystemStatusResponse); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(xxxarr.RadarrSystemStatusResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURL provides a mock function with given fields:
func (_m *RadarrAPI) GetURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewRadarrAPI creates a new instance of RadarrAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRadarrAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *RadarrAPI {
	mock := &RadarrAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
