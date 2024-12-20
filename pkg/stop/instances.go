package stop

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/pkg/client"
)

type instanceOptions struct {
	cluster clusterOptions

	awsConfig client.AWSConfig
}

func NewStopInstanceCommand() *cobra.Command {
	o := &instanceOptions{}

	c := &cobra.Command{
		Use:   "instances",
		Short: "'instances' stops ECS Container Instances",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	addClusterFlags(c, &o.cluster)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	return c
}

func (o *instanceOptions) stop(ctx context.Context) error {
	ecsCli, err := o.awsConfig.NewECSClient(ctx)
	if err != nil {
		return err
	}

	clusters, err := o.cluster.DecideClusters(ctx, ecsCli)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		log.Println("No cluster found")
		return nil
	}

	ec2Cli, err := o.awsConfig.NewEC2Client(ctx)
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		if err := stopInstances(ctx, ecsCli, ec2Cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func stopInstances(ctx context.Context, ecsCli *client.ECSClient, ec2Cli *client.EC2Client, cluster string) error {
	instanceArns, err := ecsCli.ListContainerInstances(ctx, cluster)
	if err != nil {
		return err
	}
	if len(instanceArns) == 0 {
		log.Printf("[%s] No instance found in cluster\n", cluster)
		return nil
	}

	printPreSummaryInstance(cluster, instanceArns)

	if err := ec2Cli.StopInstances(ctx, instanceArns); err != nil {
		return fmt.Errorf("failed to stop instances: %w", err)
	}
	log.Printf(" -> Successfully stopped %d instances\n", len(instanceArns))
	return nil
}

func printPreSummaryInstance(cluster string, instanceArns []string) {
	log.Printf("[%s] Instances: %d\n", cluster, len(instanceArns))
	if len(instanceArns) > 0 {
		log.Println("Instances to stop:")
		for i, inst := range instanceArns {
			log.Printf(" [%d] %s\n", i+1, inst)
		}
	} else {
		log.Println(" -> No instance to stop")
	}
}
