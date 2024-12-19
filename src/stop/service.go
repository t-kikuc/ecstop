package stop

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/client"
	"github.com/t-kikuc/ecstop/src/flag"
)

// serviceOptions is the options for scaling-in ECS services
type serviceOptions struct {
	cluster     string
	allClusters bool

	awsConfig client.AWSConfig
}

func NewStopServiceCommand() *cobra.Command {
	o := &serviceOptions{}

	c := &cobra.Command{
		Use:   "services",
		Short: "Scale-in ECS Services by updating desiredCount to 0",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	flag.AddClusterFlags(c, &o.cluster, &o.allClusters)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	return c
}

func (o *serviceOptions) stop(ctx context.Context) error {
	cli, err := o.awsConfig.NewECSClient(ctx)
	if err != nil {
		return err
	}

	if o.allClusters {
		return scaleinServicesInClusters(ctx, cli)
	} else {
		return scaleinServicesInCluster(ctx, cli, o.cluster)
	}
}

func scaleinServicesInClusters(ctx context.Context, cli *client.ECSClient) error {
	clusters, err := cli.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}
	for _, cluster := range clusters {
		if err = scaleinServicesInCluster(ctx, cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func scaleinServicesInCluster(ctx context.Context, cli *client.ECSClient, cluster string) error {
	services, e := cli.DescribeServices(ctx, cluster)
	if e != nil {
		return fmt.Errorf("failed to list services of cluster %s: %w", cluster, e)
	}
	if len(services) == 0 {
		fmt.Printf("No service found in cluster %s\n", cluster)
		return nil
	}

	runningServices := filterRunning(services)
	printPreSummary(cluster, services, runningServices)

	// Scale-in services
	for i, s := range runningServices {
		e := cli.ScaleinService(ctx, cluster, *s.ServiceName)
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
		// Sometimes RunningCount>0 although DesiredCount is already 0
		if s.DesiredCount > 0 || s.RunningCount > 0 {
			runningServices = append(runningServices, s)
		}
	}
	return runningServices
}

func printPreSummary(cluster string, services []types.Service, runningServices []types.Service) {
	total := len(services)
	running := len(runningServices)

	fmt.Printf("[%s] Total Services: %d, Running Services: %d\n", cluster, total, running)
	if running > 0 {
		fmt.Println("Running Services:")
		for i, s := range runningServices {
			fmt.Printf(" [%d]  %s:: running %d, desired: %d\n", i+1, *s.ServiceName, s.RunningCount, s.DesiredCount)
		}
	} else {
		fmt.Println(" -> No service to scale-in")
	}
}
