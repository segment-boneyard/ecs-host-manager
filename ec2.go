package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func getInstanceID() (id string, err error) {
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
