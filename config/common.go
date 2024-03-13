package config

var appName string = "MyAwesomeApp"

var MyGreatApp = map[string]string{
	"AppID":                appName,
	"LaunchTemplateName":   appName + "-launchTemplate",
	"AutoScalingGroupName": appName + "-autoscaling",
}
