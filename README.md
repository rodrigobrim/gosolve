# gosolve
Go-AWS Lab repository

### My toughts

I enjoyed the challenge and eventually, I spent more time developing code, than the pseudo one. I hope it doesn't count against me =D.

For unit tests, I think the best approach is achieved by mocking the HTTP requests (changing aws.HTTPClient), but I didn't domain it, so I opted for an approach I'm faster.

### Authentication

Uses the standard execution flow described [here](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/).

### Testing

#### Requirements

- Create an EC2 key with the name `Default` or set the proper name on staging/prod files for each AWS account (staging/prod).
- To connect to the created machines, the default security group in the default VPC should have the ssh ports opened.
- All modules have some unit tests. To check it, set the AWS authentication environment variables:

#### Deploy example:

Run:

```shell
export AWS_ACCESS_KEY_ID="replace-with-real-values"
export AWS_SECRET_ACCESS_KEY="replace-with-real-values"
export AWS_DEFAULT_REGION="replace-with-real-values"
go run . apply staging
go run . apply prod
go run . destroy staging
go run . destroy prod
```

> [!WARNING]  
> Don't use productive AWS accounts to test it. Please create a new VPC and set it in the config to use a shared account.


Execution example:
```
go run . apply staging
go run . apply prod
go run . destroy staging
go run . destroy prod
LaunchTemplate: creating
LaunchTemplate: created
AutoScalingGroup: creating
AutoScalingGroup: created
MyAwesomeApp, all resources created
LaunchTemplate: creating
LaunchTemplate: created
AutoScalingGroup: creating
AutoScalingGroup: created
MyAwesomeApp, all resources created
AutoScalingGroup: destroying
AutoScalingGroup: destroyed
LaunchTemplate: destroying
LaunchTemplate: destroyed
MyAwesomeApp, all resources destroyed
AutoScalingGroup: destroying
AutoScalingGroup: destroyed
LaunchTemplate: destroying
LaunchTemplate: destroyed
MyAwesomeApp, all resources destroyed
```