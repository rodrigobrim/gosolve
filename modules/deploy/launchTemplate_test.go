package deploy

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

func (c *Ec2ClientMock) CreateLaunchTemplate(ctx context.Context, params *ec2.CreateLaunchTemplateInput, optFns ...func(*ec2.Options)) (*ec2.CreateLaunchTemplateOutput, error) {
	var meta middleware.Metadata
	meta.Set(&AppIdKey, &AppIdValue)
	var out *ec2.CreateLaunchTemplateOutput = &ec2.CreateLaunchTemplateOutput{
		LaunchTemplate: &types.LaunchTemplate{LaunchTemplateName: &AppIdValue},
		ResultMetadata: meta,
	}
	return out, nil
}

func (c *Ec2ClientMock) DeleteLaunchTemplate(ctx context.Context, params *ec2.DeleteLaunchTemplateInput, optFns ...func(*ec2.Options)) (*ec2.DeleteLaunchTemplateOutput, error) {
	var meta middleware.Metadata
	meta.Set(&AppIdKey, &AppIdValue)
	var out *ec2.DeleteLaunchTemplateOutput = &ec2.DeleteLaunchTemplateOutput{
		LaunchTemplate: &types.LaunchTemplate{LaunchTemplateName: &AppIdValue},
		ResultMetadata: meta,
	}
	return out, nil
}

func TestApplyLaunchTemplate(t *testing.T) {
	client := &Ec2ClientMock{t, ec2.Options{}}

	optFns := func(opts *ec2.Options) {
		opts.AppID = AppIdValue
	}

	output := ApplyLaunchTemplate(client, context.TODO(), loadTestLaunchTemplateInput(), optFns)
	assert.Equal(t, AppIdValue, *output.LaunchTemplate.LaunchTemplateName)
}

func TestDestroyLaunchTemplate(t *testing.T) {
	client := &Ec2ClientMock{t, ec2.Options{}}

	output := DestroyLaunchTemplate(client, context.TODO(), loadTestLaunchTemplateInput())
	assert.Equal(t, AppIdValue, *output.LaunchTemplate.LaunchTemplateName)
}

func loadTestLaunchTemplateInput() *ec2.CreateLaunchTemplateInput {
	return &ec2.CreateLaunchTemplateInput{LaunchTemplateName: &AppIdValue, DryRun: aws.Bool(true)}
}
