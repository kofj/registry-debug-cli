package docker

import (
	"fmt"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/sirupsen/logrus"
)

// New new docker registry client.
func New(endpoint, username, password string, tls, insecure bool) (hub *registry.Registry, err error) {
	var url = fmt.Sprintf("http://%s", endpoint)
	if tls {
		url = fmt.Sprintf("https://%s", endpoint)
	}

	if insecure || !tls {
		hub, err = registry.NewInsecure(url, username, password)
	} else {
		hub, err = registry.New(url, username, password)
	}
	if err != nil {
		logrus.WithError(err).Error("hub failed")
		return
	}
	hub.Logf = func(format string, args ...interface{}) {
		logrus.Infof(format, args...)
	}
	return
}
