package stop

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestNewStopTaskCommand(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title   string
		flags   map[string]string
		wantErr bool
	}{
		{
			title: "valid flags",
			flags: map[string]string{
				"cluster": "xxx-cluster",
				"group":   "xxx-group",
			},
			wantErr: false,
		},
		{
			title:   "no flag causes error of OneRequired",
			flags:   map[string]string{},
			wantErr: true,
		},
		{
			title: "all flags causes error of MutuallyExclusive",
			flags: map[string]string{
				"group":        "xxx-group",
				"group-prefix": "xxx-group-prefix",
				"standalone":   "true",
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			cmd := NewStopTaskCommand()
			for k, v := range tc.flags {
				cmd.Flags().Set(k, v)
			}
			err1 := cmd.ValidateRequiredFlags()
			err2 := cmd.ValidateFlagGroups()

			assert.Equal(t, tc.wantErr, err1 != nil || err2 != nil)
		})
	}
}

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
		{"group", "_group", false},
		{"group", "", false},
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
