// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SetStatusCodeFunc is an autogenerated mock type for the setStatusCodeFunc type
type SetStatusCodeFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, statusCode
func (_m *SetStatusCodeFunc) Execute(ctx context.Context, statusCode int) {
	_m.Called(ctx, statusCode)
}

type NewSetStatusCodeFuncT interface {
	mock.TestingT
	Cleanup(func())
}

// NewSetStatusCodeFunc creates a new instance of SetStatusCodeFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSetStatusCodeFunc(t NewSetStatusCodeFuncT) *SetStatusCodeFunc {
	mock := &SetStatusCodeFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
