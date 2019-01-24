package main

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func (s *CodebotServer) TriggerJenkins(xGitlabEvent string, body io.Reader) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["X-Gitlab-Event"] = xGitlabEvent
	b, err := ioutil.ReadAll(body)
	if err != nil {
		glog.Errorf("fail to read request body: %v", err)
	}
	glog.Infof("request body: %s", string(b))
	req, err := NewRequest(http.MethodPost, s.JenkinsMap[0].JenkinsAddr, headers, body)
	if err != nil {
		glog.Errorf("fail to make new request: %v", err)
	}
	client := &http.Client{}
	statusCode, buf, _, err := DoRequest(client, req)
	if err != nil {
		glog.Errorf("fail to do request: %v", err)
	}
	glog.Infof("status code: %d; buf: %s", statusCode, string(buf))

}

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
