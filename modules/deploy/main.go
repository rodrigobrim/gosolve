package deploy

import (
	envConfig "config"
	"fmt"
	. "helpers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func Apply(cfgInputs envConfig.Config) {

	ctx := cfgInputs.Ctx
	cfg := PanicOnError(config.LoadDefaultConfig(cfgInputs.Ctx, cfgInputs.ConfigOpts))

	// TODO
	// Set cfgInputs.NetworkingInput to ApplyNetworking(cfgInputs.NetworkingConfig)
	// Set cfgInputs.HardeningInputs to ApplyHardening(cfgInputs.HardeningConfig)
	// Set cfgInputs.Certificate to ApplySSLCertificate(cfgInputs.ApplySSLCertificateInput)
	// Set cfgInputs.MyLoadBalancerConfig.DNS to ApplyRoute53(cfgInputs.MyRoute53Config)
	// Set cfgInputs.AutoScalingGroupInput to ApplyLoadBalancer(cfgInputs.MyLoadBalancerConfig)

	if cfgInputs.LaunchTemplateInput != nil {
		ec2Client := ec2.NewFromConfig(cfg, cfgInputs.Ec2OptFns)
		fmt.Println("LaunchTemplate: creating")
		ApplyLaunchTemplate(ec2Client, ctx, cfgInputs.LaunchTemplateInput, cfgInputs.Ec2OptFns)
		fmt.Println("LaunchTemplate: created")
	}

	if cfgInputs.AutoScalingGroupInput != nil {
		autoScalingClient := autoscaling.NewFromConfig(cfg, cfgInputs.AutoScalingOptFns)
		fmt.Println("AutoScalingGroup: creating")
		_, created := ApplyAutoScalingGroup(autoScalingClient, ctx, cfgInputs.AutoScalingGroupInput, cfgInputs.AutoScalingOptFns)
		if created {
			fmt.Println("AutoScalingGroup: created")
		} else {
			fmt.Println("AutoScalingGroup: already exists")
		}
	}

	fmt.Printf("%s, all resources created\n", cfg.AppID)
}

func Destroy(cfgInputs envConfig.Config) {
	ctx := cfgInputs.Ctx
	cfg := PanicOnError(config.LoadDefaultConfig(cfgInputs.Ctx, cfgInputs.ConfigOpts))

	if cfgInputs.AutoScalingGroupInput != nil {
		autoScalingClient := autoscaling.NewFromConfig(cfg, cfgInputs.AutoScalingOptFns)
		fmt.Println("AutoScalingGroup: destroying")
		_, destroyed := DestroyAutoScalingGroup(autoScalingClient, ctx, cfgInputs.AutoScalingGroupInput)
		if destroyed {
			fmt.Println("AutoScalingGroup: destroyed")
		} else {
			fmt.Println("AutoScalingGroup: already destroyed")
		}
	}
	if cfgInputs.LaunchTemplateInput != nil {
		ec2Client := ec2.NewFromConfig(cfg, cfgInputs.Ec2OptFns)
		fmt.Println("LaunchTemplate: destroying")
		DestroyLaunchTemplate(ec2Client, ctx, cfgInputs.LaunchTemplateInput)
		fmt.Println("LaunchTemplate: destroyed")
	}

	fmt.Printf("%s, all resources destroyed\n", cfg.AppID)

	// Improve the return message (idempotence)
}
