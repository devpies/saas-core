// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/devpies/saas-core/internal/admin/model"
)

// TenantService is an autogenerated mock type for the tenantService type
type TenantService struct {
	mock.Mock
}

// CancelSubscription provides a mock function with given fields: ctx, subID
func (_m *TenantService) CancelSubscription(ctx context.Context, subID string) (int, error) {
	ret := _m.Called(ctx, subID)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, subID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, subID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, subID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionInfo provides a mock function with given fields: ctx, tenantID
func (_m *TenantService) GetSubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, int, error) {
	ret := _m.Called(ctx, tenantID)

	var r0 model.SubscriptionInfo
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.SubscriptionInfo, int, error)); ok {
		return rf(ctx, tenantID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.SubscriptionInfo); ok {
		r0 = rf(ctx, tenantID)
	} else {
		r0 = ret.Get(0).(model.SubscriptionInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) int); ok {
		r1 = rf(ctx, tenantID)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, tenantID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ListTenants provides a mock function with given fields: ctx
func (_m *TenantService) ListTenants(ctx context.Context) ([]model.Tenant, int, error) {
	ret := _m.Called(ctx)

	var r0 []model.Tenant
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.Tenant, int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.Tenant); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Tenant)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) int); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RefundUser provides a mock function with given fields: ctx, subID
func (_m *TenantService) RefundUser(ctx context.Context, subID string) (int, error) {
	ret := _m.Called(ctx, subID)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, subID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, subID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, subID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTenantService creates a new instance of TenantService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTenantService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TenantService {
	mock := &TenantService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
