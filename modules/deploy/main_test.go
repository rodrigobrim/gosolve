package deploy

import (
	"config"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

var AppIdKey string = "AppID"
var AppIdValue string = "app-test"

type Ec2OptionsMock struct {
	AppID string
}

type Ec2ClientMock struct {
	t       *testing.T
	options ec2.Options
}

func TestApply(t *testing.T) {

	assert.Panics(t, func() {
		Apply(config.Config{
			Ctx:                   context.TODO(),
			Ec2OptFns:             nil,
			AWSConfig:             aws.Config{},
			ConfigOpts:            nil,
			LaunchTemplateInput:   nil,
			AutoScalingGroupInput: nil,
			AutoScalingOptFns:     nil,
		})
	})

	assert.Panics(t, func() {
		cfg := config.GetConfig("staging")
		cfg.LaunchTemplateInput.DryRun = aws.Bool(true)
		cfg.ConfigOpts = func(opt *awsConfig.LoadOptions) error {
			opt.AppID = "mocked"
			return nil
		}
		Apply(cfg)
	})

	assert.Panics(t, func() {
		cfg := config.GetConfig("staging")
		cfg.LaunchTemplateInput.DryRun = aws.Bool(true)
		cfg.ConfigOpts = func(opt *awsConfig.LoadOptions) error {
			opt.AppID = "mocked"
			return nil
		}
		Apply(cfg)
	})

	assert.Panics(t, func() {
		cfg := config.GetConfig("staging")
		cfg.LaunchTemplateInput = nil
		cfg.ConfigOpts = func(opt *awsConfig.LoadOptions) error {
			opt.AppID = "mocked-create"
			return nil
		}
		Apply(cfg)
	})

	assert.Panics(t, func() {
		cfg := config.Config{
			Ctx: context.TODO(),
			AWSConfig: aws.Config{
				AppID: "mocked",
			},
			AutoScalingGroupInput: &autoscaling.CreateAutoScalingGroupInput{AutoScalingGroupName: aws.String("name")},
			ConfigOpts: func(opt *awsConfig.LoadOptions) error {
				opt.AppID = "mocked"
				return nil
			},
			AutoScalingOptFns: func(opt *autoscaling.Options) {
				opt.AppID = "mocked"
			},
			LaunchTemplateInput: nil,
			Ec2OptFns:           nil,
		}
		Apply(cfg)
	})
}

func TestDestroy(t *testing.T) {
	assert.Panics(t, func() {
		Destroy(config.Config{
			Ctx:                   context.TODO(),
			Ec2OptFns:             nil,
			AWSConfig:             aws.Config{},
			ConfigOpts:            nil,
			LaunchTemplateInput:   nil,
			AutoScalingGroupInput: nil,
			AutoScalingOptFns:     nil,
		})
	})

	// assert.Panics(t, func() {
	cfg := config.GetConfig("staging")
	Destroy(config.Config{
		Ctx:        context.TODO(),
		Ec2OptFns:  cfg.Ec2OptFns,
		AWSConfig:  cfg.AWSConfig,
		ConfigOpts: cfg.ConfigOpts,
		LaunchTemplateInput: &ec2.CreateLaunchTemplateInput{
			LaunchTemplateName: aws.String("LaunchTemplateName"),
			LaunchTemplateData: &types.RequestLaunchTemplateData{
				ImageId:        cfg.LaunchTemplateInput.LaunchTemplateData.ImageId,
				InstanceType:   types.InstanceTypeT2Nano,
				KeyName:        aws.String("default"),
				SecurityGroups: []string{"default"},
			},
		},
		AutoScalingGroupInput: nil,
		AutoScalingOptFns:     nil,
	})
	// })

	assert.Panics(t, func() {
		Destroy(config.Config{
			Ctx: context.TODO(),
			AWSConfig: aws.Config{
				AppID: "mocked",
			},
			AutoScalingGroupInput: &autoscaling.CreateAutoScalingGroupInput{AutoScalingGroupName: aws.String("name")},
			ConfigOpts: func(opt *awsConfig.LoadOptions) error {
				opt.AppID = "mocked"
				return nil
			},
			AutoScalingOptFns: func(opt *autoscaling.Options) {
				opt.AppID = "mocked"
			},
			LaunchTemplateInput: nil,
			Ec2OptFns:           nil,
		})
	})

	assert.Panics(t, func() {
		Destroy(config.Config{
			Ctx: context.TODO(),
			AWSConfig: aws.Config{
				AppID: "mocked",
			},
			AutoScalingGroupInput: &autoscaling.CreateAutoScalingGroupInput{AutoScalingGroupName: aws.String("mocked")},
			ConfigOpts: func(opt *awsConfig.LoadOptions) error {
				opt.AppID = "mocked"
				return nil
			},
			AutoScalingOptFns: func(opt *autoscaling.Options) {
				opt.AppID = "mocked"
			},
			LaunchTemplateInput: nil,
			Ec2OptFns:           nil,
		})
	})

	assert.Panics(t, func() {
		Destroy(config.Config{
			Ctx: context.TODO(),
			AWSConfig: aws.Config{
				AppID: "mocked",
			},
			AutoScalingGroupInput: &autoscaling.CreateAutoScalingGroupInput{AutoScalingGroupName: aws.String("mocked-error")},
			ConfigOpts: func(opt *awsConfig.LoadOptions) error {
				opt.AppID = "mocked"
				return nil
			},
			AutoScalingOptFns: func(opt *autoscaling.Options) {
				opt.AppID = "mocked"
			},
			LaunchTemplateInput: nil,
			Ec2OptFns:           nil,
		})
	})
}
