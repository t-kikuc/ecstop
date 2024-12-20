package stop

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterRunning(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title        string
		desiredCount int32
		runningCount int32
		wantRunning  bool
	}{
		{"explicitly running", 1, 1, true},
		{"terminated or launching", 1, 0, true},
		{"to be drained", 0, 1, true},
		{"already stopped", 0, 0, false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			running := filterRunning([]types.Service{{
				DesiredCount: tc.desiredCount,
				RunningCount: tc.runningCount,
			}})
			assert.Equal(t, tc.wantRunning, len(running) > 0)
		})
	}
}
