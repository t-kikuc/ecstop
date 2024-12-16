package stop

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"

	"github.com/t-kikuc/ecstop/src/client"
)

type taskOptions struct {
	cluster     string
	allClusters bool

	group       string
	groupPrefix string
	standalone  bool
	// TODO: Add --all-group? (dangerous...)
}

func NewStopTaskCommand() *cobra.Command {
	o := &taskOptions{}

	c := &cobra.Command{
		Use:   "tasks",
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
	c.Flags().BoolVar(&o.allClusters, flag_allClusters, false, "Scale-in tasks of all clusters in the region")

	c.MarkFlagsOneRequired(flag_cluster, flag_allClusters)
	c.MarkFlagsMutuallyExclusive(flag_cluster, flag_allClusters)

	// Group
	c.Flags().StringVar(&o.group, flag_group, "", "Group name to scale-in tasks")
	c.Flags().StringVar(&o.groupPrefix, flag_groupPrefix, "", "Group name prefix to scale-in tasks")
	c.Flags().BoolVar(&o.standalone, flag_standalone, false, "Scale-in standalone tasks")

	c.MarkFlagsOneRequired(flag_group, flag_groupPrefix, flag_standalone)
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
	if len(tasks) == 0 {
		fmt.Printf("[%s] No tasks found in cluster\n", cluster)
		return nil
	}

	matchedTasks := o.filterByGroup(tasks)
	printPreSummaryTask(cluster, tasks, matchedTasks)

	for _, task := range matchedTasks {
		if err = cli.StopTask(ctx, cluster, *task.TaskArn); err != nil {
			return err
		}
		fmt.Printf(" -> Successfully stopped Task: %s\n", *task.TaskArn)
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

		// C. Standalone Tasks. The group of service tasks is "service:<service-name>".
		if o.standalone && !strings.HasPrefix(*task.Group, "service:") {
			filtered = append(filtered, task)
			continue
		}
	}
	return filtered
}
func printPreSummaryTask(cluster string, all, matched []types.Task) {
	fmt.Printf("[%s] All Tasks: %d, Tasks to stop: %d\n", cluster, len(all), len(matched))
	if len(matched) > 0 {
		fmt.Println("Tasks to stop:")
		for i, task := range matched {
			fmt.Printf(" [%d] Group: %s, Arn: %s\n", i+1, *task.Group, *task.TaskArn)
		}
	} else {
		fmt.Println(" -> No tasks to stop")
	}
}
