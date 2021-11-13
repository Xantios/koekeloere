package moffel

// @TODO: Parse out {{EVENT_NAME}} and {{FILE_NAME}} as variables, so we can replace them
// @TODO: support "exec" driver or something, so you can use it for other things than webhooks only

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"

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
var drivers = map[string]interface{}{
	"http":  http,
	"https": https,
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

	if len(uris) <= 0 {
		log.Warn("No URI's defined, logging to CLI only")
	}

	for _, uri := range uris {
		client, err := parseUri(uri)

		if err != nil {
			log.Errorf("Cant parse %s")
			continue
		}

		if drivers[client.Name] == nil {
			log.Warnf("There is no handler for %s, feel free to PR one", client.Name)
			continue
		}

		clients = append(clients, client)
	}
}

func GetClients() []MoffelClient {
	return clients
}

func parseUri(uri string) (MoffelClient, error) {

	c := MoffelClient{}
	parsedUrl, err := url.Parse(uri)

	if err != nil {
		return c, err
	}

	port, err := strconv.Atoi(parsedUrl.Port())
	if err != nil || port == 0 {
		if parsedUrl.Scheme == "http" {
			port = 80
		}

		if parsedUrl.Scheme == "https" {
			port = 443
		}
	}

	c.Name = parsedUrl.Scheme
	c.Protocol = parsedUrl.Scheme
	c.Server = parsedUrl.Host
	c.Port = port
	c.Path = parsedUrl.Path
	c.Query = "" // @todo: pull from parsedUrl.Query()

	return c, nil
}

func Emit(event string, filename string) {
	log.Infof("Stuff: %s %s\n", event, filename)

	for _, client := range clients {
		__call(client.Name, client.Server, 1)
	}
}

func __call(funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(drivers[funcName])

	if len(params) != f.Type().NumIn() {
		err = errors.New("number of params is out of index")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	var res []reflect.Value = f.Call(in)
	result = res[0].Interface()
	return
}
