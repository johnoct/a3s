package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSClient struct {
	Config  aws.Config
	Profile string
	Region  string
}

func New(ctx context.Context, profile, region string) (*AWSClient, error) {
	var opts []func(*config.LoadOptions) error

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &AWSClient{
		Config:  cfg,
		Profile: profile,
		Region:  cfg.Region,
	}, nil
}

func (c *AWSClient) SwitchProfile(ctx context.Context, profile string) error {
	newClient, err := New(ctx, profile, c.Region)
	if err != nil {
		return err
	}
	*c = *newClient
	return nil
}

func (c *AWSClient) SwitchRegion(ctx context.Context, region string) error {
	newClient, err := New(ctx, c.Profile, region)
	if err != nil {
		return err
	}
	*c = *newClient
	return nil
}