package v1alpha1

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MutableStruct struct {
	isDirty bool
}

func (m *MutableStruct) SetDirty() {
	m.isDirty = true
}

func (m MutableStruct) IsDirty() bool {
	return m.isDirty
}

type BranchNodeStatus struct {
	MutableStruct
	Phase           BranchNodePhase `json:"phase"`
	FinalizedNodeID *NodeID         `json:"finalNodeId"`
}

func (in *BranchNodeStatus) GetPhase() BranchNodePhase {
	return in.Phase
}

func (in *BranchNodeStatus) SetBranchNodeError() {
	in.SetDirty()
	in.Phase = BranchNodeError
}

func (in *BranchNodeStatus) SetBranchNodeSuccess(id NodeID) {
	in.SetDirty()
	in.Phase = BranchNodeSuccess
	in.FinalizedNodeID = &id
}

func (in *BranchNodeStatus) GetFinalizedNode() *NodeID {
	return in.FinalizedNodeID
}

func (in *BranchNodeStatus) Equals(other *BranchNodeStatus) bool {
	if in == nil && other == nil {
		return true
	}
	if in != nil && other != nil {
		phaseEqual := in.Phase == other.Phase
		if phaseEqual {
			if in.FinalizedNodeID == nil && other.FinalizedNodeID == nil {
				return true
			}
			if in.FinalizedNodeID != nil && other.FinalizedNodeID != nil {
				return *in.FinalizedNodeID == *other.FinalizedNodeID
			}
			return false
		}
		return false
	}
	return false
}

type DynamicNodePhase int

const (
	DynamicNodePhaseNone DynamicNodePhase = iota
	DynamicNodePhaseExecuting
	DynamicNodePhaseFailing
)

type DynamicNodeStatus struct {
	MutableStruct
	Phase  DynamicNodePhase `json:"phase"`
	Reason string           `json:"reason"`
}

func (in *DynamicNodeStatus) GetDynamicNodePhase() DynamicNodePhase {
	return in.Phase
}

func (in *DynamicNodeStatus) GetDynamicNodeReason() string {
	return in.Reason
}

func (in *DynamicNodeStatus) SetDynamicNodeReason(reason string) {
	in.SetDirty()
	in.Reason = reason
}

func (in *DynamicNodeStatus) SetDynamicNodePhase(phase DynamicNodePhase) {
	in.Phase = phase
}

func (in *DynamicNodeStatus) Equals(o *DynamicNodeStatus) bool {
	if in == nil && o == nil {
		return true
	}
	if in == nil || o == nil {
		return false
	}
	return in.Phase == o.Phase && in.Reason == o.Reason
}

type WorkflowNodePhase int

const (
	WorkflowNodePhaseUndefined WorkflowNodePhase = iota
	WorkflowNodePhaseExecuting
)

type WorkflowNodeStatus struct {
	MutableStruct
	Phase WorkflowNodePhase `json:"phase"`
}

func (in *WorkflowNodeStatus) GetWorkflowNodePhase() WorkflowNodePhase {
	return in.Phase
}

func (in *WorkflowNodeStatus) SetWorkflowNodePhase(phase WorkflowNodePhase) {
	in.Phase = phase
}

type NodeStatus struct {
	MutableStruct
	Phase                NodePhase     `json:"phase"`
	QueuedAt             *metav1.Time  `json:"queuedAt,omitempty"`
	StartedAt            *metav1.Time  `json:"startedAt,omitempty"`
	StoppedAt            *metav1.Time  `json:"stoppedAt,omitempty"`
	LastUpdatedAt        *metav1.Time  `json:"lastUpdatedAt,omitempty"`
	LastAttemptStartedAt *metav1.Time  `json:"laStartedAt,omitempty"`
	Message              string        `json:"message,omitempty"`
	DataDir              DataReference `json:"dataDir,omitempty"`
	Attempts             uint32        `json:"attempts"`
	Cached               bool          `json:"cached"`

	// This is useful only for branch nodes. If this is set, then it can be used to determine if execution can proceed
	ParentNode    *NodeID                  `json:"parentNode,omitempty"`
	ParentTask    *TaskExecutionIdentifier `json:"-,omitempty"`
	BranchStatus  *BranchNodeStatus        `json:"branchStatus,omitempty"`
	SubNodeStatus map[NodeID]*NodeStatus   `json:"subNodeStatus,omitempty"`
	// We can store the outputs at this layer

	// TODO not used delete
	WorkflowNodeStatus *WorkflowNodeStatus `json:"workflowNodeStatus,omitempty"`

	TaskNodeStatus    *TaskNodeStatus    `json:",omitempty"`
	DynamicNodeStatus *DynamicNodeStatus `json:"dynamicNodeStatus,omitempty"`
}

func (in *NodeStatus) IsDirty() bool {
	isDirty := in.MutableStruct.IsDirty() ||
		(in.TaskNodeStatus != nil && in.TaskNodeStatus.IsDirty()) ||
		(in.DynamicNodeStatus != nil && in.DynamicNodeStatus.IsDirty()) ||
		(in.WorkflowNodeStatus != nil && in.WorkflowNodeStatus.IsDirty()) ||
		(in.BranchStatus != nil && in.BranchStatus.IsDirty())
	if isDirty {
		return true
	}

	for _, sub := range in.SubNodeStatus {
		if sub.IsDirty() {
			return true
		}
	}

	return false
}

func (in *NodeStatus) GetBranchStatus() MutableBranchNodeStatus {
	if in.BranchStatus == nil {
		return nil
	}
	return in.BranchStatus
}

func (in *NodeStatus) GetWorkflowStatus() MutableWorkflowNodeStatus {
	if in.WorkflowNodeStatus == nil {
		return nil
	}
	return in.WorkflowNodeStatus
}

func (in *NodeStatus) GetTaskStatus() MutableTaskNodeStatus {
	if in.TaskNodeStatus == nil {
		return nil
	}
	return in.TaskNodeStatus
}

func (in NodeStatus) VisitNodeStatuses(visitor NodeStatusVisitFn) {
	for n, s := range in.SubNodeStatus {
		visitor(n, s)
	}
}

func (in NodeStatus) GetDynamicNodeStatus() MutableDynamicNodeStatus {
	if in.DynamicNodeStatus == nil {
		return nil
	}
	return in.DynamicNodeStatus
}

func (in *NodeStatus) ClearWorkflowStatus() {
	in.WorkflowNodeStatus = nil
	in.SetDirty()
}

func (in *NodeStatus) ClearTaskStatus() {
	in.TaskNodeStatus = nil
	in.SetDirty()
}

func (in *NodeStatus) ClearLastAttemptStartedAt() {
	in.LastAttemptStartedAt = nil
	in.SetDirty()
}

func (in *NodeStatus) GetLastUpdatedAt() *metav1.Time {
	return in.LastUpdatedAt
}

func (in *NodeStatus) GetLastAttemptStartedAt() *metav1.Time {
	return in.LastAttemptStartedAt
}

func (in *NodeStatus) GetAttempts() uint32 {
	return in.Attempts
}

func (in *NodeStatus) SetCached() {
	in.Cached = true
	in.SetDirty()
}

func (in *NodeStatus) IsCached() bool {
	return in.Cached
}

// ResetDirty is for unit tests, shouldn't be used in actual logic.
func (in *NodeStatus) ResetDirty() {
	in.MutableStruct.isDirty = false
}

func (in *NodeStatus) IncrementAttempts() uint32 {
	in.Attempts++
	in.SetDirty()
	return in.Attempts
}

func (in *NodeStatus) GetOrCreateDynamicNodeStatus() MutableDynamicNodeStatus {
	if in.DynamicNodeStatus == nil {
		in.SetDirty()
		in.DynamicNodeStatus = &DynamicNodeStatus{
			MutableStruct: MutableStruct{},
		}
	}

	return in.DynamicNodeStatus
}

func (in *NodeStatus) ClearDynamicNodeStatus() {
	in.DynamicNodeStatus = nil
	in.SetDirty()
}

func (in *NodeStatus) GetOrCreateBranchStatus() MutableBranchNodeStatus {
	if in.BranchStatus == nil {
		in.SetDirty()
		in.BranchStatus = &BranchNodeStatus{
			MutableStruct: MutableStruct{},
		}
	}

	return in.BranchStatus
}

func (in *NodeStatus) GetWorkflowNodeStatus() ExecutableWorkflowNodeStatus {
	if in.WorkflowNodeStatus == nil {
		return nil
	}

	return in.WorkflowNodeStatus
}

func (in *NodeStatus) GetPhase() NodePhase {
	return in.Phase
}

func (in *NodeStatus) GetMessage() string {
	return in.Message
}

func IsPhaseTerminal(phase NodePhase) bool {
	return phase == NodePhaseSucceeded || phase == NodePhaseFailed || phase == NodePhaseSkipped || phase == NodePhaseTimedOut
}

func (in *NodeStatus) GetOrCreateTaskStatus() MutableTaskNodeStatus {
	if in.TaskNodeStatus == nil {
		in.SetDirty()
		in.TaskNodeStatus = &TaskNodeStatus{
			MutableStruct: MutableStruct{},
		}
	}

	return in.TaskNodeStatus
}

func (in *NodeStatus) UpdatePhase(p NodePhase, occurredAt metav1.Time, reason string) {
	if in.Phase == p {
		// We will not update the phase multiple times. This prevents the comparison from returning false positive
		return
	}

	in.Phase = p
	in.Message = reason
	if len(reason) > maxMessageSize {
		in.Message = reason[:maxMessageSize]
	}

	n := occurredAt
	if occurredAt.IsZero() {
		n = metav1.Now()
	}

	if p == NodePhaseQueued && in.QueuedAt == nil {
		in.QueuedAt = &n
	} else if p == NodePhaseRunning {
		if in.StartedAt == nil {
			in.StartedAt = &n
		}
		if in.LastAttemptStartedAt == nil {
			in.LastAttemptStartedAt = &n
		}
	} else if IsPhaseTerminal(p) && in.StoppedAt == nil {
		if in.StartedAt == nil {
			in.StartedAt = &n
		}

		in.StoppedAt = &n
	}

	in.LastUpdatedAt = &n
	in.SetDirty()
}

func (in *NodeStatus) GetStartedAt() *metav1.Time {
	return in.StartedAt
}

func (in *NodeStatus) GetStoppedAt() *metav1.Time {
	return in.StoppedAt
}

func (in *NodeStatus) GetQueuedAt() *metav1.Time {
	return in.QueuedAt
}

func (in *NodeStatus) GetParentNodeID() *NodeID {
	return in.ParentNode
}

func (in *NodeStatus) GetParentTaskID() *core.TaskExecutionIdentifier {
	if in.ParentTask != nil {
		return in.ParentTask.TaskExecutionIdentifier
	}
	return nil
}

func (in *NodeStatus) SetParentNodeID(n *NodeID) {
	in.ParentNode = n
	in.SetDirty()
}

func (in *NodeStatus) SetParentTaskID(t *core.TaskExecutionIdentifier) {
	in.ParentTask = &TaskExecutionIdentifier{
		TaskExecutionIdentifier: t,
	}
	in.SetDirty()
}

func (in *NodeStatus) GetOrCreateWorkflowStatus() MutableWorkflowNodeStatus {
	if in.WorkflowNodeStatus == nil {
		in.SetDirty()
		in.WorkflowNodeStatus = &WorkflowNodeStatus{
			MutableStruct: MutableStruct{},
		}
	}

	return in.WorkflowNodeStatus
}

func (in NodeStatus) GetTaskNodeStatus() ExecutableTaskNodeStatus {
	// Explicitly return nil here to avoid a misleading non-nil interface.
	if in.TaskNodeStatus == nil {
		return nil
	}

	return in.TaskNodeStatus
}

func (in *NodeStatus) GetNodeExecutionStatus(id NodeID) ExecutableNodeStatus {
	n, ok := in.SubNodeStatus[id]
	if ok {
		return n
	}

	if in.SubNodeStatus == nil {
		in.SubNodeStatus = make(map[NodeID]*NodeStatus)
	}

	newNodeStatus := &NodeStatus{
		MutableStruct: MutableStruct{},
	}
	newNodeStatus.SetParentTaskID(in.GetParentTaskID())
	newNodeStatus.SetParentNodeID(in.GetParentNodeID())

	in.SubNodeStatus[id] = newNodeStatus
	in.SetDirty()
	return newNodeStatus
}

func (in *NodeStatus) IsTerminated() bool {
	return in.GetPhase() == NodePhaseFailed || in.GetPhase() == NodePhaseSkipped || in.GetPhase() == NodePhaseSucceeded
}

func (in *NodeStatus) GetDataDir() DataReference {
	return in.DataDir
}

func (in *NodeStatus) SetDataDir(d DataReference) {
	in.DataDir = d
	if in.DataDir != d {
		in.SetDirty()
	}
}

func (in *NodeStatus) Equals(other *NodeStatus) bool {
	// Assuming in is never nil
	if other == nil {
		return false
	}

	if in.Phase == other.Phase {
		if in.Phase == NodePhaseSucceeded || in.Phase == NodePhaseFailed {
			return true
		}
	}

	if in.Attempts != other.Attempts {
		return false
	}

	if in.Phase != other.Phase {
		return false
	}

	if !in.TaskNodeStatus.Equals(other.TaskNodeStatus) {
		return false
	}

	if in.DataDir != other.DataDir {
		return false
	}

	if in.ParentNode != nil && other.ParentNode != nil {
		if *in.ParentNode != *other.ParentNode {
			return false
		}
	} else if !(in.ParentNode == other.ParentNode) {
		// Both are not nil
		return false
	}

	if !reflect.DeepEqual(in.ParentTask, other.ParentTask) {
		return false
	}

	if len(in.SubNodeStatus) != len(other.SubNodeStatus) {
		return false
	}

	for k, v := range in.SubNodeStatus {
		otherV, ok := other.SubNodeStatus[k]
		if !ok {
			return false
		}
		if !v.Equals(otherV) {
			return false
		}
	}

	return in.BranchStatus.Equals(other.BranchStatus) && in.DynamicNodeStatus.Equals(other.DynamicNodeStatus)
}

// THIS IS NOT AUTO GENERATED
func (in *CustomState) DeepCopyInto(out *CustomState) {
	if in == nil || *in == nil {
		return
	}

	raw, err := json.Marshal(in)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, out)
	if err != nil {
		return
	}
}

func (in *CustomState) DeepCopy() *CustomState {
	if in == nil || *in == nil {
		return nil
	}

	out := &CustomState{}
	in.DeepCopyInto(out)
	return out
}

type TaskNodeStatus struct {
	MutableStruct
	Phase              int    `json:"phase,omitempty"`
	PhaseVersion       uint32 `json:"phaseVersion,omitempty"`
	PluginState        []byte `json:"pState,omitempty"`
	PluginStateVersion uint32 `json:"psv,omitempty"`
	BarrierClockTick   uint32 `json:"tick,omitempty"`
}

func (in *TaskNodeStatus) GetBarrierClockTick() uint32 {
	return in.BarrierClockTick
}

func (in *TaskNodeStatus) SetBarrierClockTick(tick uint32) {
	in.BarrierClockTick = tick
}

func (in *TaskNodeStatus) SetPluginState(s []byte) {
	in.PluginState = s
	in.SetDirty()
}

func (in *TaskNodeStatus) SetPluginStateVersion(v uint32) {
	in.PluginStateVersion = v
	in.SetDirty()
}

func (in *TaskNodeStatus) GetPluginState() []byte {
	return in.PluginState
}

func (in *TaskNodeStatus) GetPluginStateVersion() uint32 {
	return in.PluginStateVersion
}

func (in *TaskNodeStatus) SetPhase(phase int) {
	in.Phase = phase
	in.SetDirty()
}

func (in *TaskNodeStatus) SetPhaseVersion(version uint32) {
	in.PhaseVersion = version
	in.SetDirty()
}

func (in TaskNodeStatus) GetPhase() int {
	return in.Phase
}

func (in TaskNodeStatus) GetPhaseVersion() uint32 {
	return in.PhaseVersion
}

func (in *TaskNodeStatus) UpdatePhase(phase int, phaseVersion uint32) {
	in.Phase = phase
	in.PhaseVersion = phaseVersion
	in.SetDirty()
}

func (in *TaskNodeStatus) DeepCopyInto(out *TaskNodeStatus) {
	if in == nil {
		return
	}

	raw, err := json.Marshal(in)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, out)
	if err != nil {
		return
	}
}

func (in *TaskNodeStatus) DeepCopy() *TaskNodeStatus {
	if in == nil {
		return nil
	}

	out := &TaskNodeStatus{}
	in.DeepCopyInto(out)
	return out
}

func (in *TaskNodeStatus) Equals(other *TaskNodeStatus) bool {
	if in == nil && other == nil {
		return true
	}
	if in == nil || other == nil {
		return false
	}
	return in.Phase == other.Phase && in.PhaseVersion == other.PhaseVersion && in.PluginStateVersion == other.PluginStateVersion && bytes.Equal(in.PluginState, other.PluginState) && in.BarrierClockTick == other.BarrierClockTick
}
