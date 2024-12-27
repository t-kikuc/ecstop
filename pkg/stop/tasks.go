package stop

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/spf13/cobra"

	"github.com/t-kikuc/ecstop/pkg/client"
)

type taskOptions struct {
	cluster clusterOptions

	group       string
	groupPrefix string
	standalone  bool
	// TODO: Add --all-group? (dangerous...)
	awsConfig client.AWSConfig
}

func NewStopTaskCommand() *cobra.Command {
	o := &taskOptions{}

	c := &cobra.Command{
		Use:   "tasks",
		Short: "'tasks' stops ECS Tasks",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.stop(context.Background())
		},
	}

	const (
		flag_group       = "group"
		flag_groupPrefix = "group-prefix"
		flag_standalone  = "standalone"
	)

	addClusterFlags(c, &o.cluster)
	client.AddAWSConfigFlags(c, &o.awsConfig)

	// Group
	c.Flags().StringVar(&o.group, flag_group, "", "Group name to stop tasks")
	c.Flags().StringVar(&o.groupPrefix, flag_groupPrefix, "", "Group name prefix to stop tasks")
	c.Flags().BoolVar(&o.standalone, flag_standalone, false, "Stop standalone tasks, whose group prefix is not 'service:'")

	c.MarkFlagsOneRequired(flag_group, flag_groupPrefix, flag_standalone)
	c.MarkFlagsMutuallyExclusive(flag_group, flag_groupPrefix, flag_standalone)

	return c
}

func (o *taskOptions) stop(ctx context.Context) error {
	cli, err := o.awsConfig.NewECSClient(ctx)
	if err != nil {
		return err
	}

	clusters, err := o.cluster.DecideClusters(ctx, cli)
	if err != nil {
		return nil
	}
	if len(clusters) == 0 {
		log.Println("No cluster found")
		return nil
	}

	for _, cluster := range clusters {
		if err = o.stopTasks(ctx, cli, cluster); err != nil {
			return err
		}
	}
	return nil
}

func (o *taskOptions) stopTasks(ctx context.Context, cli *client.ECSClient, cluster string) error {
	tasks, err := cli.DescribeTasks(ctx, cluster)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		log.Printf("[%s] No tasks found in cluster\n", cluster)
		return nil
	}

	matchedTasks := o.filterByGroup(tasks)
	printPreSummaryTask(cluster, tasks, matchedTasks)

	for _, task := range matchedTasks {
		if err = cli.StopTask(ctx, cluster, *task.TaskArn); err != nil {
			return err
		}
		log.Printf(" -> ✅Successfully stopped Task: %s\n", *task.TaskArn)
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
	txt := fmt.Sprintf("[%s] All Tasks: %d, Tasks to stop: %d", cluster, len(all), len(matched))
	if len(matched) <= 0 {
		log.Printf("%s -> No tasks to stop\n", txt)
	} else {
		log.Println(txt)
		log.Printf("\nTasks to stop:\n")
		for i, task := range matched {
			log.Printf(" [%d] Group: %s, Arn: %s\n", i+1, *task.Group, *task.TaskArn)
		}
	}
}
