package scalein

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecscale0/src/client"
)

// serviceOptions is the options for scaling-in ECS services
type serviceOptions struct {
	cluster string
}

func NewScaleinServiceCommand() *cobra.Command {
	o := &serviceOptions{}

	c := &cobra.Command{
		Use:   "service",
		Short: "Scale-in ECS Services of the cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.scaleinServices(context.Background())
		},
	}

	c.Flags().StringVar(&o.cluster, "cluster", "", "Cluster name to scale-in services")
	c.MarkFlagRequired("cluster")

	return c
}

func (o *serviceOptions) scaleinServices(ctx context.Context) error {
	cli, e := client.NewDefaultClient()
	if e != nil {
		return e
	}

	services, e := cli.DescribeServices(ctx, o.cluster)
	if e != nil {
		return fmt.Errorf("failed to list services of cluster %s: %w", o.cluster, e)
	}
	if len(services) == 0 {
		fmt.Printf("No service found in cluster %s\n", o.cluster)
		return nil
	}

	runningServices := filterRunning(services)
	printPreSummary(services, runningServices)

	// Scale-in services
	for i, s := range runningServices {
		e := cli.ScaleinService(ctx, o.cluster, *s.ServiceName)
		if e != nil {
			return fmt.Errorf("failed to scale-in [%d]%s: %w", i+1, *s.ServiceName, e)
		} else {
			fmt.Printf(" -> successfully scaled-in [%d]%s \n", i+1, *s.ServiceName)
		}
	}

	return nil
}

// filterRunning filters running services from the given services
func filterRunning(services []types.Service) []types.Service {
	var runningServices []types.Service
	for _, s := range services {
		if s.RunningCount > 0 {
			runningServices = append(runningServices, s)
		}
	}
	return runningServices
}

func printPreSummary(services []types.Service, runningServices []types.Service) {
	total := len(services)
	running := len(runningServices)

	fmt.Printf("Total Services: %d, Running Services: %d\n", total, running)
	if running > 0 {
		fmt.Println("Running Services:")
		for i, s := range runningServices {
			fmt.Printf(" [%d]  %s:: running %d, desired: %d\n", i+1, *s.ServiceName, s.RunningCount, s.DesiredCount)
		}
	} else {
		fmt.Println("No service to scale-in")
	}
}
