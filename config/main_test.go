package config

import (
	// . "helpers"
	"reflect"
	"testing"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {

	for _, env := range []string{"staging", "prod"} {
		testConfigs(t, GetConfig(env))
	}

	assert.True(t, reflect.DeepEqual(Config{}, GetConfig("")))

}

func TestLoadEnvConfig(t *testing.T) {

	for _, env := range []string{"staging", "prod"} {
		testConfigs(t, loadEnvConfig(env))
	}

	assert.True(t, reflect.DeepEqual(Config{}, loadEnvConfig("")))

}

func testConfigs(t *testing.T, config interface{}) {
	c := config.(Config)
	assert.NotNil(t, c.Ctx)
	assert.NotNil(t, c.AWSConfig)
	assert.NotNil(t, c.ConfigOpts)
	assert.NotNil(t, c.Ec2OptFns)
	assert.NotNil(t, c.LaunchTemplateInput)
	assert.NotNil(t, c.LaunchTemplateInput.LaunchTemplateName)
	assert.NotNil(t, c.AutoScalingOptFns)
	assert.NotNil(t, c.AutoScalingGroupInput)
	assert.NotNil(t, c.AutoScalingGroupInput.AutoScalingGroupName)

	rawConfig := c.AWSConfig

	cfg, err := awsConfig.LoadDefaultConfig(c.Ctx, c.ConfigOpts)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	ec2Client := ec2.NewFromConfig(rawConfig, c.Ec2OptFns)
	assert.NotNil(t, ec2Client)
	assert.NotEmpty(t, ec2Client.Options().Region)
	assert.NotEmpty(t, ec2Client.Options().AppID)

	autoScalingClient := autoscaling.NewFromConfig(rawConfig, c.AutoScalingOptFns)
	assert.NotNil(t, autoScalingClient)
	assert.NotEmpty(t, autoScalingClient.Options().Region)
	assert.NotEmpty(t, autoScalingClient.Options().AppID)

}

func Ec2OptFns(fn func(*ec2.Options)) ec2.Options {

	return ec2.Options{}
}
