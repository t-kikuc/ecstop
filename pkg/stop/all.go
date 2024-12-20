package stop

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/pkg/client"
)

type allOptions struct {
	cluster clusterOptions

	awsConfig client.AWSConfig
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

	addClusterFlags(c, &o.cluster)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	return c
}

func (o *allOptions) stop(ctx context.Context) error {
	log.Println("[1] Start stopping ECS Services")
	srcOpts := &serviceOptions{
		cluster:   o.cluster,
		awsConfig: o.awsConfig,
	}
	if err := srcOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Services: %w", err)
	}
	log.Println("[1] Successfully finished stopping ECS Services")

	log.Println("[2] Start stopping ECS Standalone Tasks")
	taskOpts := &taskOptions{
		cluster:    o.cluster,
		standalone: true,
		awsConfig:  o.awsConfig,
	}
	if err := taskOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Standalone Tasks: %w", err)
	}
	log.Println("[2] Successfully finished stopping ECS Standalone Tasks")

	log.Println("[3] Start stopping ECS Container Instances")
	instOpts := &instanceOptions{
		cluster:   o.cluster,
		awsConfig: o.awsConfig,
	}
	if err := instOpts.stop(ctx); err != nil {
		return fmt.Errorf("failed while stopping ECS Container Instances: %w", err)
	}
	log.Println("[3] Successfully finished stopping ECS Container Instances")

	return nil
}
