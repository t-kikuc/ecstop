package stop

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/client"
	"github.com/t-kikuc/ecstop/src/flag"
)

type instanceOptions struct {
	cluster     string
	allClusters bool

	awsConfig client.AWSConfig
}

func NewStopInstanceCommand() *cobra.Command {
	o := &instanceOptions{}

	c := &cobra.Command{
		Use:   "instances",
		Short: "Stop ECS Container Instances",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	flag.AddClusterFlags(c, &o.cluster, &o.allClusters)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	return c
}

func (o *instanceOptions) stop(ctx context.Context) error {
	ecsCli, err := o.awsConfig.NewECSClient(ctx)
	if err != nil {
		return err
	}

	if o.allClusters {
		return o.stopInClusters(ctx, ecsCli)
	} else {
		return o.stopInstances(ctx, ecsCli, o.cluster)
	}
}

func (o *instanceOptions) stopInClusters(ctx context.Context, ecsCli *client.ECSClient) error {
	clusters, err := ecsCli.ListClusters(ctx)
	if err != nil {
		return err
	}
	for _, cluster := range clusters {
		if err = o.stopInstances(ctx, ecsCli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func (o *instanceOptions) stopInstances(ctx context.Context, ecsCli *client.ECSClient, cluster string) error {
	instanceArns, err := ecsCli.ListContainerInstances(ctx, cluster)
	if err != nil {
		return err
	}
	if len(instanceArns) == 0 {
		fmt.Printf("[%s] No instance found in cluster\n", cluster)
		return nil
	}

	printPreSummaryInstance(cluster, instanceArns)

	ec2cli, err := o.awsConfig.NewEC2Client(ctx)
	if err != nil {
		return err
	}

	if err := ec2cli.StopInstances(ctx, instanceArns); err != nil {
		return fmt.Errorf("failed to stop instances: %w", err)
	}
	fmt.Printf(" -> Successfully stopped %d instances\n", len(instanceArns))
	return nil
}

func printPreSummaryInstance(cluster string, instanceArns []string) {
	fmt.Printf("[%s] Instances: %d\n", cluster, len(instanceArns))
	if len(instanceArns) > 0 {
		fmt.Println("Instances to stop:")
		for i, inst := range instanceArns {
			fmt.Printf(" [%d] %s\n", i+1, inst)
		}
	} else {
		fmt.Println(" -> No instance to stop")
	}
}
