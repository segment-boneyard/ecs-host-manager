package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type autoScalingInfo struct {
	group               *autoscaling.Group
	launchConfiguration *autoscaling.LaunchConfiguration
}

var (
	errNoAutoScalingInstance = errors.New("the instance is not part of an autoscaling group")
	errNoAutoScalingGroup    = errors.New("the autoscaling group doesn't exist")
)

func getAutoScalingInstanceInfo(session *session.Session, instanceId string) (info autoScalingInfo, err error) {
	var client = autoscaling.New(session)
	var output *autoscaling.DescribeAutoScalingInstancesOutput

	if output, err = client.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{aws.String(instanceId)},
	}); err != nil {
		return
	}

	if len(output.AutoScalingInstances) == 0 {
		err = errNoAutoScalingInstance
		return
	}

	if info, err = getAutoScalingGroupInfo(
		session,
		*output.AutoScalingInstances[0].AutoScalingGroupName,
	); err != nil {
		return
	}

	// It's fine to error, we'll assume it doesn't exist anymore.
	info.launchConfiguration, _ = getLaunchConfigurationByName(
		session,
		*output.AutoScalingInstances[0].LaunchConfigurationName,
	)
	return
}

func getAutoScalingGroupInfo(session *session.Session, group string) (info autoScalingInfo, err error) {
	var client = autoscaling.New(session)
	var output *autoscaling.DescribeAutoScalingGroupsOutput

	if output, err = client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(group)},
	}); err != nil {
		return
	}

	if len(output.AutoScalingGroups) == 0 {
		err = errNoAutoScalingGroup
		return
	}

	info.group = output.AutoScalingGroups[0]

	// It's fine to error, we'll assume it doesn't exist anymore.
	info.launchConfiguration, _ = getLaunchConfigurationByName(
		session,
		*info.group.LaunchConfigurationName,
	)
	return
}

func getLaunchConfigurationByName(session *session.Session, name string) (launchConfiguration *autoscaling.LaunchConfiguration, err error) {
	var client = autoscaling.New(session)
	var output *autoscaling.DescribeLaunchConfigurationsOutput

	if output, err = client.DescribeLaunchConfigurations(&autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{aws.String(name)},
	}); err != nil {
		return
	}

	launchConfiguration = output.LaunchConfigurations[0]
	return
}
