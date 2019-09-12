// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import core "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
import executors "github.com/lyft/flytepropeller/pkg/controller/executors"
import mock "github.com/stretchr/testify/mock"
import v1alpha1 "github.com/lyft/flytepropeller/pkg/apis/flyteworkflow/v1alpha1"

// Node is an autogenerated mock type for the Node type
type Node struct {
	mock.Mock
}

// AbortHandler provides a mock function with given fields: ctx, w, currentNode
func (_m *Node) AbortHandler(ctx context.Context, w v1alpha1.ExecutableWorkflow, currentNode v1alpha1.ExecutableNode) error {
	ret := _m.Called(ctx, w, currentNode)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, v1alpha1.ExecutableWorkflow, v1alpha1.ExecutableNode) error); ok {
		r0 = rf(ctx, w, currentNode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Initialize provides a mock function with given fields: ctx
func (_m *Node) Initialize(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecursiveNodeHandler provides a mock function with given fields: ctx, w, currentNode
func (_m *Node) RecursiveNodeHandler(ctx context.Context, w v1alpha1.ExecutableWorkflow, currentNode v1alpha1.ExecutableNode) (executors.NodeStatus, error) {
	ret := _m.Called(ctx, w, currentNode)

	var r0 executors.NodeStatus
	if rf, ok := ret.Get(0).(func(context.Context, v1alpha1.ExecutableWorkflow, v1alpha1.ExecutableNode) executors.NodeStatus); ok {
		r0 = rf(ctx, w, currentNode)
	} else {
		r0 = ret.Get(0).(executors.NodeStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, v1alpha1.ExecutableWorkflow, v1alpha1.ExecutableNode) error); ok {
		r1 = rf(ctx, w, currentNode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetInputsForStartNode provides a mock function with given fields: ctx, w, inputs
func (_m *Node) SetInputsForStartNode(ctx context.Context, w v1alpha1.BaseWorkflowWithStatus, inputs *core.LiteralMap) (executors.NodeStatus, error) {
	ret := _m.Called(ctx, w, inputs)

	var r0 executors.NodeStatus
	if rf, ok := ret.Get(0).(func(context.Context, v1alpha1.BaseWorkflowWithStatus, *core.LiteralMap) executors.NodeStatus); ok {
		r0 = rf(ctx, w, inputs)
	} else {
		r0 = ret.Get(0).(executors.NodeStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, v1alpha1.BaseWorkflowWithStatus, *core.LiteralMap) error); ok {
		r1 = rf(ctx, w, inputs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}