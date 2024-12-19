package stop

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterByGroup_match(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		group       string
		taskGroup   string
		wantMatched bool
	}{
		{"group", "group", true},
		{"group", "group1", false},
		{"group", "Group", false},
		{"group", "", false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			o := &taskOptions{group: tc.group}
			filtered := o.filterByGroup([]types.Task{{Group: &tc.taskGroup}})
			assert.Equal(t, tc.wantMatched, len(filtered) > 0)
		})
	}
}

func TestFilterByGroup_prefix(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		prefix      string
		taskGroup   string
		wantMatched bool
	}{
		{"group", "group", true},
		{"group", "group1", true},
		{"group", "Group", false},
		{"group", "", false},
		{"group", "_group", false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			o := &taskOptions{groupPrefix: tc.prefix}
			filtered := o.filterByGroup([]types.Task{{Group: &tc.taskGroup}})
			assert.Equal(t, tc.wantMatched, len(filtered) > 0)
		})
	}
}

func TestFilterByGroup_standalone(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		taskGroup   string
		wantMatched bool
	}{
		{"service:", false},
		{"xxx:service:", true},
		{"group", true},
		{"", true},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			o := &taskOptions{standalone: true}
			filtered := o.filterByGroup([]types.Task{{Group: &tc.taskGroup}})
			assert.Equal(t, tc.wantMatched, len(filtered) > 0)
		})
	}
}
