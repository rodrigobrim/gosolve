package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoScalingTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func loadStagingConfig() Config {

	c := &autoscaling.Client{}
	c.Options()

	appId := MyGreatApp["AppID"]
	imageId := "ami-0a8227ca4bf7367c3"
	instanceType := ec2Types.InstanceTypeT2Nano
	securityGroups := []string{"default"}
	keyName := "default"
	region := "sa-east-1"
	availabilityZones := []string{region + "a", region + "b", region + "c"}

	return Config{
		Ec2OptFns: func(opt *ec2.Options) {
			opt.AppID = appId
			opt.Region = region
		},
		AWSConfig: aws.Config{
			Region: region,
			AppID:  appId,
		},
		Ctx: context.TODO(),
		ConfigOpts: func(opt *config.LoadOptions) error {
			opt.AppID = appId
			opt.Region = region
			return nil
		},
		LaunchTemplateInput: &ec2.CreateLaunchTemplateInput{
			LaunchTemplateName: aws.String(MyGreatApp["LaunchTemplateName"]),
			LaunchTemplateData: &ec2Types.RequestLaunchTemplateData{
				ImageId:        aws.String(imageId),
				InstanceType:   instanceType,
				KeyName:        aws.String(keyName),
				SecurityGroups: securityGroups,
			},
		},
		AutoScalingGroupInput: &autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: aws.String(MyGreatApp["AutoScalingGroupName"]),
			MinSize:              aws.Int32(1),
			MaxSize:              aws.Int32(3),
			AvailabilityZones:    availabilityZones,
			LaunchTemplate: &autoScalingTypes.LaunchTemplateSpecification{
				LaunchTemplateName: aws.String(MyGreatApp["LaunchTemplateName"]),
			},
		},
		AutoScalingOptFns: func(opt *autoscaling.Options) {
			opt.AppID = appId
		},
	}

}
