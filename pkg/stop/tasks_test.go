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
		title       string
		group       string
		taskGroup   string
		wantMatched bool
	}{
		{"complete match", "group", "group", true},
		{"prefix match", "group", "group1", false},
		{"different case", "group", "Group", false},
		{"empty", "group", "", false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
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
		title       string
		prefix      string
		taskGroup   string
		wantMatched bool
	}{
		{"complete match", "group", "group", true},
		{"prefix match", "group", "group1", true},
		{"different case", "group", "Group", false},
		{"suffix match", "group", "_group", false},
		{"empty", "group", "", false},
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
		title       string
		taskGroup   string
		wantMatched bool
	}{
		{"service task", "service:", false},
		{"suffix is service task", "xxx:service:", true},
		{"not service task", "group", true},
		{"empty", "", true},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			o := &taskOptions{standalone: true}
			filtered := o.filterByGroup([]types.Task{{Group: &tc.taskGroup}})
			assert.Equal(t, tc.wantMatched, len(filtered) > 0)
		})
	}
}
