package stop

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/pkg/client"
)

// serviceOptions is the options for scaling-in ECS services
type serviceOptions struct {
	cluster clusterOptions

	awsConfig client.AWSConfig
}

func NewStopServiceCommand() *cobra.Command {
	o := &serviceOptions{}

	c := &cobra.Command{
		Use:   "services",
		Short: "'services' scale-ins ECS Services by updating desiredCount to 0",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	addClusterFlags(c, &o.cluster)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	return c
}

func (o *serviceOptions) stop(ctx context.Context) error {
	cli, err := o.awsConfig.NewECSClient(ctx)
	if err != nil {
		return err
	}

	clusters, err := o.cluster.DecideClusters(ctx, cli)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		log.Println("No cluster found")
		return nil
	}

	for _, cluster := range clusters {
		if err = stopServices(ctx, cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func stopServices(ctx context.Context, cli *client.ECSClient, cluster string) error {
	services, err := cli.DescribeServices(ctx, cluster)
	if err != nil {
		return fmt.Errorf("failed to list services of cluster %s: %w", cluster, err)
	}
	if len(services) == 0 {
		log.Printf("[%s] No service found in cluster\n", cluster)
		return nil
	}

	// We don't need to scale-in services which are already stopped.
	runningServices := filterRunning(services)
	printPreSummary(cluster, services, runningServices)

	// Scale-in services
	for i, s := range runningServices {
		err := cli.ScaleinService(ctx, cluster, *s.ServiceName)
		if err != nil {
			return fmt.Errorf("failed to scale-in [%d]%s: %w", i+1, *s.ServiceName, err)
		} else {
			log.Printf(" -> successfully scaled-in [%d]%s \n", i+1, *s.ServiceName)
		}
	}

	return nil
}

// filterRunning filters running services from the given services.
func filterRunning(services []types.Service) []types.Service {
	var running []types.Service
	for _, s := range services {
		// Sometimes RunningCount>0 although DesiredCount is already 0. (e.g. while draining)
		if s.DesiredCount > 0 || s.RunningCount > 0 {
			running = append(running, s)
		}
	}
	return running
}

func printPreSummary(cluster string, services []types.Service, runningServices []types.Service) {
	total := len(services)
	running := len(runningServices)

	log.Printf("[%s] Total Services: %d, Running Services: %d", cluster, total, running)
	if running > 0 {
		log.Printf("\nRunning Services:\n")
		for i, s := range runningServices {
			log.Printf(" [%d] %s) running: %d, desired: %d\n", i+1, *s.ServiceName, s.RunningCount, s.DesiredCount)
		}
	} else {
		log.Printf(" -> No service to scale-in\n")
	}
}
