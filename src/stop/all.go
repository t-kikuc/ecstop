package stop

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type allOptions struct {
	cluster     string
	allClusters bool
}

func NewStopAllCommand() *cobra.Command {
	o := &allOptions{}

	c := &cobra.Command{
		Use:   "all",
		Short: "Stop ECS Services, Standalone Tasks", // and Container Instances",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stopAll(context.Background())
		},
	}

	const (
		flag_cluster     = "cluster"
		flag_allClusters = "all-clusters"
	)

	c.Flags().StringVar(&o.cluster, flag_cluster, "", "Name or ARN of the cluster to stop resources")
	c.Flags().BoolVar(&o.allClusters, flag_allClusters, false, "Stop resources of all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)

	return c
}

func (o *allOptions) stopAll(ctx context.Context) error {
	fmt.Println("[1] Start stopping ECS Services")
	srvOpts := &serviceOptions{
		cluster:     o.cluster,
		allClusters: o.allClusters,
	}
	if err := srvOpts.scaleinServices(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Services: %w", err)
	}
	fmt.Println("[1] Successfully finished stopping ECS Services")

	fmt.Println("[2] Start stopping ECS Standalone Tasks")
	taskOpts := &taskOptions{
		cluster:     o.cluster,
		allClusters: o.allClusters,
		standalone:  true,
	}
	if err := taskOpts.stopTasks(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Standalone Tasks: %w", err)
	}
	fmt.Println("[2] Successfully finished stopping ECS Standalone Tasks")

	return nil
}
