package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/golang/glog"
)

func NewRequest(method, url string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		// req.Header.Add() ?
		req.Header.Set(k, v)
	}

	return req, nil
}

func DoRequest(client *http.Client, req *http.Request) (statusCode int, buf []byte, headers map[string][]string, err error) {
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	headers = resp.Header

	buf, err = ioutil.ReadAll(resp.Body)

	return
}

func FormatDateToRFC3399(data []byte) []byte {
	dataStr := string(data)
	UTCReg := regexp.MustCompile("[1-9][0-9]{3}-[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-6][0-9]:[0-6][0-9] UTC")
	UTCStr := UTCReg.FindAllString(dataStr, -1)
	for _, item := range UTCStr {
		t, err := dateparse.ParseLocal(item)
		if err != nil {
			glog.Errorf("fail to parse UTC date string: %v", err)
		}
		dataStr = strings.Replace(dataStr, item, t.Format(time.RFC3339), -1)
	}
	return []byte(dataStr)
}
