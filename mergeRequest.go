package main

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/xanzy/go-gitlab"
)

func (s *CodebotServer) handleMergeRequest(data []byte) {
	var mergeRequestEvent gitlab.MergeEvent
	data = FormatDateToRFC3399(data)
	err := json.Unmarshal(data, &mergeRequestEvent)
	if err != nil {
		glog.Errorf("fail to unmarshal as MergeEvent: %v", err)
		return
	}
	s.TriggerJenkins(MergeRequestHookEvent, data)
}
