package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSClient interface {
	ListServices(ctx context.Context, cluster string) (seviceArns []string, e error)

	ScaleinService(ctx context.Context, cluster string, service string) error

	// DeleteActiveTaskDefinitions(ctx context.Context, taskDefArns []string) error
	// DeleteInactiveTaskDefinitions(ctx context.Context, taskDefArns []string) error
}

// ecsClient is a real implementation of ECSClient
type ecsClient struct {
	client *ecs.Client
}

// Create a new ECS Client with default configuration
func NewDefaultClient() (ECSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	cli := &ecsClient{
		client: ecs.NewFromConfig(cfg),
	}

	return cli, nil
}
