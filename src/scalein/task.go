package scalein

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"

	"github.com/t-kikuc/ecscale0/src/client"
)

type taskOptions struct {
	cluster     string
	allClusters bool

	group       string
	groupPrefix string
	standalone  bool
}

func NewStopTaskCommand() *cobra.Command {
	o := &taskOptions{}

	c := &cobra.Command{
		Use:   "task",
		Short: "Stop ECS Tasks",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stopTasks(context.Background())
		},
	}

	const (
		flag_cluster     = "cluster"
		flag_allClusters = "all-clusters"
		flag_group       = "group"
		flag_groupPrefix = "group-prefix"
		flag_standalone  = "standalone"
	)

	// Cluster
	c.Flags().StringVar(&o.cluster, flag_cluster, "", "Cluster name/arn to scale-in tasks")
	c.Flags().BoolVar(&o.allClusters, flag_allClusters, false, "Scale-in services of all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)

	// Group
	c.Flags().StringVar(&o.group, flag_group, "", "Group name to scale-in tasks")
	c.Flags().StringVar(&o.groupPrefix, flag_groupPrefix, "", "Group name prefix to scale-in tasks")
	c.Flags().BoolVar(&o.standalone, flag_standalone, false, "Scale-in standalone tasks")

	c.MarkFlagsMutuallyExclusive(flag_group, flag_groupPrefix, flag_standalone)

	return c
}

func (o *taskOptions) stopTasks(ctx context.Context) error {
	cli, err := client.NewDefaultClient()
	if err != nil {
		return err
	}

	if o.allClusters {
		return o.stopTasksInClusters(ctx, cli)
	} else {
		return o.stopTasksInCluster(ctx, cli, o.cluster)
	}
}

func (o *taskOptions) stopTasksInClusters(ctx context.Context, cli client.ECSClient) error {
	clusters, err := cli.ListClusters(ctx)
	if err != nil {
		return err
	}
	for _, cluster := range clusters {
		if err = o.stopTasksInCluster(ctx, cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func (o *taskOptions) stopTasksInCluster(ctx context.Context, cli client.ECSClient, cluster string) error {
	tasks, err := cli.DescribeTasks(ctx, cluster)
	if err != nil {
		return err
	}

	tasks = o.filterByGroup(tasks)

	for _, task := range tasks {
		if err = cli.StopTask(ctx, cluster, *task.TaskArn); err != nil {
			return err
		}
	}
	return nil
}

func (o *taskOptions) filterByGroup(tasks []types.Task) []types.Task {
	var filtered []types.Task
	for _, task := range tasks {
		// A. Complete Match
		if len(o.group) > 0 && *task.Group == o.group {
			filtered = append(filtered, task)
			continue
		}
		// B. Prefix
		if len(o.groupPrefix) > 0 && strings.HasPrefix(*task.Group, o.groupPrefix) {
			filtered = append(filtered, task)
			continue
		}

		// C. Standalone Tasks. Service Tasks have "service:<service-name>" prefix.
		if o.standalone && !strings.HasPrefix(*task.Group, "service:") {
			filtered = append(filtered, task)
			continue
		}
	}
	return filtered
}
