package stop

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddClusterFlags(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title   string
		flags   map[string]string
		wantErr bool
	}{
		{
			title: "only cluster is valid",
			flags: map[string]string{
				"cluster": "xxx-cluster",
			},
			wantErr: false,
		},
		{
			title: "only all-clusters is valid",
			flags: map[string]string{
				"all-clusters": "true",
			},
		},
		{
			title:   "no flag causes error of OneRequired",
			flags:   map[string]string{},
			wantErr: true,
		},
		{
			title: "all flags causes error of MutuallyExclusive",
			flags: map[string]string{
				"cluster":      "xxx-cluster",
				"all-clusters": "true",
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			cmd := &cobra.Command{}
			addClusterFlags(cmd, &clusterOptions{})
			for k, v := range tc.flags {
				cmd.Flags().Set(k, v)
			}
			err1 := cmd.ValidateRequiredFlags()
			err2 := cmd.ValidateFlagGroups()

			assert.Equal(t, tc.wantErr, err1 != nil || err2 != nil)
		})
	}
}
