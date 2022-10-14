// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/devpies/saas-core/internal/project/model"

	time "time"
)

// ColumnService is an autogenerated mock type for the columnService type
type ColumnService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, nc, now
func (_m *ColumnService) Create(ctx context.Context, nc model.NewColumn, now time.Time) (model.Column, error) {
	ret := _m.Called(ctx, nc, now)

	var r0 model.Column
	if rf, ok := ret.Get(0).(func(context.Context, model.NewColumn, time.Time) model.Column); ok {
		r0 = rf(ctx, nc, now)
	} else {
		r0 = ret.Get(0).(model.Column)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.NewColumn, time.Time) error); ok {
		r1 = rf(ctx, nc, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateColumns provides a mock function with given fields: ctx, pid, now
func (_m *ColumnService) CreateColumns(ctx context.Context, pid string, now time.Time) error {
	ret := _m.Called(ctx, pid, now)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, pid, now)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, columnID
func (_m *ColumnService) Delete(ctx context.Context, columnID string) error {
	ret := _m.Called(ctx, columnID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, columnID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, projectID
func (_m *ColumnService) List(ctx context.Context, projectID string) ([]model.Column, error) {
	ret := _m.Called(ctx, projectID)

	var r0 []model.Column
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Column); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Column)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Retrieve provides a mock function with given fields: ctx, columnID
func (_m *ColumnService) Retrieve(ctx context.Context, columnID string) (model.Column, error) {
	ret := _m.Called(ctx, columnID)

	var r0 model.Column
	if rf, ok := ret.Get(0).(func(context.Context, string) model.Column); ok {
		r0 = rf(ctx, columnID)
	} else {
		r0 = ret.Get(0).(model.Column)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, columnID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, columnID, update, now
func (_m *ColumnService) Update(ctx context.Context, columnID string, update model.UpdateColumn, now time.Time) (model.Column, error) {
	ret := _m.Called(ctx, columnID, update, now)

	var r0 model.Column
	if rf, ok := ret.Get(0).(func(context.Context, string, model.UpdateColumn, time.Time) model.Column); ok {
		r0 = rf(ctx, columnID, update, now)
	} else {
		r0 = ret.Get(0).(model.Column)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.UpdateColumn, time.Time) error); ok {
		r1 = rf(ctx, columnID, update, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewColumnServiceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewColumnService creates a new instance of ColumnService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewColumnService(t NewColumnServiceT) *ColumnService {
	mock := &ColumnService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}