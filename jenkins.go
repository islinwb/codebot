package main

import (
	"bytes"
	"net/http"

	"github.com/golang/glog"
)

func (s *CodebotServer) TriggerJenkins(xGitlabEvent string, data []byte) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["X-Gitlab-Event"] = xGitlabEvent
	req, err := NewRequest(http.MethodPost, s.JenkinsMap[0].JenkinsAddr, headers, bytes.NewReader(data))
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
