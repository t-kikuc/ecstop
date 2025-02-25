package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSClient struct {
	client *ecs.Client
}

// Create a new ECSClient with default configuration
func (ac AWSConfig) NewECSClient(ctx context.Context) (*ECSClient, error) {
	cfg, err := ac.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	cli := &ECSClient{
		client: ecs.NewFromConfig(cfg),
	}

	return cli, nil
}

func (c *ECSClient) ListClusters(ctx context.Context) (clusterArns []string, err error) {
	in := &ecs.ListClustersInput{}
	for {
		out, err := c.client.ListClusters(ctx, in)
		if err != nil {
			return nil, err
		}
		clusterArns = append(clusterArns, out.ClusterArns...)

		if out.NextToken == nil {
			break
		}
		in.NextToken = out.NextToken
	}

	return clusterArns, nil
}

func (c *ECSClient) listServices(ctx context.Context, cluster string) (seviceArns []string, err error) {
	in := &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(100),
	}
	for {
		out, err := c.client.ListServices(ctx, in)
		if err != nil {
			return nil, err
		}
		seviceArns = append(seviceArns, out.ServiceArns...)
		if out.NextToken == nil {
			break
		}
		in.NextToken = out.NextToken
	}
	return seviceArns, nil
}

func (c *ECSClient) DescribeServices(ctx context.Context, cluster string) ([]types.Service, error) {
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

func (c *ECSClient) ScaleinService(ctx context.Context, cluster string, service string) error {
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

func (c *ECSClient) DescribeTasks(ctx context.Context, cluster string) (tasks []types.Task, err error) {
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

func (c *ECSClient) listTasks(ctx context.Context, cluster string) (taskArns []string, err error) {
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

func (c *ECSClient) StopTask(ctx context.Context, cluster, taskArn string) error {
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

func (c *ECSClient) DescribeContainerInstances(ctx context.Context, cluster string) (instances []types.ContainerInstance, err error) {
	in := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(100), // DescribeContainerInstances also supports up to 100 instances
	}
	for {
		listOut, err := c.client.ListContainerInstances(ctx, in)
		if err != nil {
			return nil, err
		}
		if len(listOut.ContainerInstanceArns) == 0 {
			break
		}

		descOut, err := c.client.DescribeContainerInstances(ctx, &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(cluster),
			ContainerInstances: listOut.ContainerInstanceArns,
		})
		if err != nil {
			return nil, err
		}

		instances = append(instances, descOut.ContainerInstances...)

		if listOut.NextToken == nil {
			break
		}
		in.NextToken = listOut.NextToken
	}
	return instances, nil
}
