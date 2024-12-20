package stop

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterRunning(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		desiredCount int32
		runningCount int32
		wantRunning  bool
	}{
		{1, 1, true},
		{1, 0, true},
		{0, 1, true},
		{0, 0, false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			running := filterRunning([]types.Service{{
				DesiredCount: tc.desiredCount,
				RunningCount: tc.runningCount,
			}})
			assert.Equal(t, tc.wantRunning, len(running) > 0)
		})
	}
}
