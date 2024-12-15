package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func (c *ecsClient) ListServices(ctx context.Context, cluster string) (seviceArns []string, e error) {
	in := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	out, e := c.client.ListServices(ctx, in)
	if e != nil {
		return nil, e
	}
	return out.ServiceArns, nil
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
