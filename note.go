package main

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

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

func (s *CodebotServer) handleNote(data []byte) {
	var noteHook NoteHook
	err := json.Unmarshal(data, &noteHook)
	if err != nil {
		glog.Errorf("fail to unmarshal as NoteHook: %v", err)
		return
	}
	noteableType := noteHook.ObjectAttributes.NoteableType
	switch noteableType {
	case "MergeRequest":
		data = FormatDateToRFC3399(data)
		var mergeRequestNote gitlab.MergeCommentEvent
		err := json.Unmarshal(data, &mergeRequestNote)
		if err != nil {
			glog.Errorf("fail to unmarshal as MergeCommentEvent: %v", err)
			return
		}

		mr := mergeRequestNote.MergeRequest
		mergeEvent := gitlab.MergeEvent{
			ObjectKind: "merge_request",
			User:       mergeRequestNote.User,
			ObjectAttributes: struct {
				ID              int              `json:"id"`
				TargetBranch    string           `json:"target_branch"`
				SourceBranch    string           `json:"source_branch"`
				SourceProjectID int              `json:"source_project_id"`
				AuthorID        int              `json:"author_id"`
				AssigneeID      int              `json:"assignee_id"`
				Title           string           `json:"title"`
				CreatedAt       string           `json:"created_at"` // Should be *time.Time (see Gitlab issue #21468)
				UpdatedAt       string           `json:"updated_at"` // Should be *time.Time (see Gitlab issue #21468)
				StCommits       []*gitlab.Commit `json:"st_commits"`
				StDiffs         []*gitlab.Diff   `json:"st_diffs"`
				MilestoneID     int              `json:"milestone_id"`
				State           string           `json:"state"`
				MergeStatus     string           `json:"merge_status"`
				TargetProjectID int              `json:"target_project_id"`
				Iid             int              `json:"iid"`
				Description     string           `json:"description"`
				Position        int              `json:"position"`
				LockedAt        string           `json:"locked_at"`
				UpdatedByID     int              `json:"updated_by_id"`
				MergeError      string           `json:"merge_error"`
				MergeParams     struct {
					ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
				} `json:"merge_params"`
				MergeWhenBuildSucceeds   bool               `json:"merge_when_build_succeeds"`
				MergeUserID              int                `json:"merge_user_id"`
				MergeCommitSha           string             `json:"merge_commit_sha"`
				DeletedAt                string             `json:"deleted_at"`
				ApprovalsBeforeMerge     string             `json:"approvals_before_merge"`
				RebaseCommitSha          string             `json:"rebase_commit_sha"`
				InProgressMergeCommitSha string             `json:"in_progress_merge_commit_sha"`
				LockVersion              int                `json:"lock_version"`
				TimeEstimate             int                `json:"time_estimate"`
				Source                   *gitlab.Repository `json:"source"`
				Target                   *gitlab.Repository `json:"target"`
				LastCommit               struct {
					ID        string         `json:"id"`
					Message   string         `json:"message"`
					Timestamp *time.Time     `json:"timestamp"`
					URL       string         `json:"url"`
					Author    *gitlab.Author `json:"author"`
				} `json:"last_commit"`
				WorkInProgress bool   `json:"work_in_progress"`
				URL            string `json:"url"`
				Action         string `json:"action"`
				Assignee       struct {
					Name      string `json:"name"`
					Username  string `json:"username"`
					AvatarURL string `json:"avatar_url"`
				} `json:"assignee"`
			}{
				ID:              mr.ID,
				TargetBranch:    mr.TargetBranch,
				SourceBranch:    mr.SourceBranch,
				SourceProjectID: mr.SourceProjectID,
				AuthorID:        mr.Author.ID,
				AssigneeID:      mr.Assignee.ID,
				Title:           mr.Title,
			},
		}

		// "Note Hook" fails to trigger Jenkins ?
		noteContent := mergeRequestNote.ObjectAttributes.Note
		if testThisReg.MatchString(strings.ToLower(noteContent)) {
			// trigger CI
			b, err := json.Marshal(mergeEvent)
			if err != nil {
				glog.Errorf("fail to marshal: %v", err)
			}
			s.TriggerJenkins(MergeRequestHookEvent, b)
		}

	case "Issue":
	default:
		glog.Errorf("Bot don't handle %s for now", noteableType)

	}
}
