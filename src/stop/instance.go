package stop

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/client"
	"github.com/t-kikuc/ecstop/src/ec2client"
)

type instanceOptions struct {
	cluster     string
	allClusters bool
}

func NewStopInstanceCommand() *cobra.Command {
	o := &instanceOptions{}

	c := &cobra.Command{
		Use:   "instance",
		Short: "Stop ECS Container Instances",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stopInstances(context.Background())
		},
	}

	const (
		flag_cluster     = "cluster"
		flag_allClusters = "all-clusters"
	)

	// Cluster
	c.Flags().StringVar(&o.cluster, flag_cluster, "", "Cluster name/arn to stop instances")
	c.Flags().BoolVar(&o.allClusters, flag_allClusters, false, "Stop instances of all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)

	return c
}

func (o *instanceOptions) stopInstances(ctx context.Context) error {
	cli, err := client.NewDefaultClient()
	if err != nil {
		return err
	}

	if o.allClusters {
		return o.stopInstancesInClusters(ctx, cli)
	} else {
		return o.stopInstancesInCluster(ctx, cli, o.cluster)
	}
}

func (o *instanceOptions) stopInstancesInClusters(ctx context.Context, cli client.ECSClient) error {
	clusters, err := cli.ListClusters(ctx)
	if err != nil {
		return err
	}
	for _, cluster := range clusters {
		if err = o.stopInstancesInCluster(ctx, cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func (o *instanceOptions) stopInstancesInCluster(ctx context.Context, cli client.ECSClient, cluster string) error {
	instanceArns, err := cli.ListContainerInstances(ctx, cluster)
	if err != nil {
		return err
	}
	if len(instanceArns) == 0 {
		fmt.Printf("[%s] No instance found in cluster\n", cluster)
		return nil
	}

	printPreSummaryInstance(cluster, instanceArns)

	ec2client, err := ec2client.NewDefaultClient()
	if err != nil {
		return err
	}

	if err := ec2client.StopInstances(ctx, instanceArns); err != nil {
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
