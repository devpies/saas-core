// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SessionManager is an autogenerated mock type for the sessionManager type
type SessionManager struct {
	mock.Mock
}

// Destroy provides a mock function with given fields: ctx
func (_m *SessionManager) Destroy(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenewToken provides a mock function with given fields: ctx
func (_m *SessionManager) RenewToken(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewSessionManagerT interface {
	mock.TestingT
	Cleanup(func())
}

// NewSessionManager creates a new instance of SessionManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSessionManager(t NewSessionManagerT) *SessionManager {
	mock := &SessionManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
