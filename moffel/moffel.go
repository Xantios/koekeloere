package moffel

// @TODO: Parse out {{EVENT_NAME}} and {{FILE_NAME}} as variables, so we can replace them
// @TODO: support "exec" driver or something, so you can use it for other things than webhooks only

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type MoffelClient struct {
	Name     string
	Protocol string
	Server   string
	Port     int
	Path     string
	Query    string
}

var clients []MoffelClient
var log logrus.Logger

// @TODO: Use this when sending webhooks so the CLI user knows whats happening
var verbose bool = false

// @TODO: Add std support for Slack/Discord maybe?
var protocols = []string{
	"http",
	"https",
}

func SetVerbose(verb *bool) {
	if *verb {
		verbose = true
	}
}

func SetLogger(instance *logrus.Logger) {
	log = *instance
}

func Init(uris []string) {
	for _, uri := range uris {
		client, err := parseUri(uri)

		if err != nil {
			log.Errorf("Cant parse %s")
		}

		clients = append(clients, client)
	}
}

func GetClients() []MoffelClient {
	return clients
}

func parseUri(uri string) (MoffelClient, error) {

	c := MoffelClient{}
	protocol := strings.SplitN(uri, ":", 1)[0]
	server := strings.SplitN(uri, ":", 1)[1]

	if !stringInStringSlice(protocol, protocols) {
		log.Errorf("Cant register protocol %s.\n", protocols)
		return c, errors.New("cant register protocol")
	}

	if c.Protocol == "http" || c.Name == "https" {
		c.Name = c.Protocol
	}

	if c.Protocol == "http" {
		c.Port = 80
	}

	if c.Protocol == "https" {
		c.Port = 443
	}

	c.Server = server

	return c, nil
}

func stringInStringSlice(q string, slices []string) bool {
	for _, slice := range slices {
		if slice == q {
			return true
		}
	}

	return false
}

func Emit(event string, filename string) {
	log.Infof("Stuff: %s %s\n", event, filename)
}
