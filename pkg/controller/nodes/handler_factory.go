package nodes

import (
	"context"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/dynamic"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/catalog"

	"github.com/lyft/flytestdlib/promutils"

	"github.com/pkg/errors"

	"github.com/lyft/flytepropeller/pkg/apis/flyteworkflow/v1alpha1"
	"github.com/lyft/flytepropeller/pkg/controller/executors"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/branch"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/end"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/handler"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/start"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/subworkflow"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/subworkflow/launchplan"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/task"
)

//go:generate mockery -all -case=underscore

type HandlerFactory interface {
	GetHandler(kind v1alpha1.NodeKind) (handler.Node, error)
	Setup(ctx context.Context, setup handler.SetupContext) error
}

type handlerFactory struct {
	handlers map[v1alpha1.NodeKind]handler.Node
}

func (f handlerFactory) GetHandler(kind v1alpha1.NodeKind) (handler.Node, error) {
	h, ok := f.handlers[kind]
	if !ok {
		return nil, errors.Errorf("Handler not registered for NodeKind [%v]", kind)
	}
	return h, nil
}

func (f handlerFactory) Setup(ctx context.Context, setup handler.SetupContext) error {
	for _, v := range f.handlers {
		if err := v.Setup(ctx, setup); err != nil {
			return err
		}
	}
	return nil
}

func NewHandlerFactory(ctx context.Context, executor executors.Node, workflowLauncher launchplan.Executor, kubeClient executors.Client, client catalog.Client, scope promutils.Scope) (HandlerFactory, error) {

	f := &handlerFactory{
		handlers: map[v1alpha1.NodeKind]handler.Node{
			v1alpha1.NodeKindBranch: branch.New(executor, scope),
			v1alpha1.NodeKindTask: dynamic.New(
				task.New(ctx, kubeClient, client, scope),
				executor,
				scope),
			v1alpha1.NodeKindWorkflow: subworkflow.New(executor, workflowLauncher, scope),
			v1alpha1.NodeKindStart:    start.New(),
			v1alpha1.NodeKindEnd:      end.New(),
		},
	}

	return f, nil
}
