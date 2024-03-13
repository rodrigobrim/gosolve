package deploy

import (
	"context"
	. "helpers"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/smithy-go/middleware"
)

// TODO:
// 	- Configure loadbalancer (health check, DNS, certificates)

type CreateAutoScalingGroup interface {
	CreateAutoScalingGroup(ctx context.Context, params *autoscaling.CreateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CreateAutoScalingGroupOutput, error)
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	Options() autoscaling.Options
}

func ApplyAutoScalingGroup(client CreateAutoScalingGroup, ctx context.Context, input *autoscaling.CreateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CreateAutoScalingGroupOutput, bool) {

	describeParams := autoScalingGetDescribeParams(input)
	resourceExists := PanicOnError(client.DescribeAutoScalingGroups(ctx, describeParams))
	if resourceExists == nil {
		return nil, false
	}

	obj := PanicOnError(client.CreateAutoScalingGroup(ctx, input, optFns...))

	client = autoScalingConvertMockToInterface(client)

	waiter := autoscaling.NewGroupInServiceWaiter(client)
	skipMock := 0
	_, err := waiter.WaitForOutput(
		ctx,
		describeParams,
		(5 * time.Minute),
		func(opt *autoscaling.GroupInServiceWaiterOptions) {
			opt.MinDelay = (10 * time.Second)
			opt.MaxDelay = (60 * time.Second)
			opt.Retryable = func(
				c context.Context,
				i *autoscaling.DescribeAutoScalingGroupsInput,
				o *autoscaling.DescribeAutoScalingGroupsOutput,
				e error) (bool, error) {
				if client.Options().AppID == "mocked" {
					mockOut := autoScalingMockOutput()
					o = mockOut
					if skipMock > 0 {
						o = nil
					}
					skipMock++
				}
				if o != nil {
					for _, group := range o.AutoScalingGroups {
						if len(group.Instances) == 0 {
							return true, nil
						}
					}
				}
				return false, nil
			}
		},
	)
	if err != nil {
		return nil, false
	}

	return obj, true
}

type DeleteAutoScalingGroup interface {
	DeleteAutoScalingGroup(ctx context.Context, params *autoscaling.DeleteAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteAutoScalingGroupOutput, error)
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	Options() autoscaling.Options
}

func DestroyAutoScalingGroup(client DeleteAutoScalingGroup, ctx context.Context, input *autoscaling.CreateAutoScalingGroupInput) (*autoscaling.DeleteAutoScalingGroupOutput, bool) {

	res := PanicOnError(client.DescribeAutoScalingGroups(ctx, autoScalingGetDescribeParams(input)))
	if res == nil {
		return nil, false
	}

	i := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: input.AutoScalingGroupName,
		ForceDelete:          aws.Bool(true),
	}
	obj := PanicOnError(client.DeleteAutoScalingGroup(ctx, i))

	skipMock := 0
	waiter := autoscaling.NewGroupNotExistsWaiter(autoScalingConvertMockToInterface(client))
	_, err := waiter.WaitForOutput(
		ctx,
		&autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []string{*input.AutoScalingGroupName}},
		(5 * time.Minute),
		func(opt *autoscaling.GroupNotExistsWaiterOptions) {
			opt.MinDelay = (10 * time.Second)
			opt.MaxDelay = (60 * time.Second)
			opt.Retryable = func(
				c context.Context,
				i *autoscaling.DescribeAutoScalingGroupsInput,
				o *autoscaling.DescribeAutoScalingGroupsOutput,
				e error) (bool, error) {
				if *input.AutoScalingGroupName == "mocked" {
					mockOut := autoScalingMockOutput()
					o = mockOut
					if skipMock > 0 {
						o = nil
					}
					skipMock++
				}
				if o != nil {
					for _, group := range o.AutoScalingGroups {
						if len(group.Instances) != 0 {
							return true, nil
						}
					}
				}
				return false, nil
			}
		},
	)

	if err != nil {
		return nil, false
	}

	return obj, true
}

func autoScalingGetDescribeParams(input *autoscaling.CreateAutoScalingGroupInput) *autoscaling.DescribeAutoScalingGroupsInput {

	return &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{*input.AutoScalingGroupName},
	}
}

type AutoScalingClientMock struct {
	t       *testing.T
	options autoscaling.Options
}

func autoScalingConvertMockToInterface(in interface{}) *autoscaling.Client {
	switch t := in.(type) {
	case *AutoScalingClientMock:
		return autoscaling.New(autoscaling.Options{AppID: t.options.AppID})
	case *autoscaling.Client:
		return t
	default:
		return nil
	}
}

func autoScalingMockOutput() *autoscaling.DescribeAutoScalingGroupsOutput {
	var g types.AutoScalingGroup
	g.Instances = []types.Instance{}
	return &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []types.AutoScalingGroup{g},
		NextToken:         nil,
		ResultMetadata:    middleware.Metadata{},
	}
}
