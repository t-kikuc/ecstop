package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

type AWSConfig struct {
	profile string
	region  string
}

func AddAWSConfigFlags(c *cobra.Command, configP *AWSConfig) {
	const (
		flag_profile = "profile"
		flag_region  = "region"
	)

	c.Flags().StringVarP(&configP.profile, flag_profile, "p", "", "AWS profile")
	c.Flags().StringVarP(&configP.region, flag_region, "r", "", "AWS region")
}

func (ac AWSConfig) loadConfig(ctx context.Context) (aws.Config, error) {
	var (
		cfg aws.Config
		err error
	)

	if ac.profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(ac.profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return cfg, err
	}

	if ac.region != "" {
		cfg.Region = ac.region
	}
	return cfg, nil
}
