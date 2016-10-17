package main

import (
	"flag"
	"os"
	"reflect"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/segmentio/ecs-logs-go/apex"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var debug bool = terminal.IsTerminal(0)
	var instanceId string
	var region string
	var err error

	// Initialize the agent:
	// - parse command line arguments
	// - configure the global logger
	// - get the EC2 instance immutable configuration
	// - create the AWS session used to communicate with the AWS API
	flag.BoolVar(&debug, "debug", debug, "turn on debug mode")
	flag.Parse()

	if debug {
		log.SetHandler(text.New(os.Stderr))
	} else {
		log.SetHandler(apex_ecslogs.NewHandler(os.Stderr))
	}

	if instanceId, err = getInstanceId(); err != nil {
		log.Fatalf("failed to retrieve EC2 instance Id: %s", err)
	}

	if region, err = getRegion(); err != nil {
		log.Fatalf("failed to retrieve EC2 region: %s", err)
	}

	log.WithFields(log.Fields{
		"id":     instanceId,
		"region": region,
	}).Info("managing ec2 instance")

	session := session.New(&aws.Config{
		Region: aws.String(region),
	})

	ticker := time.NewTicker(10 * time.Second)

	// Run the agent:
	// - get the autoscaling configuration
	// - compare the current launch configurations and the set on the instance
	// - when the configs are different,
	//   * spawn a new host
	//   * waits for the host to start and pass health checks
	//   * move tasks off of the old EC2 instance
	//   * replace the instance in the autoscaling group
	//   * terminate
	for range ticker.C {
		var instanceInfo autoScalingInfo
		var groupInfo autoScalingInfo

		if instanceInfo, err = getAutoScalingInstanceInfo(session, instanceId); err != nil {
			log.WithError(err).Error("failed to get autoscaling instance information")
			continue
		}

		if groupInfo, err = getAutoScalingGroupInfo(session, *instanceInfo.group.AutoScalingGroupName); err != nil {
			log.WithError(err).Error("failed to get autoscaling group information")
			continue
		}

		if groupInfo.launchConfiguration == nil {
			log.WithField("group", *groupInfo.group.AutoScalingGroupName).Error("the launch configuration is nil")
			continue
		}

		if reflect.DeepEqual(instanceInfo.launchConfiguration, groupInfo.launchConfiguration) {
			continue
		}

		log.Info("launch configurations differ, starting upgrade process")

		if err = upgrade(session, instanceId, *groupInfo.group.AutoScalingGroupName); err != nil {
			log.WithError(err).Error("failed to replace the outdated EC2 instance")
		} else {
			ticker.Stop()
			ticker.C = nil // nothing to do anymore, wait for termination
		}
	}
}

func upgrade(session *session.Session, instanceId string, autoScalingGroup string) (err error) {
	// TODO:
	return
}
