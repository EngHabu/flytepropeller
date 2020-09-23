package task

import (
	"testing"
	"time"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/event"

	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	pluginCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/io/mocks"
	"github.com/lyft/flytestdlib/storage"
	"github.com/stretchr/testify/assert"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/handler"
	handlerMocks "github.com/lyft/flytepropeller/pkg/controller/nodes/handler/mocks"
)

func TestToTaskEventPhase(t *testing.T) {
	assert.Equal(t, core.TaskExecution_UNDEFINED, ToTaskEventPhase(pluginCore.PhaseUndefined))
	assert.Equal(t, core.TaskExecution_SUCCEEDED, ToTaskEventPhase(pluginCore.PhaseSuccess))
	assert.Equal(t, core.TaskExecution_RUNNING, ToTaskEventPhase(pluginCore.PhaseRunning))
	assert.Equal(t, core.TaskExecution_FAILED, ToTaskEventPhase(pluginCore.PhasePermanentFailure))
	assert.Equal(t, core.TaskExecution_FAILED, ToTaskEventPhase(pluginCore.PhaseRetryableFailure))
	assert.Equal(t, core.TaskExecution_WAITING_FOR_RESOURCES, ToTaskEventPhase(pluginCore.PhaseWaitingForResources))
	assert.Equal(t, core.TaskExecution_INITIALIZING, ToTaskEventPhase(pluginCore.PhaseInitializing))
	assert.Equal(t, core.TaskExecution_UNDEFINED, ToTaskEventPhase(pluginCore.PhaseNotReady))
	assert.Equal(t, core.TaskExecution_QUEUED, ToTaskEventPhase(pluginCore.PhaseQueued))
}

func Test_trimErrorMessage(t *testing.T) {
	const inputStr = "0123456789"
	t.Run("Length less or equal than max", func(t *testing.T) {
		input := inputStr
		assert.Equal(t, input, trimErrorMessage(input, 10))
	})

	t.Run("Length > max", func(t *testing.T) {
		input := inputStr
		assert.Equal(t, "01236789", trimErrorMessage(input, 8))
	})

	t.Run("Odd Max", func(t *testing.T) {
		input := inputStr
		assert.Equal(t, "01236789", trimErrorMessage(input, 9))
	})

	t.Run("Odd input", func(t *testing.T) {
		input := "012345678"
		assert.Equal(t, "012345678", trimErrorMessage(input, 9))
	})
}

func TestToTaskExecutionEvent(t *testing.T) {
	tkID := &core.Identifier{}
	nodeID := &core.NodeExecutionIdentifier{}
	id := &core.TaskExecutionIdentifier{
		TaskId:          tkID,
		NodeExecutionId: nodeID,
	}
	n := time.Now()
	np, _ := ptypes.TimestampProto(n)

	in := &mocks.InputFilePaths{}
	const inputPath = "in"
	in.On("GetInputPath").Return(storage.DataReference(inputPath))

	out := &mocks.OutputFilePaths{}
	const outputPath = "out"
	out.On("GetOutputPath").Return(storage.DataReference(outputPath))

	nodeExecutionMetadata := handlerMocks.NodeExecutionMetadata{}
	nodeExecutionMetadata.OnIsInterruptible().Return(true)
	tev, err := ToTaskExecutionEvent(id, in, out, pluginCore.PhaseInfoWaitingForResources(n, 0, "reason"),
		&nodeExecutionMetadata)
	assert.NoError(t, err)
	assert.Nil(t, tev.Logs)
	assert.Equal(t, core.TaskExecution_WAITING_FOR_RESOURCES, tev.Phase)
	assert.Equal(t, uint32(0), tev.PhaseVersion)
	assert.Equal(t, np, tev.OccurredAt)
	assert.Equal(t, tkID, tev.TaskId)
	assert.Equal(t, nodeID, tev.ParentNodeExecutionId)
	assert.Equal(t, inputPath, tev.InputUri)
	assert.Nil(t, tev.OutputResult)
	assert.Equal(t, event.TaskExecutionMetadata_INTERRUPTIBLE, tev.Metadata.InstanceType)

	l := []*core.TaskLog{
		{Uri: "x", Name: "y", MessageFormat: core.TaskLog_JSON},
	}
	c := &structpb.Struct{}
	tev, err = ToTaskExecutionEvent(id, in, out, pluginCore.PhaseInfoRunning(1, &pluginCore.TaskInfo{
		OccurredAt: &n,
		Logs:       l,
		CustomInfo: c,
	}), &nodeExecutionMetadata)
	assert.NoError(t, err)
	assert.Equal(t, core.TaskExecution_RUNNING, tev.Phase)
	assert.Equal(t, uint32(1), tev.PhaseVersion)
	assert.Equal(t, l, tev.Logs)
	assert.Equal(t, c, tev.CustomInfo)
	assert.Equal(t, np, tev.OccurredAt)
	assert.Equal(t, tkID, tev.TaskId)
	assert.Equal(t, nodeID, tev.ParentNodeExecutionId)
	assert.Equal(t, inputPath, tev.InputUri)
	assert.Nil(t, tev.OutputResult)
	assert.Equal(t, event.TaskExecutionMetadata_INTERRUPTIBLE, tev.Metadata.InstanceType)

	defaultNodeExecutionMetadata := handlerMocks.NodeExecutionMetadata{}
	defaultNodeExecutionMetadata.OnIsInterruptible().Return(false)
	tev, err = ToTaskExecutionEvent(id, in, out, pluginCore.PhaseInfoSuccess(&pluginCore.TaskInfo{
		OccurredAt: &n,
		Logs:       l,
		CustomInfo: c,
	}), &defaultNodeExecutionMetadata)
	assert.NoError(t, err)
	assert.Equal(t, core.TaskExecution_SUCCEEDED, tev.Phase)
	assert.Equal(t, uint32(0), tev.PhaseVersion)
	assert.Equal(t, l, tev.Logs)
	assert.Equal(t, c, tev.CustomInfo)
	assert.Equal(t, np, tev.OccurredAt)
	assert.Equal(t, np, tev.OccurredAt)
	assert.Equal(t, tkID, tev.TaskId)
	assert.Equal(t, nodeID, tev.ParentNodeExecutionId)
	assert.NotNil(t, tev.OutputResult)
	assert.Equal(t, inputPath, tev.InputUri)
	assert.Equal(t, outputPath, tev.GetOutputUri())
	assert.Empty(t, event.TaskExecutionMetadata_DEFAULT, tev.Metadata.InstanceType)
}

func TestToTransitionType(t *testing.T) {
	assert.Equal(t, handler.TransitionTypeEphemeral, ToTransitionType(pluginCore.TransitionTypeEphemeral))
	assert.Equal(t, handler.TransitionTypeBarrier, ToTransitionType(pluginCore.TransitionTypeBarrier))
}
