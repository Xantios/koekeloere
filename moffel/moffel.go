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
	Parser   *url.URL
}

var clients []MoffelClient
var log logrus.Logger

var verbose bool = false

// Add drivers here
var drivers = map[string]interface{}{
	"http":  http,
	"https": http,
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
			log.Warnf("There is no handler for \"%s\", feel free to PR one", client.Name)
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
	c.Parser = parsedUrl

	return c, nil
}

func Emit(event string, filename string) {
	log.Infof("Stuff: %s %s\n", event, filename)

	for _, client := range clients {

		if verbose {
			log.Printf("Emitting %s event <%s,%s>", client.Name, event, filename)
		}

		// Function definition should match
		// __call("MyFunction",<fileEvent>,<fileName>,[*url.URL]) error
		errInterface, callErr := __call(client.Name, event, filename, client.Parser)
		err := errInterface.(error)
		if callErr != nil {
			log.Errorf("Error while calling %s: %s", client.Name, callErr.Error())
		}

		if err != nil {
			log.Errorf("Driver error: %s", err.Error())
		} else {
			if verbose {
				log.Info("Emitted %s event OK")
			}
		}
	}
}

// Call a function by a variable eg: __call("myFunction","someParam1","someParam2")
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
