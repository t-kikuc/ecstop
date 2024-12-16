package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (c *ecsClient) ListClusters(ctx context.Context) (clusterArns []string, err error) {
	out, err := c.client.ListClusters(ctx, &ecs.ListClustersInput{}) // TODO: pagination (up to 100 clusters by default)
	if err != nil {
		return nil, err
	}
	return out.ClusterArns, nil
}

func (c *ecsClient) listServices(ctx context.Context, cluster string) (seviceArns []string, err error) {
	out, err := c.client.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}
	return out.ServiceArns, nil
}

func (c *ecsClient) DescribeServices(ctx context.Context, cluster string) ([]types.Service, error) {
	serviceArns, err := c.listServices(ctx, cluster)
	if err != nil {
		return nil, err
	}

	var services []types.Service
	for i := 0; i < len(serviceArns); i += 10 {
		end := i + 10
		if end > len(serviceArns) {
			end = len(serviceArns)
		}

		out, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: serviceArns[i:end],
		})
		if err != nil {
			return nil, err
		}
		services = append(services, out.Services...)
	}
	return services, nil
}

func (c *ecsClient) ScaleinService(ctx context.Context, cluster string, service string) error {
	_, err := c.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:      aws.String(cluster),
		Service:      aws.String(service),
		DesiredCount: aws.Int32(0),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *ecsClient) DescribeTasks(ctx context.Context, cluster string) (tasks []types.Task, err error) {
	taskArns, err := c.listTasks(ctx, cluster)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(taskArns); i += 100 {
		end := i + 100
		if end > len(taskArns) {
			end = len(taskArns)
		}

		out, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Cluster: aws.String(cluster),
			Tasks:   taskArns[i:end],
		})
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, out.Tasks...)
	}

	return tasks, nil
}

func (c *ecsClient) listTasks(ctx context.Context, cluster string) (taskArns []string, err error) {
	listTaskIn := &ecs.ListTasksInput{
		Cluster: aws.String(cluster),
	}
	for {
		out, err := c.client.ListTasks(ctx, listTaskIn)
		if err != nil {
			return nil, err
		}
		taskArns = append(taskArns, out.TaskArns...)

		if out.NextToken == nil {
			break
		}
		listTaskIn.NextToken = out.NextToken
	}

	return taskArns, nil
}

func (c *ecsClient) StopTask(ctx context.Context, cluster, taskArn string) error {
	_, err := c.client.StopTask(ctx, &ecs.StopTaskInput{
		Cluster: aws.String(cluster),
		Task:    aws.String(taskArn),
		Reason:  aws.String("task was stopped by ecstop"),
	})
	if err != nil {
		return err
	}
	return nil
}
