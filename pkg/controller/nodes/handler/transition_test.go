package handler

import (
	"testing"

	"github.com/lyft/flytestdlib/storage"
	"github.com/stretchr/testify/assert"
)

func TestDoTransition(t *testing.T) {
	t.Run("ephemeral", func(t *testing.T) {
		tr := DoTransition(TransitionTypeEphemeral, PhaseInfoQueued("queued"))
		assert.Equal(t, TransitionTypeEphemeral, tr.Type())
		assert.Equal(t, EPhaseQueued, tr.Info().Phase)
	})

	t.Run("barrier", func(t *testing.T) {
		tr := DoTransition(TransitionTypeBarrier, PhaseInfoSuccess(&ExecutionInfo{
			OutputInfo: &OutputInfo{OutputURI: "uri"},
		}))
		assert.Equal(t, TransitionTypeBarrier, tr.Type())
		assert.Equal(t, EPhaseSuccess, tr.Info().Phase)
		assert.Equal(t, storage.DataReference("uri"), tr.Info().Info.OutputInfo.OutputURI)
	})
}
