package moffel

import (
	"errors"
	client "net/http"
	"net/url"
)

// Http Client
// One could argue that the client is broken for not implementing a timeout by defualt.

// http helper to do a Http request (first POST, if that doesnt work GET)
func http(eventType string, fileName string, url *url.URL) error {

	// Post req

	// Get Req
	_, err := httpGet(url, eventType, fileName)
	if err != nil {
		return err
	}

	return nil
}

func httpGet(url *url.URL, event string, file string) (string, error) {

	q := url.Query()
	q.Set("type", event)
	q.Set("file", file)
	url.RawQuery = q.Encode()

	resp, err := client.Get(url.String())
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	var out []byte
	_, e := resp.Body.Read(out)
	if e != nil {
		return "", err
	}

	return string(out), nil
}

func httpPost() {
	//
}
