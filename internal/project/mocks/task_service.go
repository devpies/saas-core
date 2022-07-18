// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/devpies/saas-core/internal/project/model"

	time "time"
)

// TaskService is an autogenerated mock type for the taskService type
type TaskService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, task, projectID, now
func (_m *TaskService) Create(ctx context.Context, task model.NewTask, projectID string, now time.Time) (model.Task, error) {
	ret := _m.Called(ctx, task, projectID, now)

	var r0 model.Task
	if rf, ok := ret.Get(0).(func(context.Context, model.NewTask, string, time.Time) model.Task); ok {
		r0 = rf(ctx, task, projectID, now)
	} else {
		r0 = ret.Get(0).(model.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.NewTask, string, time.Time) error); ok {
		r1 = rf(ctx, task, projectID, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, taskID
func (_m *TaskService) Delete(ctx context.Context, taskID string) error {
	ret := _m.Called(ctx, taskID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, taskID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, projectID
func (_m *TaskService) List(ctx context.Context, projectID string) ([]model.Task, error) {
	ret := _m.Called(ctx, projectID)

	var r0 []model.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Task); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Task)
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

// Retrieve provides a mock function with given fields: ctx, taskID
func (_m *TaskService) Retrieve(ctx context.Context, taskID string) (model.Task, error) {
	ret := _m.Called(ctx, taskID)

	var r0 model.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) model.Task); ok {
		r0 = rf(ctx, taskID)
	} else {
		r0 = ret.Get(0).(model.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, taskID, update, now
func (_m *TaskService) Update(ctx context.Context, taskID string, update model.UpdateTask, now time.Time) (model.Task, error) {
	ret := _m.Called(ctx, taskID, update, now)

	var r0 model.Task
	if rf, ok := ret.Get(0).(func(context.Context, string, model.UpdateTask, time.Time) model.Task); ok {
		r0 = rf(ctx, taskID, update, now)
	} else {
		r0 = ret.Get(0).(model.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.UpdateTask, time.Time) error); ok {
		r1 = rf(ctx, taskID, update, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewTaskServiceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewTaskService creates a new instance of TaskService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTaskService(t NewTaskServiceT) *TaskService {
	mock := &TaskService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
