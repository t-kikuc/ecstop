package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client struct {
	client *ec2.Client
}

// Create a new EC2Client with default configuration
func NewEC2Client() (*EC2Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("unable to load SDK config, %v", err)
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
