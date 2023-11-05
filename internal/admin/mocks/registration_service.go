// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	cognitoidentityprovider "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	mock "github.com/stretchr/testify/mock"

	model "github.com/devpies/saas-core/internal/admin/model"
)

// RegistrationService is an autogenerated mock type for the registrationService type
type RegistrationService struct {
	mock.Mock
}

// RegisterTenant provides a mock function with given fields: ctx, tenant
func (_m *RegistrationService) RegisterTenant(ctx context.Context, tenant model.NewTenant) (int, error) {
	ret := _m.Called(ctx, tenant)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.NewTenant) (int, error)); ok {
		return rf(ctx, tenant)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.NewTenant) int); ok {
		r0 = rf(ctx, tenant)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.NewTenant) error); ok {
		r1 = rf(ctx, tenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResendTemporaryPassword provides a mock function with given fields: ctx, username
func (_m *RegistrationService) ResendTemporaryPassword(ctx context.Context, username string) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	ret := _m.Called(ctx, username)

	var r0 *cognitoidentityprovider.AdminCreateUserOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*cognitoidentityprovider.AdminCreateUserOutput, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *cognitoidentityprovider.AdminCreateUserOutput); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cognitoidentityprovider.AdminCreateUserOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRegistrationService creates a new instance of RegistrationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRegistrationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RegistrationService {
	mock := &RegistrationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
