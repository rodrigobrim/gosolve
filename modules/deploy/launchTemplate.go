package deploy

import (
	"context"
	"fmt"
	. "helpers"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// TODO:
// 	- Make Creation/Deletion idempotent
//	- Wait for resources creation

type CreateLaunchTemplate interface {
	CreateLaunchTemplate(ctx context.Context, params *ec2.CreateLaunchTemplateInput, optFns ...func(*ec2.Options)) (*ec2.CreateLaunchTemplateOutput, error)
}

func ApplyLaunchTemplate(client CreateLaunchTemplate, ctx context.Context, input *ec2.CreateLaunchTemplateInput, optFns ...func(*ec2.Options)) *ec2.CreateLaunchTemplateOutput {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Launchtemplate: doesn't exists")
		}
	}()
	return PanicOnError(client.CreateLaunchTemplate(ctx, input, optFns...))
}

type DeleteLaunchTemplate interface {
	DeleteLaunchTemplate(ctx context.Context, params *ec2.DeleteLaunchTemplateInput, optFns ...func(*ec2.Options)) (*ec2.DeleteLaunchTemplateOutput, error)
}

func DestroyLaunchTemplate(client DeleteLaunchTemplate, ctx context.Context, input *ec2.CreateLaunchTemplateInput) *ec2.DeleteLaunchTemplateOutput {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Launchtemplate: doesn't exists")
		}
	}()
	i := &ec2.DeleteLaunchTemplateInput{LaunchTemplateName: input.LaunchTemplateName}
	return PanicOnError(client.DeleteLaunchTemplate(ctx, i))
}
