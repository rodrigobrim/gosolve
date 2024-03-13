package deploy

import (
	"config"
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

var AutoScalingGroupNameKey string = "AutoScalingGroupName"

func (c *AutoScalingClientMock) CreateAutoScalingGroup(ctx context.Context, params *autoscaling.CreateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CreateAutoScalingGroupOutput, error) {
	var meta middleware.Metadata
	meta.Set(&AppIdKey, &AppIdValue)
	meta.Set(&AutoScalingGroupNameKey, &AppIdValue)
	var out *autoscaling.CreateAutoScalingGroupOutput = &autoscaling.CreateAutoScalingGroupOutput{
		ResultMetadata: meta,
	}
	return out, nil
}

func (c *AutoScalingClientMock) DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	if params.AutoScalingGroupNames[0] == "mocked" {
		return autoScalingMockOutput(), nil
	}
	return nil, nil
}

func (c *AutoScalingClientMock) Options() autoscaling.Options {
	return autoscaling.Options{}
}

func (c *AutoScalingClientMock) DeleteAutoScalingGroup(ctx context.Context, params *autoscaling.DeleteAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteAutoScalingGroupOutput, error) {
	var meta middleware.Metadata
	v := reflect.ValueOf(c.options)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if typeOfS.Field(i).Name == "AppID" {
			meta.Set(typeOfS.Field(i).Name, v.Field(i).String())
			c.options.AppID = v.Field(i).String()
		}
	}
	var out *autoscaling.DeleteAutoScalingGroupOutput = &autoscaling.DeleteAutoScalingGroupOutput{
		ResultMetadata: meta,
	}
	return out, nil
}

func TestApplyAutoScalingGroup(t *testing.T) {

	client := &AutoScalingClientMock{t, autoscaling.Options{
		AppID: "mocked",
	}}
	cfg := config.GetConfig("staging")
	optFns := func(opts *autoscaling.Options) {
		opts.AppID = cfg.AWSConfig.AppID
	}

	output, created := ApplyAutoScalingGroup(client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: aws.String("mocked"),
		},
	)
	assert.True(t, created)
	assert.NotNil(t, output)
	p := output.ResultMetadata.Get(&AutoScalingGroupNameKey)
	name := p.(*string)
	assert.Equal(t, AppIdValue, *name)

	p = output.ResultMetadata.Get(&AppIdKey)
	appId := p.(*string)
	assert.Equal(t, AppIdValue, *appId)

	output, created = ApplyAutoScalingGroup(
		client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: cfg.AutoScalingGroupInput.AutoScalingGroupName},
		optFns)
	assert.False(t, created)
	assert.Nil(t, output)

	newClient := &AutoScalingClientMock{t, autoscaling.Options{
		AppID: "mocked-create",
	}}
	output, created = ApplyAutoScalingGroup(
		newClient,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: aws.String("test-brim"),
		})
	assert.False(t, created)

	client = &AutoScalingClientMock{t, autoscaling.Options{
		AppID: "mocked-error",
	}}
	output, created = ApplyAutoScalingGroup(
		client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: cfg.AutoScalingGroupInput.AutoScalingGroupName},
		optFns)
	assert.False(t, created)
	assert.Nil(t, output)

}

func TestDestroyAutoScalingGroup(t *testing.T) {
	client := &AutoScalingClientMock{t, autoscaling.Options{}}
	cfg := config.GetConfig("staging")

	output, destroyed := DestroyAutoScalingGroup(
		client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: cfg.AutoScalingGroupInput.AutoScalingGroupName},
	)
	assert.False(t, destroyed)
	assert.Nil(t, output)

	output, destroyed = DestroyAutoScalingGroup(
		client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: aws.String("mocked")},
	)
	assert.True(t, destroyed)
	assert.NotNil(t, output)

	output, destroyed = DestroyAutoScalingGroup(
		client,
		context.TODO(),
		&autoscaling.CreateAutoScalingGroupInput{
			AutoScalingGroupName: aws.String("mocked-destroyed")},
	)
	assert.False(t, destroyed)
	assert.Nil(t, output)

}

func TestAutoScalingConvertMockToInterface(t *testing.T) {

	i := autoScalingConvertMockToInterface(&autoscaling.Client{})
	assert.IsType(t, &autoscaling.Client{}, i)

	i = autoScalingConvertMockToInterface(nil)
	assert.Nil(t, i)
}
