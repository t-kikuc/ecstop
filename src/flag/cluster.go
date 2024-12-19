package flag

import "github.com/spf13/cobra"

// AddClusterFlags adds flags `--cluster` and `--all-clusters` to the command
func AddClusterFlags(c *cobra.Command, clusterP *string, allClusterP *bool) {
	const (
		flag_cluster     = "cluster"
		flag_allClusters = "all-clusters"
	)

	c.Flags().StringVarP(clusterP, flag_cluster, "c", "", "Name or ARN of the cluster")
	c.Flags().BoolVarP(allClusterP, flag_allClusters, "a", false, "Stop in all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)
}
