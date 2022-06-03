// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StatusSetter is an autogenerated mock type for the statusSetter type
type StatusSetter struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, statusCode
func (_m *StatusSetter) Execute(ctx context.Context, statusCode int) {
	_m.Called(ctx, statusCode)
}

type NewStatusSetterT interface {
	mock.TestingT
	Cleanup(func())
}

// NewStatusSetter creates a new instance of StatusSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStatusSetter(t NewStatusSetterT) *StatusSetter {
	mock := &StatusSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
