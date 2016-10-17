package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func getInstanceId() (id string, err error) {
	var res *http.Response

	if res, err = http.Get("http://169.254.169.254/latest/meta-data/instance-id"); err != nil {
		return
	}

	defer res.Body.Close()
	var b []byte

	if b, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}

	id = strings.TrimSpace(string(b))
	return
}

func getRegion() (region string, err error) {
	var az string

	if az, err = getAvailabilityZone(); err != nil {
		return
	}

	region = strings.TrimRightFunc(az, unicode.IsLetter)
	return
}

func getAvailabilityZone() (az string, err error) {
	var res *http.Response

	if res, err = http.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone"); err != nil {
		return
	}

	defer res.Body.Close()
	var b []byte

	if b, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}

	az = strings.TrimSpace(string(b))
	return
}

func launchInstance(launchConfiguration string) (instance *ec2.Instance, err error) {
	return
}
