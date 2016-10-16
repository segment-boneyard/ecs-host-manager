package main

import (
	"flag"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/segmentio/ecs-logs-go/apex"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var debug bool = terminal.IsTerminal(0)
	var ec2InstanceID string
	var err error

	flag.BoolVar(&debug, "debug", debug, "turn on debug mode")
	flag.Parse()

	if debug {
		log.SetHandler(text.New(os.Stderr))
	} else {
		log.SetHandler(apex_ecslogs.NewHandler(os.Stderr))
	}

	if ec2InstanceID, err = getInstanceID(); err != nil {
		log.Fatalf("failed to retrieve EC2 instance ID: %s", err)
	}

	log.WithFields(log.Fields{
		"ec2-instance-id": ec2InstanceID,
	}).Info("started")

	for range time.Tick(1 * time.Minute) {

	}
}
