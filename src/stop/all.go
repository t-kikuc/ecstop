package stop

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/flag"
)

type allOptions struct {
	cluster     string
	allClusters bool
}

func NewStopAllCommand() *cobra.Command {
	o := &allOptions{}

	c := &cobra.Command{
		Use:   "all",
		Short: "Stop ECS Services, Standalone Tasks, and Container Instances",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	flag.AddClusterFlags(c, &o.cluster, &o.allClusters)

	return c
}

func (o *allOptions) stop(ctx context.Context) error {
	fmt.Println("[1] Start stopping ECS Services")
	srvOpts := &serviceOptions{
		cluster:     o.cluster,
		allClusters: o.allClusters,
	}
	if err := srvOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Services: %w", err)
	}
	fmt.Println("[1] Successfully finished stopping ECS Services")

	fmt.Println("[2] Start stopping ECS Standalone Tasks")
	taskOpts := &taskOptions{
		cluster:     o.cluster,
		allClusters: o.allClusters,
		standalone:  true,
	}
	if err := taskOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Standalone Tasks: %w", err)
	}
	fmt.Println("[2] Successfully finished stopping ECS Standalone Tasks")

	fmt.Println("[3] Start stopping ECS Container Instances")
	instOpts := &instanceOptions{
		cluster:     o.cluster,
		allClusters: o.allClusters,
	}
	if err := instOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Container Instances: %w", err)
	}
	fmt.Println("[3] Successfully finished stopping ECS Container Instances")

	return nil
}
