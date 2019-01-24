package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/xanzy/go-gitlab"
)

var testThisReg = regexp.MustCompile("^test this$")

type NoteHook struct {
	ObjectKind       string `json:"object_kind"`
	ProjectID        int    `json:"project_id"`
	ObjectAttributes struct {
		ID           int    `json:"id"`
		Note         string `json:"note"`
		NoteableType string `json:"noteable_type"`
		AuthorID     int    `json:"author_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		ProjectID    int    `json:"project_id"`
		Attachment   string `json:"attachment"`
		LineCode     string `json:"line_code"`
		CommitID     string `json:"commit_id"`
		NoteableID   int    `json:"noteable_id"`
		System       bool   `json:"system"`
		URL          string `json:"url"`
	} `json:"object_attributes"`
}

func (s *CodebotServer) handleNote(r *http.Request) {
	var noteHook NoteHook
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("fail to read request body: %v", err)
		return
	}
	err = json.Unmarshal(b, &noteHook)
	if err != nil {
		glog.Errorf("fail to unmarshal as NoteHook: %v", err)
		return
	}
	noteableType := noteHook.ObjectAttributes.NoteableType
	switch noteableType {
	case "MergeRequest":
		var mergeRequestNote gitlab.MergeCommentEvent
		err := json.Unmarshal(b, &mergeRequestNote)
		if err != nil {
			glog.Errorf("fail to unmarshal as MergeCommentEvent: %v", err)
			return
		}
		noteContent := mergeRequestNote.ObjectAttributes.Note
		if testThisReg.MatchString(strings.ToLower(noteContent)) {
			// trigger CI
			s.TriggerJenkins(NoteHookEvent, r.Body)
		}

	case "Issue":
	default:
		glog.Errorf("Bot don't handle %s for now", noteableType)

	}
}
