// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	xxxarr "github.com/clambin/mediaclients/xxxarr"
)

// RadarrGetter is an autogenerated mock type for the RadarrGetter type
type RadarrGetter struct {
	mock.Mock
}

type RadarrGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *RadarrGetter) EXPECT() *RadarrGetter_Expecter {
	return &RadarrGetter_Expecter{mock: &_m.Mock}
}

// GetCalendar provides a mock function with given fields: ctx
func (_m *RadarrGetter) GetCalendar(ctx context.Context) ([]xxxarr.RadarrCalendarResponse, error) {
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

// RadarrGetter_GetCalendar_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCalendar'
type RadarrGetter_GetCalendar_Call struct {
	*mock.Call
}

// GetCalendar is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RadarrGetter_Expecter) GetCalendar(ctx interface{}) *RadarrGetter_GetCalendar_Call {
	return &RadarrGetter_GetCalendar_Call{Call: _e.mock.On("GetCalendar", ctx)}
}

func (_c *RadarrGetter_GetCalendar_Call) Run(run func(ctx context.Context)) *RadarrGetter_GetCalendar_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RadarrGetter_GetCalendar_Call) Return(response []xxxarr.RadarrCalendarResponse, err error) *RadarrGetter_GetCalendar_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetCalendar_Call) RunAndReturn(run func(context.Context) ([]xxxarr.RadarrCalendarResponse, error)) *RadarrGetter_GetCalendar_Call {
	_c.Call.Return(run)
	return _c
}

// GetHealth provides a mock function with given fields: ctx
func (_m *RadarrGetter) GetHealth(ctx context.Context) ([]xxxarr.RadarrHealthResponse, error) {
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

// RadarrGetter_GetHealth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHealth'
type RadarrGetter_GetHealth_Call struct {
	*mock.Call
}

// GetHealth is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RadarrGetter_Expecter) GetHealth(ctx interface{}) *RadarrGetter_GetHealth_Call {
	return &RadarrGetter_GetHealth_Call{Call: _e.mock.On("GetHealth", ctx)}
}

func (_c *RadarrGetter_GetHealth_Call) Run(run func(ctx context.Context)) *RadarrGetter_GetHealth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RadarrGetter_GetHealth_Call) Return(response []xxxarr.RadarrHealthResponse, err error) *RadarrGetter_GetHealth_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetHealth_Call) RunAndReturn(run func(context.Context) ([]xxxarr.RadarrHealthResponse, error)) *RadarrGetter_GetHealth_Call {
	_c.Call.Return(run)
	return _c
}

// GetMovieByID provides a mock function with given fields: ctx, movieID
func (_m *RadarrGetter) GetMovieByID(ctx context.Context, movieID int) (xxxarr.RadarrMovieResponse, error) {
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

// RadarrGetter_GetMovieByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMovieByID'
type RadarrGetter_GetMovieByID_Call struct {
	*mock.Call
}

// GetMovieByID is a helper method to define mock.On call
//   - ctx context.Context
//   - movieID int
func (_e *RadarrGetter_Expecter) GetMovieByID(ctx interface{}, movieID interface{}) *RadarrGetter_GetMovieByID_Call {
	return &RadarrGetter_GetMovieByID_Call{Call: _e.mock.On("GetMovieByID", ctx, movieID)}
}

func (_c *RadarrGetter_GetMovieByID_Call) Run(run func(ctx context.Context, movieID int)) *RadarrGetter_GetMovieByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *RadarrGetter_GetMovieByID_Call) Return(response xxxarr.RadarrMovieResponse, err error) *RadarrGetter_GetMovieByID_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetMovieByID_Call) RunAndReturn(run func(context.Context, int) (xxxarr.RadarrMovieResponse, error)) *RadarrGetter_GetMovieByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetMovies provides a mock function with given fields: ctx
func (_m *RadarrGetter) GetMovies(ctx context.Context) ([]xxxarr.RadarrMovieResponse, error) {
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

// RadarrGetter_GetMovies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMovies'
type RadarrGetter_GetMovies_Call struct {
	*mock.Call
}

// GetMovies is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RadarrGetter_Expecter) GetMovies(ctx interface{}) *RadarrGetter_GetMovies_Call {
	return &RadarrGetter_GetMovies_Call{Call: _e.mock.On("GetMovies", ctx)}
}

func (_c *RadarrGetter_GetMovies_Call) Run(run func(ctx context.Context)) *RadarrGetter_GetMovies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RadarrGetter_GetMovies_Call) Return(response []xxxarr.RadarrMovieResponse, err error) *RadarrGetter_GetMovies_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetMovies_Call) RunAndReturn(run func(context.Context) ([]xxxarr.RadarrMovieResponse, error)) *RadarrGetter_GetMovies_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueue provides a mock function with given fields: ctx
func (_m *RadarrGetter) GetQueue(ctx context.Context) (xxxarr.RadarrQueueResponse, error) {
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

// RadarrGetter_GetQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueue'
type RadarrGetter_GetQueue_Call struct {
	*mock.Call
}

// GetQueue is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RadarrGetter_Expecter) GetQueue(ctx interface{}) *RadarrGetter_GetQueue_Call {
	return &RadarrGetter_GetQueue_Call{Call: _e.mock.On("GetQueue", ctx)}
}

func (_c *RadarrGetter_GetQueue_Call) Run(run func(ctx context.Context)) *RadarrGetter_GetQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RadarrGetter_GetQueue_Call) Return(response xxxarr.RadarrQueueResponse, err error) *RadarrGetter_GetQueue_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetQueue_Call) RunAndReturn(run func(context.Context) (xxxarr.RadarrQueueResponse, error)) *RadarrGetter_GetQueue_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueuePage provides a mock function with given fields: ctx, pageNr
func (_m *RadarrGetter) GetQueuePage(ctx context.Context, pageNr int) (xxxarr.RadarrQueueResponse, error) {
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

// RadarrGetter_GetQueuePage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueuePage'
type RadarrGetter_GetQueuePage_Call struct {
	*mock.Call
}

// GetQueuePage is a helper method to define mock.On call
//   - ctx context.Context
//   - pageNr int
func (_e *RadarrGetter_Expecter) GetQueuePage(ctx interface{}, pageNr interface{}) *RadarrGetter_GetQueuePage_Call {
	return &RadarrGetter_GetQueuePage_Call{Call: _e.mock.On("GetQueuePage", ctx, pageNr)}
}

func (_c *RadarrGetter_GetQueuePage_Call) Run(run func(ctx context.Context, pageNr int)) *RadarrGetter_GetQueuePage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *RadarrGetter_GetQueuePage_Call) Return(response xxxarr.RadarrQueueResponse, err error) *RadarrGetter_GetQueuePage_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetQueuePage_Call) RunAndReturn(run func(context.Context, int) (xxxarr.RadarrQueueResponse, error)) *RadarrGetter_GetQueuePage_Call {
	_c.Call.Return(run)
	return _c
}

// GetSystemStatus provides a mock function with given fields: ctx
func (_m *RadarrGetter) GetSystemStatus(ctx context.Context) (xxxarr.RadarrSystemStatusResponse, error) {
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

// RadarrGetter_GetSystemStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSystemStatus'
type RadarrGetter_GetSystemStatus_Call struct {
	*mock.Call
}

// GetSystemStatus is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RadarrGetter_Expecter) GetSystemStatus(ctx interface{}) *RadarrGetter_GetSystemStatus_Call {
	return &RadarrGetter_GetSystemStatus_Call{Call: _e.mock.On("GetSystemStatus", ctx)}
}

func (_c *RadarrGetter_GetSystemStatus_Call) Run(run func(ctx context.Context)) *RadarrGetter_GetSystemStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RadarrGetter_GetSystemStatus_Call) Return(response xxxarr.RadarrSystemStatusResponse, err error) *RadarrGetter_GetSystemStatus_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *RadarrGetter_GetSystemStatus_Call) RunAndReturn(run func(context.Context) (xxxarr.RadarrSystemStatusResponse, error)) *RadarrGetter_GetSystemStatus_Call {
	_c.Call.Return(run)
	return _c
}

// GetURL provides a mock function with given fields:
func (_m *RadarrGetter) GetURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RadarrGetter_GetURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURL'
type RadarrGetter_GetURL_Call struct {
	*mock.Call
}

// GetURL is a helper method to define mock.On call
func (_e *RadarrGetter_Expecter) GetURL() *RadarrGetter_GetURL_Call {
	return &RadarrGetter_GetURL_Call{Call: _e.mock.On("GetURL")}
}

func (_c *RadarrGetter_GetURL_Call) Run(run func()) *RadarrGetter_GetURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RadarrGetter_GetURL_Call) Return(url string) *RadarrGetter_GetURL_Call {
	_c.Call.Return(url)
	return _c
}

func (_c *RadarrGetter_GetURL_Call) RunAndReturn(run func() string) *RadarrGetter_GetURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewRadarrGetter creates a new instance of RadarrGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRadarrGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *RadarrGetter {
	mock := &RadarrGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
