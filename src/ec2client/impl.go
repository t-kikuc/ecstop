package ec2client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (c *ec2Client) StopInstances(ctx context.Context, instanceIDs []string) error {
	_, err := c.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return err
	}
	return nil
}
