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

func (c *ecsClient) ListServices(ctx context.Context, cluster string) (seviceArns []string, err error) {
	out, err := c.client.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}
	return out.ServiceArns, nil
}

func (c *ecsClient) DescribeServices(ctx context.Context, cluster string) ([]types.Service, error) {
	serviceArns, err := c.ListServices(ctx, cluster)
	if err != nil {
		return nil, err
	}

	var services []types.Service
	for i := 0; i < len(serviceArns); i += 10 {
		end := i + 10
		if end > len(serviceArns) {
			end = len(serviceArns)
		}

		out, e := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: serviceArns[i:end],
		})
		if e != nil {
			return nil, e
		}
		services = append(services, out.Services...)
	}
	return services, nil
}

func (c *ecsClient) ScaleinService(ctx context.Context, cluster string, service string) error {
	in := &ecs.UpdateServiceInput{
		Cluster:      aws.String(cluster),
		Service:      aws.String(service),
		DesiredCount: aws.Int32(0),
	}

	_, e := c.client.UpdateService(ctx, in)
	if e != nil {
		return e
	}

	return nil
}

func (c *ecsClient) ListStandaloneTasks(ctx context.Context, cluster string) (taskArns []string, err error) {
	out, err := c.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster: aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}
	return out.TaskArns, nil
}
