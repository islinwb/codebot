package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/xanzy/go-gitlab"
)

func (s *CodebotServer) handleMergeRequest(r *http.Request) {
	var mergeRequestEvent gitlab.MergeEvent
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("fail to read request body: %v", err)
		return
	}
	err = json.Unmarshal(b, &mergeRequestEvent)
	if err != nil {
		glog.Errorf("fail to unmarshal as MergeEvent: %v", err)
		return
	}
	s.TriggerJenkins(MergeRequestHookEvent, r)
}
