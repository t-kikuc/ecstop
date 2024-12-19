package stop

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/client"
)

type clusterOptions struct {
	cluster     string
	allClusters bool
}

// addClusterFlags adds flags `--cluster` and `--all-clusters` to the command
func addClusterFlags(c *cobra.Command, clusterP *clusterOptions) {
	const (
		flag_cluster     = "cluster"
		flag_allClusters = "all-clusters"
	)

	c.Flags().StringVarP(&clusterP.cluster, flag_cluster, "c", "", "Name or ARN of the cluster")
	c.Flags().BoolVarP(&clusterP.allClusters, flag_allClusters, "a", false, "Stop in all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)
}

func (co clusterOptions) DecideClusters(ctx context.Context, cli *client.ECSClient) ([]string, error) {
	if co.cluster != "" {
		return []string{co.cluster}, nil
	}

	// Since at least on of `--cluster` or `--all-clusters` is required, we can assume that `--all-clusters` is true
	clusters, err := cli.ListClusters(ctx)
	if err != nil {
		return nil, err
	}
	if len(clusters) == 0 {
		return nil, nil
	}
	return clusters, nil
}
