package ec2client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client interface {
	StopInstances(ctx context.Context, instanceIDs []string) error
}

// ec2Client is a real implementation of EC2Client
type ec2Client struct {
	client *ec2.Client
}

// Create a new EC2 Client with default configuration
func NewDefaultClient() (EC2Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	cli := &ec2Client{
		client: ec2.NewFromConfig(cfg),
	}

	return cli, nil
}
