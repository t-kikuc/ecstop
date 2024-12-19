package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client struct {
	client *ec2.Client
}

// Create a new EC2Client with default configuration
func (ac AWSConfig) NewEC2Client(ctx context.Context) (*EC2Client, error) {
	cfg, err := ac.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	cli := &EC2Client{
		client: ec2.NewFromConfig(cfg),
	}

	return cli, nil
}

func (c *EC2Client) StopInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return err
	}
	return nil
}
