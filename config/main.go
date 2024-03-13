package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Config struct {
	Ec2OptFns             func(*ec2.Options)
	AWSConfig             aws.Config
	Ctx                   context.Context
	ConfigOpts            func(*config.LoadOptions) error
	LaunchTemplateInput   *ec2.CreateLaunchTemplateInput
	AutoScalingGroupInput *autoscaling.CreateAutoScalingGroupInput
	AutoScalingOptFns     func(*autoscaling.Options)
}

func GetConfig(env string) Config {
	return loadEnvConfig(env)
}

func loadEnvConfig(env string) Config {

	if env == "staging" {
		return loadStagingConfig()
	}

	if env == "prod" {
		return loadProdConfig()
	}

	return Config{}
}
