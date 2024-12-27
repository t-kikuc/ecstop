package stop

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
	instances, err := ecsCli.DescribeContainerInstances(ctx, cluster)
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		log.Printf("[%s] No instances found in cluster\n", cluster)
		return nil
	}

	toStop, stopped, toSkip := categorizeInstances(instances)
	printPreSummaryInstance(cluster, toStop, stopped, toSkip)
	if len(toStop) <= 0 {
		return nil
	}

	instanceIDs := make([]string, 0, len(toStop))
	for _, inst := range toStop {
		instanceIDs = append(instanceIDs, *inst.Ec2InstanceId)
	}

	if err := ec2Cli.StopInstances(ctx, instanceIDs); err != nil {
		return fmt.Errorf("failed to stop instances: %w", err)
	}
	log.Printf(" -> ✅Successfully stopped %d instances\n", len(instances))
	return nil
}

func printPreSummaryInstance(cluster string, toStop, notConnected, toSkip []types.ContainerInstance) {
	txt := fmt.Sprintf("[%s] Total Instances: %d (to Stop: %d, Not Connected: %d, External or something: %d)", cluster, len(toStop)+len(notConnected)+len(toSkip), len(toStop), len(notConnected), len(toSkip))
	if len(toStop) <= 0 {
		log.Printf("%s -> Not instances to stop", txt)
	} else {
		log.Println(txt)
		log.Println("Instances to stop:")
		for i, inst := range toStop {
			log.Printf(" [%d] %s\n", i+1, *inst.Ec2InstanceId)
		}
	}
}

func categorizeInstances(instances []types.ContainerInstance) (toStop, notConnected, toSkip []types.ContainerInstance) {
	for _, i := range instances {
		if isExternalInstance(i) {
			// External instances cannot be stopped.
			toSkip = append(toSkip, i)
			continue
		}

		if !i.AgentConnected {
			// Instances that are not connected to the ECS agent cannot be stopped.
			notConnected = append(notConnected, i)
			continue
		}

		// i.AgentConnected
		toStop = append(toStop, i)
	}
	return toStop, notConnected, toSkip
}

func isExternalInstance(instance types.ContainerInstance) bool {
	for _, a := range instance.Attributes {
		if *a.Name == "ecs.capability.external" {
			return true
		}
	}
	return false
}
